package controller

import (
	"encoding/json"
	"io"
	"os"
)

type Reader struct {
	filename string
	result   []byte
}

func NewReader(filename string) *Reader {
	return &Reader{
		filename: filename,
	}
}

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

func (r *Reader) GetConfig() (error, *Config) {
	err := r.readJson()
	if err != nil {
		return err, nil
	}

	var config Config
	err = json.Unmarshal(r.result, &config)
	if err != nil {
		return err, nil
	}
	return nil, &config
}
