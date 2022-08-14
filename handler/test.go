package handler

import (
	"encoding/json"
	"fmt"
	"net/http"
	"test-assignment-cookie-sync/service"
)

type cookie struct {
	*service.CookieS
}

func NewCookieHandler(c *service.CookieS) cookie {
	return cookie{c}
}

func (c cookie) Cookie(w http.ResponseWriter, r *http.Request) {
	cookie, err := c.CookieS.ProcessCookie(r.Context())
	if err != nil {
		responseBuilder(w, http.StatusInternalServerError)
		return
	}
	b, err := json.Marshal(cookie)
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Write(b)
}

type syncCookie struct {
	*service.SyncS
}

func NewCookieSyncHandler(s *service.SyncS) syncCookie {
	return syncCookie{s}
}

func (s syncCookie) Sync(w http.ResponseWriter, r *http.Request) {
	m, err := readData(r, "dsp_cookie_id", "ssp_cookie_id")
	if err != nil {
		responseBuilder(w, http.StatusBadRequest, fmt.Errorf("parse body: %v", err))
		return
	}
	err = s.SyncS.ProcessSyncCookie(r.Context(), m["dsp_cookie_id"].(string), m["ssp_cookie_id"].(string))
	if err != nil {
		responseBuilder(w, http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Write([]byte("success"))
}
