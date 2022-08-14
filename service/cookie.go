package service

import (
	"context"
	"net/http"
	"test-assignment-cookie-sync/connector"

	satoriuuid "github.com/satori/go.uuid"
)

const (
	COOKIE = "cookie"
	SYNC   = "sync"
)

type CookieS struct {
	conn connector.DB
}

func NewCookie(conn connector.DB) *CookieS {
	return &CookieS{conn}
}

func (c *CookieS) ProcessCookie(ctx context.Context) (*http.Cookie, error) {
	cid := satoriuuid.NewV4().String()
	return &http.Cookie{
		Name:   "dsp",
		Value:  cid,
		MaxAge: 0,
	}, c.conn.Cookie().PersistCookie(ctx, satoriuuid.NewV4().String(), cid)
}

type SyncS struct {
	conn connector.DB
}

func NewSync(conn connector.DB) *SyncS {
	return &SyncS{conn}
}

func (s *SyncS) ProcessSyncCookie(ctx context.Context, dspCookieId, sspCookieId string) error {
	return s.conn.Cookie().UpdateCookie(ctx, dspCookieId, sspCookieId)
}
