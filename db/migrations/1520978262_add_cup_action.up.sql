CREATE TYPE maker.ACT AS ENUM (
  'give',
  'open',
  'join',
  'exit',
  'lock',
  'free',
  'draw',
  'wipe',
  'shut',
  'bite'
);

CREATE TABLE maker.cup_action (
  log_id  INTEGER,
  id      INTEGER,
  tx      CHARACTER VARYING(66) NOT NULL,
  act     maker.ACT             NOT NULL,
  arg     CHARACTER VARYING(66),
  lad     CHARACTER VARYING(66) NOT NULL,
  ink     NUMERIC DEFAULT 0     NOT NULL,
  art     NUMERIC DEFAULT 0     NOT NULL,
  ire     NUMERIC DEFAULT 0     NOT NULL,
  block   INTEGER               NOT NULL,
  deleted BOOLEAN DEFAULT FALSE,
  guy     CHARACTER VARYING(66),
  CONSTRAINT tx_act_arg_constraint UNIQUE (tx, act, arg)
);