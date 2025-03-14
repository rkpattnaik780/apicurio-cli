package delete

import (
	"context"

	"github.com/apicurio/apicurio-cli/pkg/cmd/context/contextcmdutil"
	"github.com/apicurio/apicurio-cli/pkg/core/ioutil/icon"
	"github.com/apicurio/apicurio-cli/pkg/core/ioutil/iostreams"
	"github.com/apicurio/apicurio-cli/pkg/core/localize"
	"github.com/apicurio/apicurio-cli/pkg/core/logging"
	"github.com/apicurio/apicurio-cli/pkg/core/servicecontext"
	"github.com/apicurio/apicurio-cli/pkg/shared/contextutil"
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

	name string
}

// NewDeleteCommand command for deleting service contexts
func NewDeleteCommand(f *factory.Factory) *cobra.Command {
	opts := &options{
		Connection:     f.Connection,
		Logger:         f.Logger,
		IO:             f.IOStreams,
		localizer:      f.Localizer,
		Context:        f.Context,
		ServiceContext: f.ServiceContext,
	}

	cmd := &cobra.Command{
		Use:     "delete",
		Short:   f.Localizer.MustLocalize("context.delete.cmd.shortDescription"),
		Long:    f.Localizer.MustLocalize("context.delete.cmd.longDescription"),
		Example: f.Localizer.MustLocalize("context.delete.cmd.example"),
		Args:    cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			return runDelete(opts)
		},
	}

	flags := contextcmdutil.NewFlagSet(cmd, f)

	flags.AddContextName(&opts.name)

	return cmd
}

func runDelete(opts *options) error {

	svcContext, err := opts.ServiceContext.Load()
	if err != nil {
		return err
	}

	currCtx := svcContext.CurrentContext

	if opts.name == "" || opts.name == currCtx {

		if currCtx == "" {
			return opts.localizer.MustLocalizeError("context.common.error.notSet")
		}

		opts.name = currCtx

		svcContext.CurrentContext = ""

		opts.Logger.Info(opts.localizer.MustLocalize("context.delete.log.warning.currentUnset"))
	}

	if _, err = contextutil.GetContext(svcContext, opts.localizer, opts.name); err != nil {
		return err
	}

	delete(svcContext.Contexts, opts.name)

	err = opts.ServiceContext.Save(svcContext)
	if err != nil {
		return err
	}

	opts.Logger.Info(icon.SuccessPrefix(), opts.localizer.MustLocalize("context.delete.log.successMessage"))

	return nil

}
