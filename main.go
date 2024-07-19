package main

import (
	"io"
	"log"
	"net/http"
)

// handleRequestAndRedirect handles incoming requests and forwards them to the target server
func handleRequestAndRedirect(w http.ResponseWriter, r *http.Request) {
	// Create a new request for the target server
	req, err := http.NewRequest(r.Method, "http://google.com"+r.RequestURI, r.Body)
	if err != nil {
		http.Error(w, "Error creating request", http.StatusInternalServerError)
		return
	}

	// Copy headers from the original request
	for name, values := range r.Header {
		for _, value := range values {
			req.Header.Add(name, value)
		}
	}

	// Send the request to the target server
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		http.Error(w, "Error forwarding request", http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	// Copy the response headers
	for name, values := range resp.Header {
		for _, value := range values {
			w.Header().Add(name, value)
		}
	}

	// Write the status code and body of the response
	w.WriteHeader(resp.StatusCode)
	io.Copy(w, resp.Body)
}

func main() {
	// Set up the proxy server
	http.HandleFunc("/", handleRequestAndRedirect)
	port := "8080"
	log.Printf("Starting proxy server on port %s\n", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
