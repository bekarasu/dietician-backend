CREATE TABLE IF NOT EXISTS tracking_metrics (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL,
    metric_type VARCHAR(50) NOT NULL, -- weight, body_fat, waist_cm, etc.
    target_value DECIMAL(8,2),
    current_value DECIMAL(8,2) NOT NULL,
    deadline DATE,
    status VARCHAR(50) NOT NULL DEFAULT 'active',
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_tracking_metrics_user_id ON tracking_metrics (user_id);
