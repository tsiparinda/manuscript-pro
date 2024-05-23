SELECT value FROM
  (SELECT  skeys(dict) key, svals(dict) value 
  FROM imagedata cd
  inner join collections c on c.id=cd.id_col
    where c.id=($1) 
    and iscollectionreadable(c.id, ($3))=true 
    ) AS stat
  where key= ($2);