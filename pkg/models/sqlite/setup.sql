-- USERS TABLE --
CREATE TABLE IF NOT EXISTS "users" (
	"ID"	INTEGER NOT NULL UNIQUE,
	"Login"	TEXT NOT NULL UNIQUE,
	"Email"	TEXT NOT NULL UNIQUE,
	"Created"	TEXT NOT NULL,
	"Password"	TEXT NOT NULL,
	PRIMARY KEY("ID" AUTOINCREMENT)
);

-- POSTS TABLE --
CREATE TABLE IF NOT EXISTS "posts" (
	"ID"	INTEGER NOT NULL UNIQUE,
	"UserID"	INTEGER NOT NULL,
	"UserLogin" TEXT NOT NULL,
	"Title"	TEXT NOT NULL,
	"Text"	TEXT NOT NULL,
	"Created"	TEXT NOT NULL,
	PRIMARY KEY("ID" AUTOINCREMENT),
	FOREIGN KEY("UserID") REFERENCES "users"("ID")
);

-- TAGS TABLE --
CREATE TABLE IF NOT EXISTS "tags" (
	"ID"	INTEGER NOT NULL UNIQUE,
	"Tag"	TEXT NOT NULL UNIQUE,
	PRIMARY KEY("ID" AUTOINCREMENT)
);

-- POSTS AND TAGS TABLE --
CREATE TABLE IF NOT EXISTS "postsAndTags" (
	"ID"	INTEGER NOT NULL UNIQUE,
	"PostID"	INTEGER NOT NULL,
	"TagID"	INTEGER NOT NULL,
	PRIMARY KEY("ID" AUTOINCREMENT),
	FOREIGN KEY("PostID") REFERENCES "posts"("ID"),
	FOREIGN KEY("TagID") REFERENCES "tags"("ID")
);

-- COMMENTS TABLE --
CREATE TABLE IF NOT EXISTS "comments" (
	"ID"	INTEGER NOT NULL UNIQUE,
	"UserID"	INTEGER NOT NULL,
	"PostID"	INTEGER NOT NULL,
	"Text"	TEXT NOT NULL,
	"Created"	TEXT NOT NULL,
	PRIMARY KEY("ID" AUTOINCREMENT),
	FOREIGN KEY("UserID") REFERENCES "users"("ID"),
	FOREIGN KEY("PostID") REFERENCES "posts"("ID")
);