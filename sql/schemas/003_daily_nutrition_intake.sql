-- +goose Up 

CREATE TABLE daily_nutrition_intake(
  id UUID PRIMARY KEY,
  created_at TIMESTAMP NOT NULL,
  user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
  program VARCHAR(16) NOT NULL,
  calories FLOAT NOT NULL,
  carbohydrates FLOAT NOT NULL,
  protein FLOAT NOT NULL,
  fat FLOAT NOT NULL,
  fiber FLOAT NOT NULL
);

-- +goose Down
DROP TABLE daily_nutrition_intake;
