package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"net/url"
	"reflect"
	"testing"
)

func TestNewService(t *testing.T) {
	// test valid site configurations
	for i, test := range []struct {
		sites []*Site
	}{
		{nil},
		{[]*Site{}},
		{[]*Site{
			{Name: "test", Bind: "/test", URL: "http://example.com/test"},
		}},
		{[]*Site{
			{Name: "test1", Bind: "/test1", URL: "http://example.com/test1"},
			{Name: "test2", Bind: "/test2", URL: "http://example.com/test2"},
		}},
		{[]*Site{
			// duplicate destinations are ok.
			{Name: "test1", Bind: "/test1", URL: "http://example.com/test"},
			{Name: "test2", Bind: "/test2", URL: "http://example.com/test"},
		}},
	} {
		s, err := NewService(test.sites)
		if err != nil {
			t.Errorf("test %d: %v", i, err)
			continue
		}
		if s == nil {
			t.Errorf("test %d: nil service", i)
			continue
		}
	}

	// test invalid site configurations
	for i, test := range []struct {
		sites []*Site
	}{
		{[]*Site{new(Site)}},
		{[]*Site{
			// relative bind url
			{Name: "test", Bind: "test", URL: "http://example.com/test"},
		}},
		{[]*Site{
			// scheme-relative destination url
			{Name: "test", Bind: "/test", URL: "//example.com/test"},
		}},
		{[]*Site{
			// destination url path
			{Name: "test", Bind: "/test", URL: "/test"},
		}},
		{[]*Site{
			// duplicate site names
			{Name: "test1", Bind: "/test1", URL: "http://example.com/test1"},
			{Name: "test1", Bind: "/test2", URL: "http://example.com/test2"},
		}},
		{[]*Site{
			// duplicate bind urls
			{Name: "test1", Bind: "/test1", URL: "http://example.com/test1"},
			{Name: "test2", Bind: "/test1", URL: "http://example.com/test2"},
		}},
	} {
		s, err := NewService(test.sites)
		if err == nil {
			t.Errorf("test %d: expected error", i)
			continue
		}
		if s != nil {
			t.Errorf("test %d: non-nil service", i)
			continue
		}
	}
}

func TestServiceServeHTTP_redirect(t *testing.T) {
	_, server := newSimpleServiceTest(t)
	c := newServiceClient()
	defer server.Close()
	defer server.CloseClientConnections()

	resp, err := c.Get(server.URL + "/test")
	if !isredirect(err) {
		if err == nil {
			t.Errorf("did not redirect (%v)", resp.Status)
			return
		}
		t.Error(err)
		return
	}
	if resp.StatusCode != http.StatusTemporaryRedirect {
		t.Errorf("status: %v", http.StatusText(resp.StatusCode))
	}
	loc := resp.Header.Get("Location")
	if loc != "http://example.com/test" {
		t.Errorf("redirect location %q", loc)
	}
}

func TestServiceServeHTTP_aliasNotFound(t *testing.T) {
	_, server := newSimpleServiceTest(t)
	c := newServiceClient()
	defer server.Close()
	defer server.CloseClientConnections()

	resp, err := c.Get(server.URL + "/unknown-alias")
	if isredirect(err) {
		t.Errorf("%s %q", http.StatusText(resp.StatusCode), resp.Header.Get("Location"))
		return
	}
	if err != nil {
		t.Error(err)
		return
	}
	if resp.StatusCode != http.StatusNotFound {
		t.Error(http.StatusText(resp.StatusCode))
	}
}

func TestServiceServeHTTP_apiNotFound(t *testing.T) {
	_, server := newSimpleServiceTest(t)
	c := newServiceClient()
	defer server.Close()
	defer server.CloseClientConnections()

	resp, err := c.Get(server.URL + "/.api/v999")
	if isredirect(err) {
		t.Errorf("%s %q", http.StatusText(resp.StatusCode), resp.Header.Get("Location"))
		return
	}
	if err != nil {
		t.Error(err)
		return
	}
	if resp.StatusCode != http.StatusNotFound {
		t.Error(http.StatusText(resp.StatusCode))
	}
}

func TestServiceServeHTTP_apiV1NotFound(t *testing.T) {
	_, server := newSimpleServiceTest(t)
	c := newServiceClient()
	defer server.Close()
	defer server.CloseClientConnections()

	resp, err := c.Get(server.URL + "/.api/v1/unknown-resource")
	if isredirect(err) {
		t.Errorf("%s %q", http.StatusText(resp.StatusCode), resp.Header.Get("Location"))
		return
	}
	if err != nil {
		t.Error(err)
		return
	}
	if resp.StatusCode != http.StatusNotFound {
		t.Error(http.StatusText(resp.StatusCode))
	}
}

func TestServiceServeHTTP_apiV1Aliases(t *testing.T) {
	_, server := newSimpleServiceTest(t)
	c := newServiceClient()
	defer server.Close()
	defer server.CloseClientConnections()

	resp, err := c.Get(server.URL + "/.api/v1/aliases")
	if isredirect(err) {
		t.Errorf("%s %q", http.StatusText(resp.StatusCode), resp.Header.Get("Location"))
		return
	}
	if err != nil {
		t.Error(err)
		return
	}
	if resp.StatusCode != http.StatusOK {
		t.Error(http.StatusText(resp.StatusCode))
	}
	p, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Errorf("read: %v", err)
	}
	var body []map[string]interface{}
	err = json.Unmarshal(p, &body)
	if err != nil {
		t.Errorf("unmarshal: %v", err)
		return
	}
	if len(body) != 1 {
		var names []interface{}
		for i := range body {
			names = append(names, body[i]["name"])
		}
		t.Errorf("length %d %q", len(body), names)
		return
	}
	if !reflect.DeepEqual(body[0]["name"], "test-site") {
		t.Errorf("name %q", body[0]["name"])
	}
}

func newServiceClient() *http.Client {
	c := new(http.Client)
	c.CheckRedirect = nofollow
	return c
}

func newSimpleServiceTest(t *testing.T) (*Service, *httptest.Server) {
	site := &Site{
		Name:        "test-site",
		Bind:        "/test",
		URL:         "http://example.com/test",
		Description: "A site to test basic service behavior",
	}
	s, err := NewService([]*Site{site})
	if err != nil {
		t.Fatalf("new: %v", err)
	}
	handler, err := s.Handler()
	if err != nil {
		t.Fatalf("handler: %v", err)
	}
	server := httptest.NewServer(handler)
	return s, server
}

func isredirect(err error) bool {
	uerr, ok := err.(*url.Error)
	return ok && uerr.Err == errNofollow
}

// nofollow can be used as an http.Client's CheckRedirect action.
// nofollow never lets a client follow redirects.
func nofollow(r *http.Request, via []*http.Request) error {
	return errNofollow
}

var errNofollow = fmt.Errorf("nofollow")
