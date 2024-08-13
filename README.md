# GoNetWorker

GoNetWorker는 웹 서버 운영 연습을 위한 가상의 사용자입니다. config 파일에 endpoint를 설정하면 설정값 범위 내로 랜덤으로 request를 보내게 됩니다. 

```json
{
  "settings": {
    "sleep_range" : 5
  },
  "works": [
    {
      "uri": "http://localhost",
      "port": 8080,
      "info": [
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

