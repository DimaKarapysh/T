-- +goose Up
-- +goose StatementBegin
INSERT INTO subscriptions (id, service_name, price, user_id, start_date, end_date)
VALUES ('11111111-1111-1111-1111-111111111111', 'Yandex Plus', 299, 'aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaaa',
        '2024-01-01', NULL),
       ('22222222-2222-2222-2222-222222222222', 'Netflix', 899, 'bbbbbbbb-bbbb-bbbb-bbbb-bbbbbbbbbbbb', '2024-03-01',
        NULL),
       ('33333333-3333-3333-3333-333333333333', 'Spotify', 169, 'cccccccc-cccc-cccc-cccc-cccccccccccc', '2023-12-01',
        '2024-06-01'),
       ('44444444-4444-4444-4444-444444444444', 'IVI', 399, 'dddddddd-dddd-dddd-dddd-dddddddddddd', '2024-04-01', NULL),
       ('55555555-5555-5555-5555-555555555555', 'Kinopoisk', 399, 'eeeeeeee-eeee-eeee-eeee-eeeeeeeeeeee', '2024-02-01',
        '2024-08-01');

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DELETE
FROM subscriptions
WHERE id = '11111111-1111-1111-1111-111111111111'
   OR id = '22222222-2222-2222-2222-222222222222'
   OR id = '33333333-3333-3333-3333-333333333333'
   OR id = '44444444-4444-4444-4444-444444444444'
   OR id = '55555555-5555-5555-5555-555555555555';
-- +goose StatementEnd
