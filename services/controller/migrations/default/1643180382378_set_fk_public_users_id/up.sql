alter table "public"."users"
  add constraint "users_id_fkey"
  foreign key ("id")
  references "public"."account"
  ("id") on update restrict on delete restrict;
