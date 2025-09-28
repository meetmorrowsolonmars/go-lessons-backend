SELECT id,
       user_id,
       account_id,
       type,
       category_id,
       amount,
       description,
       create_time
FROM operations
WHERE id = $1;
