package delete

import (
	"context"

	"github.com/apicurio/apicurio-cli/pkg/cmd/serviceaccount/svcaccountcmdutil/validation"
	"github.com/apicurio/apicurio-cli/pkg/core/cmdutil/flagutil"
	"github.com/apicurio/apicurio-cli/pkg/core/config"
	"github.com/apicurio/apicurio-cli/pkg/core/ioutil/icon"
	"github.com/apicurio/apicurio-cli/pkg/core/ioutil/iostreams"
	"github.com/apicurio/apicurio-cli/pkg/core/localize"
	"github.com/apicurio/apicurio-cli/pkg/core/logging"
	"github.com/apicurio/apicurio-cli/pkg/shared/factory"

	svcacctmgmterrors "github.com/redhat-developer/app-services-sdk-core/app-services-sdk-go/serviceaccountmgmt/apiv1/error"

	"github.com/AlecAivazis/survey/v2"
	"github.com/spf13/cobra"
)

type options struct {
	IO         *iostreams.IOStreams
	Config     config.IConfig
	Connection factory.ConnectionFunc
	Logger     logging.Logger
	localizer  localize.Localizer
	Context    context.Context

	id    string
	force bool
}

// NewDeleteCommand creates a new command to delete a service account
func NewDeleteCommand(f *factory.Factory) *cobra.Command {
	opts := &options{
		Config:     f.Config,
		Connection: f.Connection,
		Logger:     f.Logger,
		IO:         f.IOStreams,
		localizer:  f.Localizer,
		Context:    f.Context,
	}

	cmd := &cobra.Command{
		Use:     "delete",
		Short:   opts.localizer.MustLocalize("serviceAccount.delete.cmd.shortDescription"),
		Long:    opts.localizer.MustLocalize("serviceAccount.delete.cmd.longDescription"),
		Example: opts.localizer.MustLocalize("serviceAccount.delete.cmd.example"),
		Args:    cobra.NoArgs,
		RunE: func(cmd *cobra.Command, _ []string) error {
			if !opts.IO.CanPrompt() && !opts.force {
				return flagutil.RequiredWhenNonInteractiveError("yes")
			}

			validator := &validation.Validator{
				Localizer: opts.localizer,
			}

			validID := validator.ValidateUUID(opts.id)
			if validID != nil {
				return validID
			}

			return runDelete(opts)
		},
	}

	cmd.Flags().StringVar(&opts.id, "id", "", opts.localizer.MustLocalize("serviceAccount.delete.flag.id.description"))
	cmd.Flags().BoolVarP(&opts.force, "yes", "y", false, opts.localizer.MustLocalize("serviceAccount.delete.flag.yes.description"))

	_ = cmd.MarkFlagRequired("id")

	return cmd
}

func runDelete(opts *options) (err error) {
	conn, err := opts.Connection()
	if err != nil {
		return err
	}

	_, httpRes, err := conn.API().ServiceAccountMgmt().GetServiceAccount(opts.Context, opts.id).Execute()
	if httpRes != nil {
		defer httpRes.Body.Close()
	}

	if apiErr := svcacctmgmterrors.GetAPIError(err); apiErr != nil {
		switch apiErr.GetError() {
		case "service_account_not_found":
			return opts.localizer.MustLocalizeError("serviceAccount.common.error.notFoundError", localize.NewEntry("ID", opts.id))
		default:
			return err
		}
	}

	if !opts.force {
		var confirmDelete bool
		promptConfirmDelete := &survey.Confirm{
			Message: opts.localizer.MustLocalize("serviceAccount.delete.input.confirmDelete.message", localize.NewEntry("ID", opts.id)),
		}

		err = survey.AskOne(promptConfirmDelete, &confirmDelete)
		if err != nil {
			return err
		}

		if !confirmDelete {
			opts.Logger.Debug(opts.localizer.MustLocalize("serviceAccount.delete.log.debug.deleteNotConfirmed"))
			return nil
		}
	}

	return deleteServiceAccount(opts)
}

func deleteServiceAccount(opts *options) error {
	conn, err := opts.Connection()
	if err != nil {
		return err
	}

	httpRes, err := conn.API().ServiceAccountMgmt().DeleteServiceAccount(opts.Context, opts.id).Execute()
	if httpRes != nil {
		defer httpRes.Body.Close()
	}

	if apiErr := svcacctmgmterrors.GetAPIError(err); apiErr != nil {
		switch apiErr.GetError() {
		case "service_account_access_invalid":
			return opts.localizer.MustLocalizeError("serviceAccount.common.error.forbidden", localize.NewEntry("Operation", "delete"))
		case "service_account_not_found":
			return opts.localizer.MustLocalizeError("serviceAccount.common.error.notFoundError", localize.NewEntry("ID", opts.id))
		default:
			return err
		}
	}

	opts.Logger.Info(icon.SuccessPrefix(), opts.localizer.MustLocalize("serviceAccount.delete.log.info.deleteSuccess"))

	return nil
}
