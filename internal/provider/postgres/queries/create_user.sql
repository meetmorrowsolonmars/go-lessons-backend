INSERT INTO users (email, password, full_name, create_time)
VALUES ($1, $2, $3, $4)
RETURNING id, create_time;
