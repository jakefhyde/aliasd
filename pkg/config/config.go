/*
Copyright © 2021 dragozir

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in
all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
THE SOFTWARE.
*/
package config

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/mitchellh/go-homedir"
)

type ServerSpec struct {
	Server Server `yaml:"spec"`
}

type Server struct {
	Address  string `yaml:"address"`
	Port     uint16 `yaml:"port"`
	AllowAll bool   `yaml:"allowAll"`
}

type ResourcesSpec struct {
	Resources map[string]Resource `yaml:"resources"`
}

type Resource struct {
	Image        string        `yaml:"image"`
	VolumeMounts []VolumeMount `yaml:"volumeMounts"`
}

type VolumeMount struct {
	MountPath string `yaml:"mountPath"`
	HostPath  string `yaml:"hostPath"`
}

var HomeDir string
var Prefix string
var BinDir string
var ConfigDir string
var Dirs []string
var AliasdPath string

func init() {
	home, err := homedir.Dir()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	HomeDir = home
	Prefix = HomeDir + "/.aliasd"
	BinDir = Prefix + "/bin"
	ConfigDir = Prefix + "/config"
	Dirs = []string{BinDir, ConfigDir}

	exe, err := os.Executable()
	if err != nil {
		return
	}

	aliasdPath, err := filepath.EvalSymlinks(exe)
	if err != nil {
		log.Fatalf("Could not retrieve runtime name %v", err)
	}
	AliasdPath = aliasdPath
}
