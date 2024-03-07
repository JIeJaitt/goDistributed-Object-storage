package main

import (
	"goDistributed-Object-storage/apiServer/heartbeat"
	"goDistributed-Object-storage/apiServer/locate"
	"goDistributed-Object-storage/apiServer/objects"
	"goDistributed-Object-storage/apiServer/versions"
	"goDistributed-Object-storage/dataServer/temp"
	"log"
	"net/http"
	"os"
)

func main() {
	go heartbeat.ListenHeartbeat()
	http.HandleFunc("/objects/", objects.Handler)
	http.HandleFunc("/temp/", temp.Handler)
	http.HandleFunc("/locate/", locate.Handler)
	http.HandleFunc("/versions/", versions.Handler)
	log.Fatal(http.ListenAndServe(os.Getenv("LISTEN_ADDRESS"), nil))
}
