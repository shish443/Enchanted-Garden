-- ну чтож практика написания SQLлок никогда не будет лишней да?
-- +goose Up
CREATE TABLE branches (
    id BIGSERIAL PRIMARY KEY,
    name VARCHAR(200) NOT NULL, -- VARCHAR(200)-говорит что нельзя записать более 200 символов
    parent_id BIGINT,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    FOREIGN KEY(parent_id) REFERENCES branches(id), --говорит может быть только ключ или пусто
    UNIQUE (parent_id,name) -- проверяет не занято ли название
);

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
