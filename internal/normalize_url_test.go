// SPDX-FileCopyrightText: Copyright (c) 2015-2025 go-swagger maintainers
// SPDX-License-Identifier: Apache-2.0

package internal

import (
	"fmt"
	"iter"
	"net/url"
	"slices"
	"testing"

	"github.com/go-openapi/testify/v2/assert"
	"github.com/go-openapi/testify/v2/require"
)

func TestUrlnorm(t *testing.T) {
	for testCase := range urlTestCases() {
		t.Run(fmt.Sprintf("should normalize URL: %q", testCase.url), func(t *testing.T) {
			u, err := url.Parse(testCase.url)
			if testCase.expectErr {
				require.Error(t, err, "expected url.Parse to fail")
				return
			}

			require.NoError(t, err)

			NormalizeURL(u)
			normalized := u.String()
			assert.EqualTf(t, testCase.expected, normalized, "got an unexpected normalization: %s", normalized)

			_, err = url.Parse(normalized)
			require.NoErrorf(t, err, "normalizing %q yielded an URL that no longer parses", testCase.url)
		})
	}
}

type urlTestCase struct {
	url       string
	expected  string // expected normalized url
	expectErr bool   // expected error when true
}

const (
	// url.Parse reads these as the host ":a" on a default port, so dropping the port would
	// leave "https://:a", which no longer parses.
	degenerateHostHTTPS = "https://:a:443" // since go1.26 this degenerate case is no longer tolerated
	degenerateHostHTTP  = "http://:a:80"

	ipv6NonDefaultPort = "https://[2001:db8::1]:8443/folder"

	// a default port, an upper-cased scheme and host, and a duplicate slash, all at once.
	mixedCaseDefaultPort = "HTTPs://xYz.cOm:443/folder//file"

	// userinfo holds a colon of its own, which the port removal must not read as a port.
	//nolint:gosec // test URLs carrying userinfo, not credentials
	userinfoDefaultPort = "https://user:pw@xYz.cOm:443/folder"
	//nolint:gosec // test URLs carrying userinfo, not credentials
	userinfoNormalized = "https://user:pw@xyz.com/folder"
)

func urlTestCases() iter.Seq[urlTestCase] {
	return slices.Values([]urlTestCase{
		{
			url:      mixedCaseDefaultPort,
			expected: "https://xyz.com/folder/file",
		},
		{
			url:      "HTTP://xYz.cOm:80/folder//file",
			expected: "http://xyz.com/folder/file",
		},
		{
			url:      "postGRES://xYz.cOm:5432/folder//file",
			expected: "postgres://xyz.com:5432/folder/file",
		},
		{
			url:      userinfoDefaultPort,
			expected: userinfoNormalized,
		},
		{
			url:      "https://[2001:DB8::1]:443/folder",
			expected: "https://[2001:db8::1]/folder",
		},
		{
			url:      ipv6NonDefaultPort,
			expected: ipv6NonDefaultPort,
		},
		{
			// Since go1.26, url.Parse rejects a colon outside a bracketed IPv6 host on an http or https URL
			// (GODEBUG urlstrictcolons=1, the default from a go.mod declaring go 1.26 or later). Both $refs below
			// used to parse - the first as the host ":a" on port 443, the second as "0:443" on port 443 - and both
			// now fail. normalizeURI logs a warning, repairs the $ref to the empty URI and resolves it against the
			// base, so the base itself comes back.
			url:       degenerateHostHTTPS,
			expected:  degenerateHostHTTPS,
			expectErr: true,
		},
		{
			url:       degenerateHostHTTP,
			expected:  degenerateHostHTTP,
			expectErr: true,
		},
		{
			// a run of slashes collapses to one, however long the run is
			url:      "https://xyz.com/a//b///c////d",
			expected: "https://xyz.com/a/b/c/d",
		},
		{
			url:      "https://xyz.com////",
			expected: "https://xyz.com/",
		},
		{
			url:       "https://:]:443",
			expected:  "https://:]:443",
			expectErr: true,
		},
		{
			// the host is ":80" on port 80, so the removal has to run twice. Emptying the host
			// drops the "//" as well, which url.URL.String has always done.
			url:       "http://:80:80",
			expected:  "http:",
			expectErr: true,
		},
	})
}
