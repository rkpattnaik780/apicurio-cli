package root

import (
	"flag"

	"github.com/apicurio/apicurio-cli/pkg/cmd/completion"
	contextcmd "github.com/apicurio/apicurio-cli/pkg/cmd/context"
	"github.com/apicurio/apicurio-cli/pkg/cmd/login"
	"github.com/apicurio/apicurio-cli/pkg/cmd/logout"
	"github.com/apicurio/apicurio-cli/pkg/cmd/registry"
	"github.com/apicurio/apicurio-cli/pkg/cmd/request"
	"github.com/apicurio/apicurio-cli/pkg/cmd/serviceaccount"
	"github.com/apicurio/apicurio-cli/pkg/core/cmdutil/flagutil"
	"github.com/apicurio/apicurio-cli/pkg/shared/factory"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"

	"github.com/apicurio/apicurio-cli/pkg/cmd/registry/artifact"
	"github.com/apicurio/apicurio-cli/pkg/cmd/registry/artifact/role"
	"github.com/apicurio/apicurio-cli/pkg/cmd/registry/rule"
)

func NewRootCommand(f *factory.Factory, version string) *cobra.Command {

	cmd := &cobra.Command{
		SilenceUsage:  true,
		SilenceErrors: true,
		Use:           "apicr",
		Short:         "apicurio service registry cli",
		Long:          "",
		Example:       "",
	}
	fs := cmd.PersistentFlags()
	flagutil.VerboseFlag(fs)

	// this flag comes out of the box, but has its own basic usage text, so this overrides that
	var help bool

	fs.BoolVarP(&help, "help", "h", false, "Prints help information")

	pflag.CommandLine.AddGoFlagSet(flag.CommandLine)
	cmd.AddCommand(login.NewLoginCmd(f))
	cmd.AddCommand(logout.NewLogoutCommand(f))
	cmd.AddCommand(completion.NewCompletionCommand(f))

	// Plugin command
	cmd.AddCommand(registry.NewServiceRegistryCommand(f))
	cmd.AddCommand(contextcmd.NewContextCmd(f))
	cmd.AddCommand(serviceaccount.NewServiceAccountCommand(f))
	cmd.AddCommand(request.NewCallCmd(f))

	cmd.AddCommand(artifact.NewArtifactsCommand(f))
	cmd.AddCommand(role.NewRoleCommand(f))
	cmd.AddCommand(rule.NewRuleCommand(f))
	return cmd
}
