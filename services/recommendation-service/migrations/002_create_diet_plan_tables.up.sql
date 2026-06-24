CREATE TABLE IF NOT EXISTS diet_plans (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL,
    start_date DATE NOT NULL,
    end_date DATE,
    goals TEXT,
    status VARCHAR(50) NOT NULL DEFAULT 'active',
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_diet_plans_user_id ON diet_plans (user_id);
CREATE INDEX idx_diet_plans_status ON diet_plans (status);

CREATE TABLE IF NOT EXISTS diet_plan_meals (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    diet_plan_id UUID NOT NULL REFERENCES diet_plans(id) ON DELETE CASCADE,
    day_of_week INT, -- 1=Monday, 7=Sunday. Or specific date if not a weekly recurring plan.
    meal_type VARCHAR(50) NOT NULL, -- e.g., breakfast, lunch, dinner, snack
    recipe_id UUID, -- Optional reference to a recipe service if exists
    name VARCHAR(255) NOT NULL,
    description TEXT,
    calories INT,
    protein_g DECIMAL(5,1),
    carbs_g DECIMAL(5,1),
    fat_g DECIMAL(5,1),
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_diet_plan_meals_plan_id ON diet_plan_meals (diet_plan_id);
