-- Create "invited_users" table
CREATE TABLE "invited_users" (
 "id" uuid NOT NULL DEFAULT gen_random_uuid(),
 "created_at" timestamptz NULL,
 "updated_at" timestamptz NULL,
 "deleted_at" timestamptz NULL,
 "email" text NULL,
 "password_hash" character varying(255) NULL,
 PRIMARY KEY ("id")
);
-- Create index "idx_invited_users_email" to table: "invited_users"
CREATE UNIQUE INDEX "idx_invited_users_email" ON "invited_users" ("email");
-- Create "users" table
CREATE TABLE "users" (
 "id" uuid NOT NULL DEFAULT gen_random_uuid(),
 "created_at" timestamptz NULL,
 "updated_at" timestamptz NULL,
 "deleted_at" timestamptz NULL,
 "username" text NULL,
 "email" text NULL,
 "role" text NULL DEFAULT 'user',
 "password_hash" character varying(255) NULL,
 "first_name" text NULL,
 "last_name" text NULL,
 PRIMARY KEY ("id"),
 CONSTRAINT "uni_users_password_hash" UNIQUE ("password_hash")
);
-- Create index "idx_users_email" to table: "users"
CREATE UNIQUE INDEX "idx_users_email" ON "users" ("email");
-- Create index "idx_users_username" to table: "users"
CREATE UNIQUE INDEX "idx_users_username" ON "users" ("username");
