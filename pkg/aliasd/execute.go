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
package aliasd

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/dragozir/aliasd/pkg/config"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v2"
)

var cwd string
var alias = filepath.Base(os.Args[0])

func execute() {
	wd, err := os.Getwd()
	if err != nil {
		log.Fatalf("Could not get cwd %s : %v", wd, err)
	}
	cwd = wd
	log.Infof("Running %s from %s", os.Args[0], wd)

	path := fmt.Sprintf("%s/%s.yaml", config.ConfigDir, alias)

	if _, err := os.Stat(path); os.IsNotExist(err) {
		log.Fatalf("%s proxy does not exist, aborting", path)
	}

	bytes, err := ioutil.ReadFile(path)
	if err != nil {
		log.Fatalf("Error loading config file : %v", err)
	}

	var spec config.ResourcesSpec

	if err := yaml.Unmarshal(bytes, &spec); err != nil {
		log.Fatalf("Error unmarshaling config file : %v", err)
	}
	fmt.Printf("Yaml: %+v\n", spec)

	mounts := ""
	image := spec.Resources[alias].Image
	args := strings.Join(os.Args[1:], " ")
	for _, v := range spec.Resources[alias].VolumeMounts {
		if strings.Contains(v.HostPath, "$(pwd)") {
			args = strings.Replace(args, wd, v.MountPath, -1)
			mounts += fmt.Sprintf("-v %s/:%s", wd, v.MountPath)
		} else {
			args = strings.Replace(args, v.HostPath, v.MountPath, -1)
			mounts += fmt.Sprintf("-v %s:%s", v.HostPath, v.MountPath)
		}
	}

	cmdStr := fmt.Sprintf("run --rm %s %s %s", mounts, image, args)

	log.Infof("Running %s", cmdStr)

	dockerCmd := exec.Command("docker", strings.Split(cmdStr, " ")...)

	if err := dockerCmd.Run(); err != nil {
		log.Fatalf("Could not execute %s with %v", dockerCmd.Path, dockerCmd.Args)
	}
}

var execCmd = &cobra.Command{
	Use:   "execute",
	Short: "Execute a command proxy",
	Run: func(cmd *cobra.Command, _ []string) {
		execute()
	},
}

func init() {
	rootCmd.AddCommand(execCmd)
}
