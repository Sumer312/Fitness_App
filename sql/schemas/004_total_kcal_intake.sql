-- +goose Up 

CREATE TABLE total_calorie_intake(
  id UUID PRIMARY KEY,
  created_at TIMESTAMP NOT NULL,
  calories FLOAT NOT NULL,
  total_deficit FLOAT NOT NULL,
  total_surplus FLOAT NOT NULL,
  user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE
);

-- +goose Down
DROP TABLE total_calorie_intake;

