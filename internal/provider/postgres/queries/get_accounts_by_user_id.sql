SELECT id, user_id, title, is_default, create_time
FROM accounts
WHERE user_id = $1
ORDER BY is_default DESC, title ASC;
