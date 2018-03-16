package graphql_server_test

import (
	"context"
	"log"

	"encoding/json"

	"math/big"

	"github.com/8thlight/sai_watcher/everyblock"
	"github.com/8thlight/sai_watcher/graphql_server"
	"github.com/8thlight/sai_watcher/utils"
	"github.com/neelance/graphql-go"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/vulcanize/vulcanizedb/pkg/config"
	"github.com/vulcanize/vulcanizedb/pkg/core"
	"github.com/vulcanize/vulcanizedb/pkg/datastore/postgres"
)

func formatJSON(data []byte) []byte {
	var v interface{}
	if err := json.Unmarshal(data, &v); err != nil {
		log.Fatalf("invalid JSON: %s", err)
	}
	formatted, err := json.Marshal(v)
	if err != nil {
		log.Fatal(err)
	}
	return formatted
}

func convertHelper(method string, input [32]byte, precsion int) string {
	converted := big.NewInt(0)
	converted.SetBytes(input[:])
	return utils.Convert(method, converted.String(), precsion)
}

var _ bool = Describe("GraphQL", func() {
	var graphQLRepositories graphql_server.GraphQLRepositories
	var zero [32]byte
	var one [32]byte
	var pep everyblock.Peek
	var pip everyblock.Peek
	var per everyblock.Per

	BeforeEach(func() {
		node := core.Node{GenesisBlock: "GENESIS", NetworkID: 1, ID: "x123", ClientName: "geth"}
		db, err := postgres.NewDB(config.Database{
			Hostname: "localhost",
			Name:     "vulcanize_private",
			Port:     5432,
		}, node)
		Expect(err).NotTo(HaveOccurred())
		db.Query(`DELETE FROM maker.cups`)
		db.Query(`DELETE FROM maker.peps`)
		db.Query(`DELETE FROM maker.peps_everyblock`)
		db.Query(`DELETE FROM logs`)
		db.Query(`DELETE FROM log_filters`)
		ebds := everyblock.DataStore{DB: db}
		graphQLRepositories = graphql_server.GraphQLRepositories{
			Everyblock: ebds,
		}
		zero = [32]byte{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0}
		one = [32]byte{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1}
		pep = everyblock.Peek{Value: one, OK: false}
		pip = everyblock.Peek{Value: zero, OK: false}
		per = everyblock.Per{Value: big.NewInt(1)}
		ebds.Create(1, pep, pip, per)

	})

	It("Queries schema everyblock for a specific block number", func() {
		var variables map[string]interface{}
		resolver := graphql_server.NewResolver(graphQLRepositories)
		var schema = graphql.MustParseSchema(graphql_server.Schema, resolver)
		response := schema.Exec(context.Background(),
			`{
                          everyblock(blockNumber: 1) {
                             total
                             rows {
                                    blockNumber
                                    per
                                    pep
                                    pip
                                 }
                           }
                       }`,
			"",
			variables)
		expected := `{
          "everyblock": {
            "total": 1,
            "rows": [
              {
                "blockNumber": 1,
                "per": "0.000000000000000000000000001",
                "pep": "0.000000000000000001",
                "pip": "0"
              }
            ]
          }
        }`
		var v interface{}
		if len(response.Errors) != 0 {
			log.Fatal(response.Errors)
		}
		err := json.Unmarshal(response.Data, &v)
		Expect(err).ToNot(HaveOccurred())
		actualJSON := formatJSON(response.Data)
		expectedJSON := formatJSON([]byte(expected))
		Expect(actualJSON).To(Equal(expectedJSON))
	})

})
