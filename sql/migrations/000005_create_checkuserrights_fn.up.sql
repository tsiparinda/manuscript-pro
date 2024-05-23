CREATE OR REPLACE FUNCTION iscollectionreadable(colid integer, userid integer)
 RETURNS boolean
 LANGUAGE plpgsql
AS 
$$
  declare
  	fAllow boolean = false;
begin
	select true
	into fAllow
	from collections c
	 where c.id=colid and 
	 (c.is_public=true  
		or c.id in (select id_col from colusers where id_user=userid)
		or c.id_author =userid
	);
return fallow;
end;
$$;


CREATE OR REPLACE FUNCTION public.iscollectionwriteble(colid integer, userid integer)
 RETURNS boolean
 LANGUAGE plpgsql
AS $$
begin
	return 
	exists(select 1
	from collections c
	 where c.id=colid and 
	 ( c.id in (select id_col from colusers where id_user=userid and colusers.is_write=true)
		or c.id_author =userid
	));

end;
$$;