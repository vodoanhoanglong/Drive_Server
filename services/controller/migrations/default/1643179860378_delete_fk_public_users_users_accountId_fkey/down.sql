alter table "public"."users"
  add constraint "users_accountId_fkey"
  foreign key ("accountId")
  references "public"."account"
  ("id") on update restrict on delete restrict;
