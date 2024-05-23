with updimage as (
UPDATE imagedata SET dict = case when dict is null then  hstore(($2),($3))
else dict || hstore(($2),($3)) end 
where imagedata.id_col=($1)
returning *
)
insert into imagedata (id_col, dict)
select ($1), hstore(($2),($3))
where not exists (SELECT * FROM imagedata WHERE imagedata.id_col=($1));