-- +goose Up
-- +goose StatementBegin
SELECT 'up SQL query';
-- +goose StatementEnd
CREATE TABLE messages (
                            id int(11) NOT NULL AUTO_INCREMENT,
                            phone_number text NOT NULL,
                            text LONGTEXT not null,
                            created_at date DEFAULT current_timestamp(),
                            updated_at date DEFAULT current_timestamp(),
                            deleted_at date DEFAULT NULL,
                            PRIMARY KEY (id)
);
-- +goose Down
-- +goose StatementBegin
SELECT 'down SQL query';
-- +goose StatementEnd
drop table messages;
