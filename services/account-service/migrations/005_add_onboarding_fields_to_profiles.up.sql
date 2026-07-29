ALTER TABLE user_profiles 
ADD COLUMN age INT,
ADD COLUMN display_name VARCHAR(100),
ADD COLUMN target_weight_kg DECIMAL(5,1),
ADD COLUMN daily_calorie_target INT,
ADD COLUMN target_water_ml INT,
ADD COLUMN target_coffee_cups INT;
