package server

import (
	"V2RayFree/internal/model"
	"encoding/base64"
	"io"
	"log"
	"net/http"
	"strings"
	"time"
)

type FreeConverter interface {
	Converter([]byte) []string
	~string
}

type SubscribeLink string

func (SubscribeLink) Converter(data []byte) []string {
	return strings.Split(string(data), "\n")
}

type SubscribeBase64 string

func (SubscribeBase64) Converter(data []byte) []string {
	var err error
	if data, err = base64.StdEncoding.DecodeString(string(data)); err != nil {
		return nil
	}
	return new(SubscribeLink).Converter(data)
}

var (
	client = http.Client{
		Timeout: time.Second * 2,
	}
	retryNumber = 5
	retryTime   = time.Second * 2
)

func NewFreeConverter[T FreeConverter](f T) []model.Node {
	var err error
	var resp *http.Response
	log.Printf("开始请求: %s", f)
	for retry := 0; ; retry++ {
		if retry >= retryNumber {
			return nil
		} else if resp, err = client.Get(string(f)); err == nil && resp.StatusCode == http.StatusOK {
			break
		}
		log.Printf("请求 %s 失败, 进行第%d次尝试: %v", f, retry, err)
		time.Sleep(retryTime)
	}
	defer resp.Body.Close()
	var body []byte
	if body, err = io.ReadAll(resp.Body); err != nil {
		return nil
	}
	return model.NewNodes(f.Converter(body))
}
