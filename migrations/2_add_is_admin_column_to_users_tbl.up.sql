-- Добавляем столбец is_admin в таблицу пользователей
ALTER TABLE users
    ADD COLUMN is_admin BOOLEAN NOT NULL DEFAULT FALSE;
