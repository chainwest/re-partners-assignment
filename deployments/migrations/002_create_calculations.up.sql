-- Create calculations table
CREATE TABLE IF NOT EXISTS calculations (
    id BIGSERIAL PRIMARY KEY,
    pack_set_id BIGINT REFERENCES pack_sets(id) ON DELETE SET NULL,
    pack_sizes JSONB NOT NULL,
    amount INTEGER NOT NULL,
    breakdown JSONB NOT NULL,
    total_packs INTEGER NOT NULL,
    overage INTEGER NOT NULL,
    calculated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT calculations_amount_check CHECK (amount > 0),
    CONSTRAINT calculations_total_packs_check CHECK (total_packs >= 0),
    CONSTRAINT calculations_overage_check CHECK (overage >= 0),
    CONSTRAINT calculations_pack_sizes_check CHECK (jsonb_typeof(pack_sizes) = 'array'),
    CONSTRAINT calculations_breakdown_check CHECK (jsonb_typeof(breakdown) = 'object')
);

-- Create index on pack_set_id for faster lookups
CREATE INDEX idx_calculations_pack_set_id ON calculations(pack_set_id);

-- Create index on calculated_at for sorting and filtering
CREATE INDEX idx_calculations_calculated_at ON calculations(calculated_at DESC);

-- Create index on amount for analytics
CREATE INDEX idx_calculations_amount ON calculations(amount);

-- Create composite index for common queries
CREATE INDEX idx_calculations_pack_set_amount ON calculations(pack_set_id, amount);

-- Add comments to table
COMMENT ON TABLE calculations IS 'История расчётов упаковок для аудита и демонстрации';
COMMENT ON COLUMN calculations.id IS 'Уникальный идентификатор расчёта';
COMMENT ON COLUMN calculations.pack_set_id IS 'Ссылка на набор размеров (опционально)';
COMMENT ON COLUMN calculations.pack_sizes IS 'Массив размеров упаковок, использованных в расчёте';
COMMENT ON COLUMN calculations.amount IS 'Требуемое количество элементов';
COMMENT ON COLUMN calculations.breakdown IS 'Разбивка решения: размер -> количество упаковок';
COMMENT ON COLUMN calculations.total_packs IS 'Общее количество упаковок в решении';
COMMENT ON COLUMN calculations.overage IS 'Превышение (количество элементов сверх требуемого)';
COMMENT ON COLUMN calculations.calculated_at IS 'Время выполнения расчёта';

