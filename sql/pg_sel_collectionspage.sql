select *
from 
(SELECT c.Id, c.title, c.is_public, c.id_author,
  u.Id id_user, u.Name author, iscollectionwriteble(c.id, ($1)) isWritable
FROM collections c
inner join users u on c.id_author=u.id
where iscollectionreadable(c.id, ($1))=true 
and ((($2)>0 and c.id_author=($2)) or ($2)=0)
) rrr
ORDER BY Id desc
LIMIT ($3) OFFSET ($4)