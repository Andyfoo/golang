// Copyright 2018 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package fmt_test

import (
	"path/filepath"
	"regexp"

	"github.com/Andyfoo/golang/x/exp/errors"
	"github.com/Andyfoo/golang/x/exp/errors/fmt"
)

func baz() error { return errors.New("baz flopped") }
func bar() error { return fmt.Errorf("bar(nameserver 139): %v", baz()) }
func foo() error { return fmt.Errorf("foo: %s", bar()) }

func Example_formatting() {
	err := foo()
	fmt.Println("Error:")
	fmt.Printf("%v\n", err)
	fmt.Println()
	fmt.Println("Detailed error:")
	fmt.Println(stripPath(fmt.Sprintf("%+v\n", err)))
	// Output:
	// Error:
	// foo: bar(nameserver 139): baz flopped
	//
	// Detailed error:
	// foo:
	//     github.com/Andyfoo/golang/x/exp/errors/fmt_test.foo
	//         github.com/Andyfoo/golang/x/exp/errors/fmt/format_example_test.go:17
	//   - bar(nameserver 139):
	//     github.com/Andyfoo/golang/x/exp/errors/fmt_test.bar
	//         github.com/Andyfoo/golang/x/exp/errors/fmt/format_example_test.go:16
	//   - baz flopped:
	//     github.com/Andyfoo/golang/x/exp/errors/fmt_test.baz
	//         github.com/Andyfoo/golang/x/exp/errors/fmt/format_example_test.go:15
}

func stripPath(s string) string {
	rePath := regexp.MustCompile(`( [^ ]*)golang.org`)
	s = rePath.ReplaceAllString(s, " golang.org")
	s = filepath.ToSlash(s)
	return s
}
