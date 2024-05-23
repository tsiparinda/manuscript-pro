create extension IF NOT EXISTS  hstore;
create extension IF NOT EXISTS  lo;

CREATE TABLE imagedata (
	id serial4 NOT NULL,
	id_col serial4 not null,
	dict "hstore" NULL,
	created_at timestamptz NOT NULL DEFAULT CURRENT_TIMESTAMP,
	updated_at timestamptz NOT NULL DEFAULT CURRENT_TIMESTAMP,
	CONSTRAINT imagedata_pkey PRIMARY KEY (id)
);

-- citedata foreign keys
ALTER TABLE imagedata ADD CONSTRAINT imagedata_collections_fk FOREIGN KEY (id_col) REFERENCES collections(id) ON DELETE CASCADE;

CREATE INDEX idx_imagedata_dict ON imagedata USING hash (dict);
CREATE INDEX idx_imagedata_id_col ON imagedata USING hash (id_col);