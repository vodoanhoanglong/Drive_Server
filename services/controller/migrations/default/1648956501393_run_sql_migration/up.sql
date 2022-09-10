CREATE
OR REPLACE FUNCTION public.check_file_name(path_input text, name_input text, extension_input text) RETURNS SETOF files LANGUAGE sql STABLE AS $function$
SELECT *
FROM files
WHERE status <> 'deleted' AND path SIMILAR TO (path_input || '/%') AND name LIKE name_input AND extension LIKE extension_input AND length(path) = length(path_input) + 37 $function$;
