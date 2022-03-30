alter table "public"."users"
  add constraint "users_departmentId_fkey"
  foreign key (departmentId)
  references "public"."departments"
  (id) on update restrict on delete restrict;
alter table "public"."users" alter column "departmentId" drop not null;
alter table "public"."users" add column "departmentId" uuid;
