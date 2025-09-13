package fsys

import (
	"encoding/json"
	"os"
)

func ReadJson[T any](path string, data T) error {
	file, err := os.Open(path)
	if err != nil {
		return err
	}
	defer file.Close()

	decoder := json.NewDecoder(file)
	if err := decoder.Decode(data); err != nil {
		return err
	}

	if err := file.Close(); err != nil {
		return err
	}
	return nil
}

func WriteJson[T any](path string, data T) error {
	file, err := os.Create(path)
	if err != nil {
		return err
	}
	defer file.Close()

	encoder := json.NewEncoder(file)

	if err := encoder.Encode(data); err != nil {
		return err
	}

	if err := file.Sync(); err != nil {
		return err
	}
	return nil
}

func WriteJsonIndent[T any](path string, data T) error {
	file, err := os.Create(path)
	if err != nil {
		return err
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ")

	if err := encoder.Encode(data); err != nil {
		return err
	}

	if err := file.Sync(); err != nil {
		return err
	}
	return nil
}
