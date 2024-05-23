select id, "name", email, coalesce("description", ''), '' as avatar
from users
where (($1)>0 and  id=($1))
or (($1)=0);