INSERT INTO accounts (user_id, title, is_default, create_time)
VALUES ($1, $2, $3, $4)
RETURNING id, create_time;
