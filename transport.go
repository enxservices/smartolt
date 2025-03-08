package smartolt

import "net/http"

type TransportWithToken struct {
	Token     string
	Transport http.RoundTripper
}

func (t *TransportWithToken) RoundTrip(req *http.Request) (*http.Response, error) {
	req.Header.Set("X-Token", t.Token)
	if t.Transport == nil {
		t.Transport = http.DefaultTransport
	}
	return t.Transport.RoundTrip(req)
}
