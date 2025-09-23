SELECT id, email, password, full_name, create_time
FROM users
WHERE id = $1
LIMIT 1;
