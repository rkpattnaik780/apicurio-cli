package create

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

// NewCreateCommand creates a new command to create contexts
func NewCreateCommand(f *factory.Factory) *cobra.Command {

	opts := &options{
		Connection:     f.Connection,
		IO:             f.IOStreams,
		Logger:         f.Logger,
		localizer:      f.Localizer,
		ServiceContext: f.ServiceContext,
	}

	cmd := &cobra.Command{
		Use:     "create",
		Short:   f.Localizer.MustLocalize("context.create.cmd.shortDescription"),
		Long:    f.Localizer.MustLocalize("context.create.cmd.longDescription"),
		Example: f.Localizer.MustLocalize("context.create.cmd.example"),
		Args:    cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {

			return runCreate(opts)
		},
	}

	flags := contextcmdutil.NewFlagSet(cmd, f)

	flags.StringVar(
		&opts.name,
		"name",
		"",
		opts.localizer.MustLocalize("context.common.flag.name"),
	)

	return cmd

}

func runCreate(opts *options) error {

	svcContext, err := opts.ServiceContext.Load()
	if err != nil {
		return err
	}

	profileValidator := &contextcmdutil.Validator{
		Localizer:  opts.localizer,
		SvcContext: svcContext,
	}

	svcContextsMap := svcContext.Contexts

	if svcContextsMap == nil {
		svcContextsMap = make(map[string]servicecontext.ServiceConfig)
	}

	err = profileValidator.ValidateName(opts.name)
	if err != nil {
		return err
	}

	err = profileValidator.ValidateNameIsAvailable(opts.name)
	if err != nil {
		return err
	}

	context, _ := contextutil.GetContext(svcContext, opts.localizer, opts.name)
	if context != nil {
		return opts.localizer.MustLocalizeError("context.create.log.alreadyExists", localize.NewEntry("Name", opts.name))
	}

	svcContextsMap[opts.name] = servicecontext.ServiceConfig{}
	svcContext.CurrentContext = opts.name

	svcContext.Contexts = svcContextsMap

	err = opts.ServiceContext.Save(svcContext)
	if err != nil {
		return err
	}

	opts.Logger.Info(icon.SuccessPrefix(), opts.localizer.MustLocalize("context.create.log.successMessage", localize.NewEntry("Name", opts.name)))

	return nil
}
