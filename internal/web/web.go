package web

import (
	"embed"
	"html/template"
	"io/fs"
	"net/http"
	"os"
)

//go:embed templates/*.html
var templatesFS embed.FS

//go:embed static
var staticFS embed.FS

const siteURL = "https://cuorpugnale.com"
const defaultSpotifyURL = "https://open.spotify.com/search/Cuorpugnale"
const spotifyEmbedURL = "https://open.spotify.com/embed/show/033nh6SpB5aFTDuQ6FMkSI?utm_source=generator"
const youtubeURL = "https://youtube.com/@cuorpugnale?si=GMnp_eG1ujakmclG"
const instagramURL = "https://www.instagram.com/cuorpugnale"

type indexData struct {
	SpotifyEmbed string
	SpotifyURL   string
	YouTubeURL   string
	InstagramURL string
	SiteURL      string
	OGImageURL   string
}

type Server struct {
	tmpl *template.Template
}

func New() (*Server, error) {
	tmpl, err := template.ParseFS(templatesFS, "templates/*.html")
	if err != nil {
		return nil, err
	}
	return &Server{tmpl: tmpl}, nil
}

func (s *Server) Handler() http.Handler {
	mux := http.NewServeMux()

	staticSub, err := fs.Sub(staticFS, "static")
	if err != nil {
		panic(err)
	}
	mux.Handle("GET /static/", securityHeaders(cacheStatic(http.StripPrefix("/static/", http.FileServerFS(staticSub)))))

	mux.HandleFunc("GET /healthz", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("ok"))
	})

	mux.HandleFunc("GET /{$}", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		w.Header().Set("Cache-Control", "no-cache")
		data := indexData{
			SpotifyEmbed: spotifyEmbedURL,
			SpotifyURL:   spotifyURL(),
			YouTubeURL:   youtubeURL,
			InstagramURL: instagramURL,
			SiteURL:      siteURL,
			OGImageURL:   siteURL + "/static/img/campaign/age-of-umbra-1920.jpg",
		}
		if err := s.tmpl.ExecuteTemplate(w, "index.html", data); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	})

	return securityHeaders(mux)
}

func securityHeaders(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Security-Policy", "default-src 'self'; script-src 'self'; style-src 'self'; img-src 'self' data:; frame-src https://open.spotify.com; base-uri 'self'; form-action 'self'; frame-ancestors 'none'")
		w.Header().Set("Referrer-Policy", "strict-origin-when-cross-origin")
		w.Header().Set("X-Content-Type-Options", "nosniff")
		w.Header().Set("X-Frame-Options", "DENY")
		next.ServeHTTP(w, r)
	})
}

func spotifyURL() string {
	if url := os.Getenv("SPOTIFY_URL"); url != "" {
		return url
	}
	return defaultSpotifyURL
}

func cacheStatic(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Cache-Control", "public, max-age=31536000, immutable")
		next.ServeHTTP(w, r)
	})
}
