// SPDX-FileCopyrightText: Copyright (c) 2015-2025 go-swagger maintainers
// SPDX-License-Identifier: Apache-2.0

package jsonreference_test

import (
	"fmt"
	"log"

	"github.com/go-openapi/jsonreference"
)

func ExampleRef_GetURL() {
	fragRef := jsonreference.MustCreateRef("#/definitions/Pet")

	fmt.Printf("URL: %s\n", fragRef.GetURL())

	// Output: URL: #/definitions/Pet
}

func ExampleRef_Inherits() {
	parent := jsonreference.MustCreateRef("http://example.com/base.json")
	child, err := jsonreference.New("#/definitions/Pet")
	if err != nil {
		log.Printf("%v", err)

		return
	}

	resolved, err := parent.Inherits(child)
	if err != nil {
		log.Printf("%v", err)

		return
	}

	fmt.Printf("URL: %v\n", resolved)

	// Output: URL: http://example.com/base.json#/definitions/Pet
}
