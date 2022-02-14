-- USERS TABLE --
CREATE TABLE IF NOT EXISTS "users" (
	"id"	INTEGER NOT NULL UNIQUE,
	"login"	TEXT NOT NULL UNIQUE,
	"email"	TEXT NOT NULL UNIQUE,
	"created"	TEXT NOT NULL,
	"password"	TEXT NOT NULL,
	PRIMARY KEY("id" AUTOINCREMENT)
);

-- POSTS TABLE --
CREATE TABLE IF NOT EXISTS "posts" (
	"id"	INTEGER NOT NULL UNIQUE,
	"user_id"	INTEGER NOT NULL,
	"user_login" TEXT NOT NULL,
	"title"	TEXT NOT NULL,
	"text"	TEXT NOT NULL,
	"created"	TEXT NOT NULL,
	PRIMARY KEY("id" AUTOINCREMENT),
	FOREIGN KEY("user_id") REFERENCES "users"("id")
);

-- TAGS TABLE --
CREATE TABLE IF NOT EXISTS "tags" (
	"id"	INTEGER NOT NULL UNIQUE,
	"tag"	TEXT NOT NULL UNIQUE,
	PRIMARY KEY("id" AUTOINCREMENT)
);

-- POSTS AND TAGS TABLE --
CREATE TABLE IF NOT EXISTS "posts_and_tags" (
	"id"	INTEGER NOT NULL UNIQUE,
	"post_id"	INTEGER NOT NULL,
	"tag_id"	INTEGER NOT NULL,
	PRIMARY KEY("id" AUTOINCREMENT),
	FOREIGN KEY("post_id") REFERENCES "posts"("id"),
	FOREIGN KEY("tag_id") REFERENCES "tags"("id")
);

-- COMMENTS TABLE --
CREATE TABLE IF NOT EXISTS "comments" (
	"id"	INTEGER NOT NULL UNIQUE,
	"user_id"	INTEGER NOT NULL,
	"post_id"	INTEGER NOT NULL,
	"login" TEXT NOT NULL,
	"text"	TEXT NOT NULL,
	"created"	TEXT NOT NULL,
	PRIMARY KEY("id" AUTOINCREMENT),
	FOREIGN KEY("user_id") REFERENCES "users"("id"),
	FOREIGN KEY("post_id") REFERENCES "posts"("id")
);

-- LIKE POST TABLE --
CREATE TABLE IF NOT EXISTS "like_post" (
	"id"	INTEGER NOT NULL UNIQUE,
	"post_id"	INTEGER NOT NULL,
	"user_id"	INTEGER NOT NULL,
	"is_like"	INTEGER NOT NULL,
	FOREIGN KEY("post_id") REFERENCES "posts"("id"),
	FOREIGN KEY("user_id") REFERENCES "users"("id"),
	PRIMARY KEY("id" AUTOINCREMENT)
);