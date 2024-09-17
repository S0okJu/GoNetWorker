package core

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
)

type reader struct {
	filepath string
	data     []byte
}

func newreader(filepath string) *reader {
	return &reader{filepath: filepath}
}

func (r *reader) read() error {
	// Open the JSON file:%s/
	jsonFile, err := os.Open(r.filepath)
	if err != nil {
		return fmt.Errorf("Error opening JSON file: %v", err)
	}
	defer jsonFile.Close()

	// Read the file's content
	byteValue, err := io.ReadAll(jsonFile)
	if err != nil {
		return fmt.Errorf("Error reading JSON file: %v", err)
	}
	r.data = byteValue

	return nil
}

func (r *reader) convertTo() (*Config, error) {
	var config Config
	err := json.Unmarshal(r.data, &config)
	if err != nil {
		return nil, err
	}
	return &config, nil
}

func LoadJSON(filepath string) (*Config, error) {
	r := newreader(filepath)
	err := r.read()
	if err != nil {
		return nil, err
	}

	cfg, err := r.convertTo()
	if err != nil {
		return nil, err
	}
	return cfg, nil
}
