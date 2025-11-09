// SPDX-FileCopyrightText: Copyright (c) 2015-2025 go-swagger maintainers
// SPDX-License-Identifier: Apache-2.0

package jsonreference

import (
	"testing"

	"github.com/go-openapi/jsonpointer"
	"github.com/go-openapi/testify/v2/assert"
	"github.com/go-openapi/testify/v2/require"
)

func TestIsRoot(t *testing.T) {
	t.Run("with empty fragment", func(t *testing.T) {
		in := "#"
		r1, err := New(in)
		require.NoError(t, err)
		assert.True(t, r1.IsRoot())
	})

	t.Run("with fragment", func(t *testing.T) {
		in := "#/ok"
		r1 := MustCreateRef(in)
		assert.False(t, r1.IsRoot())
	})

	t.Run("with invalid ref", func(t *testing.T) {
		assert.Panics(t, assert.PanicTestFunc(func() {
			MustCreateRef("%2")
		}))
	})
}

func TestFullURL(t *testing.T) {
	t.Run("with fragment", func(t *testing.T) {
		const (
			in = "http://host/path/a/b/c#/f/a/b"
		)

		r1, err := New(in)
		require.NoError(t, err)
		assert.Equal(t, in, r1.String())
		require.False(t, r1.HasFragmentOnly)
		require.True(t, r1.HasFullURL)
		require.False(t, r1.HasURLPathOnly)
		require.False(t, r1.HasFileScheme)
		require.Equal(t, "/f/a/b", r1.GetPointer().String())
	})

	t.Run("with empty fragment", func(t *testing.T) {
		const in = "http://host/path/a/b/c"

		r1, err := New(in)
		require.NoError(t, err)
		assert.Equal(t, in, r1.String())
		require.False(t, r1.HasFragmentOnly)
		require.True(t, r1.HasFullURL)
		require.False(t, r1.HasURLPathOnly)
		require.False(t, r1.HasFileScheme)
		require.Empty(t, r1.GetPointer().String())
	})
}

func TestFragmentOnly(t *testing.T) {
	const in = "#/fragment/only"

	r1, err := New(in)
	require.NoError(t, err)
	assert.Equal(t, in, r1.String())

	require.True(t, r1.HasFragmentOnly)
	require.False(t, r1.HasFullURL)
	require.False(t, r1.HasURLPathOnly)
	require.False(t, r1.HasFileScheme)
	require.Equal(t, "/fragment/only", r1.GetPointer().String())

	p, err := jsonpointer.New(r1.referenceURL.Fragment)
	require.NoError(t, err)

	t.Run("Ref with fragmentOnly", func(t *testing.T) {
		r2 := Ref{referencePointer: p, HasFragmentOnly: true}
		assert.Equal(t, in, r2.String())
	})

	t.Run("Ref without fragmentOnly", func(t *testing.T) {
		r3 := Ref{referencePointer: p, HasFragmentOnly: false}
		assert.Equal(t, in[1:], r3.String())
	})
}

func TestURLPathOnly(t *testing.T) {
	const in = "/documents/document.json"

	r1, err := New(in)
	require.NoError(t, err)
	assert.Equal(t, in, r1.String())
	require.False(t, r1.HasFragmentOnly)
	require.False(t, r1.HasFullURL)
	require.True(t, r1.HasURLPathOnly)
	require.False(t, r1.HasFileScheme)
	require.Empty(t, r1.GetPointer().String())
}

func TestURLRelativePathOnly(t *testing.T) {
	const in = "document.json"

	r1, err := New(in)
	require.NoError(t, err)
	assert.Equal(t, in, r1.String())
	require.False(t, r1.HasFragmentOnly)
	require.False(t, r1.HasFullURL)
	require.True(t, r1.HasURLPathOnly)
	require.False(t, r1.HasFileScheme)
	require.Empty(t, r1.GetPointer().String())
}

func TestInheritsInValid(t *testing.T) {
	const (
		in1 = "http://www.test.com/doc.json"
		in2 = "#/a/b"
	)

	r1, err := New(in1)
	require.NoError(t, err)

	t.Run("inherits from empty Ref", func(t *testing.T) {
		r2 := Ref{}
		result, err := r1.Inherits(r2)
		require.Error(t, err)
		assert.Nil(t, result)
	})

	t.Run("inherits from non-empty Ref", func(t *testing.T) {
		r1 = Ref{}
		r2, err := New(in2)
		require.NoError(t, err)

		result, err := r1.Inherits(r2)
		require.NoError(t, err)
		require.NotNil(t, result)
		assert.Equal(t, r2, *result)
	})
}

func TestInheritsValid(t *testing.T) {
	const (
		in1 = "http://www.test.com/doc.json"
		in2 = "#/a/b"
		out = in1 + in2
	)

	r1, err := New(in1)
	require.NoError(t, err)
	r2, err := New(in2)
	require.NoError(t, err)

	result, err := r1.Inherits(r2)
	require.NoError(t, err)
	require.NotNil(t, result)

	assert.Equal(t, out, result.String())
	assert.Equal(t, "/a/b", result.GetPointer().String())
}

func TestInheritsDifferentHost(t *testing.T) {
	const (
		in1 = "http://www.test.com/doc.json"
		in2 = "http://www.test2.com/doc.json#bla"
	)

	r1, err := New(in1)
	require.NoError(t, err)
	r2, err := New(in2)
	require.NoError(t, err)

	result, err := r1.Inherits(r2)
	require.NoError(t, err)
	require.NotNil(t, result)

	assert.Equal(t, in2, result.String())
	assert.Empty(t, result.GetPointer().String())
}

func TestFileScheme(t *testing.T) {
	const (
		in1 = "file:///Users/mac/1.json#a"
		in2 = "file:///Users/mac/2.json#b"
	)

	r1, err := New(in1)
	require.NoError(t, err)
	r2, err := New(in2)
	require.NoError(t, err)

	require.False(t, r1.HasFragmentOnly)
	require.True(t, r1.HasFileScheme)
	require.True(t, r1.HasFullFilePath)
	require.True(t, r1.IsCanonical())
	assert.Empty(t, r1.GetPointer().String())

	result, err := r1.Inherits(r2)
	require.NoError(t, err)
	assert.Equal(t, in2, result.String())
	assert.Empty(t, result.GetPointer().String())
}

func TestReferenceResolution(t *testing.T) {
	// 5.4. Reference Resolution Examples
	// http://tools.ietf.org/html/rfc3986#section-5.4
	const base = "http://a/b/c/d;p?q"

	baseRef, err := New(base)
	require.NoError(t, err)
	require.Equal(t, base, baseRef.String())

	checks := []string{
		// 5.4.1. Normal Examples
		// http://tools.ietf.org/html/rfc3986#section-5.4.1

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

		// 5.4.2. Abnormal Examples
		// http://tools.ietf.org/html/rfc3986#section-5.4.2

		"../../../g", "http://a/g",
		"../../../../g", "http://a/g",

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
		// "http:g", "http://a/b/c/g", // for backward compatibility

	}
	for i := 0; i < len(checks); i += 2 {
		child := checks[i]
		expected := checks[i+1]

		childRef, e := New(child)
		require.NoErrorf(t, e, "test: %d: New(%s) failed error: %v", i/2, child, e)

		res, e := baseRef.Inherits(childRef)
		require.NoErrorf(t, e, "test: %d", i/2)
		require.NotNilf(t, res, "test: %d", i/2)
		assert.Equalf(t, expected, res.String(), "test: %d", i/2)
	}
}

func TestIdenticalURLEncoded(t *testing.T) {
	expected, err := New("https://localhost/🌭#/🍔")
	require.NoErrorf(t, err, "failed to create jsonreference: %v", err)

	actual, err := New("https://localhost/%F0%9F%8C%AD#/%F0%9F%8D%94")
	require.NoErrorf(t, err, "failed to create jsonreference: %v", err)
	require.Equal(t, expected, actual)
}
