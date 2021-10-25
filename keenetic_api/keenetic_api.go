package keenetic_api

import (
	"bytes"
	"crypto/md5"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"strings"
)

const (
	authPath      = "/auth"
	loginPath     = "/login"
	dashboardPath = "/dashboard"
	rciPath       = "/rci/"
)

var errBadCode = errors.New("keenetic return bad status code")

type StatRQ interface {
	GetRqBody() io.Reader
	Unmarshal(b io.Reader) error
}

type api struct {
	endpoint string
	login    string
	password string

	cl http.Client

	ndmChallenge string
	ndmRealm     string
	cookie       []*http.Cookie
}

func NewApi(endpoint string, login string, password string) *api {
	return &api{
		endpoint: endpoint,
		login:    login,
		password: password,
		cl:       http.Client{},
	}
}

// Auth Авторизация в keenetic, Требуется выполнить перед получением метрик
func (a *api) Auth() error {
	var err error
	if err = a.getAuth(); err != nil {
		return err
	}

	if err = a.doAuth(); err != nil {
		return err
	}

	if err = a.getAuth(); err != nil {
		return err
	}

	return nil
}

// getAuth проверка авторизации у keenetic
func (a *api) getAuth() error {
	var err error
	var rq *http.Request
	if rq, err = http.NewRequest(http.MethodGet, a.endpoint+authPath, nil); err != nil {
		return err
	}

	for i := range a.cookie {
		rq.AddCookie(a.cookie[i])
	}

	rq.Header.Set("Accept", "application/json, text/plain, */*")
	rq.Header.Set("Referer", a.endpoint+loginPath)
	rq.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/95.0.4638.54 Safari/537.36")

	var rs *http.Response
	if rs, err = a.cl.Do(rq); err != nil {
		return err
	}

	switch rs.StatusCode {
	case http.StatusOK:
		a.cookie = append(a.cookie, &http.Cookie{
			Name:  "_authorized",
			Value: "admin",
			Path:  "/",
			Raw:   "_authorized=admin; Path=/",
		}, &http.Cookie{
			Name:  "sysmode",
			Value: "router",
			Path:  "/",
			Raw:   "sysmode=router; Path=/",
		})
		return nil
	case http.StatusUnauthorized:
		a.ndmChallenge = rs.Header.Get("X-NDM-Challenge")
		a.ndmRealm = rs.Header.Get("X-NDM-Realm")
		a.cookie = rs.Cookies()
		return nil
	}

	return errBadCode
}

type AuthJson struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}

// getAuth авторизация в keenetic
func (a *api) doAuth() error {
	var m5 = md5.New()
	m5.Write([]byte(a.login + ":" + a.ndmRealm + ":" + a.password))

	var sh256 = sha256.New()
	sh256.Write([]byte(a.ndmChallenge))
	sh256.Write([]byte(hex.EncodeToString(m5.Sum(nil))))
	sh256.Sum(nil)

	var s = strings.Builder{}
	s.Grow(sha256.BlockSize * 4)
	s.Write([]byte(hex.EncodeToString(sh256.Sum(nil))))

	var b = bytes.NewBuffer(nil)
	b.Grow(s.Len())

	_ = json.NewEncoder(b).Encode(&AuthJson{
		Login:    a.login,
		Password: s.String(),
	})

	var err error
	var rq *http.Request
	if rq, err = http.NewRequest(http.MethodPost, a.endpoint+authPath, b); err != nil {
		return err
	}

	for i := range a.cookie {
		rq.AddCookie(a.cookie[i])
	}

	rq.Header.Set("Accept", "application/json, text/plain, */*")
	rq.Header.Set("Content-Type", "application/json;charset=UTF-8")
	rq.Header.Set("Origin", a.endpoint)
	rq.Header.Set("Referer", a.endpoint+loginPath)
	rq.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/95.0.4638.54 Safari/537.36")

	var rs *http.Response
	if rs, err = a.cl.Do(rq); err != nil {
		return err
	}

	if rs.StatusCode != http.StatusOK {
		return errBadCode
	}

	return nil
}

func (a *api) Metric(q StatRQ) error {
	var err error
	var rq *http.Request
	if rq, err = http.NewRequest(http.MethodPost, a.endpoint+rciPath, q.GetRqBody()); err != nil {
		return err
	}

	for i := range a.cookie {
		rq.AddCookie(a.cookie[i])
	}

	rq.Header.Set("Accept", "application/json, text/plain, */*")
	rq.Header.Set("Origin", a.endpoint)
	rq.Header.Set("Referer", a.endpoint+dashboardPath)
	rq.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/95.0.4638.54 Safari/537.36")

	var rs *http.Response
	if rs, err = a.cl.Do(rq); err != nil {
		return err
	}

	if rs.StatusCode != http.StatusOK {
		return errBadCode
	}
	defer rs.Body.Close()

	if err = q.Unmarshal(rs.Body); err != nil {
		return err
	}

	return nil
}
