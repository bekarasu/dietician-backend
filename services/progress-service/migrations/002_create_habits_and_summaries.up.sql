CREATE TABLE IF NOT EXISTS habit_logs (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL,
    habit_name VARCHAR(100) NOT NULL,
    completed BOOLEAN NOT NULL DEFAULT false,
    notes TEXT,
    logged_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_habit_logs_user_id ON habit_logs (user_id);

CREATE TABLE IF NOT EXISTS weekly_progress_summaries (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL,
    week_start DATE NOT NULL,
    week_end DATE NOT NULL,
    start_weight DECIMAL(5,1),
    end_weight DECIMAL(5,1),
    weight_change DECIMAL(5,1),
    habits_completed INT DEFAULT 0,
    habits_total INT DEFAULT 0,
    summary TEXT,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_weekly_progress_user_id ON weekly_progress_summaries (user_id);
CREATE UNIQUE INDEX idx_weekly_progress_user_week ON weekly_progress_summaries (user_id, week_start);
