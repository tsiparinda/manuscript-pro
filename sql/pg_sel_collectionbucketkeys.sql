select (each(dict)).key
FROM citedata cd
inner join collections c on c.id=cd.id_col
where c.id=($1) and cd.bucket=($2)
and iscollectionreadable(c.id, ($3))=true ;