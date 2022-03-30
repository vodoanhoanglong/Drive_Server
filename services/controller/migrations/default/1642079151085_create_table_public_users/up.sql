CREATE TABLE "public"."users" (
    "id" uuid NOT NULL DEFAULT gen_random_uuid(),
    "fullName" text NOT NULL,
    "phone" text NOT NULL,
    "status" integer NOT NULL,
    "createdAt" timestamptz NOT NULL DEFAULT now(),
    "createdBy" text,
    "updatedAt" timestamptz NOT NULL DEFAULT now(),
    "updatedBy" text,
    "accountId" text NOT NULL,
    PRIMARY KEY ("id"),
    FOREIGN KEY ("accountId") REFERENCES "public"."account"("id") ON UPDATE restrict ON DELETE restrict,
);

CREATE EXTENSION IF NOT EXISTS pgcrypto;