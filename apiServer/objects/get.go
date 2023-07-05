package objects

import (
	"io"
	"log"
	"net/http"
	"strings"
)

func get(w http.ResponseWriter, r *http.Request) {
	object := strings.Split(r.URL.EscapedPath(), "/")[2]
	// 从getStream中拿到文件的读取流
	stream,err := getStream(object)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusNotFound)
		return
	}
	// 把文件的读取流写入到HTTP响应中
	io.Copy(w,stream)
}