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
	"bytes"
	"io/ioutil"
	"regexp"
	"text/template"

	"gopkg.in/yaml.v2"
)

type Scheme struct {
	Name             string
	Description      string        `yaml:"desc"`
	Version          string        `yaml:"ver"`
	Release          int           `yaml:"rel"`
	Subpackages      []interface{} `yaml:"subpkgs"`
	Sources          []string      `yaml:"srcs"`
	Entries          []Entry
	Defines          map[string]string
	Dependencies     []string `yaml:"deps"`
	MakeDependencies []string `yaml:"makedeps"`
}

type Entry struct {
	Before           string
	After            string
	Folder           string
	Type             string
	InstallOptions   []string `yaml:"makeopts"`
	ConfigureOptions []string `yaml:"cfgopts"`
	Build            string
	Environment      []string `yaml:"env"`
}

func (s *Scheme) Load(filename string) error {
	content, err := ioutil.ReadFile(filename)
	if err != nil {
		return err
	}
	r := regexp.MustCompile("{{(.*)}}")
	p := r.ReplaceAll(content, []byte(""))

	if err = yaml.Unmarshal([]byte(p), s); err != nil {
		return err
	}

	mp := s.schemeToMap()

	tmpl, err := template.New("scheme").Parse(string(content))
	if err != nil {
		return err
	}

	var buff bytes.Buffer
	if err = tmpl.Execute(&buff, mp); err != nil {
		return err
	}

	return yaml.Unmarshal(buff.Bytes(), s)
}
