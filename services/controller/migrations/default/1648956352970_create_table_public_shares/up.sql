CREATE TABLE "public"."shares"
(
    "accountId" Text        NOT NULL,
    "fileId"    text        NOT NULL,
    "status"    text     NOT NULL DEFAULT 'active',
    "createdAt" timestamptz NOT NULL DEFAULT now(),
    "createdBy" text,
    "updatedAt" timestamptz NOT NULL DEFAULT now(),
    "updatedBy" text,
    PRIMARY KEY ("accountId", "fileId"),
    FOREIGN KEY ("accountId") REFERENCES "public"."account" ("id") ON UPDATE restrict ON DELETE restrict,
    FOREIGN KEY ("fileId") REFERENCES "public"."files" ("id") ON UPDATE restrict ON DELETE restrict,
    FOREIGN KEY ("createdBy") REFERENCES "public"."account" ("id") ON UPDATE restrict ON DELETE restrict,
    FOREIGN KEY ("updatedBy") REFERENCES "public"."account" ("id") ON UPDATE restrict ON DELETE restrict
);
