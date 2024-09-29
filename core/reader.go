package core

import (
	"encoding/json"
	"io"
	"os"
)

// Reader
// NOTICE: JSON 파일만 읽기 가능
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

	jsonFile, err := os.Open("./test-server.json")
	if err != nil {
		return err
	}
	defer jsonFile.Close()

	byteValue, err := io.ReadAll(jsonFile)
	if err != nil {
		return err
	}
	r.result = byteValue
	return nil
}

// GetConfig 읽은 값을 Config 구조체로 반환
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
