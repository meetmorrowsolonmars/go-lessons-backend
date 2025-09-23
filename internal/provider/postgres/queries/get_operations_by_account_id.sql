SELECT id, user_id, account_id, type, amount, description, create_time
FROM operations
WHERE account_id = $1
ORDER BY create_time DESC
LIMIT $2 OFFSET $3;
