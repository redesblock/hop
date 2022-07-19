package jsonhttp_test

import (
	"encoding/json"
	"errors"
	"net"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
)

func TestRespond_defaults(t *testing.T) {
	w := httptest.NewRecorder()

	Respond(w, 0, nil)

	statusCode := w.Result().StatusCode
	wantCode := http.StatusOK
	if statusCode != wantCode {
		t.Errorf("got status code %d, want %d", statusCode, wantCode)
	}

	var m *StatusResponse

	if err := json.Unmarshal(w.Body.Bytes(), &m); err != nil {
		t.Errorf("json unmarshal response body: %s", err)
	}

	if m.Code != wantCode {
		t.Errorf("got message code %d, want %d", m.Code, wantCode)
	}

	wantMessage := http.StatusText(wantCode)
	if m.Message != wantMessage {
		t.Errorf("got message message %q, want %q", m.Message, wantMessage)
	}

	testContentType(t, w)
}

func TestRespond_statusResponse(t *testing.T) {
	for _, tc := range []struct {
		code int
	}{
		{code: http.StatusContinue},
		{code: http.StatusSwitchingProtocols},
		{code: http.StatusOK},
		{code: http.StatusCreated},
		{code: http.StatusAccepted},
		{code: http.StatusNonAuthoritativeInfo},
		{code: http.StatusResetContent},
		{code: http.StatusPartialContent},
		{code: http.StatusMultipleChoices},
		{code: http.StatusMovedPermanently},
		{code: http.StatusFound},
		{code: http.StatusSeeOther},
		{code: http.StatusNotModified},
		{code: http.StatusUseProxy},
		{code: http.StatusTemporaryRedirect},
		{code: http.StatusPermanentRedirect},
		{code: http.StatusBadRequest},
		{code: http.StatusUnauthorized},
		{code: http.StatusPaymentRequired},
		{code: http.StatusForbidden},
		{code: http.StatusNotFound},
		{code: http.StatusMethodNotAllowed},
		{code: http.StatusNotAcceptable},
		{code: http.StatusProxyAuthRequired},
		{code: http.StatusRequestTimeout},
		{code: http.StatusConflict},
		{code: http.StatusGone},
		{code: http.StatusLengthRequired},
		{code: http.StatusPreconditionFailed},
		{code: http.StatusRequestEntityTooLarge},
		{code: http.StatusRequestURITooLong},
		{code: http.StatusUnsupportedMediaType},
		{code: http.StatusRequestedRangeNotSatisfiable},
		{code: http.StatusExpectationFailed},
		{code: http.StatusTeapot},
		{code: http.StatusUpgradeRequired},
		{code: http.StatusPreconditionRequired},
		{code: http.StatusTooManyRequests},
		{code: http.StatusRequestHeaderFieldsTooLarge},
		{code: http.StatusUnavailableForLegalReasons},
		{code: http.StatusInternalServerError},
		{code: http.StatusNotImplemented},
		{code: http.StatusBadGateway},
		{code: http.StatusServiceUnavailable},
		{code: http.StatusGatewayTimeout},
		{code: http.StatusHTTPVersionNotSupported},
	} {
		w := httptest.NewRecorder()

		Respond(w, tc.code, nil)

		statusCode := w.Result().StatusCode
		if statusCode != tc.code {
			t.Errorf("got status code %d, want %d", statusCode, tc.code)
		}

		var m *StatusResponse

		if err := json.Unmarshal(w.Body.Bytes(), &m); err != nil {
			t.Errorf("json unmarshal response body: %s", err)
		}

		if m.Code != tc.code {
			t.Errorf("got message code %d, want %d", m.Code, tc.code)
		}

		wantMessage := http.StatusText(tc.code)
		if m.Message != wantMessage {
			t.Errorf("got message message %q, want %q", m.Message, wantMessage)
		}

		testContentType(t, w)
	}
}

func TestRespond_special(t *testing.T) {
	for _, tc := range []struct {
		name        string
		code        int
		response    interface{}
		wantMessage string
	}{
		{
			name:        "string 200",
			code:        http.StatusOK,
			response:    "custom message",
			wantMessage: "custom message",
		},
		{
			name:        "string 404",
			code:        http.StatusNotFound,
			response:    "element not found",
			wantMessage: "element not found",
		},
		{
			name:        "error 400",
			code:        http.StatusBadRequest,
			response:    errors.New("test error"),
			wantMessage: "test error",
		},
		{
			name:        "error 500",
			code:        http.StatusInternalServerError,
			response:    errors.New("test error"),
			wantMessage: "test error",
		},
		{
			name:        "stringer 200",
			code:        http.StatusOK,
			response:    net.IPv4(127, 0, 0, 1), // net.IP implements Stringer interface
			wantMessage: "127.0.0.1",
		},
		{
			name:        "stringer 403",
			code:        http.StatusForbidden,
			response:    net.IPv4(2, 4, 8, 16), // net.IP implements Stringer interface
			wantMessage: "2.4.8.16",
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			w := httptest.NewRecorder()

			Respond(w, tc.code, tc.response)

			statusCode := w.Result().StatusCode
			if statusCode != tc.code {
				t.Errorf("got status code %d, want %d", statusCode, tc.code)
			}

			var m *StatusResponse

			if err := json.Unmarshal(w.Body.Bytes(), &m); err != nil {
				t.Errorf("json unmarshal response body: %s", err)
			}

			if m.Code != tc.code {
				t.Errorf("got message code %d, want %d", m.Code, tc.code)
			}

			if m.Message != tc.wantMessage {
				t.Errorf("got message message %q, want %q", m.Message, tc.wantMessage)
			}

			testContentType(t, w)
		})
	}
}

func TestRespond_custom(t *testing.T) {
	w := httptest.NewRecorder()

	wantCode := http.StatusTeapot

	type response struct {
		Field1 string `json:"field1"`
		Field2 int    `json:"field2"`
	}

	r := response{
		Field1: "custom message",
		Field2: 42,
	}
	Respond(w, wantCode, r)

	statusCode := w.Result().StatusCode
	if statusCode != wantCode {
		t.Errorf("got status code %d, want %d", statusCode, wantCode)
	}

	var m response

	if err := json.Unmarshal(w.Body.Bytes(), &m); err != nil {
		t.Errorf("json unmarshal response body: %s", err)
	}

	if !reflect.DeepEqual(m, r) {
		t.Errorf("got response %+v, want %+v", m, r)
	}

	testContentType(t, w)
}

func TestStandardHTTPResponds(t *testing.T) {
	for _, tc := range []struct {
		f    func(w http.ResponseWriter, response interface{})
		code int
	}{
		{f: Continue, code: http.StatusContinue},
		{f: SwitchingProtocols, code: http.StatusSwitchingProtocols},
		{f: OK, code: http.StatusOK},
		{f: Created, code: http.StatusCreated},
		{f: Accepted, code: http.StatusAccepted},
		{f: NonAuthoritativeInfo, code: http.StatusNonAuthoritativeInfo},
		{f: ResetContent, code: http.StatusResetContent},
		{f: PartialContent, code: http.StatusPartialContent},
		{f: MultipleChoices, code: http.StatusMultipleChoices},
		{f: MovedPermanently, code: http.StatusMovedPermanently},
		{f: Found, code: http.StatusFound},
		{f: SeeOther, code: http.StatusSeeOther},
		{f: NotModified, code: http.StatusNotModified},
		{f: UseProxy, code: http.StatusUseProxy},
		{f: TemporaryRedirect, code: http.StatusTemporaryRedirect},
		{f: PermanentRedirect, code: http.StatusPermanentRedirect},
		{f: BadRequest, code: http.StatusBadRequest},
		{f: Unauthorized, code: http.StatusUnauthorized},
		{f: PaymentRequired, code: http.StatusPaymentRequired},
		{f: Forbidden, code: http.StatusForbidden},
		{f: NotFound, code: http.StatusNotFound},
		{f: MethodNotAllowed, code: http.StatusMethodNotAllowed},
		{f: NotAcceptable, code: http.StatusNotAcceptable},
		{f: ProxyAuthRequired, code: http.StatusProxyAuthRequired},
		{f: RequestTimeout, code: http.StatusRequestTimeout},
		{f: Conflict, code: http.StatusConflict},
		{f: Gone, code: http.StatusGone},
		{f: LengthRequired, code: http.StatusLengthRequired},
		{f: PreconditionFailed, code: http.StatusPreconditionFailed},
		{f: RequestEntityTooLarge, code: http.StatusRequestEntityTooLarge},
		{f: RequestURITooLong, code: http.StatusRequestURITooLong},
		{f: UnsupportedMediaType, code: http.StatusUnsupportedMediaType},
		{f: RequestedRangeNotSatisfiable, code: http.StatusRequestedRangeNotSatisfiable},
		{f: ExpectationFailed, code: http.StatusExpectationFailed},
		{f: Teapot, code: http.StatusTeapot},
		{f: UpgradeRequired, code: http.StatusUpgradeRequired},
		{f: PreconditionRequired, code: http.StatusPreconditionRequired},
		{f: TooManyRequests, code: http.StatusTooManyRequests},
		{f: RequestHeaderFieldsTooLarge, code: http.StatusRequestHeaderFieldsTooLarge},
		{f: UnavailableForLegalReasons, code: http.StatusUnavailableForLegalReasons},
		{f: InternalServerError, code: http.StatusInternalServerError},
		{f: NotImplemented, code: http.StatusNotImplemented},
		{f: BadGateway, code: http.StatusBadGateway},
		{f: ServiceUnavailable, code: http.StatusServiceUnavailable},
		{f: GatewayTimeout, code: http.StatusGatewayTimeout},
		{f: HTTPVersionNotSupported, code: http.StatusHTTPVersionNotSupported},
	} {
		w := httptest.NewRecorder()
		tc.f(w, nil)
		var m *StatusResponse

		if err := json.Unmarshal(w.Body.Bytes(), &m); err != nil {
			t.Errorf("json unmarshal response body: %s", err)
		}

		if m.Code != tc.code {
			t.Errorf("expected message code %d, got %d", tc.code, m.Code)
		}

		if m.Message != http.StatusText(tc.code) {
			t.Errorf("expected message message \"%s\", got \"%s\"", http.StatusText(tc.code), m.Message)
		}

		testContentType(t, w)
	}
}

func TestPanicRespond(t *testing.T) {
	w := httptest.NewRecorder()

	defer func() {
		err := recover()
		if _, ok := err.(*json.UnsupportedTypeError); !ok {
			t.Errorf("expected error from recover json.UnsupportedTypeError, got %#v", err)
		}
	}()

	Respond(w, http.StatusNotFound, map[bool]string{
		true: "",
	})
}

func testContentType(t *testing.T, r *httptest.ResponseRecorder) {
	t.Helper()

	if got := r.Header().Get("Content-Type"); got != DefaultContentTypeHeader {
		t.Errorf("got content type %q, want %q", got, DefaultContentTypeHeader)
	}
}
