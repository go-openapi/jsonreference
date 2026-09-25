// SPDX-FileCopyrightText: Copyright (c) 2015-2025 go-swagger maintainers
// SPDX-License-Identifier: Apache-2.0

package internal

import (
	"net/url"
	"strings"
	"testing"
	"time"

	"github.com/go-openapi/testify/v2/assert"
	"github.com/go-openapi/testify/v2/require"
)

// TestRemoveDefaultPortLinear pins the cost of removeDefaultPort on a host spelling the default
// port over and over. Removing one repetition at a time and re-parsing what remained was
// quadratic: 40 000 repetitions took seconds.
//
// The host is set directly on the url.URL: since go1.26, url.Parse rejects such a host unless
// GODEBUG urlstrictcolons=0, but callers opting out of the strict parsing still get there.
func TestRemoveDefaultPortLinear(t *testing.T) {
	const repeats = 100_000

	for _, tc := range []struct{ scheme, port string }{
		{"http", defaultHTTPPort},
		{"https", defaultHTTPSPort},
	} {
		t.Run(tc.scheme, func(t *testing.T) {
			u := &url.URL{Scheme: tc.scheme, Host: strings.Repeat(":"+tc.port, repeats), Path: "/a.json"}

			start := time.Now()
			removeDefaultPort(u)
			elapsed := time.Since(start)

			assert.Empty(t, u.Host)
			assert.Lessf(t, elapsed, time.Second, "removing %d default ports took %v", repeats, elapsed)
		})
	}
}

func TestRemoveDefaultPort(t *testing.T) {
	for _, tc := range []struct{ scheme, host, expected string }{
		{"http", "xyz.com:80", "xyz.com"},
		{"https", "xyz.com:443", "xyz.com"},
		{"http", "xyz.com:443", "xyz.com:443"},
		{"https", "xyz.com:80", "xyz.com:80"},
		{"ftp", "xyz.com:80", "xyz.com:80"},
		{"http", "xyz.com:8080", "xyz.com:8080"},
		{"http", "xyz.com", "xyz.com"},
		{"https", "[2001:db8::1]:443", "[2001:db8::1]"},
		{"http", ":80", ""},
		{"http", ":80:80:80", ""},
		// the host would become ":a", which is an invalid port: the last repetition stays
		{"http", ":a:80", ":a:80"},
		{"http", ":a:80:80:80", ":a:80"},
	} {
		t.Run(tc.scheme+"://"+tc.host, func(t *testing.T) {
			u := &url.URL{Scheme: tc.scheme, Host: tc.host}
			removeDefaultPort(u)
			assert.EqualT(t, tc.expected, u.Host)
		})
	}
}

// FuzzRemoveDefaultPort checks removeDefaultPort against the former implementation, which
// removed one repetition of the default port at a time and re-parsed the host after each.
//
// Like NormalizeURL's only caller, the fuzz target gets its host from url.Parse: removeDefaultPort
// relies on url.Parse having accepted the host it is given. A host no parse yields, such as
// "[::]::80:80", may come out differently.
func FuzzRemoveDefaultPort(f *testing.F) {
	for _, seed := range []string{
		"xyz.com:80", "xyz.com:443", ":80", ":80:80", ":443:443", ":a:80", ":a:80:80",
		":a:443:443:443", "[::1]:80", "[::1]:80:80", ":]:443", "a%20b:80", "80:80:80",
	} {
		f.Add(true, seed)
		f.Add(false, seed)
	}

	f.Fuzz(func(t *testing.T, https bool, host string) {
		if len(host) > 1024 {
			return // the reference implementation is quadratic
		}

		scheme := "http"
		if https {
			scheme = "https"
		}

		got, err := url.Parse(scheme + "://" + host + "/")
		if err != nil || got.Host != host {
			return // not a host url.Parse yields, or not a host alone (userinfo, path, query...)
		}

		want := *got
		removeDefaultPort(got)
		removeDefaultPortOneByOne(&want)

		require.EqualTf(t, want.Host, got.Host, "removing the default port from %s://%s", scheme, host)
	})
}

func removeDefaultPortOneByOne(u *url.URL) {
	for {
		port := u.Port()
		if port == "" || port != defaultPortForScheme(strings.ToLower(u.Scheme)) {
			return
		}

		host := strings.TrimSuffix(u.Host, ":"+port)
		if _, err := url.Parse("//" + host); err != nil {
			return
		}

		u.Host = host
	}
}
