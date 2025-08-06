package redis

import (
	"T/internal/config"
	"T/internal/errwrap"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
	"time"
)

type Session struct {
	db   *redis.Client
	cfg  *config.Config
	logg *zap.Logger
}

func NewSessionRepository(cfg *config.Config, db *redis.Client, logg *zap.Logger) *Session {
	return &Session{
		db:   db,
		cfg:  cfg,
		logg: logg,
	}
}

func (s *Session) SetRefreshToken(ctx context.Context, userID uuid.UUID, token string, sessionID uuid.UUID, userAgent, ip string) (err error) {
	const op = "SetRefreshToken"
	logger := s.logg.With(
		zap.String("Operation", op),
		zap.String("user_id", userID.String()),
		zap.String("session_id", sessionID.String()),
	)

	//ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	//defer cancel()

	if userID == uuid.Nil {
		return errwrap.ErrUserIDEmpty(op)
	}

	if sessionID == uuid.Nil {
		return errwrap.ErrSessionIDEmpty(op)
	}

	if token == "" {
		return errwrap.InvalidReqError(op, "token", errors.New("token cannot be empty"))
	}

	fields := map[string]any{
		"token":      token,
		"user_agent": userAgent,
		"ip":         ip,
	}

	data, err := json.Marshal(fields)
	if err != nil {
		return errwrap.WrapWithReason(op, errwrap.CodeInternal, err, "failed_to_marshal_token_data")
	}

	logger.Info("сохранение токена")

	parsedDuration, err := time.ParseDuration(s.cfg.JWT.RefreshTTL.String())
	if err != nil {
		return errwrap.WrapWithReason(op, errwrap.CodeInvalidArgument, err, "failed_to_parse_duration")
	}

	err = s.db.Set(ctx, buildKey(userID, sessionID), data, parsedDuration).Err()
	if err != nil {
		return errwrap.WrapWithReason(op, errwrap.CodeInternal, err, "failed_to_set_token")
	}

	logger.Info("токен успешно сохранен")

	return nil

}

func (s *Session) GetRefreshToken(ctx context.Context, userID uuid.UUID, sessionID uuid.UUID) (sessionData *SessionData, err error) {
	const op = "GetRefreshToken"
	logger := s.logg.With(
		zap.String("Operation", op),
		zap.String("user_id", userID.String()),
		zap.String("session_id", sessionID.String()),
	)

	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	if userID == uuid.Nil {
		return nil, errwrap.ErrUserIDEmpty(op)
	}

	if sessionID == uuid.Nil {
		return nil, errwrap.ErrSessionIDEmpty(op)
	}

	logger.Info("получение токена")

	data, err := s.db.Get(ctx, buildKey(userID, sessionID)).Result()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return nil, errwrap.WrapWithReason(op, errwrap.CodeNotFound, errors.New("token does not exist"), "token_not_found")
		}

		return nil, errwrap.WrapErrorWithReason(op, errwrap.CodeInternal, "failed_to_get_token")
	}

	var session SessionData
	if err := json.Unmarshal([]byte(data), &session); err != nil {
		return nil, errwrap.WrapWithReason(op, errwrap.CodeInternal, err, "failed_to_unmarshal_token_data")
	}

	logger.Info("токен успешно получен")

	return &session, nil

}

func (s *Session) RevokeRefreshToken(ctx context.Context, userID uuid.UUID, sessionID uuid.UUID) (err error) {
	const op = "RevokeRefreshToken"
	logger := s.logg.With(
		zap.String("Operation", op),
		zap.String("user_id", userID.String()),
		zap.String("session_id", sessionID.String()),
	)

	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	if userID == uuid.Nil {
		return errwrap.ErrSessionIDEmpty(op)
	}

	if sessionID == uuid.Nil {
		return errwrap.ErrSessionIDEmpty(op)
	}

	logger.Info("удаление токена")

	err = s.db.Del(ctx, buildKey(userID, sessionID)).Err()
	if err != nil {
		return errwrap.WrapErrorWithReason(op, errwrap.CodeInternal, "failed_to_set_token")
	}

	logger.Info("токен успешно удален")

	return nil

}

func (s *Session) RevokeAllRefreshTokens(ctx context.Context, userID uuid.UUID) (err error) {
	const op = "RevokeAllRefreshTokens"
	logger := s.logg.With(
		zap.String("Operation", op),
		zap.String("user_id", userID.String()),
	)

	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	if userID == uuid.Nil {
		return errwrap.ErrUserIDEmpty(op)
	}

	iter := s.db.Scan(ctx, 0, fmt.Sprintf("auth_tokens:%s:*", userID.String()), 0).Iterator()

	for iter.Next(ctx) {
		err := s.db.Del(ctx, iter.Val()).Err()
		if err != nil {
			return errwrap.WrapWithReason(op, errwrap.CodeInternal, err, "failed_to_delete_tokens")
		}
	}

	logger.Info("токены успешно удалены")

	return nil

}

func buildKey(userID uuid.UUID, sessionID uuid.UUID) string {
	return "auth_tokens:" + userID.String() + ":" + sessionID.String()
}
