package pcli

import (
	"github.com/pterm/pterm"
	"github.com/spf13/cobra"
)

// generateTitleString generates a pretty looking title string.
func generateTitleString(rootCmd *cobra.Command) string {
	return pterm.Sprintf("\n# %s | %s\n", pterm.Cyan(rootCmd.Name()), pterm.Green(rootCmd.Version))
}
