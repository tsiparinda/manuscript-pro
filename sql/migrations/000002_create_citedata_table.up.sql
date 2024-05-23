create extension IF NOT EXISTS  hstore;
create extension IF NOT EXISTS  lo;

CREATE TABLE collections (
	id serial4 NOT NULL,
	id_author int4 NOT NULL,
	title text not null,
	is_public boolean not null DEFAULT FALSE,
	created_at timestamptz NOT NULL DEFAULT CURRENT_TIMESTAMP,
	updated_at timestamptz NOT NULL DEFAULT CURRENT_TIMESTAMP,
	CONSTRAINT collections_pkey PRIMARY KEY (id)
);

CREATE TABLE citedata (
	id serial4 NOT NULL,
	id_col serial4 not null,
	bucket text NULL,
	dict "hstore" NULL,
	created_at timestamptz NOT NULL DEFAULT CURRENT_TIMESTAMP,
	updated_at timestamptz NOT NULL DEFAULT CURRENT_TIMESTAMP,
	CONSTRAINT citedata_pkey PRIMARY KEY (id)
);

-- citedata foreign keys
ALTER TABLE collections ADD CONSTRAINT collections_author_fk FOREIGN KEY (id_author) REFERENCES users(id);
ALTER TABLE citedata ADD CONSTRAINT citedata_collections_fk FOREIGN KEY (id_col) REFERENCES collections(id) ON DELETE CASCADE;

CREATE INDEX idx_h_dict ON citedata USING hash (dict);
CREATE INDEX idx_h_bucket ON citedata USING hash (bucket);
CREATE INDEX idx_h_id_author ON collections USING hash (id_author);
CREATE INDEX idx_h_id_col ON citedata USING hash (id_col);