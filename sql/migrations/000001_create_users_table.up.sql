CREATE TABLE IF NOT EXISTS users (
	id serial4 NOT NULL,
	"name" text NOT NULL,
	email text NOT NULL,
	"description" text null,
	avatar bytea null,
	password_hash text  NULL,
	is_verified bool NOT NULL DEFAULT false,
	verification_code text NULL,
	created_at timestamptz NOT NULL DEFAULT CURRENT_TIMESTAMP,
	updated_at timestamptz NOT NULL DEFAULT CURRENT_TIMESTAMP,
	CONSTRAINT users_pkey PRIMARY KEY (id),
	CONSTRAINT users_email_is_lowercase CHECK (email = lower(email)),
	CONSTRAINT users_email_key UNIQUE (email)
);


CREATE TABLE IF NOT EXISTS groups (
	id serial4 NOT NULL,
	"name" text not null,
	created_at timestamptz NOT NULL DEFAULT CURRENT_TIMESTAMP,
	updated_at timestamptz NOT NULL DEFAULT CURRENT_TIMESTAMP,
	CONSTRAINT groups_pkey PRIMARY KEY (id),
	CONSTRAINT groups_name_key UNIQUE ("name")
);

CREATE TABLE IF NOT EXISTS usergroups (
	id serial4 NOT NULL,
	id_user serial4 NOT NULL,
	id_group serial4 NOT NULL,
	created_at timestamptz NOT NULL DEFAULT CURRENT_TIMESTAMP,
	updated_at timestamptz NOT NULL DEFAULT CURRENT_TIMESTAMP,
	CONSTRAINT usergroups_pkey PRIMARY KEY (id),
	CONSTRAINT usergroups_unique_usergroups UNIQUE (id_user, id_group)
);

ALTER TABLE usergroups ADD CONSTRAINT usergroups_groups_fk FOREIGN KEY (id_group) REFERENCES groups(id)  ON DELETE CASCADE;

ALTER TABLE usergroups ADD CONSTRAINT usergroups_users_fk FOREIGN KEY (id_user) REFERENCES users(id)  ON DELETE CASCADE;


COMMENT ON COLUMN public.users.id IS 'User''s ID';
COMMENT ON COLUMN public.users."name" IS 'Use''r name';
COMMENT ON COLUMN public.users.email IS 'User''s email';