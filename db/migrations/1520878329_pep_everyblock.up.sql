BEGIN;
CREATE SCHEMA maker;
CREATE TABLE maker.peps_everyblock (
  id           SERIAL,
  block_number BIGINT  NOT NULL,
  block_id     INTEGER NOT NULL,
  block_time   BIGINT,
  pep          NUMERIC,
  pip          NUMERIC,
  per          NUMERIC,
  CONSTRAINT blocks_fk FOREIGN KEY (block_id)
  REFERENCES blocks (id)
  ON DELETE CASCADE
);

END;
