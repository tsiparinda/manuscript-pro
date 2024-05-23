SELECT count(c.Id)
FROM collections c
inner join users u on c.id_author=u.id
where u.id= ($1)
