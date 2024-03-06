package versions

import (
	"encoding/json"
	"goDistributed-Object-storage/src/lib/es"
	"log"
	"net/http"
	"strings"
)

func Handler(w http.ResponseWriter, r *http.Request) {
	// 检查 HTTP 方法是否为GET
	m := r.Method
	if m != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	from := 0    // 从最早的版本开始获取元数据
	size := 1000 // 每次获取1000个版本的元数据
	// 获取 URL 中<object_name＞的部分
	// 调用 strings.Split 函数将 URL 以“/”分隔符切成数组并取第三个元素赋值给name
	// 如果客户端的 HTTP 请求的 URL 是“/versions/”而不含<object_name>部分，那么 name 就是空字符串
	name := strings.Split(r.URL.EscapedPath(), "/")[2]
	for {
		metas, err := es.SearchAllVersions(name, from, size)
		if err != nil {
			log.Println(err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		for i := range metas {
			b, _ := json.Marshal(metas[i])
			w.Write(b)
			w.Write([]byte("\n"))
		}
		if len(metas) != size {
			return
		}
		from += size
	}
}
