CREATE DOMAIN email AS TEXT
CHECK(
   VALUE ~ '^[a-zA-Z0-9_.+-]+@[a-zA-Z0-9-]+[.][a-zA-Z0-9-]+([.][a-zA-Z0-9-]+)?$'
);

CREATE EXTENSION IF NOT EXISTS pgcrypto;
CREATE EXTENSION IF NOT EXISTS pg_trgm;
CREATE EXTENSION IF NOT EXISTS btree_gin;


CREATE OR REPLACE FUNCTION "public"."set_current_timestamp_updated_at"()
RETURNS TRIGGER AS $$
DECLARE
  _new record;
BEGIN
  _new := NEW;
  _new."updated_at" = NOW();
  RETURN _new;
END;
$$ LANGUAGE plpgsql;

--------------------------------------------------
-- TABLE public.account
--------------------------------------------------
CREATE TABLE "public"."account" (
  "id"                    TEXT NOT NULL DEFAULT gen_random_uuid(), 
  "email"                 email, 
  "display_name"          TEXT, 
  "password"              TEXT,
  "avatar_url"            TEXT,
  "role"                  TEXT NOT NULL, 
  "birthday"              DATE,
  "created_at"            TIMESTAMPTZ NOT NULL DEFAULT now(), 
  "updated_at"            TIMESTAMPTZ NOT NULL DEFAULT now(),  
  "created_by"            TEXT NULL, 
  "updated_by"            TEXT NULL,  
  PRIMARY KEY ("id")
);

CREATE UNIQUE INDEX account_email_unique
  ON "public"."account"(email);

CREATE TRIGGER "account_set_current_timestamp_updated_at"
BEFORE INSERT OR UPDATE ON "public"."account"
FOR EACH ROW
EXECUTE PROCEDURE "public"."set_current_timestamp_updated_at"();

CREATE VIEW public.me AS
  SELECT * FROM public.account;
--------------------------------------------------
-- END TABLE public.account
--------------------------------------------------
