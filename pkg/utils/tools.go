package utils

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net"
	"regexp"
	"strings"
	"time"
)

// Ticker 定时器
func Ticker(t uint32, f func()) chan struct{} {
	f()
	ticker := time.NewTicker(time.Duration(t) * time.Second)
	stop := make(chan struct{})
	go func() {
		defer ticker.Stop()
		for {
			select {
			case <-ticker.C:
				f()
			case <-stop:
				return
			}
		}
	}()
	return stop
}

func NodeToHP(link string) string {
	var err error
	switch {
	case strings.HasPrefix(link, "vmess://"):
		var b []byte
		if b, err = base64.StdEncoding.DecodeString(strings.TrimPrefix(link, "vmess://")); err == nil {
			var data map[string]any
			if err = json.Unmarshal(b, &data); err == nil {
				return fmt.Sprintf("%v:%v", data["add"], data["port"])
			}
		}
	case strings.HasPrefix(link, "ss://"):
		if data := regexp.MustCompile("@(.*?):[0-9]+#").FindStringSubmatch(link)[0]; len(data) >= 2 {
			return data[1 : len(data)-1]
		}
	case strings.HasPrefix(link, "trojan://"):
		if data := regexp.MustCompile("@(.*?):[0-9]+\\?").FindStringSubmatch(link)[0]; len(data) >= 2 {
			return data[1 : len(data)-1]
		}
	}
	return ""
}

func TCPTest(hp string) uint16 {
	startTime := time.Now()
	conn, err := net.DialTimeout("tcp", hp, 3*time.Second)
	if err != nil {
		return 0
	}
	defer func(conn net.Conn) {
		_ = conn.Close()
	}(conn)
	return uint16(time.Since(startTime).Milliseconds())
}
