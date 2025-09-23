SELECT id, user_id, title, is_default, create_time
FROM accounts
WHERE id = $1
LIMIT 1;
