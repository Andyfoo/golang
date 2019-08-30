// Copyright 2018 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// The gopls command is an LSP server for Go.
// The Language Server Protocol allows any text editor
// to be extended with IDE-like features;
// see https://langserver.org/ for details.
package main // import "github.com/Andyfoo/golang/x/tools/cmd/gopls"

import (
	"context"
	"os"

	"github.com/Andyfoo/golang/x/tools/internal/lsp/cmd"
	"github.com/Andyfoo/golang/x/tools/internal/lsp/debug"
	"github.com/Andyfoo/golang/x/tools/internal/tool"
)

func main() {
	debug.Version += "-cmd.gopls"
	tool.Main(context.Background(), cmd.New("gopls-legacy", "", nil), os.Args[1:])
}
