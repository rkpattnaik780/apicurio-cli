package versions

import (
	"context"
	"github.com/apicurio/apicurio-cli/pkg/cmd/registry/artifact/util"

	"github.com/apicurio/apicurio-cli/pkg/cmd/registry/registrycmdutil"
	"github.com/apicurio/apicurio-cli/pkg/core/cmdutil/flagutil"
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
	artifact     string
	group        string
	outputFormat string

	registryID string

	IO             *iostreams.IOStreams
	Logger         logging.Logger
	Connection     factory.ConnectionFunc
	localizer      localize.Localizer
	Context        context.Context
	ServiceContext servicecontext.IContext
}

func NewVersionsCommand(f *factory.Factory) *cobra.Command {
	opts := &options{
		Connection:     f.Connection,
		IO:             f.IOStreams,
		localizer:      f.Localizer,
		Logger:         f.Logger,
		Context:        f.Context,
		ServiceContext: f.ServiceContext,
	}

	cmd := &cobra.Command{
		Use:     "versions",
		Short:   f.Localizer.MustLocalize("artifact.cmd.versions.description.short"),
		Long:    f.Localizer.MustLocalize("artifact.cmd.versions.description.long"),
		Example: f.Localizer.MustLocalize("artifact.cmd.versions.example"),
		Args:    cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {

			if opts.artifact == "" {
				return f.Localizer.MustLocalizeError("artifact.common.message.artifactIdRequired")
			}

			if opts.registryID != "" {
				return runGet(opts)
			}

			registryInstance, err := contextutil.GetCurrentRegistryInstance(f)
			if err != nil {
				return err
			}

			opts.registryID = registryInstance.GetId()
			return runGet(opts)
		},
	}

	cmd.Flags().StringVar(&opts.artifact, "artifact-id", "", opts.localizer.MustLocalize("artifact.common.id"))
	cmd.Flags().StringVarP(&opts.group, "group", "g", registrycmdutil.DefaultArtifactGroup, opts.localizer.MustLocalize("artifact.common.group"))
	cmd.Flags().StringVar(&opts.registryID, "instance-id", "", opts.localizer.MustLocalize("registry.common.flag.instance.id"))
	cmd.Flags().StringVarP(&opts.outputFormat, "output", "o", "json", opts.localizer.MustLocalize("artifact.common.message.output.formatNoTable"))

	flagutil.EnableOutputFlagCompletion(cmd)

	return cmd
}

func runGet(opts *options) error {
	format := util.OutputFormatFromString(opts.outputFormat)
	if format == util.UnknownOutputFormat || format == util.TableOutputFormat {
		return opts.localizer.MustLocalizeError("artifact.common.error.invalidOutputFormat")
	}

	conn, err := opts.Connection()
	if err != nil {
		return err
	}

	dataAPI, _, err := conn.API().ServiceRegistryInstance(opts.registryID)
	if err != nil {
		return err
	}

	if opts.group == registrycmdutil.DefaultArtifactGroup {
		opts.Logger.Info(opts.localizer.MustLocalize("registry.artifact.common.message.no.group", localize.NewEntry("DefaultArtifactGroup", registrycmdutil.DefaultArtifactGroup)))
	}

	opts.Logger.Info(opts.localizer.MustLocalize("artifact.common.message.artifact.versions.fetching"))

	request := dataAPI.VersionsApi.ListArtifactVersions(opts.Context, opts.group, opts.artifact)
	response, _, err := request.Execute()
	if err != nil {
		return registrycmdutil.TransformInstanceError(err)
	}

	opts.Logger.Info(icon.SuccessPrefix(), opts.localizer.MustLocalize("artifact.common.message.artifact.versions.fetched"))

	return util.Dump(opts.IO.Out, format, response, nil)
}
