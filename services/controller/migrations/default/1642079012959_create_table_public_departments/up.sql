CREATE TABLE "public"."departments" ("id" uuid NOT NULL DEFAULT gen_random_uuid(), "name" text NOT NULL, "des" text, "status" integer NOT NULL, "createdAt" timestamptz NOT NULL DEFAULT now(), "createdBy" text, "updatedAt" timestamptz NOT NULL DEFAULT now(), "updatedBy" text, PRIMARY KEY ("id") , UNIQUE ("name"));
CREATE EXTENSION IF NOT EXISTS pgcrypto;
