CREATE EXTENSION IF NOT EXISTS pgcrypto;

CREATE TABLE chapters (
	id uuid NOT NULL DEFAULT gen_random_uuid(),
	novel_id uuid NOT NULL,
	title text NOT NULL,
	author text NOT NULL,
	description text NOT NULL,
	content text NOT NULL,
	chapter_index bigint DEFAULT 0,
	deleted boolean DEFAULT false,
	created_at timestamptz NOT NULL DEFAULT now(),
	updated_at timestamptz NOT NULL DEFAULT now(),
	PRIMARY KEY (novel_id, id)
) PARTITION BY HASH (novel_id);

CREATE TABLE users (
	id uuid PRIMARY KEY DEFAULT gen_random_uuid(),
	username text NOT NULL,
	type text NOT NULL,
	email text NOT NULL,
	password text NOT NULL,
	created_at timestamptz NOT NULL DEFAULT now(),
	updated_at timestamptz NOT NULL DEFAULT now()
);

CREATE TABLE novels (
	id uuid PRIMARY KEY DEFAULT gen_random_uuid(),
	title text NOT NULL,
	author text NOT NULL,
	description text NOT NULL,
	chapter_count bigint NOT NULL DEFAULT 0,
	created_at timestamptz NOT NULL DEFAULT now(),
	updated_at timestamptz NOT NULL DEFAULT now()
);

CREATE INDEX idx_chapters_novel_index_id ON chapters (novel_id, chapter_index, id);
CREATE INDEX idx_users_email_id ON users (email, id);
CREATE INDEX idx_novels_title_id ON novels (title, id);
