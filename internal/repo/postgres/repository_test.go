package postgres_test

import (
	"T/migrations"
	"context"
	"database/sql"
	"fmt"
	"os"
	"testing"
	"time"

	"T/internal/entity"
	repopg "T/internal/repo/postgres"
	"T/internal/repo/postgres/sqlc"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	_ "github.com/lib/pq"
	"github.com/pressly/goose/v3"
	"go.opentelemetry.io/otel/trace/noop"
	"go.uber.org/zap"
)

func setup(t *testing.T) (context.Context, *repopg.RepoSub) {
	t.Helper()

	ctx := context.Background()

	pool := mustPGXPool(t, ctx)
	t.Cleanup(func() { pool.Close() })

	mustMigrate(t)
	mustClean(t, pool)

	repo := repopg.NewRepoSub(
		sqlc.New(pool),
		zap.NewNop(),
		noop.NewTracerProvider().Tracer("test"),
	)
	return ctx, repo
}

func TestRepoSub_CRUD(t *testing.T) {
	ctx, repo := setup(t)

	// --- Create
	desc := "my desc"
	id, err := repo.Create(ctx, &entity.Task{
		Title:       "hello",
		Description: &desc,
		Status:      entity.StatusNew,
	})
	if err != nil {
		t.Fatalf("Create: %v", err)
	}
	if id == uuid.Nil {
		t.Fatalf("Create: got nil uuid")
	}

	// --- GetByID
	got, err := repo.GetByID(ctx, id)
	if err != nil {
		t.Fatalf("GetByID: %v", err)
	}
	if got.Title != "hello" {
		t.Fatalf("GetByID: title mismatch: %q", got.Title)
	}
	if got.Description == nil || *got.Description != "my desc" {
		t.Fatalf("GetByID: description mismatch: %#v", got.Description)
	}
	if got.Status != entity.StatusNew {
		t.Fatalf("GetByID: status mismatch: %q", got.Status)
	}

	// --- Update
	newDesc := "updated"
	got.Title = "hello2"
	got.Description = &newDesc
	got.Status = entity.StatusDone

	if err := repo.Update(ctx, got); err != nil {
		t.Fatalf("Update: %v", err)
	}

	got2, err := repo.GetByID(ctx, id)
	if err != nil {
		t.Fatalf("GetByID after Update: %v", err)
	}
	if got2.Title != "hello2" || got2.Description == nil || *got2.Description != "updated" || got2.Status != entity.StatusDone {
		t.Fatalf("Update check mismatch: %+v", got2)
	}

	// --- List
	list, err := repo.List(ctx, 10, 0)
	if err != nil {
		t.Fatalf("List: %v", err)
	}
	if len(list) == 0 {
		t.Fatalf("List: expected >=1, got %d", len(list))
	}

	// --- Delete
	if err := repo.Delete(ctx, id); err != nil {
		t.Fatalf("Delete: %v", err)
	}

	// после delete NotFound)
	_, err = repo.GetByID(ctx, id)
	if err == nil {
		t.Fatalf("GetByID after Delete: expected error, got nil")
	}
}

func mustPGXPool(t *testing.T, ctx context.Context) *pgxpool.Pool {
	t.Helper()

	dsn := pgDSN()

	cfg, err := pgxpool.ParseConfig(dsn)
	if err != nil {
		t.Fatalf("ParseConfig: %v", err)
	}

	pool, err := pgxpool.NewWithConfig(ctx, cfg)
	if err != nil {
		t.Fatalf("pgxpool.New: %v", err)
	}

	// ждём, пока Postgres реально примет коннекты
	waitForPostgres(t, ctx, pool)

	return pool
}

func waitForPostgres(t *testing.T, parent context.Context, pool *pgxpool.Pool) {
	t.Helper()

	ctx, cancel := context.WithTimeout(parent, 10*time.Second)
	defer cancel()

	ticker := time.NewTicker(200 * time.Millisecond)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			t.Fatalf("postgres not ready: %v", ctx.Err())
		case <-ticker.C:
			if err := pool.Ping(ctx); err == nil {
				return
			}
		}
	}
}

func mustMigrate(t *testing.T) {
	t.Helper()

	dsn := pgDSN()

	db, err := sql.Open("postgres", dsn)
	if err != nil {
		t.Fatalf("sql.Open: %v", err)
	}
	defer db.Close()

	goose.SetBaseFS(migrations.FS)
	if err := goose.Up(db, "."); err != nil {
		t.Fatalf("goose.Up: %v", err)
	}
}

func mustClean(t *testing.T, pool *pgxpool.Pool) {
	t.Helper()

	// чистим
	_, err := pool.Exec(context.Background(), `TRUNCATE TABLE tasks RESTART IDENTITY CASCADE;`)
	if err != nil {
		t.Fatalf("clean tasks: %v", err)
	}
}

func getenv(key, def string) string {
	v := os.Getenv(key)
	if v == "" {
		return def
	}
	return v
}

func pgDSN() string {
	host := getenv("POSTGRES_HOST", "localhost")
	port := getenv("POSTGRES_PORT", "5433")
	user := getenv("POSTGRES_USER", "postgres")
	pass := getenv("POSTGRES_PASSWORD", "qqqq")
	db := getenv("POSTGRES_DATABASE", "auth")

	return fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable",
		user, pass, host, port, db)
}
