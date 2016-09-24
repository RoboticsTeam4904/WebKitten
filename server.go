package main

import "net/http"

func StartServer(port string, hub *Hub) {
	serveSingle("/", "./www/index.html")
	serveSingle("/index.js", "./www/index.js")
	http.HandleFunc("/assets/", http.StripPrefix("/assets/", http.FileServer(http.Dir("./www/assets"))))
	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.RequestWriter) {
		serveWS(hub, w, r)
	})

	panic(http.ListenAndServe(":"+port, nil))
}

func serveSingle(pattern string, filename string) {
	http.HandleFunc(pattern, func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, filename)
	})
}
