BEGIN;
CREATE SCHEMA maker;
CREATE TABLE maker.peps_everyblock (
  id           SERIAL,
  block_number INTEGER NOT NULL,
  pep        NUMERIC,
  pip        NUMERIC,
  per        NUMERIC
);

END;
