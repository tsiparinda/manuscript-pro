select bucket 
FROM citedata cd
where cd.id_col=($1)
and iscollectionreadable(cd.id_col, ($2))=true ;