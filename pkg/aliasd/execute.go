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

func execute(args []string) {
	wd, err := os.Getwd()
	if err != nil {
		log.Fatalf("Could not get cwd %s : %v", wd, err)
	}
	cwd = wd
	log.Debugf("Running %s from %s", os.Args[0], wd)

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

	mounts := ""
	image := spec.Resources[alias].Image
	argString := strings.Join(args, " ")
	for _, v := range spec.Resources[alias].VolumeMounts {
		if strings.Contains(v.HostPath, "$(pwd)") {
			argString = strings.Replace(argString, wd, v.MountPath, -1)
			mounts += fmt.Sprintf("-v %s/:%s", wd, v.MountPath)
		} else {
			argString = strings.Replace(argString, v.HostPath, v.MountPath, -1)
			mounts += fmt.Sprintf("-v %s:%s", v.HostPath, v.MountPath)
		}
	}

	cmdStr := fmt.Sprintf("run --rm %s %s %s", mounts, image, argString)

	log.Debugf("Running %s", cmdStr)

	dockerCmd := exec.Command("docker", strings.Split(cmdStr, " ")...)

	if out, err := dockerCmd.Output(); err != nil {
		log.Fatalf("Could not execute %s with %v with %v", dockerCmd.Path, dockerCmd.Args, err)
	} else {
		fmt.Printf("%s", out)
	}
}

var execCmd = &cobra.Command{
	Use:   "execute",
	Short: "Execute a command proxy",
	Long: `Executes a specific comman proxy by name, which can be useful if your
environment makes it difficult to create and maintain symlinks.`,
	Args: cobra.ArbitraryArgs,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) < 2 || args[0] == "-h" || args[0] == "--help" || args[1] == "-h" || args[2] == "--help" {
			cmd.Help()
			return
		}
		if args[0] == "-n" || args[0] == "--name" {
			alias = args[1]
			execute(args[2:]) // forward all subsequent args to proxy
		}
	},
}

func init() {
	rootCmd.AddCommand(execCmd)
	execCmd.SetUsageTemplate("aliasd execute -n [name] [proxy_flags]")

	// parse manually
	execCmd.DisableFlagParsing = true
}
