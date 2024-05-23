CREATE TABLE colusers (
	id serial4 NOT NULL,
	id_col int4 NOT NULL,
	id_user int4 NOT NULL,
	is_write boolean not null DEFAULT FALSE,
	created_at timestamptz NOT NULL DEFAULT CURRENT_TIMESTAMP,
	updated_at timestamptz NOT NULL DEFAULT CURRENT_TIMESTAMP,
	CONSTRAINT colusers_pkey PRIMARY KEY (id),
	CONSTRAINT col_unique_key UNIQUE ("id_col", "id_user")
);

-- foreign keys
ALTER TABLE colusers ADD CONSTRAINT colusers_users_fk FOREIGN KEY (id_user) REFERENCES users(id) ON DELETE CASCADE;
ALTER TABLE colusers ADD CONSTRAINT colusers_collections_fk FOREIGN KEY (id_col) REFERENCES collections(id) ON DELETE CASCADE;
