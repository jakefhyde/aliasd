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
	"path/filepath"

	"github.com/dragozir/aliasd/pkg/config"
	"github.com/spf13/cobra"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

var CfgFile string

var rootCmd = &cobra.Command{
	Use:   "aliasd",
	Short: "Docker CLI proxy",
	Long:  `aliasd is a CLI utility for running docker images aliased to the commands they normally provide`,
	RunE: func(cmd *cobra.Command, args []string) error {
		return cmd.Help()
	},
}

func Execute() {
	if filepath.Base(os.Args[0]) != "aliasd" {
		execute()
		return
	}
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}

func initDotfiles() {

	createDir := func(path string) {
		if _, err := os.Stat(path); os.IsNotExist(err) {
			log.Infof("%s does not exist, creating now", path)
			err := os.MkdirAll(path, os.FileMode(0755))
			if err != nil {
				log.Fatalf("Could not create %s : %v", path, err)
			}
		}
	}

	for _, dir := range config.Dirs {
		createDir(dir)
	}
}

func init() {
	initDotfiles()

	cobra.OnInitialize(initConfig)

	rootCmd.PersistentFlags().StringVar(&CfgFile, "config", "", "config file (default is $HOME/.aliasd.yaml)")
}

func initConfig() {
	if CfgFile != "" {
		viper.SetConfigFile(CfgFile)
	} else {

		viper.AddConfigPath(config.HomeDir)
		viper.SetConfigName(".aliasd")
	}

	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err == nil {
		fmt.Println("Using config file:", viper.ConfigFileUsed())
	}
}
