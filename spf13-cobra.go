package pcli

import (
	"strings"

	"github.com/pterm/pterm"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"github.com/urfave/cli/v2"
)

// Spf13CobraHelpFunc is a drop in replacement for spf13/cobra `HelpFunc`
//
// Usage:
//     rootCmd.SetHelpFunc(pcli.Spf13CobraHelpFunc(rootCmd))
func Spf13CobraHelpFunc(rootCmd *cobra.Command) func(*cobra.Command, []string) {
	return func(cmd *cobra.Command, strings []string) {
		var ret string

		ret += GenerateTitleString(rootCmd.Name(), rootCmd.Version)
		ret += spf13CobraGenerateUsageTemplate(cmd, rootCmd)
		ret += spf13CobraGenerateDescriptionTemplate(cmd.Long)
		ret += spf13CobraGenerateCommandsTemplate(cmd.Commands())
		ret += spf13CobraGenerateFlagsTemplate(cmd.Flags())
		ret += "\n"

		pterm.Print(ret)
	}
}

// Spf13CobraFlagErrorFunc is a drop in replacement for spf13/cobra `HelpFunc`
//
// Usage:
//     rootCmd.SetFlagErrorFunc(pcli.Spf13CobraFlagErrorFunc(rootCmd))
func Spf13CobraFlagErrorFunc(rootCmd *cobra.Command) func(*cobra.Command, error) error {
	return func(cmd *cobra.Command, err error) error {
		var ret string

		ret += GenerateTitleString(rootCmd.Name(), rootCmd.Version)
		ret += spf13CobraGenerateUsageTemplate(cmd, rootCmd)
		ret += spf13CobraGenerateDescriptionTemplate(cmd.Long)
		ret += spf13CobraGenerateCommandsTemplate(cmd.Commands())
		ret += spf13CobraGenerateFlagsTemplate(cmd.Flags())
		ret += "\n\n"
		ret += pterm.Error.WithShowLineNumber(false).Sprintln(err)

		pterm.Print(ret)
		return nil
	}
}

func spf13CobraGenerateUsageTemplate(cmd, rootCmd *cobra.Command) string {
	var ret string

	if cmd.Short != "" {
		ret += HelpSectionPrinter("Usage")
		ret += pterm.Sprintfln("%s %s", pterm.Gray(">"), pterm.Magenta(cmd.Short))
		ret += "\n"
	}

	ret += pterm.Sprintfln("%s [global options] command [options] [arguments...]", rootCmd.Use)

	return ret
}

func spf13CobraGenerateDescriptionTemplate(description string) string {
	var ret string

	if description != "" {
		ret += HelpSectionPrinter("Description")
		ret += pterm.DefaultParagraph.Sprintln(description)
	}

	return ret
}

func spf13CobraGenerateAuthorsTemplate(authors []*cli.Author) string {
	var ret string

	if len(authors) > 0 {
		ret += HelpSectionPrinter("Authors")
		data := pterm.TableData{}
		for _, author := range authors {
			data = append(data, []string{author.Name, pterm.Gray(author.Email)})
		}
		result, _ := pterm.DefaultTable.WithData(data).Srender()
		ret += result + "\n"
	}

	return ret
}

func spf13CobraGenerateCommandsTemplate(commands []*cobra.Command) string {
	var ret string

	if len(commands) > 0 {
		ret += HelpSectionPrinter("Commands")
		data := pterm.TableData{}
		for _, command := range commands {
			if command.Hidden {
				continue
			}
			data = append(data, []string{command.Use + " " + strings.Join(command.Aliases, " "), command.Short})
		}
		result, _ := pterm.DefaultTable.WithData(data).Srender()
		ret += result + "\n"
	}

	return ret
}

func spf13CobraGenerateFlagsTemplate(flags *pflag.FlagSet) string {
	if !flags.HasFlags() {
		return ""
	}

	var ret string
	ret += HelpSectionPrinter("Flags")

	flagTableData := pterm.TableData{}
	flagUsageLines := strings.Split(strings.TrimSpace(flags.FlagUsages()), "\n")
	for _, line := range flagUsageLines {
		line = strings.TrimSpace(line)
		parts := strings.Split(line, "   ")
		flagString := parts[0]
		flagUsage := strings.TrimSpace(strings.Join(parts[1:], "   "))
		flagTableData = append(flagTableData, []string{flagString, flagUsage})
	}

	table, _ := pterm.DefaultTable.WithData(flagTableData).Srender()
	ret += table

	return ret
}
