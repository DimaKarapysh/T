-- +goose Up
-- +goose StatementBegin
INSERT INTO tasks (title, description, status, created_at, updated_at)
VALUES
    ('Купить продукты', 'Сходить в магазин и купить хлеб и молоко', 'new', now(), now()),
    ('Сделать тестовое задание', 'Реализовать TODO сервис на Go + Fiber', 'in_progress', now(), now()),
    ('Позаниматься спортом', 'Пробежать 5 км вечером', 'done', now() - interval '2 days', now() - interval '1 day'),
    ('Прочитать книгу', 'Закончить читать «Clean Architecture»', 'new', now(), now()),
    ('Позвонить родителям', 'Связаться по видеозвонку', 'in_progress', now(), now());
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DELETE FROM tasks
WHERE title IN (
                'Купить продукты',
                'Сделать тестовое задание',
                'Позаниматься спортом',
                'Прочитать книгу',
                'Позвонить родителям'
    );
-- +goose StatementEnd
