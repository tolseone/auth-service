-- Удаляем столбец is_admin из таблицы пользователей
ALTER TABLE users DROP COLUMN IF EXISTS is_admin;
