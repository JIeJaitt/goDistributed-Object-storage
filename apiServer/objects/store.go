package objects

import (
	"io"
	"net/http"
)

func storeObject(r io.Reader, object string) (int, error) {
	// 调用putStream函数，得到一个stream
	// putStream是object的一个文件流用于写入
	stream,err:= putStream(object)
	if err != nil {
		return http.StatusServiceUnavailable, err
	}
	// 把HTTP请求正文（文件内容）写入这个stream
	io.Copy(stream, r)
	err = stream.Close()
	if err != nil {
		return http.StatusInternalServerError,err
	}
	return http.StatusOK,err
}
