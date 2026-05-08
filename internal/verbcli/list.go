package verbcli

import (
	"fmt"

	"github.com/spf13/cobra"
)

func newListCommand(discovered []DiscoveredVerb) *cobra.Command {
	return &cobra.Command{
		Use:   "list",
		Short: "List discovered JavaScript verbs",
		RunE: func(cmd *cobra.Command, args []string) error {
			for _, item := range discovered {
				verb := item.Verb
				if verb == nil || verb.File == nil {
					continue
				}
				_, err := fmt.Fprintf(
					cmd.OutOrStdout(),
					"%s\t%s\t%s\t%s\t%s\n",
					verb.FullPath(),
					item.Repository.Repository.Name,
					verb.File.RelPath,
					verb.FunctionName,
					verb.OutputMode,
				)
				if err != nil {
					return err
				}
			}
			return nil
		},
	}
}
