package server

import (
	"V2RayFree/internal/model"
	"V2RayFree/pkg/utils"
	"strings"
)

type FreeFilter[T any] interface {
	Filter() T
}

type NodeType []model.Node

func (n NodeType) Filter() NodeType {
	var data NodeType
	for i := range n {
		if strings.HasPrefix(n[i].Link, "vmess://") ||
			strings.HasPrefix(n[i].Link, "ss://") ||
			strings.HasPrefix(n[i].Link, "trojan://") {
			data = append(data, n[i])
		}
	}
	return data
}

type NodeTcpTest []model.Node

func (n NodeTcpTest) Filter() []model.Node {
	for i := range n {
		n[i].TCPTest = utils.TCPTest(utils.NodeToHP(n[i].Link))
	}
	return n
}
