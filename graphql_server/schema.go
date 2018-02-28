package graphql_server

import (
    cupsrepo "github.com/8thlight/sai_watcher/cup/repositories"
	"github.com/vulcanize/vulcanizedb/pkg/core"
	"github.com/vulcanize/vulcanizedb/pkg/datastore"
	"github.com/vulcanize/vulcanizedb/pkg/filters"
)

var Schema = `
	schema {
		query: Query
	}
	type Query {
        logFilter(name: String!): LogFilter
        watchedEvents(name: String!): WatchedEventList
		cupHistory(cupIndex: Int!): CupList
	}

    type LogFilter {
        name: String!
        fromBlock: Int
        toBlock: Int
        address: String!
        topics: [String]!
    }

  	type WatchedEventList {
		total: Int!
		watchedEvents: [WatchedEvent]!
	}

    type WatchedEvent {
        name: String!
        blockNumber: Int!
        address: String!
        tx_hash: String!
        topic0: String!
        topic1: String!
        topic2: String!
        topic3: String!
        data: String!
    }

	type CupList {
		total: Int!
		cups: [Cup]!
	}

	type Cup {
		cupIndex: Int!
		blockNumber: Int!
		logId: Int!
		lad: String!
		ink: String!
		art: String!
		irk: String!
		is_closed: Boolean!
	}
`

type GraphQLRepositories struct {
	datastore.BlockRepository
	cupsrepo.ICupsRepository
	datastore.LogRepository
	datastore.WatchedEventRepository
	datastore.FilterRepository
}

type Resolver struct {
	graphQLRepositories GraphQLRepositories
}

func NewResolver(repositories GraphQLRepositories) *Resolver {
	return &Resolver{graphQLRepositories: repositories}
}

func (r *Resolver) LogFilter(args struct {
	Name string
}) (*logFilterResolver, error) {
	logFilter, err := r.graphQLRepositories.GetFilter(args.Name)
	if err != nil {
		return &logFilterResolver{}, err
	}
	return &logFilterResolver{&logFilter}, nil
}

type logFilterResolver struct {
	lf *filters.LogFilter
}

func (lfr *logFilterResolver) Name() string {
	return lfr.lf.Name
}

func (lfr *logFilterResolver) FromBlock() *int32 {
	fromBlock := int32(lfr.lf.FromBlock)
	return &fromBlock
}

func (lfr *logFilterResolver) ToBlock() *int32 {
	toBlock := int32(lfr.lf.ToBlock)
	return &toBlock
}

func (lfr *logFilterResolver) Address() string {
	return lfr.lf.Address
}

func (lfr *logFilterResolver) Topics() []*string {
	var topics = make([]*string, 4)
	for i := range topics {
		if lfr.lf.Topics[i] != "" {
			topics[i] = &lfr.lf.Topics[i]
		}
	}
	return topics
}

func (r *Resolver) WatchedEvents(args struct {
	Name string
}) (*watchedEventsResolver, error) {
	watchedEvents, err := r.graphQLRepositories.GetWatchedEvents(args.Name)
	if err != nil {
		return &watchedEventsResolver{}, err
	}
	return &watchedEventsResolver{watchedEvents: watchedEvents}, err
}

type watchedEventsResolver struct {
	watchedEvents []*core.WatchedEvent
}

func (wesr watchedEventsResolver) WatchedEvents() []*watchedEventResolver {
	return resolveWatchedEvents(wesr.watchedEvents)
}

func (wesr watchedEventsResolver) Total() int32 {
	return int32(len(wesr.watchedEvents))
}

func resolveWatchedEvents(watchedEvents []*core.WatchedEvent) []*watchedEventResolver {
	watchedEventResolvers := make([]*watchedEventResolver, 0)
	for _, watchedEvent := range watchedEvents {
		watchedEventResolvers = append(watchedEventResolvers, &watchedEventResolver{watchedEvent})
	}
	return watchedEventResolvers
}

type watchedEventResolver struct {
	we *core.WatchedEvent
}

func (wer watchedEventResolver) Name() string {
	return wer.we.Name
}

func (wer watchedEventResolver) BlockNumber() int32 {
	return int32(wer.we.BlockNumber)
}

func (wer watchedEventResolver) Address() string {
	return wer.we.Address
}

func (wer watchedEventResolver) TxHash() string {
	return wer.we.TxHash
}

func (wer watchedEventResolver) Topic0() string {
	return wer.we.Topic0
}

func (wer watchedEventResolver) Topic1() string {
	return wer.we.Topic1
}

func (wer watchedEventResolver) Topic2() string {
	return wer.we.Topic2
}

func (wer watchedEventResolver) Topic3() string {
	return wer.we.Topic3
}

func (wer watchedEventResolver) Data() string {
	return wer.we.Data
}

func (r *Resolver) CupHistory(args struct {
	CupIndex int32
}) (*cupsResolver, error) {
	cups, err := r.graphQLRepositories.GetCupsByIndex(int(args.CupIndex))
	if err != nil {
		return &cupsResolver{}, err
	}
	return &cupsResolver{cups: cups}, err
}

type cupsResolver struct {
	cups []cupsrepo.DBCup
}

func (csr cupsResolver) Cups() []*cupResolver {
	return resolveCups(csr.cups)
}

func (csr cupsResolver) Total() int32 {
	return int32(len(csr.cups))
}

func resolveCups(cups []cupsrepo.DBCup) []*cupResolver {
	cupResolvers := make([]*cupResolver, 0)
	for _, cup := range cups {
		cupResolvers = append(cupResolvers, &cupResolver{cup})
	}
	return cupResolvers
}

type cupResolver struct {
	c cupsrepo.DBCup
}

func (cr cupResolver) LogId() int32 {
	return int32(cr.c.LogID)
}

func (cr cupResolver) BlockNumber() int32 {
	return int32(cr.c.BlockNumber)
}

func (cr cupResolver) CupIndex() int32 {
	return int32(cr.c.CupIndex)
}

func (cr cupResolver) Lad() string {
	return cr.c.Lad
}

func (cr cupResolver) Ink() string {
	return string(cr.c.Ink)
}

func (cr cupResolver) Art() string {
	return string(cr.c.Art)
}

func (cr cupResolver) Irk() string {
	return string(cr.c.Irk)
}

func (cr cupResolver) IsClosed() bool {
	return cr.c.IsClosed
}
