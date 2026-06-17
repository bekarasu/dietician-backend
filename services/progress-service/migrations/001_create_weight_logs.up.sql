CREATE TABLE IF NOT EXISTS weight_logs (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL,
    weight_kg DECIMAL(5,1) NOT NULL,
    notes TEXT,
    logged_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_weight_logs_user_id ON weight_logs (user_id);
CREATE INDEX idx_weight_logs_logged_at ON weight_logs (user_id, logged_at);
