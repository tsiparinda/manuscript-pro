SELECT value FROM
  (SELECT  skeys(dict) key, svals(dict) value 
  FROM citedata cd
  inner join collections c on c.id=cd.id_col
    where c.id=($1) and bucket=($2)
    and iscollectionreadable(c.id, ($4))=true 
    ) AS stat
  where key= ($3);