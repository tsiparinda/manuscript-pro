select cd.id, c.id_author, bucket, (each(dict)).key, (each(dict)).value 
FROM citedata cd
inner join collections c on c.id=cd.id_col;