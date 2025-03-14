package bash

import (
	"github.com/apicurio/apicurio-cli/pkg/shared/factory"
	"github.com/spf13/cobra"
)

func NewCommand(f *factory.Factory) *cobra.Command {
	cmd := &cobra.Command{
		Use:                   "bash",
		Short:                 f.Localizer.MustLocalize("completion.bash.cmd.shortDescription"),
		Long:                  f.Localizer.MustLocalize("completion.bash.cmd.longDescription"),
		Example:               f.Localizer.MustLocalize("completion.bash.cmd.example"),
		DisableFlagsInUseLine: true,
		Args:                  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			return cmd.Parent().Parent().GenBashCompletion(f.IOStreams.Out)
		},
	}

	return cmd
}
