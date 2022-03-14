package settingStorage

import (
	"gitlab.com/xiayesuifeng/gopanel/core/storage"
	"sync"
)

var (
	once           sync.Once
	currentStorage *SettingStorage
)

type SettingStorage struct {
	baseStorage storage.Storage
}

func GetStorage() *SettingStorage {
	once.Do(func() {
		currentStorage = &SettingStorage{
			baseStorage: storage.GetBaseStorage(),
		}
	})

	return currentStorage
}

func (s SettingStorage) Get(key string, defaultValue []byte) []byte {
	bytes, err := s.baseStorage.Get("setting", key)
	if err != nil || bytes == nil {
		return defaultValue
	}

	return bytes
}

func (s SettingStorage) Set(key string, value []byte) error {
	return s.baseStorage.Set("setting", key, value)
}

func (s SettingStorage) Has(key string) bool {
	return s.baseStorage.Has("setting", key)
}

func (s SettingStorage) Delete(key string) error {
	return s.baseStorage.Delete("setting", key)
}
