package main

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"strconv"

	"github.com/gcrtnst/sw-term-server/internal/vterm"
)

type ServiceHandler struct {
	Service Service
}

func (h *ServiceHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" && r.Method != "" {
		resp := &ServiceResponse{
			Code: http.StatusMethodNotAllowed,
			Body: []byte("method not allowed"),
		}
		_ = resp.WriteResponse(w)
		return
	}

	query, err := url.ParseQuery(r.URL.RawQuery)
	if err != nil {
		resp := &ServiceResponse{
			Code: http.StatusBadRequest,
			Body: []byte("invalid url query"),
		}
		_ = resp.WriteResponse(w)
		return
	}

	resp := h.Service.ServeAPI(query)
	_ = resp.WriteResponse(w)
}

type Service interface {
	ServeAPI(url.Values) *ServiceResponse
}

type KeyboardService struct {
	TermSlot *TermSlot
	Logger   *log.Logger
}

func (srv *KeyboardService) ServeAPI(query url.Values) *ServiceResponse {
	queryKey := query.Get("key")
	if queryKey == "" {
		return &ServiceResponse{
			Code: http.StatusBadRequest,
			Body: []byte(`missing parameter "key"`),
		}
	}
	key := Key(queryKey)

	mod := vterm.ModNone
	queryMod := query.Get("mod")
	if queryMod != "" {
		n, err := strconv.ParseUint(queryMod, 10, 8)
		if err != nil {
			s := fmt.Sprintf(`failed to parse parameter "mod": %s`, err.Error())
			return &ServiceResponse{
				Code: http.StatusBadRequest,
				Body: []byte(s),
			}
		}

		mod = vterm.Modifier(n)
	}

	err := srv.TermSlot.Keyboard(key, mod)
	if errors.Is(err, ErrInvalidKey) {
		s := err.Error()
		return &ServiceResponse{
			Code: http.StatusBadRequest,
			Body: []byte(s),
		}
	}
	if err != nil {
		srv.Logger.Printf("error: %s", err.Error())
		return &ServiceResponse{
			Code: http.StatusInternalServerError,
			Body: []byte("internal server error"),
		}
	}

	return &ServiceResponse{
		Code: http.StatusOK,
		Body: []byte{},
	}
}

type ScreenService struct {
	TermSlot *TermSlot
	Logger   *log.Logger
}

func (srv *ScreenService) ServeAPI(query url.Values) *ServiceResponse {
	ss, err := srv.TermSlot.Capture()
	if err != nil {
		srv.Logger.Printf("error: %s", err.Error())
		return &ServiceResponse{
			Code: http.StatusInternalServerError,
			Body: []byte("internal server error"),
		}
	}

	var b []byte
	b = EncodeScreenShot(ss)
	b = EscapeZero(b)
	b = append([]byte("%SWTSCRN"), b...)

	return &ServiceResponse{
		Code: http.StatusOK,
		Body: b,
	}
}

type StopService struct {
	TermSlot *TermSlot
}

func (srv *StopService) ServeAPI(query url.Values) *ServiceResponse {
	srv.TermSlot.Stop()
	return &ServiceResponse{
		Code: http.StatusOK,
		Body: []byte{},
	}
}

type ServiceResponse struct {
	Code int
	Body []byte
}

func (r *ServiceResponse) WriteResponse(w http.ResponseWriter) error {
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.Header().Set("X-Content-Type-Options", "nosniff")
	w.Header().Set("Content-Length", strconv.Itoa(len(r.Body)))
	w.WriteHeader(r.Code)
	_, err := w.Write(r.Body)
	return err
}
