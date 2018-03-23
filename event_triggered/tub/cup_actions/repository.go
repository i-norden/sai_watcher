package cup_actions

import "github.com/vulcanize/vulcanizedb/pkg/datastore/postgres"

type CupActionsRepository struct {
	DB *postgres.DB
}

func (car CupActionsRepository) CreateCupAction(cupAction CupActionModel, logID int64) error {
	id := cupAction.ID
	tx := cupAction.TransactionHash
	act := cupAction.Act
	arg := cupAction.Arg
	lad := cupAction.Lad
	ink := cupAction.Ink
	art := cupAction.Art
	ire := cupAction.Ire
	block := cupAction.Block
	deleted := cupAction.Deleted
	guy := cupAction.Guy

	_, err := car.DB.Exec(
		`INSERT INTO maker.cup_action (log_id, id, tx, act, arg, lad, ink, art, ire, block, deleted, guy)
                SELECT $1, $2, $3, $4, NULLIF($5, ''), $6, $7::NUMERIC, $8::NUMERIC, $9::NUMERIC, $10, $11, $12
                WHERE NOT EXISTS (SELECT log_id FROM maker.cup_action WHERE log_id = $1)`,
		logID, id, tx, act, arg, lad, ink, art, ire, block, deleted, guy)
	if err != nil {
		return err
	}
	return nil
}
