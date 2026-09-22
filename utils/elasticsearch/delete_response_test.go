package elasticsearch

import (
	"errors"
	"io"
	"strings"
	"testing"

	"github.com/elastic/go-elasticsearch/v9/esapi"
)

func TestHandleDeleteResponse(t *testing.T) {
	cases := []struct {
		name    string
		err     error
		status  int
		body    string
		wantErr bool
	}{
		{name: "deleted", status: 200, body: `{"acknowledged":true}`},
		{name: "already gone is not an error", status: 404, body: `{"error":"resource_not_found_exception"}`},
		{name: "forbidden", status: 403, body: `{"error":"security_exception"}`, wantErr: true},
		{name: "server error", status: 500, body: `{"error":"internal"}`, wantErr: true},
		{name: "transport error", err: errors.New("dial tcp: connection refused"), wantErr: true},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			var res *esapi.Response
			if c.err == nil {
				res = &esapi.Response{
					StatusCode: c.status,
					Body:       io.NopCloser(strings.NewReader(c.body)),
				}
			}

			result, err := HandleDeleteResponse(c.err, res)

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
