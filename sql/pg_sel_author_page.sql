SELECT c.Id, c.title, c.is_public,  
    u.Id, u.Name author
FROM collections c
inner join users u on c.id_author=u.id
inner join citedata cd on cd.id_col=c.id
where u.id= ($1)
ORDER BY c.Id desc
LIMIT ($2) OFFSET ($3)