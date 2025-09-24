SELECT id, user_id, account_id, type, amount, description, create_time
FROM operations
WHERE id = $1;
