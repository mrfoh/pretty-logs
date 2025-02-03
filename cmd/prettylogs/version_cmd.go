package prettylogs

import "github.com/spf13/cobra"

var VERSION = "1.0.1"

func NewVersionCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "version",
		Short: "Print the version number of PrettyLogs",
		Run: func(cmd *cobra.Command, args []string) {
			cmd.Printf("PrettyLogs v%s\n", VERSION)
		},
	}
}
