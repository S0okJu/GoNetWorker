package core

type Config struct {
	Works []Work `json:"works"`
}

type Work struct {
	Uri  string    `json:"uri"`
	Port int       `json:"port"`
	Info []ReqInfo `json:"info"`
}

func (w *Work) Count() int {
	return len(w.Info)
}

type ReqInfo struct {
	Path   string            `json:"path"`
	Method string            `json:"method"`
	Param  map[string]string `json:"param"`
	Weight int               `json:"weight"`
}
