SELECT COUNT (c.Id)
FROM collections c
inner join users u on c.id_author = u.Id;