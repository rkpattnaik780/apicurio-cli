package list

import (
	"context"
	"fmt"

	"github.com/apicurio/apicurio-cli/pkg/cmd/context/contextcmdutil"
	"github.com/apicurio/apicurio-cli/pkg/core/ioutil/dump"
	"github.com/apicurio/apicurio-cli/pkg/core/ioutil/icon"
	"github.com/apicurio/apicurio-cli/pkg/core/ioutil/iostreams"
	"github.com/apicurio/apicurio-cli/pkg/core/localize"
	"github.com/apicurio/apicurio-cli/pkg/core/logging"
	"github.com/apicurio/apicurio-cli/pkg/core/servicecontext"
	"github.com/apicurio/apicurio-cli/pkg/shared/factory"
	"github.com/spf13/cobra"
)

type options struct {
	IO             *iostreams.IOStreams
	Logger         logging.Logger
	Connection     factory.ConnectionFunc
	localizer      localize.Localizer
	Context        context.Context
	ServiceContext servicecontext.IContext

	outputFormat string
}

// NewListCommand creates a new command to list available contexts
func NewListCommand(f *factory.Factory) *cobra.Command {

	opts := &options{
		Connection:     f.Connection,
		IO:             f.IOStreams,
		Logger:         f.Logger,
		localizer:      f.Localizer,
		ServiceContext: f.ServiceContext,
	}

	cmd := &cobra.Command{
		Use:     "list",
		Short:   f.Localizer.MustLocalize("context.list.cmd.shortDescription"),
		Long:    f.Localizer.MustLocalize("context.list.cmd.longDescription"),
		Example: f.Localizer.MustLocalize("context.list.cmd.example"),
		Args:    cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			return runList(opts)
		},
	}

	flags := contextcmdutil.NewFlagSet(cmd, f)

	flags.AddOutput(&opts.outputFormat)

	return cmd
}

func runList(opts *options) error {

	svcContext, err := opts.ServiceContext.Load()
	if err != nil {
		return err
	}

	svcContextsMap := svcContext.Contexts

	if svcContextsMap == nil {
		opts.Logger.Info(opts.localizer.MustLocalize("context.list.log.info.noContexts"))
		return nil
	}

	if opts.outputFormat == dump.EmptyFormat {
		currentCtx := svcContext.CurrentContext
		var profileList string

		for name := range svcContextsMap {
			if currentCtx != "" && name == currentCtx {
				profileList += fmt.Sprintln(name, icon.SuccessPrefix())
			} else {
				profileList += fmt.Sprintln(name)
			}
		}
		opts.Logger.Info(profileList)
		return nil
	}

	return dump.Formatted(opts.IO.Out, opts.outputFormat, svcContextsMap)

}
