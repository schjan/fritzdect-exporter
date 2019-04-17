package main

import (
	"context"
	"fmt"
	"github.com/giantswarm/microerror"
	"github.com/giantswarm/microkit/command"
	"github.com/giantswarm/micrologger"
	"github.com/schjan/fritzdect-exporter/flag"
	"github.com/schjan/fritzdect-exporter/server"
	"github.com/schjan/fritzdect-exporter/service"
	"github.com/spf13/viper"
	"os"

	microserver "github.com/giantswarm/microkit/server"
)

var (
	description string     = "Fritzdect-exporter exposes values from Fritz!DECT thermostats to prometheus"
	f           *flag.Flag = flag.New()
	gitCommit   string     = "n/a"
	name        string     = "fritzdect-exporter"
	source      string     = "https://github.com/schjan/fritzdect-exporter"
)

func main() {
	err := mainErr()
	if err != nil {
		panic(fmt.Sprintf("%#v", err))
	}
}

func mainErr() error {
	var err error

	// Create a new logger that is used by all packages.
	var newLogger micrologger.Logger
	{
		c := micrologger.Config{
			IOWriter: os.Stdout,
		}
		newLogger, err = micrologger.New(c)
		if err != nil {
			return microerror.Maskf(err, "micrologger.New")
		}
	}

	// Define server factory to create the custom server once all command line
	// flags are parsed and all microservice configuration is processed.
	newServerFactory := func(v *viper.Viper) microserver.Server {
		// New custom service implements the business logic.
		var newService *service.Service
		{
			c := service.Config{
				Flag:   f,
				Logger: newLogger,
				Viper:  v,

				Description: description,
				GitCommit:   gitCommit,
				ProjectName: name,
				Source:      source,
			}
			newService, err = service.New(c)
			if err != nil {
				// panic because we're inside embedded function
				panic(fmt.Sprintf("%#v\n", microerror.Maskf(err, "service.New")))
			}

			go newService.Boot(context.Background())
		}

		// New custom server that bundles microkit endpoints.
		var newServer microserver.Server
		{
			c := server.Config{
				Logger:      newLogger,
				Service:     newService,
				Viper:       v,
				ProjectName: name,
			}

			newServer, err = server.New(c)
			if err != nil {
				panic(fmt.Sprintf("%#v\n", microerror.Maskf(err, "server.New")))
			}
		}

		return newServer

	}

	// Create a new microkit command that manages operator daemon.
	var newCommand command.Command
	{
		c := command.Config{
			Logger:        newLogger,
			ServerFactory: newServerFactory,

			Description: description,
			GitCommit:   gitCommit,
			Name:        name,
			Source:      source,
		}

		newCommand, err = command.New(c)
		if err != nil {
			return microerror.Maskf(err, "command.New")
		}
	}

	daemonCommand := newCommand.DaemonCommand().CobraCommand()

	daemonCommand.PersistentFlags().String(f.Service.FritzBox.Url, "", "Root URL of fritzbox (Default http://fritz.box).")
	daemonCommand.PersistentFlags().String(f.Service.FritzBox.User.Name, "", "FritzBox username (optional).")
	daemonCommand.PersistentFlags().String(f.Service.FritzBox.User.Password, "", "FritzBox password.")

	return newCommand.CobraCommand().Execute()

	//c, err := collector.NewFritzDect()
	//if err != nil {
	//	return err
	//}
	//
	//prometheus.MustRegister(c)
}
