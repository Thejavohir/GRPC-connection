package repo

type InMemoryStorageI interface {
	SetWithTTL(key, value string, seconds int64) error
	Get(key string) (interface{}, error)
}
