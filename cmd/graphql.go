// Copyright © 2018 NAME HERE <EMAIL ADDRESS>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package cmd

import (

	"github.com/spf13/cobra"
	"github.com/neelance/graphql-go"
	"github.com/neelance/graphql-go/relay"
	"github.com/vulcanize/vulcanizedb/pkg/geth"
	"github.com/vulcanize/vulcanizedb/utils"
	"net/http"
	"log"
	"github.com/8thlight/sai_watcher/graphql_server"
	"github.com/vulcanize/vulcanizedb/pkg/datastore/postgres/repositories"
	cupsRepo "github.com/8thlight/sai_watcher/cup/repositories"
)

// graphqlCmd represents the graphql command
var graphqlCmd = &cobra.Command{
	Use:   "graphql",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
	    schema := parseSchema()
	    serve(schema)
	},
}

func init() {
	rootCmd.AddCommand(graphqlCmd)
}

func parseSchema() *graphql.Schema {

	blockchain := geth.NewBlockchain("/Users/mattkrump/Library/Ethereum/geth.ipc")
	db := utils.LoadPostgres(cfg, blockchain.Node())
	blockRepository := &repositories.BlockRepository{DB: &db}
	logRepository := &repositories.LogRepository{DB: &db}
	filterRepository := &repositories.FilterRepository{DB: &db}
	watchedEventRepository := &repositories.WatchedEventRepository{DB: &db}
	cupsRepository := &cupsRepo.CupsRepository{DB: &db}
	graphQLRepositories := graphql_server.GraphQLRepositories{
		WatchedEventRepository: watchedEventRepository,
		BlockRepository:        blockRepository,
		LogRepository:          logRepository,
		FilterRepository:       filterRepository,
		ICupsRepository:        cupsRepository,
	}
	schema := graphql.MustParseSchema(graphql_server.Schema, graphql_server.NewResolver(graphQLRepositories))
	return schema

}

func serve(schema *graphql.Schema) {
	http.Handle("/", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write(page)
	}))
	http.Handle("/query", &relay.Handler{Schema: schema})

	log.Fatal(http.ListenAndServe(":9090", nil))
}

var page = []byte(`
<!DOCTYPE html>
<html>
	<head>
		<link rel="stylesheet" href="https://cdnjs.cloudflare.com/ajax/libs/graphiql/0.10.2/graphiql.css" />
		<script src="https://cdnjs.cloudflare.com/ajax/libs/fetch/1.1.0/fetch.min.js"></script>
		<script src="https://cdnjs.cloudflare.com/ajax/libs/react/15.5.4/react.min.js"></script>
		<script src="https://cdnjs.cloudflare.com/ajax/libs/react/15.5.4/react-dom.min.js"></script>
		<script src="https://cdnjs.cloudflare.com/ajax/libs/graphiql/0.10.2/graphiql.js"></script>
	</head>
	<body style="width: 100%; height: 100%; margin: 0; overflow: hidden;">
		<div id="graphiql" style="height: 100vh;">Loading...</div>
		<script>
			function graphQLFetcher(graphQLParams) {
				return fetch("/query", {
					method: "post",
					body: JSON.stringify(graphQLParams),
					credentials: "include",
				}).then(function (response) {
					return response.text();
				}).then(function (responseBody) {
					try {
						return JSON.parse(responseBody);
					} catch (error) {
						return responseBody;
					}
				});
			}
			ReactDOM.render(
				React.createElement(GraphiQL, {fetcher: graphQLFetcher}),
				document.getElementById("graphiql")
			);
		</script>
	</body>
</html>
`)
