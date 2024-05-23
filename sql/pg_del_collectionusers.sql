DELETE  FROM colusers WHERE  "id_col"=($1) and 
id_col in (select id from collections where id_author=($2)) ;
