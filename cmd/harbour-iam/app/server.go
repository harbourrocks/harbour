package server

import (
	// "github.com/harbourrocks/harbour/pkg/harbour-iam/options"

	"github.com/spf13/cobra"
)

// NewIAMServerCommmand creates a *cobra.Command object with default parameters
func NewIAMServerCommmand() *cobra.Command {
	// s := options.NewDefaultOptions()
	cmd := &cobra.Command{
		Use: "harbour-iam",
		Long: `The harbour.rocks IAM server manages
authentication and authorization for the harbour environment.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			// s.OIDC_URL = ""
			return nil
		},
	}

	return cmd
}
