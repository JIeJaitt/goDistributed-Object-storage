package main

import (
	"goDistributed-Object-storage/src/lib/es"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

func main() {
	// 获取对象存储目录下的所有文件
	files, _ := filepath.Glob(os.Getenv("STORAGE_ROOT") + "/objects/*")

	for i := range files {
		// 从文件名中获得对象的散列值
		hash := strings.Split(filepath.Base(files[i]), ".")[0]
		// 判断对象的散列值是否在元数据中
		hashInMetadata, e := es.HasHash(hash)
		if e != nil {
			log.Println(e)
			return
		}
		if !hashInMetadata {
			del(hash)
		}
	}
}

func del(hash string) {
	log.Println("delete", hash)
	url := "http://" + os.Getenv("LISTEN_ADDRESS") + "/objects/" + hash
	request, _ := http.NewRequest("DELETE", url, nil)
	client := http.Client{}
	client.Do(request)
}
