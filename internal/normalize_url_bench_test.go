// SPDX-FileCopyrightText: Copyright (c) 2015-2025 go-swagger maintainers
// SPDX-License-Identifier: Apache-2.0

package internal

import (
	"net/url"
	"testing"
)

func BenchmarkNormalizeURL(b *testing.B) {
	// the shapes NormalizeURL branches on: nothing to do, a default port with an upper-cased
	// host, a fragment, a query, and a path holding runs of slashes
	benchURLs := []string{
		"https://example.com/v1/pets/{petId}/photos",
		mixedCaseDefaultPort,
		"file:///base/path.json#/definitions/a",
		"http://a/b/c/d;p?q",
		"https://example.com/a//b///c////d?x=1#/frag",
	}

	parsed := make([]url.URL, 0, len(benchURLs))
	for _, raw := range benchURLs {
		u, err := url.Parse(raw)
		if err != nil {
			b.Fatal(err)
		}
		parsed = append(parsed, *u)
	}

	b.ReportAllocs()
	for b.Loop() {
		for i := range parsed {
			u := parsed[i]
			NormalizeURL(&u)
		}
	}
}
