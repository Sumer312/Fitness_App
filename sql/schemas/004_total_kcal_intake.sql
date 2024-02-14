-- +goose Up 

CREATE TABLE total_calorie_intake(
  id UUID PRIMARY KEY,
  created_at TIMESTAMP NOT NULL,
  calories INT NOT NULL,
  total_deficit INT NOT NULL,
  total_surplus INT NOT NULL,
  user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE
);

-- +goose Down
DROP TABLE total_calorie_intake;

