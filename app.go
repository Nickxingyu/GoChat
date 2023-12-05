package main

import (
	"flag"
	"log"
	"net/http"
)

var addr = flag.String("addr", ":3000", "http server address")

func main() {
	flag.Parse()

	wsServer := GetWsServer()
	go wsServer.Run()
	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		wsServer.ServeWs(w, r)
	})

	log.Fatal(http.ListenAndServe(*addr, nil))
}
