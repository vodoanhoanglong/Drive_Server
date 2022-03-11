alter table "public"."users" alter column "accountId" drop not null;
alter table "public"."users" add column "accountId" text;
