alter table "public"."files"
  add constraint "files_userId_fkey"
  foreign key ("userId")
  references "public"."users"
  ("id") on update restrict on delete restrict;
