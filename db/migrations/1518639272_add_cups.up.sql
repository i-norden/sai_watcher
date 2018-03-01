BEGIN;
CREATE SCHEMA maker;

CREATE TABLE maker.cups (
  log_id    INTEGER,
  id           SERIAL,
  cup_index    NUMERIC,
  lad          VARCHAR(42),
  ink          NUMERIC,
  art          NUMERIC,
  irk          NUMERIC,
  block_number NUMERIC,
  is_closed    BOOLEAN,
  CONSTRAINT log_index_fk FOREIGN KEY (log_id)
  REFERENCES logs (id)
  ON DELETE CASCADE
);
COMMIT;
