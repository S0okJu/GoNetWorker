package processcontroller

type StartSignal struct {
	RequestId string `json:"request_id"`
	Status    string `json:"status"`
}

type StopSignal struct {
	RequestId string `json:"request_id"`
	Status    string `json:"status"`
}
