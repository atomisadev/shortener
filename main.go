package main

import (
	"fmt"
	"html/template"
	"net/http"
	"time"
)

type PageData struct {
	URL string
}

func logging(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		currentTime := time.Now().Format("15:04:05")
		fmt.Println("API Called. Method:", r.Method, ", Endpoint:", r.URL.Path, ", Time:", currentTime)
		next.ServeHTTP(w, r)
	})
}

func main() {
	urlService, err := New()
	if err != nil {
		fmt.Println("Error initializing URL shortener:", err)
		return
	}

	mux := http.NewServeMux()

	mux.HandleFunc("/shorten", func(w http.ResponseWriter, r *http.Request) {
		url := r.URL.Query().Get("url")
		if url == "" {
			http.Error(w, "missing url parameter", http.StatusBadRequest)
			return
		}

		alias, err := urlService.Shorten(url)
		if err != nil {
			http.Error(w, "error shortening url", http.StatusInternalServerError)
			return
		}

		shortUrl := "http://" + r.Host + "/" + alias
		w.Write([]byte(shortUrl))
	})

	mux.HandleFunc("/delete", func(w http.ResponseWriter, r *http.Request) {
		alias := r.URL.Query().Get("alias")
		if alias == "" {
			http.Error(w, "missing alias parameter", http.StatusBadRequest)
			return
		}

		err := urlService.Delete(alias)
		if err != nil {
			http.Error(w, "error deleting alias", http.StatusInternalServerError)
			return
		}

		w.Write([]byte("Alias deleted successfully"))
	})

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		alias := r.URL.Path[1:]

		if alias == "" || alias == "favicon.ico" {
			w.Write([]byte("Welcome to the URL shortener service"))
			return
		}

		url, err := urlService.Resolve(alias)
		if err != nil {
			http.Error(w, "Error resolving alias", http.StatusInternalServerError)
			return
		}

		// Insert confirm page here
		tmpl, _ := template.ParseFiles("views/confirm.html")
		data := PageData{
			URL: url,
		}
		tmpl.Execute(w, data)
	})

	loggingMux := logging(mux)

	http.ListenAndServe(":8080", loggingMux)
}
