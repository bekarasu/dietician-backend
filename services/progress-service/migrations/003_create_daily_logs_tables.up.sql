CREATE TABLE IF NOT EXISTS daily_logs (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL,
    log_date DATE NOT NULL,
    water_intake_ml INT DEFAULT 0,
    sleep_hours DECIMAL(4,1) DEFAULT 0.0,
    exercise_minutes INT DEFAULT 0,
    mood VARCHAR(50),
    notes TEXT,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW()
);

CREATE UNIQUE INDEX idx_daily_logs_user_date ON daily_logs (user_id, log_date);

CREATE TABLE IF NOT EXISTS daily_log_meals (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    daily_log_id UUID NOT NULL REFERENCES daily_logs(id) ON DELETE CASCADE,
    meal_type VARCHAR(50) NOT NULL, -- breakfast, lunch, dinner, snack
    name VARCHAR(255) NOT NULL,
    calories INT,
    protein_g DECIMAL(5,1),
    carbs_g DECIMAL(5,1),
    fat_g DECIMAL(5,1),
    logged_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_daily_log_meals_log_id ON daily_log_meals (daily_log_id);
