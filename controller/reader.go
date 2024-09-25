package controller

import (
	"encoding/json"
	"io"
	"os"
)

// Reader is the file reader struct
// The file can read only JSON files
type Reader struct {
	filename string
	result   []byte
}

func NewReader(filename string) *Reader {
	return &Reader{
		filename: filename,
	}
}

// readJson reads the JSON file
func (r *Reader) readJson() error {

	// Open the JSON file
	jsonFile, err := os.Open("./test.json")
	if err != nil {
		return err
	}
	defer jsonFile.Close()

	// Read the file's content
	byteValue, err := io.ReadAll(jsonFile)
	if err != nil {
		return err
	}
	r.result = byteValue
	return nil
}

// GetConfig reads the JSON file and returns the Config struct
func (r *Reader) GetConfig() (*Config, error) {
	err := r.readJson()
	if err != nil {
		return nil, err
	}

	var config Config
	err = json.Unmarshal(r.result, &config)
	if err != nil {
		return nil, err
	}
	return &config, nil
}
