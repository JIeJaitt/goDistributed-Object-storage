package objects

import (
	"log"
	"net/http"
	"strings"
)

func put(w http.ResponseWriter, r *http.Request) {
	// 拿到对象名
	object := strings.Split(r.URL.EscapedPath(), "/")[2]
	// 把对象的内容保存到对象的名字里面去
	c, err := storeObject(r.Body, object)
	if err != nil {
		log.Println(err)
	}
	// 返回状态码
	w.WriteHeader(c)
}