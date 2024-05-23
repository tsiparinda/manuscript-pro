create extension IF NOT EXISTS  hstore;
create extension IF NOT EXISTS  lo;

-- CREATE TABLE imagecollectionlist (
-- 	id serial4 NOT NULL,
-- 	filedata "hstore" NULL,
-- 	created_at timestamptz NOT NULL DEFAULT CURRENT_TIMESTAMP,
-- 	updated_at timestamptz NOT NULL DEFAULT CURRENT_TIMESTAMP,
-- 	CONSTRAINT imagecollectionlist_pkey PRIMARY KEY (id)
-- );

-- CREATE INDEX idx_imagecollectionlist_filedata ON imagecollectionlist USING hash (filedata);

CREATE TABLE imagecollectionlist (
	id serial4 NOT NULL,
	filepath text not NULL,
	created_at timestamptz NOT NULL DEFAULT CURRENT_TIMESTAMP,
	updated_at timestamptz NOT NULL DEFAULT CURRENT_TIMESTAMP,
	CONSTRAINT imagecollectionlist_pkey PRIMARY KEY (id)
);


CREATE unique INDEX idx_imagecollectionlist_filepath ON imagecollectionlist  (filepath);