CREATE TABLE public.users (
	id int NOT NULL,
	"name" text NOT NULL,
	passwordhash text NOT NULL,
	"role" int DEFAULT 0 NOT NULL,
	"isoauth" bool DEFAULT false NOT NULL,
	CONSTRAINT users_pk PRIMARY KEY (id)
);

CREATE TABLE public.shares (
	id text NOT NULL,
	title text NOT NULL,
	"content" text NOT NULL,
	passwordhash text NOT NULL,
	expire_at timestamp with time zone NOT NULL,
	author_id int NOT NULL,
	hide_author bool DEFAULT false NOT NULL,
	CONSTRAINT shares_pk PRIMARY KEY (id)
);

CREATE TABLE public.sessions (
	session_id text NOT NULL,
	user_id int NOT NULL,
	expire_at timestamp with time zone NOT NULL,
	CONSTRAINT sessions_pk PRIMARY KEY (session_id)
);