package kibana

import (
	"errors"
	"io"
	"net/http"
	"strings"
	"testing"
)

func TestHandleDeleteResponse(t *testing.T) {
	cases := []struct {
		name    string
		err     error
		status  int
		body    string
		wantErr bool
	}{
		{name: "deleted", status: 200, body: `{}`},
		{name: "already gone is not an error", status: 404, body: `{"statusCode":404,"error":"Not Found"}`},
		{name: "unauthorized", status: 401, body: `{"statusCode":401,"error":"Unauthorized"}`, wantErr: true},
		{name: "server error", status: 500, body: `{"statusCode":500,"error":"Internal Server Error"}`, wantErr: true},
		{name: "transport error", err: errors.New("dial tcp: connection refused"), wantErr: true},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			var res *http.Response
			if c.err == nil {
				res = &http.Response{
					StatusCode: c.status,
					Body:       io.NopCloser(strings.NewReader(c.body)),
				}
			}

			result, err := HandleDeleteResponse(res, c.err)

			if c.wantErr && err == nil {
				t.Fatalf("expected an error, got none")
			}
			if !c.wantErr && err != nil {
				t.Fatalf("expected no error, got %v", err)
			}
			if c.wantErr != result.Requeue {
				t.Errorf("expected requeue %v, got %v", c.wantErr, result.Requeue)
			}
		})
	}
}
