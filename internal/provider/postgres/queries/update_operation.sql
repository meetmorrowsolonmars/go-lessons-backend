UPDATE operations
SET amount      = $2,
    description = $3
WHERE id = $1;
