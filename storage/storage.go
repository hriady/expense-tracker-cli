package storage

import (
	"encoding/json"
	"os"
)

type StorageManager[T any] struct {
	FileName string
}

func NewStorageManager[T any](fileName string) *StorageManager[T] {
	return &StorageManager[T]{
		FileName: fileName,
	}
}

func (sm *StorageManager[T]) Save(data []T) error {
	taskBytes, err := json.MarshalIndent(data, "", "	")
	if err != nil {
		return err
	}
	
	err = os.WriteFile(sm.FileName, taskBytes, 0644)
	if err != nil {
		return err
	}
	
	return nil
}

func (sm *StorageManager[T]) Load() (resp []T, err error) {
	
	file, err := os.Open(sm.FileName)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	
	err = json.NewDecoder(file).Decode(&resp)
	if err != nil {
		return nil, err
	}
	
	return resp, nil
}
