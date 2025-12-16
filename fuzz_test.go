// SPDX-FileCopyrightText: Copyright (c) 2015-2025 go-swagger maintainers
// SPDX-License-Identifier: Apache-2.0

package jsonreference

import (
	"iter"
	"slices"
	"strings"
	"testing"

	"github.com/go-openapi/testify/v2/require"
)

func FuzzParse(f *testing.F) {
	// initial seed
	cumulated := make([]string, 0, 100)
	for generator := range generators() {
		f.Add(generator)

		cumulated = append(cumulated, generator)
		f.Add(strings.Join(cumulated, ""))
	}

	ref := Ref{}
	f.Fuzz(func(t *testing.T, input string) {
		require.NotPanics(t, func() {
			_ = ref.parse(input)
		})
	})
}

func generators() iter.Seq[string] {
	return slices.Values([]string{
		"",
		"https://localhost/%F0%9F%8C%AD#/%F0%9F%8D%94",
		"#/ok",
		"%2",
		"http://host/path/a/b/c#/f/a/b",
		"http://host/path/a/b/c",
		"http://host/path/a/b/c#",
		"#/fragment/only",
		"/documents/document.json",
		"document.json",
		"http://www.test.com/doc.json",
		"#/a/b",
		"http://www.test.com/doc.json",
		"http://www.test2.com/doc.json#bla",
		"http://a/b/c/d;p?q",
		"g:h", "g:h",
		"g", "http://a/b/c/g",
		"./g", "http://a/b/c/g",
		"g/", "http://a/b/c/g/",
		"/g", "http://a/g",
		"//g", "http://g",
		"?y", "http://a/b/c/d;p?y",
		"g?y", "http://a/b/c/g?y",
		"#s", "http://a/b/c/d;p?q#s",
		"g#s", "http://a/b/c/g#s",
		"g?y#s", "http://a/b/c/g?y#s",
		";x", "http://a/b/c/;x",
		"g;x", "http://a/b/c/g;x",
		"g;x?y#s", "http://a/b/c/g;x?y#s",
		"", "http://a/b/c/d;p?q",
		".", "http://a/b/c/",
		"./", "http://a/b/c/",
		"..", "http://a/b/",
		"../", "http://a/b/",
		"../g", "http://a/b/g",
		"../..", "http://a/",
		"../../", "http://a/",
		"../../g", "http://a/g",
		"http://tools.ietf.org/html/rfc3986#section-5.4.2",
		"../../../g", "http://a/g",
		"../../../../g", "http://a/g",
		"https://localhost/🌭#/🍔",
		"https://localhost/%F0%9F%8C%AD#/%F0%9F%8D%94",
		"/./g", "http://a/g",
		"/../g", "http://a/g",
		"g.", "http://a/b/c/g.",
		".g", "http://a/b/c/.g",
		"g..", "http://a/b/c/g..",
		"..g", "http://a/b/c/..g",
		"./../g", "http://a/b/g",
		"./g/.", "http://a/b/c/g/",
		"g/./h", "http://a/b/c/g/h",
		"g/../h", "http://a/b/c/h",
		"g;x=1/./y", "http://a/b/c/g;x=1/y",
		"g;x=1/../y", "http://a/b/c/y",
		"g?y/./x", "http://a/b/c/g?y/./x",
		"g?y/../x", "http://a/b/c/g?y/../x",
		"g#s/./x", "http://a/b/c/g#s/./x",
		"g#s/../x", "http://a/b/c/g#s/../x",
		"http:g", "http:g", // for strict parsers
		"http:g", "http://a/b/c/g",
		`a`,
		``, `/`, `/`, `/a~1b`, `/a~1b`, `/c%d`, `/e^f`, `/g|h`, `/i\j`, `/k"l`, `/ `, `/m~0n`,
		`/foo`, `/0`,
	})
}
