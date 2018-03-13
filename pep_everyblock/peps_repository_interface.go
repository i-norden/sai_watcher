package peps_everyblock

type IPepsRepository interface {
	CreatePep(value string, blockNumber int64) error
}
