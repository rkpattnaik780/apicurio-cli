package serviceaccount

import (
	"github.com/apicurio/apicurio-cli/internal/doc"
	"github.com/apicurio/apicurio-cli/pkg/cmd/serviceaccount/create"
	"github.com/apicurio/apicurio-cli/pkg/cmd/serviceaccount/delete"
	"github.com/apicurio/apicurio-cli/pkg/cmd/serviceaccount/describe"
	"github.com/apicurio/apicurio-cli/pkg/cmd/serviceaccount/list"
	"github.com/apicurio/apicurio-cli/pkg/cmd/serviceaccount/resetcredentials"
	"github.com/apicurio/apicurio-cli/pkg/shared/factory"
	"github.com/spf13/cobra"
)

// NewServiceAccountCommand creates a new command sub-group to manage service accounts
func NewServiceAccountCommand(f *factory.Factory) *cobra.Command {
	cmd := &cobra.Command{
		Use:         "service-account",
		Annotations: map[string]string{doc.AnnotationName: "Service account commands"},
		Short:       f.Localizer.MustLocalize("serviceAccount.cmd.shortDescription"),
		Long:        f.Localizer.MustLocalize("serviceAccount.cmd.longDescription"),
		Example:     f.Localizer.MustLocalize("serviceAccount.cmd.example"),
		Args:        cobra.ExactArgs(1),
	}

	cmd.AddCommand(
		create.NewCreateCommand(f),
		list.NewListCommand(f),
		delete.NewDeleteCommand(f),
		resetcredentials.NewResetCredentialsCommand(f),
		describe.NewDescribeCommand(f),
	)

	return cmd
}
