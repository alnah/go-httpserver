package main

import (
	"fmt"
	"net/http"
)

// handlerMetrics serves the metrics page for administrators.
// It sets the content type to HTML and displays the current file server hit count.
func (cfg *apiConfig) handlerMetrics(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "text/html")
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write(fmt.Appendf(nil, `
<html>

<body>
	<h1>Welcome, Chirpy Admin</h1>
	<p>Chirpy has been visited %d times!</p>
</body>

</html>
	`, cfg.fileserverHits.Load()))
}

// middlewareMetricsInc is an HTTP middleware that increments the file server hit counter
// before calling the next handler in the chain.
func (cfg *apiConfig) middlewareMetricsInc(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cfg.fileserverHits.Add(1)
		next.ServeHTTP(w, r)
	})
}
