CREATE OR REPLACE FUNCTION public.filter(search text, page integer, size integer)
 RETURNS SETOF account
 LANGUAGE sql
 STABLE
AS $function$
  SELECT * 
  FROM account 
  WHERE page > 0 AND size > 0 AND email ILIKE ('%' || search || '%')
  ORDER BY id DESC
  LIMIT size OFFSET (page - 1) * size 
$function$;
