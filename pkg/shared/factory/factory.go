package factory

import (
	"context"

	"github.com/apicurio/apicurio-cli/pkg/core/config"
	"github.com/apicurio/apicurio-cli/pkg/core/servicecontext"
	"github.com/apicurio/apicurio-cli/pkg/shared/connection"

	"github.com/apicurio/apicurio-cli/pkg/core/ioutil/iostreams"
	"github.com/apicurio/apicurio-cli/pkg/core/localize"
	"github.com/apicurio/apicurio-cli/pkg/core/logging"
)

// Factory is an abstract type which provides access to
// the root configuration and connections for the CLI
type Factory struct {
	// Type which defines the streams for the CLI
	IOStreams *iostreams.IOStreams
	// Interface to read/write to the config
	Config config.IConfig
	// Creates a connection to the API
	Connection ConnectionFunc
	// Returns a logger to create leveled logs in the application
	Logger logging.Logger
	// Localizer provides text to the commands
	Localizer localize.Localizer
	// Context returns the default context for the application
	Context context.Context
	// ServiceContext returns the identifiers for currently selected services for the context
	ServiceContext servicecontext.IContext
}

type ConnectionFunc func() (connection.Connection, error)
