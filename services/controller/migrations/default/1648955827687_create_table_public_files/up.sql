CREATE TABLE "public"."files"
(
    "id"        text        NOT NULL DEFAULT gen_random_uuid(),
    "path"      text        NOT NULL,
    "extension" text        NOT NULL,
    "name"      text        NOT NULL,
    "size"      integer     NOT NULL,
    "url"       text        NOT NULL,
    "layer"     integer     NOT NULL,
    "status"    text     NOT NULL DEFAULT 'active',
    "createdAt" timestamptz NOT NULL DEFAULT now(),
    "createdBy" text,
    "updatedAt" timestamptz NOT NULL DEFAULT now(),
    "updatedBy" text,
    PRIMARY KEY ("id"),
    FOREIGN KEY ("createdBy") REFERENCES "public"."account" ("id") ON UPDATE restrict ON DELETE restrict,
    FOREIGN KEY ("updatedBy") REFERENCES "public"."account" ("id") ON UPDATE restrict ON DELETE restrict,
    UNIQUE ("path")
);
CREATE
EXTENSION IF NOT EXISTS pgcrypto;
