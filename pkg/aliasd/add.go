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

	"github.com/dragozir/aliasd/pkg/config"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v2"
)

var specFile string

var addCmd = &cobra.Command{
	Use:   "add",
	Short: "Add a docker command proxy",
	Run: func(cmd *cobra.Command, _ []string) {
		if specFile != "" {
			bytes, err := ioutil.ReadFile(specFile)
			if err != nil {
				log.Fatalf("Error loading config file : %v", err)
			}

			var spec config.ResourcesSpec

			if err := yaml.Unmarshal(bytes, &spec); err != nil {
				log.Fatalf("Error unmarshaling config file : %v", err)
			}
			fmt.Printf("Yaml: %+v\n", spec)

			for name, resource := range spec.Resources {
				fmt.Printf("%s: %+v\n", name, resource)

				destLoc := fmt.Sprintf("%s/%s.yaml", config.ConfigDir, name)
				dest, err := os.Create(destLoc)
				if err != nil {
					log.Fatalf("Error creating config file %s : %v", destLoc, err)
				}
				defer dest.Close()

				log.Infof("Writing to %s", destLoc)

				spec := config.ResourcesSpec{}
				spec.Resources = make(map[string]config.Resource)
				spec.Resources[name] = resource

				d, err := yaml.Marshal(&spec)
				if err != nil {
					log.Fatalf("Error marshaling config file : %v", err)
				}
				dest.Write(d)

				symlink := fmt.Sprintf("%s/%s", config.BinDir, name)

				if err := os.Symlink(config.AliasdPath, symlink); err != nil {
					log.Fatalf("Could not create symlink at %s : %v", symlink, err)
				}
			}
		}
	},
}

func init() {
	rootCmd.AddCommand(addCmd)

	addCmd.PersistentFlags().StringVarP(&specFile, "file", "f", "", "resource spec file")
}
