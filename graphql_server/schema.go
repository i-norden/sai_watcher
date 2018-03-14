package graphql_server

import (
	"github.com/8thlight/sai_watcher/everyblock"
)

var Schema = `
	schema {
		query: Query
	}

	type Query {
		everyblock(blockNumber: Int!): Rows
	}

  	type Rows {
		total: Int!
		rows: [Row]!
	}
	
	type Row {
		id: Int!	
		blockNumber: Int!
		per: String!
		pep: String!
		pip: String!
	}
`

type GraphQLRepositories struct {
	Everyblock everyblock.DataStore
}

type Resolver struct {
	graphQLRepositories GraphQLRepositories
}

func NewResolver(repositories GraphQLRepositories) *Resolver {
	return &Resolver{graphQLRepositories: repositories}
}

func (r *Resolver) EveryBlock(args struct {
	BlockNumber int32
}) (*rowsResolver, error) {
	row, err := r.graphQLRepositories.Everyblock.Get(int64(args.BlockNumber))
	rows := []*everyblock.Row{row}
	if err != nil {
		return &rowsResolver{}, err
	}
	return &rowsResolver{rows: rows}, err
}

type rowsResolver struct {
	rows []*everyblock.Row
}

func (ebr rowsResolver) Rows() []*rowResolver {
	return resolveRows(ebr.rows)
}

func (ebrs rowsResolver) Total() int32 {
	return int32(len(ebrs.rows))
}

func resolveRows(rows []*everyblock.Row) []*rowResolver {
	rowResolvers := make([]*rowResolver, 0)
	for _, row := range rows {
		rowResolvers = append(rowResolvers, &rowResolver{row})
	}
	return rowResolvers
}

type rowResolver struct {
	r *everyblock.Row
}

func (rr rowResolver) Id() int32 {
	return int32(rr.r.ID)
}

func (rr rowResolver) BlockNumber() int32 {
	return int32(rr.r.BlockNumber)
}

func (rr rowResolver) Per() string {
	return rr.r.Per
}

func (rr rowResolver) Pep() string {
	return rr.r.Pep
}

func (rr rowResolver) Pip() string {
	return rr.r.Pip
}
