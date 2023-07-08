package main

import (
	"embed"
	"html/template"
	"log"
	"math/rand"
	"net/http"
	"os"
	"time"

	"github.com/devries/nlm"
)

//go:embed templates/*.html
var templateFiles embed.FS

//go:embed static
var staticFiles embed.FS

const titlesize = 5
const size = 6

func main() {
	rand.Seed(time.Now().UnixNano())
	mux := http.NewServeMux()

	var bind string
	switch len(os.Args) {
	case 1:
		bind = ":8080"
	case 2:
		bind = os.Args[1]
	default:
		log.Fatalf("Usage: %s [bind]", os.Args[0])
	}

	templates := template.Must(template.New("web").ParseFS(templateFiles, "templates/*"))

	log.Printf("Loading articles...")
	start := time.Now()
	articleBuilder, err := nlm.NewArticleBuilder(titlesize, size)
	if err != nil {
		log.Fatalf("Error creating article builder: %s", err)
	}
	elapsed := time.Now().Sub(start)
	log.Printf("Time elapsed loading articles: %s", elapsed)

	mux.Handle("/", setCacheControl("max-age=3600", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/" {
			http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
			return
		}
		indexTemplate := templates.Lookup("index.html")
		err := indexTemplate.Execute(w, nil)
		if err != nil {
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		}
	})))

	mux.Handle("/about", setCacheControl("max-age=3600", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		aboutTemplate := templates.Lookup("about.html")
		err := aboutTemplate.Execute(w, size)
		if err != nil {
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		}
	})))

	mux.Handle("/generate", NewRateLimiter("X-Forwarded-For", 0.2, 4, setCacheControl("no-cache, no-store", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		responseTemplate := templates.Lookup("article.html")
		errorTemplate := templates.Lookup("error.html")

		article := articleBuilder.GenerateArticle(120, 15000)

		err = responseTemplate.Execute(w, article)
		if err != nil {
			log.Printf("error writing article template: %s", err)
			errorTemplate.Execute(w, "Unable to render article")
		}
	}))))

	mux.Handle("/speed", setCacheControl("no-cache, no-store", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		errorTemplate := templates.Lookup("error.html")
		w.Header().Add("HX-Retarget", "#content")
		w.Header().Add("HX-Reswap", "innerHTML")
		w.Header().Add("HX-Replace-Url", "/")

		errorTemplate.Execute(w, "Too many requests, please slow down.")
	})))

	mux.Handle("/static/", setCacheControl("max-age=1209600", http.FileServer(http.FS(staticFiles))))

	server := &http.Server{
		Addr:    bind,
		Handler: loggingHandler(mux),
	}

	log.Printf("Server starting on %s", bind)
	log.Fatal(server.ListenAndServe())
}

func setCacheControl(flag string, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Cache-Control", flag)
		next.ServeHTTP(w, r)
	})
}
