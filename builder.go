/*
	Copyright (c) 2021 Nikita Nikiforov <vokestd@gmail.com>
	This software is provided 'as-is', without any express or implied
	warranty. In no event will the authors be held liable for any damages
	arising from the use of this software.
	Permission is granted to anyone to use this software for any purpose,
	including commercial applications, and to alter it and redistribute it
	freely, subject to the following restrictions:
	1. The origin of this software must not be misrepresented; you must not
		 claim that you wrote the original software. If you use this software
		 in a product, an acknowledgement in the product documentation would be
		 appreciated but is not required.
	2. Altered source versions must be plainly marked as such, and must not be
		 misrepresented as being the original software.
	3. This notice may not be removed or altered from any source distribution.
*/

package main

import (
	"fmt"
	"reflect"
	"strings"
)

const URL = "https://github.com/stdcurse/stdcurse"
const ARCH = "aarch64 x86_64"
const LICENSE = "dummy"

const (
	PREPARE = iota
	PACKAGE
)

func buildScripts(s *Script, mode int) string {
	b := strings.Builder{}
	fmt.Fprintf(&b, "\ncd \"$srcdir/%s\"\n", s.Folder)
	switch mode {
	case PREPARE:
		fmt.Fprintf(&b, "\n%s\n", s.Before)
	case PACKAGE:
		switch s.Type {
		case "autotools":
			BinutilsBuilderPackage(&b, s)
		case "custom":
			CustomBuilderPackage(&b, s)
		}
		fmt.Fprintf(&b, "\n%s\n", s.After)
	}

	return b.String()
}

func (s *Scheme) Build() string {
	b := strings.Builder{}

	fmt.Fprintf(&b, "pkgname=\"%s\"\n", s.Name)
	fmt.Fprintf(&b, "pkgver=\"%s\"\n", s.Version)
	fmt.Fprintf(&b, "pkgrel=%v\n", s.Release)
	fmt.Fprintf(&b, "pkgdesc=\"%s\"\n", s.Description)
	fmt.Fprintln(&b, "name=\"$pkgname\"\nver=\"$pkgver\"\ndesc=\"$pkgdesc\"\nrel=\"$pkgrel\"")
	fmt.Fprintf(&b, "url=\"%s\"\n", URL)
	fmt.Fprintf(&b, "arch=\"%s\"\n", ARCH)
	fmt.Fprintf(&b, "license=\"%s\"\n", LICENSE)
	fmt.Fprintf(&b, "source=\"%s\"\n", strings.Join(s.Sources, " "))
	fmt.Fprintf(&b, "depends=\"%s\"\n", strings.Join(s.Dependencies, " "))
	fmt.Fprintf(&b, "makedepends=\"%s\"\n", strings.Join(s.MakeDependencies, " "))
	fmt.Fprintln(&b, "deps=\"$depends\"\nsubpkgs=\"$subpackages\"\nmakedeps=\"$makedepends\"\nsrcs=\"$source\"\n")
	spkgs := []string{}
	for _, x := range s.Subpackages {
		if reflect.TypeOf(x).Kind() == reflect.String {
			spkgs = append(spkgs, x.(string))
		} else {
			c := x.(map[string]interface{})
			spkgs = append(spkgs, strings.ReplaceAll(c["name"].(string), "-", "_")+c["name"].(string))
		}
	}
	fmt.Fprintf(&b, "subpackages=\"%s\"\n", strings.Join(spkgs, " "))
	fmt.Fprintln(&b, "prepare(){")
	for _, x := range s.Scripts {
		fmt.Fprintln(&b, buildScripts(&x, PREPARE))
	}
	fmt.Fprintf(&b, "}\npackage(){")
	for _, x := range s.Scripts {
		fmt.Fprintln(&b, buildScripts(&x, PACKAGE))
	}
	fmt.Fprintln(&b, "}")
	return b.String()
}
