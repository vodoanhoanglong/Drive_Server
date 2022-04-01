alter table "public"."files" alter column "sharedDepartment" drop not null;
alter table "public"."files" add column "sharedDepartment" text;
