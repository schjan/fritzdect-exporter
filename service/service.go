package service

import (
	"context"
	"github.com/giantswarm/microendpoint/service/version"
	"github.com/giantswarm/microerror"
	"github.com/giantswarm/micrologger"
	"github.com/schjan/fritzdect-exporter/client"
	"github.com/schjan/fritzdect-exporter/collector"
	"github.com/schjan/fritzdect-exporter/flag"
	"github.com/spf13/viper"
	"sync"
)

// Config represents the configuration used to create a new service.
type Config struct {
	// Dependencies.
	Logger micrologger.Logger

	// Settings.
	Flag  *flag.Flag
	Viper *viper.Viper

	Description string
	GitCommit   string
	ProjectName string
	Source      string
}

// Service is a type providing implementation of microkit service interface.
type Service struct {
	Version *version.Service

	// Internals.
	bootOnce          sync.Once
	exporterCollector *collector.Set
}

// New creates a new configured service object.
func New(config Config) (*Service, error) {
	// Dependencies.
	if config.Logger == nil {
		return nil, microerror.Maskf(invalidConfigError, "%T.Logger must not be empty", config)
	}

	// Settings.
	if config.Flag == nil {
		return nil, microerror.Maskf(invalidConfigError, "%T.Flag must not be empty", config)
	}
	if config.Viper == nil {
		return nil, microerror.Maskf(invalidConfigError, "%T.Viper must not be empty", config)
	}

	var err error

	var newClient client.Client
	{
		c := client.Config{
			Url:      config.Viper.GetString(config.Flag.Service.FritzBox.Url),
			Username: config.Viper.GetString(config.Flag.Service.FritzBox.User.Name),
			Password: config.Viper.GetString(config.Flag.Service.FritzBox.User.Password),

			Logger: config.Logger,
		}

		newClient, err = client.New(c)
		if err != nil {
			return nil, microerror.Mask(err)
		}
	}

	var exporterCollector *collector.Set
	{
		c := collector.SetConfig{
			Logger: config.Logger,
			Client: newClient,
		}

		exporterCollector, err = collector.NewSet(c)
		if err != nil {
			return nil, microerror.Mask(err)
		}
	}

	var versionService *version.Service
	{
		versionConfig := version.Config{
			Description: config.Description,
			GitCommit:   config.GitCommit,
			Name:        config.ProjectName,
			Source:      config.Source,
		}

		versionService, err = version.New(versionConfig)
		if err != nil {
			return nil, microerror.Mask(err)
		}
	}

	s := &Service{
		Version:           versionService,
		exporterCollector: exporterCollector,
	}

	return s, nil
}

func (s *Service) Boot(ctx context.Context) {
	s.bootOnce.Do(func() {
		go s.exporterCollector.Boot(ctx)
	})
}
