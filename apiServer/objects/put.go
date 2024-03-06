package objects

import (
	"goDistributed-Object-storage/src/lib/es"
	"goDistributed-Object-storage/src/lib/utils"
	"log"
	"net/http"
	"net/url"
	"strings"
)

/* chapter2
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
*/

func put(w http.ResponseWriter, r *http.Request) {
	// 从请求头中获取对象的哈希值
	hash := utils.GetHashFromHeader(r.Header)
	if hash == "" {
		log.Println("missing object hash in digest header")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// 将对象存储到存储系统中
	statusCode, err := storeObject(r.Body, url.PathEscape(hash))
	if err != nil {
		log.Println(err)
		w.WriteHeader(statusCode)
		return
	}
	if statusCode != http.StatusOK {
		w.WriteHeader(statusCode)
		return
	}

	// 从 URL 中获取对象的名称
	name := strings.Split(r.URL.EscapedPath(), "/")[2]
	// 从请求头中获取对象的大小
	size := utils.GetSizeFromHeader(r.Header)
	// 将对象的元数据添加到 Elasticsearch 中
	err = es.AddVersion(name, hash, size)
	// 异步方式将对象的元数据添加到 Elasticsearch 中
	// go func() {
	// 	err = es.AddVersion(name, hash, size)
	// 	if err != nil {
	// 		log.Println(err)
	// 	}
	// }()
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
	}

	// 返回 HTTP 状态码 200 表示成功
	w.WriteHeader(http.StatusOK)
}
