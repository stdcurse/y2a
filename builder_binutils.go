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
	"strings"
)

func BinutilsBuilderPackage(b *strings.Builder, s *Script) {
	fmt.Fprintln(b, "mkdir build\ncd build")
	fmt.Fprintf(b, "../configure --prefix=/usr --libdir=/lib --sysconfdir=/etc --localstatedir=/var")
	for _, x := range s.ConfigureOptions {
		fmt.Fprintf(b, " --%s", x)
	}
	s.InstallOptions = append(s.InstallOptions, "DESTDIR=\"$pkgdir\"")
	fmt.Fprintf(b, "\nmake\nmake install %s\n", strings.Join(s.InstallOptions, " "))
}
