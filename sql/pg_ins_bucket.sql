
INSERT INTO citedata ( id_col, bucket)
select ($1), ($2)
where not exists (select id from citedata where  id_col=($1) and bucket = ($2))
