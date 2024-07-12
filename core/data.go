package core

type Config struct {
	Works []Work `json:"works"`
}

type Work struct {
	Uri     string    `json:"uri"`
	Port    int       `json:"port"`
	Request []Request `json:"request"`
}

type Request struct {
	Path   string            `json:"path"`
	Method string            `json:"method"`
	Param  map[string]string `json:"param"`
}
