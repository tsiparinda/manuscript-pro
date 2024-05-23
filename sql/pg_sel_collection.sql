select *
from 
(SELECT c.Id, c.title, c.is_public,  c.id_author 
FROM collections c
where iscollectionreadable(c.id, ($2))=true 
and c.id=($1)
) rrr
ORDER BY Id desc
