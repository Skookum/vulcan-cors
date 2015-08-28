package cors

import (
	"log"
	"strings"

	"net/http"
)

// Handler executes CORS and handles the middleware chain to the next in stack
type Handler struct {
	cfg  Middleware
	next http.Handler
}

// Runs the CORS specification on the request before passing it to the next middleware chain
func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	h.prepResponse(w)

	if r.Method == optionsMethod {
		h.handlePreflight(w, r)
	} else {
		h.handleRequest(w, r)
	}

	h.next.ServeHTTP(w, r)
}

// Runs the CORS specification for OPTION requests
func (h *Handler) handlePreflight(w http.ResponseWriter, r *http.Request) {
	method := r.Header.Get(requestMethodHeader)
	if method == "" {
		method = r.Method
	}

	h.handleCommon(w, r, method)
}

// Runs the CORS specification for standard requests
func (h *Handler) handleRequest(w http.ResponseWriter, r *http.Request) {
	method := r.Method
	h.handleCommon(w, r, method)
}

// Shares common functionality for prefilght and standard requests
func (h *Handler) handleCommon(w http.ResponseWriter, r *http.Request, method string) {
	origin := r.Header.Get(originHeader)
	if !h.cfg.isOriginAllowed(origin) {
		h.requestDenied(w, errorBadOrigin)
		return
	}

	if !h.cfg.isMethodAllowed(method, origin) {
		h.requestDenied(w, errorBadMethod)
		return
	}

	headers := strings.Split(r.Header.Get(requestHeadersHeader), ",")
	if !h.cfg.areHeadersAllowed(headers, origin) {
		h.requestDenied(w, errorBadHeader)
		return
	}

	h.buildResponse(w, r, origin, method, headers)
}

// Sets the HTTP status to forbidden and logs error message
func (h *Handler) requestDenied(w http.ResponseWriter, m string) {
	log.Println(errorRoot, ": ", m)
	w.WriteHeader(http.StatusForbidden)
	return
}

// Preconfigure headers on the response
func (h *Handler) prepResponse(w http.ResponseWriter) {
	w.Header().Add(varyHeader, originHeader)
}

// Writes the Access Control response headers
func (h *Handler) buildResponse(w http.ResponseWriter, r *http.Request, origin string, method string, headers []string) {
	w.Header().Set(allowOriginHeader, origin)
	w.Header().Set(allowMethodsHeader, method)
	w.Header().Set(allowHeadersHeader, strings.Join(headers, ","))
}