
INSERT INTO collections
( id_author, title)
Values (($1), ($2)) RETURNING id;

