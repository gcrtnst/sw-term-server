package main

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"reflect"
	"strconv"
	"testing"

	"github.com/gcrtnst/sw-term-server/internal/xpty"
)

func TestServiceHandlerServeHTTP(t *testing.T) {
	tt := []struct {
		name           string
		inReq          *http.Request
		inResp         *ServiceResponse
		wantSvcQuery   url.Values
		wantRespCode   int
		wantRespHeader http.Header
		wantRespBody   []byte
	}{
		{
			name:  "Normal",
			inReq: httptest.NewRequest("GET", "/path/to/api?key=value", nil),
			inResp: &ServiceResponse{
				Code: http.StatusOK,
				Body: []byte("test body"),
			},
			wantSvcQuery: url.Values{
				"key": []string{"value"},
			},
			wantRespCode: http.StatusOK,
			wantRespHeader: http.Header{
				"Content-Type":           []string{"text/plain; charset=utf-8"},
				"X-Content-Type-Options": []string{"nosniff"},
				"Content-Length":         []string{"9"},
			},
			wantRespBody: []byte("test body"),
		},
		{
			name:  "MethodNotAllowed",
			inReq: httptest.NewRequest("POST", "/path/to/api?key=value", nil),
			inResp: &ServiceResponse{
				Code: http.StatusOK,
				Body: []byte("test body"),
			},
			wantSvcQuery: nil,
			wantRespCode: http.StatusMethodNotAllowed,
			wantRespHeader: http.Header{
				"Content-Type":           []string{"text/plain; charset=utf-8"},
				"X-Content-Type-Options": []string{"nosniff"},
				"Content-Length":         []string{"18"},
			},
			wantRespBody: []byte("method not allowed"),
		},
		{
			name:  "InvalidQuery",
			inReq: httptest.NewRequest("GET", "/path/to/api?key=value%", nil),
			inResp: &ServiceResponse{
				Code: http.StatusOK,
				Body: []byte("test body"),
			},
			wantSvcQuery: nil,
			wantRespCode: http.StatusBadRequest,
			wantRespHeader: http.Header{
				"Content-Type":           []string{"text/plain; charset=utf-8"},
				"X-Content-Type-Options": []string{"nosniff"},
				"Content-Length":         []string{"17"},
			},
			wantRespBody: []byte("invalid url query"),
		},
	}

	for _, tc := range tt {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			srv := &MockService{Resp: tc.inResp}
			rec := httptest.NewRecorder()
			handler := &ServiceHandler{Service: srv}
			handler.ServeHTTP(rec, tc.inReq)

			if !reflect.DeepEqual(srv.Query, tc.wantSvcQuery) {
				t.Errorf("srv query: expected %#v, got %#v", tc.wantSvcQuery, srv.Query)
			}

			gotResp := rec.Result()
			defer gotResp.Body.Close()
			gotBody, _ := io.ReadAll(gotResp.Body)

			if gotResp.StatusCode != tc.wantRespCode {
				t.Errorf("resp code: expected %d, got %d", tc.wantRespCode, gotResp.StatusCode)
			}
			for key := range tc.wantRespHeader {
				if !reflect.DeepEqual(gotResp.Header[key], tc.wantRespHeader[key]) {
					t.Errorf("resp header %#v: expected %#v, got %#v", key, tc.wantRespHeader[key], gotResp.Header[key])
				}
			}
			if !bytes.Equal(gotBody, tc.wantRespBody) {
				t.Errorf("resp body: expected %#v, got %#v", string(tc.wantRespBody), string(gotBody))
			}
		})
	}
}

func TestKeyboardServiceServeAPI(t *testing.T) {
	pid := os.Getpid()

	tt := []struct {
		name      string
		inQuery   url.Values
		inErrOpen error
		wantResp  *ServiceResponse
		wantLog   []byte
		wantMTOut []byte
	}{
		{
			name: "Normal",
			inQuery: url.Values{
				"key": []string{"A"},
				"mod": []string{"6"},
			},
			inErrOpen: nil,
			wantResp: &ServiceResponse{
				Code: http.StatusOK,
				Body: []byte{},
			},
			wantLog:   []byte{},
			wantMTOut: []byte("\x1B[65;7u"),
		},
		{
			name: "OmitMod",
			inQuery: url.Values{
				"key": []string{"A"},
			},
			inErrOpen: nil,
			wantResp: &ServiceResponse{
				Code: http.StatusOK,
				Body: []byte{},
			},
			wantLog:   []byte{},
			wantMTOut: []byte("A"),
		},
		{
			name: "MissingKey",
			inQuery: url.Values{
				"mod": []string{"6"},
			},
			inErrOpen: nil,
			wantResp: &ServiceResponse{
				Code: http.StatusBadRequest,
				Body: []byte(`missing parameter "key"`),
			},
			wantLog:   []byte{},
			wantMTOut: []byte{},
		},
		{
			name: "InvalidKey",
			inQuery: url.Values{
				"key": []string{"Invalid"},
				"mod": []string{"6"},
			},
			inErrOpen: nil,
			wantResp: &ServiceResponse{
				Code: http.StatusBadRequest,
				Body: []byte(ErrInvalidKey.Error()),
			},
			wantLog:   []byte{},
			wantMTOut: []byte{},
		},
		{
			name: "InvalidMod",
			inQuery: url.Values{
				"key": []string{"A"},
				"mod": []string{"A"},
			},
			inErrOpen: nil,
			wantResp: &ServiceResponse{
				Code: http.StatusBadRequest,
				Body: []byte(fmt.Sprintf(`failed to parse parameter "mod": %s`, &strconv.NumError{
					Func: "ParseUint",
					Num:  "A",
					Err:  strconv.ErrSyntax,
				})),
			},
			wantLog:   []byte{},
			wantMTOut: []byte{},
		},
		{
			name: "ErrOpen",
			inQuery: url.Values{
				"key": []string{"A"},
				"mod": []string{"6"},
			},
			inErrOpen: errors.New("dummy error"),
			wantResp: &ServiceResponse{
				Code: http.StatusInternalServerError,
				Body: []byte("internal server error"),
			},
			wantLog:   []byte("error: dummy error\n"),
			wantMTOut: []byte{},
		},
	}

	for _, tc := range tt {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			mt := &xpty.MockTerminal{
				ErrOpen: tc.inErrOpen,
				PID:     pid,
			}
			cfg := TermConfig{
				Open: mt.Open,
				Row:  30,
				Col:  120,
				Cmd: xpty.Cmd{
					Path: "bash",
					Args: []string{"--version"},
				},
			}
			slot := NewTermSlot(cfg)

			logbuf := new(bytes.Buffer)
			logger := log.New(logbuf, "", 0)

			srv := &KeyboardService{
				TermSlot: slot,
				Logger:   logger,
			}

			gotResp := srv.ServeAPI(tc.inQuery)
			gotLog := logbuf.Bytes()

			mt.ErrOpen = nil
			slot.start()
			slot.term.pc = nil
			slot.Stop()
			gotMTOut, _ := io.ReadAll(mt.Computer())

			if gotResp.Code != tc.wantResp.Code {
				t.Errorf("resp code: expected %d, got %d", tc.wantResp.Code, gotResp.Code)
			}
			if !bytes.Equal(gotResp.Body, tc.wantResp.Body) {
				t.Errorf("resp body: expected %#v, got %#v", tc.wantResp.Body, gotResp.Body)
			}
			if !bytes.Equal(gotLog, tc.wantLog) {
				t.Errorf("log: expected %#v, got %#v", tc.wantLog, gotLog)
			}
			if !bytes.Equal(gotMTOut, tc.wantMTOut) {
				t.Errorf("mt out: expected %#v, got %#v", tc.wantMTOut, gotMTOut)
			}
		})
	}
}

func TestScreenServiceServeAPI(t *testing.T) {
	errDummy := errors.New("dummy error")
	pid := os.Getpid()

	tt := []struct {
		name        string
		inStart     bool
		inIn        []byte
		inMTErrOpen error
		wantResp    *ServiceResponse
		wantLog     []byte
	}{
		{
			name:        "Normal",
			inStart:     true,
			inIn:        []byte("AB\r\nC"),
			inMTErrOpen: nil,
			wantResp: &ServiceResponse{
				Code: http.StatusOK,
				Body: []byte{
					'%', 'S', 'W', 'T', 'S', 'C', 'R', 'N', // (signature)

					0x01,                                                                                     // CursorVisible
					0x01,                                                                                     // CursorBlink
					0x01,                                                                                     // CursorShape
					0x01, 0x5C, 0x30, 0x5C, 0x30, 0x5C, 0x30, 0x5C, 0x30, 0x5C, 0x30, 0x5C, 0x30, 0x5C, 0x30, // CursorPos.Row
					0x01, 0x5C, 0x30, 0x5C, 0x30, 0x5C, 0x30, 0x5C, 0x30, 0x5C, 0x30, 0x5C, 0x30, 0x5C, 0x30, // CursorPos.Col
					0x02, 0x5C, 0x30, 0x5C, 0x30, 0x5C, 0x30, 0x5C, 0x30, 0x5C, 0x30, 0x5C, 0x30, 0x5C, 0x30, // Stride

					0x5C, 0x30, // Cell[i].Attrs.Bold
					0x5C, 0x30, // Cell[i].Attrs.Underline
					0x5C, 0x30, // Cell[i].Attrs.Italic
					0x5C, 0x30, // Cell[i].Attrs.Blink
					0x5C, 0x30, // Cell[i].Attrs.Reverse
					0x5C, 0x30, // Cell[i].Attrs.Conceal
					0x5C, 0x30, // Cell[i].Attrs.Strike
					0x5C, 0x30, // Cell[i].Attrs.Font
					0x5C, 0x30, // Cell[i].Attrs.DWL
					0x5C, 0x30, // Cell[i].Attrs.DHL
					0x5C, 0x30, // Cell[i].Attrs.Small
					0x5C, 0x30, // Cell[i].Attrs.Baseline
					0x03,       // Cell[i].FG.Type
					0x07,       // Cell[i].FG.Idx
					0x5C, 0x30, // (padding)
					0x5C, 0x30, // (padding)
					0x05,       // Cell[i].BG.Type
					0x5C, 0x30, // Cell[i].BG.Idx
					0x5C, 0x30, // (padding)
					0x5C, 0x30, // (padding)
					0x01,                                                                                     // Cell[i].Width
					0x01, 0x5C, 0x30, 0x5C, 0x30, 0x5C, 0x30, 0x5C, 0x30, 0x5C, 0x30, 0x5C, 0x30, 0x5C, 0x30, // len(string(Cell[i].Runes))
					0x41, // string(Cell[i].Runes)

					0x5C, 0x30, // Cell[i].Attrs.Bold
					0x5C, 0x30, // Cell[i].Attrs.Underline
					0x5C, 0x30, // Cell[i].Attrs.Italic
					0x5C, 0x30, // Cell[i].Attrs.Blink
					0x5C, 0x30, // Cell[i].Attrs.Reverse
					0x5C, 0x30, // Cell[i].Attrs.Conceal
					0x5C, 0x30, // Cell[i].Attrs.Strike
					0x5C, 0x30, // Cell[i].Attrs.Font
					0x5C, 0x30, // Cell[i].Attrs.DWL
					0x5C, 0x30, // Cell[i].Attrs.DHL
					0x5C, 0x30, // Cell[i].Attrs.Small
					0x5C, 0x30, // Cell[i].Attrs.Baseline
					0x03,       // Cell[i].FG.Type
					0x07,       // Cell[i].FG.Idx
					0x5C, 0x30, // (padding)
					0x5C, 0x30, // (padding)
					0x05,       // Cell[i].BG.Type
					0x5C, 0x30, // Cell[i].BG.Idx
					0x5C, 0x30, // (padding)
					0x5C, 0x30, // (padding)
					0x01,                                                                                     // Cell[i].Width
					0x01, 0x5C, 0x30, 0x5C, 0x30, 0x5C, 0x30, 0x5C, 0x30, 0x5C, 0x30, 0x5C, 0x30, 0x5C, 0x30, // len(string(Cell[i].Runes))
					0x42, // string(Cell[i].Runes)

					0x5C, 0x30, // Cell[i].Attrs.Bold
					0x5C, 0x30, // Cell[i].Attrs.Underline
					0x5C, 0x30, // Cell[i].Attrs.Italic
					0x5C, 0x30, // Cell[i].Attrs.Blink
					0x5C, 0x30, // Cell[i].Attrs.Reverse
					0x5C, 0x30, // Cell[i].Attrs.Conceal
					0x5C, 0x30, // Cell[i].Attrs.Strike
					0x5C, 0x30, // Cell[i].Attrs.Font
					0x5C, 0x30, // Cell[i].Attrs.DWL
					0x5C, 0x30, // Cell[i].Attrs.DHL
					0x5C, 0x30, // Cell[i].Attrs.Small
					0x5C, 0x30, // Cell[i].Attrs.Baseline
					0x03,       // Cell[i].FG.Type
					0x07,       // Cell[i].FG.Idx
					0x5C, 0x30, // (padding)
					0x5C, 0x30, // (padding)
					0x05,       // Cell[i].BG.Type
					0x5C, 0x30, // Cell[i].BG.Idx
					0x5C, 0x30, // (padding)
					0x5C, 0x30, // (padding)
					0x01,                                                                                     // Cell[i].Width
					0x01, 0x5C, 0x30, 0x5C, 0x30, 0x5C, 0x30, 0x5C, 0x30, 0x5C, 0x30, 0x5C, 0x30, 0x5C, 0x30, // len(string(Cell[i].Runes))
					0x43, // string(Cell[i].Runes)

					0x5C, 0x30, // Cell[i].Attrs.Bold
					0x5C, 0x30, // Cell[i].Attrs.Underline
					0x5C, 0x30, // Cell[i].Attrs.Italic
					0x5C, 0x30, // Cell[i].Attrs.Blink
					0x5C, 0x30, // Cell[i].Attrs.Reverse
					0x5C, 0x30, // Cell[i].Attrs.Conceal
					0x5C, 0x30, // Cell[i].Attrs.Strike
					0x5C, 0x30, // Cell[i].Attrs.Font
					0x5C, 0x30, // Cell[i].Attrs.DWL
					0x5C, 0x30, // Cell[i].Attrs.DHL
					0x5C, 0x30, // Cell[i].Attrs.Small
					0x5C, 0x30, // Cell[i].Attrs.Baseline
					0x03,       // Cell[i].FG.Type
					0x07,       // Cell[i].FG.Idx
					0x5C, 0x30, // (padding)
					0x5C, 0x30, // (padding)
					0x05,       // Cell[i].BG.Type
					0x5C, 0x30, // Cell[i].BG.Idx
					0x5C, 0x30, // (padding)
					0x5C, 0x30, // (padding)
					0x01,                                                                                           // Cell[i].Width
					0x5C, 0x30, 0x5C, 0x30, 0x5C, 0x30, 0x5C, 0x30, 0x5C, 0x30, 0x5C, 0x30, 0x5C, 0x30, 0x5C, 0x30, // len(string(Cell[i].Runes))
				},
			},
			wantLog: []byte{},
		},
		{
			name:        "ErrOpen",
			inStart:     false,
			inMTErrOpen: errDummy,
			wantResp: &ServiceResponse{
				Code: http.StatusInternalServerError,
				Body: []byte("internal server error"),
			},
			wantLog: []byte("error: dummy error\n"),
		},
	}

	for _, tc := range tt {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			mt := &xpty.MockTerminal{
				ErrOpen: tc.inMTErrOpen,
				PID:     pid,
			}
			cfg := TermConfig{
				Open: mt.Open,
				Row:  2,
				Col:  2,
				Cmd: xpty.Cmd{
					Path: "bash",
					Args: []string{"--version"},
				},
			}
			slot := NewTermSlot(cfg)

			if tc.inStart {
				err := slot.start()
				if err != nil {
					t.Fatal(err)
				}

				mc := mt.Computer()
				_, err = mc.Write(tc.inIn)
				if err != nil {
					t.Fatal(err)
				}
				_, err = mc.Write([]byte{})
				if err != nil {
					t.Fatal(err)
				}
			}

			logbuf := new(bytes.Buffer)
			logger := log.New(logbuf, "", 0)

			srv := &ScreenService{
				TermSlot: slot,
				Logger:   logger,
			}

			gotResp := srv.ServeAPI(nil)
			gotLog := logbuf.Bytes()

			if slot.term != nil {
				slot.term.pc = nil
			}
			slot.Stop()

			if gotResp.Code != tc.wantResp.Code {
				t.Errorf("resp code: expected %d, got %d", tc.wantResp.Code, gotResp.Code)
			}
			if !bytes.Equal(gotResp.Body, tc.wantResp.Body) {
				t.Errorf("resp body: expected %#v, got %#v", string(tc.wantResp.Body), string(gotResp.Body))
			}
			if !bytes.Equal(gotLog, tc.wantLog) {
				t.Errorf("log: expected %#v, got %#v", string(tc.wantLog), string(gotLog))
			}
		})
	}
}

type MockService struct {
	Query url.Values
	Resp  *ServiceResponse
}

func (srv *MockService) ServeAPI(query url.Values) *ServiceResponse {
	srv.Query = query
	return srv.Resp
}
