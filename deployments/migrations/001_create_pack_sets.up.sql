-- Create pack_sets table
CREATE TABLE IF NOT EXISTS pack_sets (
    id BIGSERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    sizes JSONB NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT pack_sets_name_key UNIQUE (name),
    CONSTRAINT pack_sets_sizes_check CHECK (jsonb_typeof(sizes) = 'array')
);

-- Create index on name for faster lookups
CREATE INDEX idx_pack_sets_name ON pack_sets(name);

-- Create index on created_at for sorting
CREATE INDEX idx_pack_sets_created_at ON pack_sets(created_at DESC);

-- Add comment to table
COMMENT ON TABLE pack_sets IS 'Хранит наборы размеров упаковок для демонстрации и аудита';
COMMENT ON COLUMN pack_sets.id IS 'Уникальный идентификатор набора';
COMMENT ON COLUMN pack_sets.name IS 'Название набора размеров';
COMMENT ON COLUMN pack_sets.sizes IS 'Массив размеров упаковок в формате JSON';
COMMENT ON COLUMN pack_sets.created_at IS 'Время создания записи';
COMMENT ON COLUMN pack_sets.updated_at IS 'Время последнего обновления записи';

