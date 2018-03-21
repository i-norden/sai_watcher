CREATE VIEW public.cup_act AS
  SELECT
    act,
    arg,
    art,
    block,
    deleted,
    maker.cup_action.id,
    ink,
    ire,
    lad,
    pip,
    per,
    (pip * per * ink) / NULLIF(art, 0) * 100 AS ratio,
    (pip * per * ink)                        AS tab,
    tx
  FROM maker.cup_action
    LEFT JOIN maker.peps_everyblock ON maker.peps_everyblock.block_number = maker.cup_action.block
  ORDER BY block DESC;

CREATE VIEW public.cup AS
  SELECT
    act,
    art,
    block,
    deleted,
    id,
    ink,
    ire,
    lad,
    pip,
    per,
    (pip * per * ink) / NULLIF(art, 0) * 100 AS ratio,
    (pip * per * ink)                        AS tab
  FROM (
         SELECT DISTINCT ON (cup_action.id)
           act,
           art,
           block,
           deleted,
           id,
           ink,
           ire,
           lad,
           (SELECT pip
            FROM maker.peps_everyblock
            ORDER BY block_number DESC
            LIMIT 1),
           (SELECT per
            FROM maker.peps_everyblock
            ORDER BY block_number DESC
            LIMIT 1) AS per
         FROM maker.cup_action
         ORDER BY maker.cup_action.id DESC, maker.cup_action.block DESC
       )
       c;

