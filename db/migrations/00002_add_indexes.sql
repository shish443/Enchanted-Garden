--Enchanted-Garden/db/migrations/00002_add_indexes.sql 
-- Оптимизация выборок: добавляем индексы для частых операций поиска и фильтрации
-- +goose Up
-- +goose StatementBegin
CREATE INDEX idx_floras_branch_id ON floras (branch_id);
CREATE INDEX idx_floras_full_name ON floras (full_name);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP INDEX IF EXISTS idx_floras_branch_id;
DROP INDEX IF EXISTS idx_floras_full_name;
-- +goose StatementEnd