SELECT id, name
FROM categories
WHERE id = $1
LIMIT 1;
