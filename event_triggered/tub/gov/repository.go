package gov

import (
	"fmt"

	"github.com/vulcanize/vulcanizedb/pkg/datastore/postgres"
)

type DataStore struct {
	DB *postgres.DB
}

func (ds DataStore) CreateGov(govModel *GovModel, logID int64) error {
	tx := govModel.Tx
	var_ := govModel.Var
	arg_ := govModel.Arg
	guy := govModel.Guy
	cap_ := govModel.Cap
	mat := govModel.Mat
	tax := govModel.Tax
	fee := govModel.Fee
	axe := govModel.Axe
	gap := govModel.Gap
	block := govModel.Block
	_, err := ds.DB.Exec(
		`INSERT INTO maker.gov (log_id, tx, var, arg, guy, cap, mat, tax, fee, axe, gap, block)
                SELECT $1, $2::VARCHAR, $3, $4::NUMERIC, $5, $6::NUMERIC, $7::NUMERIC, $8::NUMERIC, $9::NUMERIC, $10::NUMERIC, $11::NUMERIC, $12 
                WHERE NOT EXISTS (SELECT tx FROM maker.gov WHERE tx = $2)`,
		logID, tx, var_, arg_, guy, cap_, mat, tax, fee, axe, gap, block)
	if err != nil {
		fmt.Println(tx)
		return err
	}
	return nil
}
