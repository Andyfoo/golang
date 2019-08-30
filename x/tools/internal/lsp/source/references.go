// Copyright 2019 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package source

import (
	"context"
	"go/ast"
	"go/types"

	"github.com/Andyfoo/golang/x/tools/internal/telemetry/trace"
	errors "github.com/Andyfoo/golang/x/xerrors"
)

// ReferenceInfo holds information about reference to an identifier in Go source.
type ReferenceInfo struct {
	Name string
	mappedRange
	ident         *ast.Ident
	obj           types.Object
	pkg           Package
	isDeclaration bool
}

// References returns a list of references for a given identifier within the packages
// containing i.File. Declarations appear first in the result.
func (i *IdentifierInfo) References(ctx context.Context, view View) ([]*ReferenceInfo, error) {
	ctx, done := trace.StartSpan(ctx, "source.References")
	defer done()
	var references []*ReferenceInfo

	// If the object declaration is nil, assume it is an import spec and do not look for references.
	if i.Declaration.obj == nil {
		return nil, errors.Errorf("no references for an import spec")
	}

	pkgs, err := i.File.GetPackages(ctx)
	if err != nil {
		return nil, err
	}
	for _, pkg := range pkgs {
		info := pkg.GetTypesInfo()
		if info == nil {
			return nil, errors.Errorf("package %s has no types info", pkg.PkgPath())
		}

		if i.Declaration.wasImplicit {
			// The definition is implicit, so we must add it separately.
			// This occurs when the variable is declared in a type switch statement
			// or is an implicit package name. Both implicits are local to a file.
			references = append(references, &ReferenceInfo{
				Name:          i.Declaration.obj.Name(),
				mappedRange:   i.Declaration.mappedRange,
				obj:           i.Declaration.obj,
				pkg:           pkg,
				isDeclaration: true,
			})
		}
		for ident, obj := range info.Defs {
			if obj == nil || !sameObj(obj, i.Declaration.obj) {
				continue
			}
			reference := &ReferenceInfo{
				Name:          ident.Name,
				ident:         ident,
				obj:           obj,
				pkg:           pkg,
				isDeclaration: true,
			}
			if reference.mappedRange, err = posToRange(ctx, view, ident.Pos(), ident.End()); err != nil {
				return nil, err
			}
			// Add the declarations at the beginning of the references list.
			references = append([]*ReferenceInfo{reference}, references...)
		}
		for ident, obj := range info.Uses {
			if obj == nil || !sameObj(obj, i.Declaration.obj) {
				continue
			}
			reference := &ReferenceInfo{
				Name:  ident.Name,
				ident: ident,
				pkg:   pkg,
				obj:   obj,
			}
			if reference.mappedRange, err = posToRange(ctx, view, ident.Pos(), ident.End()); err != nil {
				return nil, err
			}
			references = append(references, reference)
		}

	}

	return references, nil
}

// sameObj returns true if obj is the same as declObj.
// Objects are the same if they have the some Pos and Name.
func sameObj(obj, declObj types.Object) bool {
	// TODO(suzmue): support the case where an identifier may have two different
	// declaration positions.
	return obj.Pos() == declObj.Pos() && obj.Name() == declObj.Name()
}