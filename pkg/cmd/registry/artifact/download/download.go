package download

import (
	"context"
	"fmt"
	"os"

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

var unusedFlagIdValue int64 = -1

type options struct {
	group string

	contentId  int64
	globalId   int64
	hash       string
	outputFile string

	registryID string

	IO             *iostreams.IOStreams
	Logger         logging.Logger
	Connection     factory.ConnectionFunc
	localizer      localize.Localizer
	Context        context.Context
	ServiceContext servicecontext.IContext
}

// NewDownloadCommand creates a new command for downloading binary content for registry artifacts.
func NewDownloadCommand(f *factory.Factory) *cobra.Command {
	opts := &options{
		Connection:     f.Connection,
		IO:             f.IOStreams,
		localizer:      f.Localizer,
		Logger:         f.Logger,
		Context:        f.Context,
		ServiceContext: f.ServiceContext,
	}

	cmd := &cobra.Command{
		Use:     "download",
		Short:   f.Localizer.MustLocalize("artifact.cmd.download.description.short"),
		Long:    f.Localizer.MustLocalize("artifact.cmd.download.description.long"),
		Example: f.Localizer.MustLocalize("artifact.cmd.download.example"),
		Args:    cobra.RangeArgs(0, 1),
		RunE: func(cmd *cobra.Command, args []string) error {
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

	cmd.Flags().StringVarP(&opts.group, "group", "g", registrycmdutil.DefaultArtifactGroup, opts.localizer.MustLocalize("artifact.common.group"))
	cmd.Flags().StringVar(&opts.hash, "hash", "", opts.localizer.MustLocalize("artifact.common.sha"))
	cmd.Flags().Int64VarP(&opts.globalId, "global-id", "", unusedFlagIdValue, opts.localizer.MustLocalize("artifact.common.global.id"))
	cmd.Flags().Int64VarP(&opts.contentId, "content-id", "", unusedFlagIdValue, opts.localizer.MustLocalize("artifact.common.content.id"))

	cmd.Flags().StringVarP(&opts.outputFile, "output-file", "", "", opts.localizer.MustLocalize("artifact.common.message.file.location"))
	cmd.Flags().StringVar(&opts.registryID, "instance-id", "", opts.localizer.MustLocalize("artifact.common.registryIdToUse"))

	flagutil.EnableOutputFlagCompletion(cmd)

	return cmd
}

func runGet(opts *options) error {
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

	opts.Logger.Info(opts.localizer.MustLocalize("artifact.common.message.fetching.artifact"))

	var dataFile *os.File
	// nolint
	if opts.contentId != unusedFlagIdValue {
		request := dataAPI.ArtifactsApi.GetContentById(opts.Context, opts.contentId)
		dataFile, _, err = request.Execute()
	} else if opts.globalId != unusedFlagIdValue {
		request := dataAPI.ArtifactsApi.GetContentByGlobalId(opts.Context, opts.globalId)
		dataFile, _, err = request.Execute()
	} else if opts.hash != "" {
		request := dataAPI.ArtifactsApi.GetContentByHash(opts.Context, opts.hash)
		dataFile, _, err = request.Execute()
	} else {
		return opts.localizer.MustLocalizeError("artifact.cmd.common.error.specify.contentId.globalId.hash")
	}

	if err != nil {
		return registrycmdutil.TransformInstanceError(err)
	}

	fileContent, err := os.ReadFile(dataFile.Name())
	if err != nil {
		return err
	}
	if opts.outputFile != "" {
		err := os.WriteFile(opts.outputFile, fileContent, 0600)
		if err != nil {
			return err
		}
	} else {
		// Print to stdout
		fmt.Fprintf(os.Stdout, "%v\n", string(fileContent))
	}

	opts.Logger.Info(icon.SuccessPrefix(), opts.localizer.MustLocalize("artifact.common.message.fetched.successfully"))
	return nil
}
