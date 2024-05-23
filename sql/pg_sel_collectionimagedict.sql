select (each(dict)).key, (each(dict)).value 
FROM imagedata cd
inner join collections c on c.id=cd.id_col
where c.id=($1) 
and iscollectionreadable(c.id, ($2))=true ;