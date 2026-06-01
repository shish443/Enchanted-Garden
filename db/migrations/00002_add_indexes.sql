-- +goose Up
CREATE INDEX idx_floras_branch_id ON floras (branch_id);
CREATE INDEX idx_floras_full_name ON floras (full_name);

-- +goose Down
DROP INDEX IF EXISTS idx_floras_branch_id;
DROP INDEX IF EXISTS idx_floras_full_name;
