package server

import (
	"V2RayFree/internal/model"
	"V2RayFree/pkg/utils"
	"log"
	"strings"
	"sync"
)

type FreeFilter[T any] interface {
	Filter() T
}

var protocols = map[string]struct{}{ // 定义一个map来存储支持的协议
	"vmess://":  {},
	"ss://":     {},
	"trojan://": {},
}

type NodeType []model.Node

func (n NodeType) Filter() NodeType {
	log.Printf("开始过滤非 vmess ss trojan 节点, 总节点数: %d", len(n))
	var data NodeType
	for i := range n {
		if _, ok := protocols[strings.SplitN(n[i].Link, "://", 2)[0]+"://"]; ok {
			data = append(data, n[i])
		}
	}
	log.Printf("结束过滤节点, 已过滤出条节点数: %d", len(data))
	return data
}

type NodeTcpTest []model.Node

func (n NodeTcpTest) Filter() []model.Node {
	log.Println("开始TCPTest")
	var wg sync.WaitGroup
	retryMax := 3
	wg.Add(len(n))
	for i := range n {
		go func(i int) {
			for retry := 0; retry < retryMax; retry++ {
				if n[i].TCPTest = utils.TCPTest(utils.NodeToHP(n[i].Link)); n[i].TCPTest > 0 {
					break
				}
			}
			wg.Done()
		}(i)
	}
	wg.Wait()
	log.Println("完成TCPTest")
	return n
}

type NodeFailures []model.Node

func (n NodeFailures) Filter() []model.Node {
	log.Println("开始标记节点失效次数")
	for i := range n {
		if n[i].TCPTest == 0 {
			n[i].Failures++
		}
	}
	log.Println("完成标记节点失效次数")
	return n
}
