package pcli

import (
	"github.com/pterm/pterm"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

// HelpFunc is a drop in replacement for spf13/cobra `HelpFunc`
func HelpFunc(rootCmd *cobra.Command) func(*cobra.Command, []string) {
	return func(cmd *cobra.Command, strings []string) {
		var ret string

		ret += generateTitleString(rootCmd)
		ret += generateUsageTemplate(cmd, rootCmd)
		ret += generateDescriptionTemplate(cmd.Long)
		ret += generateCommandsTemplate(cmd.Commands())
		ret += generateFlagsTemplate(cmd.Flags())
		ret += "\n"

		pterm.Print(ret)
	}
}

// FlagErrorFunc is a drop in replacement for spf13/cobra `FlagErrorFunc`
func FlagErrorFunc(rootCmd *cobra.Command) func(*cobra.Command, error) error {
	return func(cmd *cobra.Command, err error) error {
		var ret string

		ret += generateTitleString(rootCmd)
		ret += generateUsageTemplate(cmd, rootCmd)
		ret += generateDescriptionTemplate(cmd.Long)
		ret += generateCommandsTemplate(cmd.Commands())
		ret += generateFlagsTemplate(cmd.Flags())
		ret += "\n\n"
		ret += pterm.Error.WithShowLineNumber(false).Sprintln(err)

		pterm.Print(ret)
		return nil
	}
}

// GlobalNormalizationFunc is a drop in replacement for spf13/cobra `GlobalNormalizationFunc`
func GlobalNormalizationFunc(rootCmd *cobra.Command) func(f *pflag.FlagSet, name string) pflag.NormalizedName {
	return rootCmd.GlobalNormalizationFunc()
}

// HelpTemplate is a drop in replacement for spf13/cobra `HelpTemplate`
func HelpTemplate(rootCmd *cobra.Command) string {
	return rootCmd.HelpTemplate()
}

// UsageFunc is a drop in replacement for spf13/cobra `UsageFunc`
func UsageFunc(rootCmd *cobra.Command) func(*cobra.Command) error {
	return rootCmd.UsageFunc()
}

// UsageTemplate is a drop in replacement for spf13/cobra `UsageTemplate`
func UsageTemplate(rootCmd *cobra.Command) string {
	return rootCmd.UsageTemplate()
}

// VersionTemplate is a drop in replacement for spf13/cobra `VersionTemplate`
func VersionTemplate(rootCmd *cobra.Command) string {
	return pterm.Info.Sprintfln("%s is on version: %s", rootCmd.Name(), pterm.Magenta(rootCmd.Version))
}
