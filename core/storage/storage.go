package storage

type Storage interface {
	Get(module, key string) ([]byte, error)
	Set(module, key string, value []byte) error
	Has(module, key string) bool
	Delete(module, key string) error
}
