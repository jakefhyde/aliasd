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
	"context"
	"fmt"
	"io/ioutil"
	"net"
	"os"

	"github.com/dragozir/aliasd/pkg/config"
	"github.com/dragozir/aliasd/pkg/proto"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"google.golang.org/grpc"
	"gopkg.in/yaml.v2"
)

var cfgFile string

type aliasdServer struct {
	proto.UnimplementedAliasdServer
}

func (s *aliasdServer) Add(ctx context.Context, request *proto.AddRequest) (*proto.AddResponse, error) {
	// todo find better way to do this
	spec := config.ResourcesSpec{
		Resources: map[string]config.Resource{
			name: {
				Image:        request.Image,
				VolumeMounts: make([]config.VolumeMount, len(request.VolumeMounts)),
			},
		},
	}

	for i, v := range request.VolumeMounts {
		spec.Resources[name].VolumeMounts[i] = config.VolumeMount{
			MountPath: v.MountPath,
			HostPath:  v.HostPath,
		}
	}

	path := fmt.Sprintf("/etc/aliasd/config/%s.yaml", request.Name)
	dest, err := os.Create(path)
	if err != nil {
		return &proto.AddResponse{Result: proto.AddResult_ConfigError}, fmt.Errorf("Error creating config file %s : %v", path, err)
	}
	defer dest.Close()

	if d, err := yaml.Marshal(&spec); err != nil {
		return &proto.AddResponse{Result: proto.AddResult_ConfigError}, fmt.Errorf("Error marshaling config file : %v", err)
	} else {
		if _, err := dest.Write(d); err != nil {
			return &proto.AddResponse{Result: proto.AddResult_ConfigError}, fmt.Errorf("Error writing config file : %v", err)
		}
		return &proto.AddResponse{Result: proto.AddResult_Success}, nil
	}
}

var serverCmd = &cobra.Command{
	Use:   "server",
	Short: "Run an aliasd server",
	RunE: func(cmd *cobra.Command, _ []string) error {
		bytes, err := ioutil.ReadFile(cfgFile)
		if err != nil {
			return fmt.Errorf("Error loading config file : %v", err)
		}

		var spec config.ServerSpec

		if err := yaml.Unmarshal(bytes, &spec); err != nil {
			return fmt.Errorf("Error unmarshaling config file : %v", err)
		}

		l, err := net.Listen("tcp", fmt.Sprintf("%s:%d", spec.Server.Address, spec.Server.Port))
		if err != nil {
			return fmt.Errorf("Failed to listen :%v", err)
		}

		s := grpc.NewServer()
		proto.RegisterAliasdServer(s, &aliasdServer{})
		log.Infof("Server listening at %v", l.Addr())

		if err := s.Serve(l); err != nil {
			return fmt.Errorf("Failed to serve: %v", err)
		}

		return nil
	},
}

func init() {
	rootCmd.AddCommand(serverCmd)

	// todo add viper configuration
	serverCmd.PersistentFlags().StringVarP(&cfgFile, "config", "c", "", "config file")
	serverCmd.MarkFlagRequired("config")
}
