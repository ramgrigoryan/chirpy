package main

import (
	"fmt"
	"net/http"
)

func (cfg *apiConfig) middlewareMetricsInc(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cfg.fileserverHits.Add(1)
		next.ServeHTTP(w, r)
	})
}

func (c *apiConfig) ReadHits(writer http.ResponseWriter, req *http.Request) {
	writer.Header().Add("application-type", "text/html")
	writer.Write([]byte(fmt.Sprintf(`
	<html>
		<body>
			<h1>Welcome, Chirpy Admin</h1>
			<p>Chirpy has been visited %d times!</p>
		</body>
	</html>
	`, c.fileserverHits.Load())))
}

func (c *apiConfig) Reset(writer http.ResponseWriter, req *http.Request) {
	c.fileserverHits.Swap(0)

	if c.platform != "dev" {
		respondWithErr(writer, http.StatusForbidden, "unable to delete all users", fmt.Errorf("forbidden operation"))
		return
	}

	_, err := c.dbQueries.DeleteUsers(req.Context())
	if err != nil {
		respondWithErr(writer, http.StatusInternalServerError, "unable to delete users", err)
	}
}
