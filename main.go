package main

import (
	"encoding/base64"
	"log/slog"
	"net/http"
	"net/url"
	"os"
	"strings"
)

const defaultPort = ":80"

func redirectHandler(l *slog.Logger) func(http.ResponseWriter, *http.Request) {
	l.Info("handler created")
	return func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/" {
			w.WriteHeader(http.StatusOK)
			w.Write([]byte("OK"))
			return
		}
		logger := l.With(
			slog.String("ip", r.RemoteAddr),
			slog.String("user-agent", r.UserAgent()),
			slog.String("referer", r.Referer()),
		)

		if r.Method != http.MethodGet {
			logger.Error("invalid method", slog.String("method", r.Method))
			http.Error(w, "Invalid method", http.StatusMethodNotAllowed)
			return
		}

		data := strings.TrimPrefix(r.URL.Path, "/")
		decodedData, err := base64.StdEncoding.DecodeString(data)
		if err != nil {
			logger.Error("invalid decoded data", slog.String("data", data), slog.String("error", err.Error()))
			http.Error(w, "Invalid data", http.StatusBadRequest)
			return
		}

		link, err := url.Parse(string(decodedData))
		if err != nil {
			logger.Error("invalid data", slog.String("data", string(decodedData)), slog.String("error", err.Error()))
			http.Error(w, "Invalid data", http.StatusBadRequest)
			return
		}

		if link.Scheme == "" && link.Host == "" {
			logger.Error("empty scheme or host")
			http.Error(w, "Invalid data", http.StatusBadRequest)
			return
		}

		logger.Info("redirecting", slog.String("link", link.String()))
		http.Redirect(w, r, link.String(), http.StatusMovedPermanently)
	}
}
func main() {
	opts := &slog.HandlerOptions{
		Level: slog.LevelDebug,
	}
	textHandler := slog.NewTextHandler(os.Stdout, opts)
	logger := slog.New(textHandler)

	logger.Info("starting server", slog.String("port", defaultPort))

	http.HandleFunc("/", redirectHandler(logger))
	if err := http.ListenAndServe(defaultPort, nil); err != nil {
		logger.Error("failed to start server", slog.String("error", err.Error()))
	}
}
