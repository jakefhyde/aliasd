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
	RunE: func(cmd *cobra.Command, _ []string) error {
		bytes, err := ioutil.ReadFile(specFile)
		if err != nil {
			return fmt.Errorf("Error loading config file : %v", err)
		}

		var spec config.ResourcesSpec

		if err := yaml.Unmarshal(bytes, &spec); err != nil {
			return fmt.Errorf("Error unmarshaling config file : %v", err)
		}

		for name, resource := range spec.Resources {

			destLoc := fmt.Sprintf("%s/%s.yaml", config.ConfigDir, name)
			dest, err := os.Create(destLoc)
			if err != nil {
				return fmt.Errorf("Error creating config file %s : %v", destLoc, err)
			}
			defer dest.Close()

			log.Infof("Writing config to %s", destLoc)

			spec := config.ResourcesSpec{
				Resources: map[string]config.Resource{
					name: resource,
				},
			}

			if d, err := yaml.Marshal(&spec); err != nil {
				return fmt.Errorf("Error marshaling config file : %v", err)
			} else {
				dest.Write(d)
			}

			symlink := fmt.Sprintf("%s/%s", config.BinDir, name)

			if err := os.Symlink(config.AliasdPath, symlink); err != nil {
				return fmt.Errorf("Could not create symlink at %s : %v", symlink, err)
			}
		}
		return nil
	},
}

func init() {
	rootCmd.AddCommand(addCmd)

	addCmd.PersistentFlags().StringVarP(&specFile, "file", "f", "", "resource spec file")
	addCmd.MarkFlagRequired("file")
}
