select id_col, id_user, is_write 
from colusers cu
inner join collections c on c.id =cu.id_col 
where c.id=($1)
and c.id_author = ($2)
and cu.id_user != ($2)
ORDER BY id_user desc
