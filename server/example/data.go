package example

// @NewNetWorker uri http://localhost port 8080
// @NetWorker path /tasks method GET
type CreateReq struct {
	Title string `json:"title"`
}

// @NetWorker path /tasks/add method POST
type AddSubReq struct {
	TaskId string `json:"task_id"`
	Title  string `json:"title"`
}
