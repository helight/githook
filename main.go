package main

import (
	"log"
	"net/http"

	"githook/webhook"
)

func test(w http.ResponseWriter, r *http.Request) {
	log.Printf("xxxxxx: %s", r.URL.Path)
	w.Write([]byte(r.URL.Path))
}

func main() {
	// To verify webhook's payload, set secret by SetSecret().
	webhook.SetSecret([]byte("abcdefgh"))

	http.HandleFunc("/test", test)

	// Add a HandlerFunc to process webhook.
	http.HandleFunc("/", webhook.HandlePush(func(ev *webhook.Event) {
		push := ev.PushEvent()
		if push == nil {
			return
		}
		log.Printf("push: verified=%v %#v", ev.Verified, push)
	}))

	// Start web server.
	log.Fatal(http.ListenAndServe(":8080", nil))
}
