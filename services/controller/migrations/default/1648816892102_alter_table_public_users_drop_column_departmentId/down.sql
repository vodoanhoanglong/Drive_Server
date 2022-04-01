alter table "public"."users" alter column "departmentId" drop not null;
alter table "public"."users" add column "departmentId" uuid;
