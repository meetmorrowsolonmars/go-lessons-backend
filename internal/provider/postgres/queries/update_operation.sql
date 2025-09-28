UPDATE operations
SET amount      = $2,
    description = $3,
    category_id = $4
WHERE id = $1;
