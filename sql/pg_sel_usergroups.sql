select g.id, g.name
from usergroups u inner join groups g on u.id_group=g.id where u.id_user = ($1);