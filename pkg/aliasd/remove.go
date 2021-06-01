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
	"os"

	"github.com/dragozir/aliasd/pkg/config"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var name string

var removeCmd = &cobra.Command{
	Use:   "remove",
	Short: "Remove a docker command proxy",
	Run: func(cmd *cobra.Command, _ []string) {
		path := fmt.Sprintf("%s/%s.yaml", config.ConfigDir, name)
		if _, err := os.Stat(path); os.IsNotExist(err) {
			log.Warnf("%s config not exist, searching for symlink", path)
		} else {
			if err := os.Remove(path); err != nil {
				log.Warnf("%s could not be removed", path)
			}
		}

		symlink := fmt.Sprintf("%s/%s", config.BinDir, name)
		if err := os.Remove(symlink); err != nil {
			log.Warnf("%s could not be removed", symlink)
		}
	},
}

func init() {
	rootCmd.AddCommand(removeCmd)

	removeCmd.Flags().StringVarP(&name, "name", "n", "", "resource name")
	removeCmd.MarkFlagRequired("name")
}
