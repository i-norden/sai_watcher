BEGIN;
ALTER TABLE maker.cup_action
  DROP CONSTRAINT log_index_fk;

ALTER TABLE maker.gov
  DROP CONSTRAINT log_index_fk;
COMMIT;
