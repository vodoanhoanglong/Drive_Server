alter table "public"."account" add column "loginType" text default 'default';
alter table "public"."account" add constraint "account_email_key" unique ("email");
