package service

import (
	"V2RayFree/internal/db"
	"V2RayFree/internal/model"
	"encoding/base64"
	"fmt"
	"net/http"
	"strings"
)

func SubscribeHandle(res http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodGet {
		fmt.Fprintf(res, "只接受Get请求")
	}
	ns := model.QueryAllNode(db.DB)
	ls := make([]string, len(ns))
	for i := range ns {
		ls[i] = ns[i].Link
	}
	fmt.Fprintf(res, base64.StdEncoding.EncodeToString([]byte(strings.Join(ls, "\n"))))
}
