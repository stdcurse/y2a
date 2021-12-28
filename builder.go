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

func interfaceMapToStringMap(m map[interface{}]interface{}) map[string]interface{} {
	ret := map[string]interface{}{}

	for k, v := range m {
		ret[fmt.Sprintf("%v", k)] = v
	}

	return ret
}

func prepareEntry(s *Entry) string {
	b := strings.Builder{}

	for _, v := range s.Environment {
		fmt.Fprintf(&b, "export %s\n", v)
	}

	fmt.Fprintf(&b, "cd \"$srcdir/%s\"\n", s.Folder)
	fmt.Fprintf(&b, "%s\n", s.Before)

	return b.String()
}

func packageEntry(s *Entry) string {
	b := strings.Builder{}

	for _, v := range s.Environment {
		fmt.Fprintf(&b, "export %s\n", v)
	}
	fmt.Fprintf(&b, "cd \"$srcdir/%s\"\n", s.Folder)

	switch s.Type {
	case "autotools":
		AutotoolsBuilderPackage(&b, s)
	case "custom":
		CustomBuilderPackage(&b, s)
	}
	fmt.Fprintf(&b, "%s\n", s.After)
	fmt.Fprintln(&b, "rm -rf \"$pkgdir/usr/share/man\"\n")

	return b.String()
}

func (s *Scheme) schemeToMap() map[string]string {
	ret := map[string]string{}
	v := reflect.ValueOf(*s)
	t := v.Type()

	for i := 0; i < v.NumField(); i++ {
		k := t.Field(i).Name
		x := v.Field(i).Interface()

		if reflect.TypeOf(x).Kind() == reflect.String {
			if j := t.Field(i).Tag.Get("yaml"); j != "" {
				k = j
			} else {
				k = strings.ToLower(k)
			}

			ret[k] = x.(string)
		}
	}

	for k, v := range s.Defines {
		ret[k] = v
	}

	return ret
}

func doSpkg(b *strings.Builder, s *map[string]interface{}) {
	c := *s

	name := strings.ReplaceAll(c["name"].(string), "-", "_")
	fmt.Fprintf(b, "%s(){\npkgdesc=\"%s\"\n%s\n}\n", name, c["desc"].(string), c["script"].(string))
}

func (s *Scheme) Build() string {
	b := strings.Builder{}

	fmt.Fprintf(&b, "pkgname=\"%s\"\n", s.Name)
	fmt.Fprintf(&b, "pkgver=\"%s\"\n", s.Version)
	fmt.Fprintf(&b, "pkgrel=\"%v\"\n", s.Release)
	fmt.Fprintf(&b, "pkgdesc=\"%s\"\n", s.Description)
	fmt.Fprintf(&b, "url=\"%s\"\n", URL)
	fmt.Fprintf(&b, "arch=\"%s\"\n", ARCH)
	fmt.Fprintf(&b, "license=\"%s\"\n", LICENSE)
	fmt.Fprintf(&b, "source=\"%s\"\n", strings.Join(s.Sources, " "))
	fmt.Fprintf(&b, "depends=\"%s\"\n", strings.Join(s.Dependencies, " "))
	fmt.Fprintf(&b, "makedepends=\"%s\"\n", strings.Join(s.MakeDependencies, " "))
	fmt.Fprintln(&b, "options=\"lib64\"")
	spkgs := []string{}
	for _, x := range s.Subpackages {
		if reflect.TypeOf(x).Kind() == reflect.String {
			spkgs = append(spkgs, x.(string))
		} else {
			c := interfaceMapToStringMap(x.(map[interface{}]interface{}))
			spkgs = append(spkgs, strings.ReplaceAll(c["name"].(string), "-", "_")+":"+c["name"].(string))
			doSpkg(&b, &c)
		}
	}
	fmt.Fprintf(&b, "subpackages=\"%s\"\n", strings.Join(spkgs, " "))
	fmt.Fprintln(&b, "prepare(){\n")
	for _, x := range s.Entries {
		fmt.Fprintln(&b, prepareEntry(&x))
	}
	fmt.Fprintf(&b, "}\npackage(){\n")
	for _, x := range s.Entries {
		fmt.Fprintln(&b, packageEntry(&x))
	}
	fmt.Fprintln(&b, "}\n")
	return b.String()
}
