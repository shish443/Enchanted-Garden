--Enchanted-Garden/db/migrations/00001_init_schema.sql 
-- +goose Up
-- +goose StatementBegin
CREATE TABLE branches (
    id BIGSERIAL PRIMARY KEY,
    name VARCHAR(200) NOT NULL, -- Ограничение длины строки до 200 символов для защиты от овертайпа
    parent_id BIGINT,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    CONSTRAINT fk_branches_parent FOREIGN KEY(parent_id) REFERENCES branches(id) ON DELETE CASCADE -- Ссылается на родительский ID; поддерживает дерево сущностей или NULL
);

-- Индекс для уникальности имен корневых веток (где parent_id IS NULL)
CREATE UNIQUE INDEX idx_branches_unique_root_name ON branches (name) WHERE parent_id IS NULL;

-- Индекс для уникальности имени внутри одного родительского подразделения
CREATE UNIQUE INDEX idx_branches_name_parent_id_unique ON branches (name, parent_id) WHERE parent_id IS NOT NULL;

CREATE TABLE floras (
    id BIGSERIAL PRIMARY KEY,
    branch_id BIGINT NOT NULL,
    full_name VARCHAR(200) NOT NULL,
    position VARCHAR(200) NOT NULL,
    hired_at DATE,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    CONSTRAINT fk_floras_branch FOREIGN KEY (branch_id) REFERENCES branches(id) ON DELETE CASCADE -- Каскадное удаление сотрудников при ликвидации ветки
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS floras;
DROP TABLE IF EXISTS branches;
-- +goose StatementEnd