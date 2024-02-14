-- +goose Up 

CREATE TABLE daily_nutrition_intake(
  id UUID PRIMARY KEY,
  created_at TIMESTAMP NOT NULL,
  user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
  calories INT NOT NULL,
  carbohydrates INT NOT NULL,
  protien INT NOT NULL,
  fat INT NOT NULL,
  fiber INT NOT NULL
);

-- +goose Down
DROP TABLE daily_nutrition_intake;
