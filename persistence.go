package main

import (
	"os"

	"encoding/gob"
)

func persisted(name string) bool {
	_, err := os.Stat(name + ".gob")
	return !os.IsNotExist(err)
}

func getPersistence[T any](name string) (T, error) {
	var data T

	file, err := os.Open(name + ".gob")
	if err != nil {
		return data, err
	}
	defer file.Close()

	decoder := gob.NewDecoder(file)
	if err := decoder.Decode(&data); err != nil {
		return data, err
	}

	return data, nil
}

func persist[T any](name string, data T) error {
	file, err := os.Create(name + ".gob")
	if err != nil {
		return err
	}
	defer file.Close()

	encoder := gob.NewEncoder(file)
	if err := encoder.Encode(data); err != nil {
		return err
	}

	return nil
}
