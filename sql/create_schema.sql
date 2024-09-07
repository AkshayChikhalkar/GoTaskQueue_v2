CREATE TABLE IF NOT EXISTS tasks (
  id SERIAL PRIMARY KEY,
  type INT NOT NULL,
  value INT NOT NULL,
  state VARCHAR(50),
  creation_time TIMESTAMP,
  last_update_time TIMESTAMP
);
