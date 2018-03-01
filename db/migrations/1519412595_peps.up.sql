CREATE TABLE maker.peps (
  id           SERIAL,
  log_id       INTEGER,
  block_number INTEGER,
  value        NUMERIC,
  CONSTRAINT log_index_fk FOREIGN KEY (log_id)
  REFERENCES logs (id)
  ON DELETE CASCADE
);


