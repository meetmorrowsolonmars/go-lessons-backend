SELECT id, email, password, full_name, create_time
FROM users
WHERE email = $1
LIMIT 1;
