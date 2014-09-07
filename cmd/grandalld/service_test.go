package main

import "testing"

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

func TestServiceServeHTTP(t *testing.T) {
}
