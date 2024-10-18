package core

import (
	"fmt"
	"math/rand"
	"net/url"
	"regexp"
	"strconv"
	"time"
)

type Selector interface {
	Select() (string, error)
}

// UrlSelector Task 데이터 기반으로 URL로 변환하는 구조체
type UrlSelector struct {
	uri  string
	port int
	task Task
}

func NewUrlSelector(uri string, port int, task Task) *UrlSelector {
	return &UrlSelector{
		uri:  uri,
		port: port,
		task: task,
	}
}

func (us *UrlSelector) Select() (string, error) {
	var baseURI string
	if us.port != 0 {
		baseURI = fmt.Sprintf("%s:%d", us.uri, us.port)
	} else {
		return "", fmt.Errorf("Invalid port number")
	}

	baseURI = baseURI + us.task.Path

	// Get Query Parameters
	if us.task.Method == "GET" {
		params := url.Values{}
		if len(us.task.Query) > 0 {
			for key, value := range us.task.Query {
				params.Add(key, value)
			}
		}

		u, err := url.Parse(baseURI)
		if err != nil {
			return "", err
		}

		u.RawQuery = params.Encode()
		finalUrl := u.String()
		return finalUrl, nil
	} else {
		return baseURI, nil
	}
}

// DynamicUrlSelector 동적 URL을 선택하는 구조체
type DynamicUrlSelector struct {
	url string
}

func NewDynamicUrlSelector(url string) *DynamicUrlSelector {
	return &DynamicUrlSelector{
		url: url,
	}
}

func (dus *DynamicUrlSelector) Select() (string, error) {
	if HasBrace(dus.url) {
		res, err := fixedUrl(dus.url)
		if err != nil {
			return "", err
		}
		return res, nil
	}
	// job이 dynamic url을 가지고 있으면
	return dus.url, nil
}

// FixedUrl
// Dynamic url 범위를 하나로 고정한다.
func fixedUrl(urlStr string) (string, error) {
	// 정규식으로 파싱
	decodedUrl, err := url.QueryUnescape(urlStr)
	if err != nil {
		return "", err
	}
	re := regexp.MustCompile(`\{\[(\d+)-(\d+)\]\}`)
	matches := re.FindStringSubmatch(decodedUrl)
	if len(matches) != 3 {
		return "", fmt.Errorf("Invalid Regex Format")
	}

	// 숫자로 반환
	minR, err1 := strconv.Atoi(matches[1])
	maxR, err2 := strconv.Atoi(matches[2])
	if err1 != nil || err2 != nil {
		return "", fmt.Errorf("Error converting to integer")
	}

	// 랜덤 시드 생성
	rand.Seed(time.Now().UnixNano())
	idx := rand.Intn(maxR-minR+1) + minR

	// 새로운 url
	res := re.ReplaceAllString(decodedUrl, fmt.Sprintf("%d", idx))

	return res, nil
}

// RandomSelector 랜덤 문자열을 선택하는 구조체
type RandomSelector struct {
	length int
}

func NewRandomSelector(length int) *RandomSelector {
	return &RandomSelector{
		length: length,
	}
}

func (rs *RandomSelector) Select() (string, error) {
	if rs.length == 0 {
		return "", fmt.Errorf("Length is 0")
	}
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	b := make([]byte, rs.length)
	for i := range b {
		b[i] = charset[rand.Intn(len(charset))]
	}
	return string(b), nil
}
