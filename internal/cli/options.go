package cli

import (
	"fmt"
	"strings"

	"github.com/kyma-incubator/reconciler/pkg/db"
	"github.com/spf13/viper"

	"github.com/kyma-incubator/reconciler/pkg/app"
	"github.com/kyma-incubator/reconciler/pkg/logger"
	"go.uber.org/zap"
)

type Options struct {
	Migrate        bool
	Verbose        bool
	InitRegistry   bool
	NonInteractive bool
	OutputFormat   string
	logger         *zap.SugaredLogger
	Registry       *app.ApplicationRegistry //will be initialized during CLI bootstrap in main.go
}

func (o *Options) String() string {
	return fmt.Sprintf("CLI options: migrate=%t verbose=%t non-interactive=%t",
		o.Migrate, o.Verbose, o.NonInteractive)
}

func (o *Options) Validate() error {
	for _, supportedFormat := range SupportedOutputFormats {
		if supportedFormat == o.OutputFormat {
			return nil
		}
	}
	return fmt.Errorf("Output format '%s' not supported - choose between '%s'", o.OutputFormat, strings.Join(SupportedOutputFormats, "', '"))
}

func (o *Options) Logger() *zap.SugaredLogger {
	if o.logger == nil {
		o.logger = logger.NewLogger(o.Verbose)
	}
	return o.logger
}

func (o *Options) InitApplicationRegistry(forceInitialization bool) error {
	if forceInitialization || o.InitRegistry {
		dbConnFact, err := db.NewConnectionFactory(viper.ConfigFileUsed(), o.Migrate, o.Verbose)
		if err != nil {
			return err
		}
		o.Registry, err = app.NewApplicationRegistry(dbConnFact, o.Verbose)
		return err
	}
	return nil
}
