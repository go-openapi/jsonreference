// SPDX-FileCopyrightText: Copyright (c) 2015-2025 go-swagger maintainers
// SPDX-License-Identifier: Apache-2.0

package internal

import (
	"net/url"
	"testing"

	"github.com/go-openapi/testify/v2/require"
)

// FuzzNormalizeURL pins the postcondition callers rely on: a URL that parsed before still
// parses after being normalized, and normalizing what came out changes nothing further.
//
// Both matter to go-openapi/spec, which turns a normalized URL back into a string, hands it
// to MustCreateRef and gets a panic when it no longer parses.
func FuzzNormalizeURL(f *testing.F) {
	for _, seed := range normalizeURLSeeds() {
		f.Add(seed)
	}

	f.Fuzz(func(t *testing.T, input string) {
		u, err := url.Parse(input)
		if err != nil {
			return
		}

		NormalizeURL(u)
		once := u.String()

		again, err := url.Parse(once)
		require.NoErrorf(t, err, "normalizing %q yielded %q, which no longer parses", input, once)

		NormalizeURL(again)
		require.EqualTf(t, once, again.String(), "normalizing %q is not idempotent", input)
	})
}

// normalizeURLSeeds covers what the normalizations branch on: the schemes carrying a default
// port, ports that are not default, hosts needing lower-casing, duplicate slashes, escaped
// paths and fragments, IPv6 literals, userinfo, and the degenerate authorities url.Parse accepts.
func normalizeURLSeeds() []string {
	return []string{
		"",
		"/folder/file",
		"HTTPs://xYz.cOm:443/folder//file",
		"HTTP://xYz.cOm:80/folder//file",
		"http://xyz.com:8080/folder",
		"postGRES://xYz.cOm:5432/folder//file",
		userinfoDefaultPort + "?q=1#/a~1b",
		"https://[2001:DB8::1]:443/folder",
		ipv6NonDefaultPort,
		degenerateHostHTTPS,
		degenerateHostHTTP,
		"https://:443",
		"file:///base/path.json#/definitions/a%20b",
		"file://",
		"https://localhost/%F0%9F%8C%AD#/%F0%9F%8D%94",
		"mailto:someone@example.com",
		"//host/path",
		"http:g",
	}
}
