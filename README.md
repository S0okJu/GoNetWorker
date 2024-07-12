# GoNetWorker

GoNetWorker는 웹 서버 운영 연습을 위한 가상의 사용자입니다. 

## 사용법
1. 적용하고자 하는 파일에 `@NewNetWorker`를 정의한 후 url와 port 정보를 기입합니다. 단, 최초로 `@NewNetWorker`를 작성할 때는 아래와 같이 `@NetWorker` 위에 기입해야 합니다.
2. `@NetWorker`를 정의한 후 path와 method 정보를 기입합니다.
3. 정보를 기반으로 `/docs/gonetworker.config.json` 파일이 생성됩니다.

#### 예시
```go
// @NewNetWorker url http://localhost port 8080
// @NetWorker path /tasks method GET
type CreateReq struct {
	Title string `json:"title"`
}

// @NetWorker path /tasks/add method POST
type AddSubReq struct {
	TaskId string `json:"task_id"`
	Title  string `json:"title"`
}

```

#### gonetworker.config.json 
```json
{
  "works": [
    {
      "url": "http://localhost",
      "port": 8080,
      "request": [
        {
          "path": "/tasks",
          "method": "GET",
          "param": {
            "title": "string"
          }
        },
        {
          "path": "/tasks/add",
          "method": "POST",
          "param": {
            "task_id": "string",
            "title": "string"
          }
        }
      ]
    }
  ]
}
```

## Worker 정보 
#### @NewNetWorker
| 주석 | 설명 | 예시                 |
|---|---|--------------------|
| base | 기본 URL | http[:]//localhost |
| port | 포트 번호 | 8080 |


#### @NetWorker
| 주석 | 설명 | 예시                 |
|---|---|--------------------|
| path | 요청 URL | /tasks |
| method | 요청 방식 | GET |
