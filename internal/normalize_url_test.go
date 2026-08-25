// SPDX-FileCopyrightText: Copyright (c) 2015-2025 go-swagger maintainers
// SPDX-License-Identifier: Apache-2.0

package internal

import (
	"net/url"
	"testing"

	"github.com/go-openapi/testify/v2/assert"
	"github.com/go-openapi/testify/v2/require"
)

// URLs spelled in more than one test in this package.
const (
	// url.Parse reads these as the host ":a" on a default port, so dropping the port would
	// leave "https://:a", which no longer parses.
	degenerateHostHTTPS = "https://:a:443"
	degenerateHostHTTP  = "http://:a:80"

	ipv6NonDefaultPort = "https://[2001:db8::1]:8443/folder"

	// userinfo holds a colon of its own, which the port removal must not read as a port.
	//nolint:gosec // test URLs carrying userinfo, not credentials
	userinfoDefaultPort = "https://user:pw@xYz.cOm:443/folder"
	//nolint:gosec // test URLs carrying userinfo, not credentials
	userinfoNormalized = "https://user:pw@xyz.com/folder"
)

func TestUrlnorm(t *testing.T) {
	testCases := []struct {
		url      string
		expected string
	}{
		{
			url:      "HTTPs://xYz.cOm:443/folder//file",
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
			url:      degenerateHostHTTPS,
			expected: degenerateHostHTTPS,
		},
		{
			url:      degenerateHostHTTP,
			expected: degenerateHostHTTP,
		},
		{
			url:      "https://:]:443",
			expected: "https://:]:443",
		},
		{
			// the host is ":80" on port 80, so the removal has to run twice. Emptying the host
			// drops the "//" as well, which url.URL.String has always done.
			url:      "http://:80:80",
			expected: "http:",
		},
	}

	for _, toPin := range testCases {
		testCase := toPin

		u, err := url.Parse(testCase.url)
		require.NoError(t, err)

		NormalizeURL(u)
		normalized := u.String()
		assert.EqualT(t, testCase.expected, normalized)

		_, err = url.Parse(normalized)
		require.NoErrorf(t, err, "normalizing %q yielded a URL that no longer parses", testCase.url)
	}
}
