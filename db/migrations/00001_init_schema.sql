--Enchanted-Garden/db/migrations/00001_init_schema.sql
-- ну чтож практика написания SQLлок никогда не будет лишней да?
-- +goose Up
CREATE TABLE branches (
    id BIGSERIAL PRIMARY KEY,
    name VARCHAR(200) NOT NULL, -- VARCHAR(200)-говорит что нельзя записать более 200 символов
    parent_id BIGINT,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    FOREIGN KEY(parent_id) REFERENCES branches(id) ON DELETE CASCADE --говорит может быть только ключ или пусто
);

-- Индекс для корневых веток (где parent_id IS NULL)
CREATE UNIQUE INDEX idx_branches_unique_root_name ON branches (name) WHERE parent_id IS NULL;

-- Индекс для дочерних веток (уникальность имени внутри одного родителя)
CREATE UNIQUE INDEX idx_branches_name_parent_id_unique ON branches (name, parent_id) WHERE parent_id IS NOT NULL;

CREATE TABLE floras (
    id BIGSERIAL PRIMARY KEY,
    branch_id BIGINT NOT NULL,
    full_name VARCHAR(200) NOT NULL,
    position VARCHAR(200) NOT NULL,
    hired_at DATE,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    FOREIGN KEY (branch_id) REFERENCES branches(id) ON DELETE CASCADE -- если айди не привязан к ветке удалчем
);

-- +goose Down
DROP TABLE IF EXISTS floras;
DROP TABLE IF EXISTS branches;
