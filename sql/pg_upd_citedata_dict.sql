UPDATE citedata SET dict = case when dict is null then hstore(($3),($4))
else dict || hstore(($3), ($4)) end 
where citedata.id_col=($1) and bucket=($2);