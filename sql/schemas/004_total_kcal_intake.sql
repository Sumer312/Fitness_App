-- +goose Up 

CREATE TABLE total_calorie_intake(
  id UUID PRIMARY KEY,
  created_at TIMESTAMP NOT NULL,
  updated_at TIMESTAMP NOT NULL,
  calories INT NOT NULL,
  program VARCHAR(16) NOT NULL,
  user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE
);

-- +goose Down
DROP TABLE total_calorie_intake;

