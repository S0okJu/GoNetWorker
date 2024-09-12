## GoNetWorker
GoNetWorker는 운영 환경을 학습하기 위한 가상의 사용자입니다. json에 endpoint를 정의하여 형식에 맞게 서버에 요청을 보낼 수 있습니다.

### Worker.json

| Key   | Value             | Describe              |
| ----- | ----------------- | --------------------- |
| works | uri, port, []task | endpoint에 대한 정보       |
| tasks | path, method      | 요청하고자 하는 경로와 메소드가 정의됨 |

```json
{  
  "settings": {  
    "sleep_range" : 5  
  },  
  "works": [
    {
      "uri": "http://localhost",
      "port": 8080,
      "tasks": [
        {
          "path": "/users/1",
          "method": "GET"
        },
        {
          "path": "/users/2",
          "method": "GET"
        }
      ]
    }
  ]
}
```


### 구조

```mermaid
flowchart LR
	subgraph GoNetWorker
	stop_ch(Stop Channel)
	req([Request])
	end
	
	user((user))
	server([server])

	user --Ctrl+C-->stop_ch
	stop_ch-->server

	user -->req
	req --Request-->server
```
