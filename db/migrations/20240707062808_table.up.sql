CREATE TABLE "users" (
  "id" SERIAL PRIMARY KEY,
  "name" varchar(255) NOT NULL,
  "email" varchar(255) UNIQUE NOT NULL,
  "password" varchar(255) NOT NULL,
  "role" char(20) NOT NULL,
  "created_at" timestamp NOT NULL DEFAULT (current_timestamp),
  "updated_at" timestamp NOT NULL DEFAULT (current_timestamp)
);

CREATE TABLE "novels" (
  "id" SERIAL PRIMARY KEY,
  "title" varchar(255) NOT NULL,
  "keyword" text NOT NULL,
  "short" text NOT NULL,
  "content" text NOT NULL,
  "user_id" int NOT NULL,
  "created_at" timestamp NOT NULL DEFAULT (current_timestamp),
  "updated_at" timestamp NOT NULL DEFAULT (current_timestamp)
);

CREATE TABLE "novel_likes" (
  "id" SERIAL PRIMARY KEY,
  "novel_id" int NOT NULL,
  "user_id" int NOT NULL,
  "created_at" timestamp NOT NULL DEFAULT (current_timestamp)
);

CREATE TABLE "comments" (
  "id" SERIAL PRIMARY KEY,
  "novel_id" INT NOT NULL,
  "user_id" INT NOT NULL,
  "parent_id" INT DEFAULT null,
  "content" TEXT NOT NULL,
  "created_at" TIMESTAMP NOT NULL DEFAULT (CURRENT_TIMESTAMP)
);

CREATE TABLE "comment_likes" (
  "id" SERIAL PRIMARY KEY,
  "comment_id" INT NOT NULL,
  "user_id" INT NOT NULL,
  "created_at" TIMESTAMP NOT NULL DEFAULT (CURRENT_TIMESTAMP)
);

CREATE TABLE "tags" (
  "id" SERIAL PRIMARY KEY,
  "name" varchar(255) UNIQUE NOT NULL
);

CREATE TABLE "novel_tags" (
  "novel_id" INT NOT NULL,
  "tag_id" INT NOT NULL,
  PRIMARY KEY ("novel_id", "tag_id")
);

ALTER TABLE "novels" ADD FOREIGN KEY ("user_id") REFERENCES "users" ("id");

ALTER TABLE "comments" ADD FOREIGN KEY ("novel_id") REFERENCES "novels" ("id");

ALTER TABLE "novel_likes" ADD FOREIGN KEY ("novel_id") REFERENCES "novels" ("id");

ALTER TABLE "novel_likes" ADD FOREIGN KEY ("user_id") REFERENCES "users" ("id");

ALTER TABLE "comments" ADD FOREIGN KEY ("user_id") REFERENCES "users" ("id");

ALTER TABLE "comments" ADD FOREIGN KEY ("parent_id") REFERENCES "comments" ("id") ON DELETE CASCADE;

ALTER TABLE "comment_likes" ADD FOREIGN KEY ("comment_id") REFERENCES "comments" ("id") ON DELETE CASCADE;

ALTER TABLE "comment_likes" ADD FOREIGN KEY ("user_id") REFERENCES "users" ("id");

ALTER TABLE "novel_tags" ADD FOREIGN KEY ("tag_id") REFERENCES "tags" ("id") ON DELETE CASCADE;

ALTER TABLE "novel_tags" ADD FOREIGN KEY ("novel_id") REFERENCES "novels" ("id") ON DELETE CASCADE;
