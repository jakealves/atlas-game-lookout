package main

import (
	"flag"
	"net/http"
)

func main() {
	webhook := flag.String("webhook", "", "*REQUIRED* the Discord webhook URL to forward requests.")
	flag.Parse()

	if *webhook == "" {
		flag.Usage()
		return
	}

	// webhook handler
	http.HandleFunc("/webhook", func(w http.ResponseWriter, r *http.Request) {
		err := PrintWebhook(r)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		err = RelayWebhook(r, *webhook)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	})
	http.ListenAndServe(":3000", nil)
}
