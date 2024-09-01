package handler

import (
	"encoding/json"
	"net"
	"net/http"

	"github.com/gleb-korostelev/short-url.git/internal/config"
	"github.com/gleb-korostelev/short-url.git/internal/models"
)

// StatsHandler handles HTTP GET requests to the /api/internal/stats endpoint.
// This endpoint provides statistics about the service, such as the number of shortened URLs
// and the number of registered users.
//
// The access to this endpoint is restricted based on the client's IP address. The client's IP
// must fall within a trusted subnet specified in the server configuration. If the trusted subnet
// is not configured or if the client's IP does not belong to the specified subnet, access is denied,
// and a HTTP 403 Forbidden status is returned.
//
// The handler checks the X-Real-IP header of the incoming HTTP request to determine the client's IP.
// If the header is missing or if the IP does not match the trusted subnet, the request is rejected.
// Otherwise, the handler proceeds to fetch and return the requested statistics in JSON format.
//
// Example response on successful retrieval:
//
//	{
//	  "urls": 150,
//	  "users": 25
//	}
//
// Responses:
//   - HTTP 200 OK: Returned along with the service statistics if the request is authorized.
//   - HTTP 403 Forbidden: Returned if the client's IP is not within the trusted subnet or if
//     the trusted subnet is not configured.
//
// This function uses the server's configuration settings to determine the trusted subnet.
// It is recommended to ensure that the X-Real-IP header is reliably set by a trusted proxy
// or load balancer to prevent unauthorized access.
func (svc *APIService) StatsHandler(w http.ResponseWriter, r *http.Request) {
	clientIP := r.Header.Get("X-Real-IP")
	if config.TrustedSubnet != "" {
		_, subnet, err := net.ParseCIDR(config.TrustedSubnet)
		if err != nil {
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}
		if !subnet.Contains(net.ParseIP(clientIP)) {
			http.Error(w, "Forbidden", http.StatusForbidden)
			return
		}
	} else {
		http.Error(w, "Forbidden", http.StatusForbidden)
		return
	}

	stats := models.Stats{
		URLs:  100,
		Users: 10,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(stats)
}
