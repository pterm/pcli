package pcli

import (
	"strings"
	"time"

	"github.com/pterm/pterm"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

// generateMarkdown generates a help document written in markdown for a command.
func generateMarkdown(cmd, rootCmd *cobra.Command) (md string) {
	pterm.DisableColor()
	md += pterm.Sprintfln("# %s", cmd.CommandPath())
	md += generateUsageTemplate(cmd, rootCmd)
	md += pterm.Sprintfln("\n## Description\n\n```\n%s\n```", cmd.Long)

	if len(cmd.Commands()) > 0 {
		md += HelpSectionPrinter("Commands")
		var data [][]string
		for _, command := range cmd.Commands() {
			data = append(data, []string{command.Use + " " + strings.Join(command.Aliases, " "), command.Short})
		}
		md += "|Command|Usage|\n"
		md += "|-------|-----|\n"
		for _, d := range data {
			md += pterm.Sprintfln("|`%s`|%s|", d[0], d[1])
		}
	}

	if cmd.Flags().HasFlags() {
		md += HelpSectionPrinter("Flags")

		var flagTableData [][]string
		flagUsageLines := strings.Split(strings.TrimSpace(cmd.Flags().FlagUsages()), "\n")
		for _, line := range flagUsageLines {
			line = strings.TrimSpace(line)
			parts := strings.Split(line, "   ")
			flagString := parts[0]
			flagUsage := strings.TrimSpace(strings.Join(parts[1:], "   "))
			flagTableData = append(flagTableData, []string{flagString, flagUsage})
		}

		md += "|Flag|Usage|\n"
		md += "|----|-----|\n"
		for _, d := range flagTableData {
			md += pterm.Sprintfln("|`%s`|%s|", d[0], d[1])
		}

		md += "---\n\n"
		md += "###### Automatically generated with [PTerm](https://github.com/pterm/cli-template) on " + time.Now().Format("02 January 2006") + "\n"
	}
	pterm.EnableColor()

	return
}

// GenerateMarkdownDocs walks trough every subcommand of rootCmd and creates a documentation written in Markdown for it.
func GenerateMarkdownDocs(command, rootCmd *cobra.Command) (markdown []string) {
	markdown = append(markdown, generateMarkdown(rootCmd, rootCmd))
	for _, cmd := range command.Commands() {
		markdown = append(markdown, generateMarkdown(cmd, rootCmd))
		GenerateMarkdownDocs(cmd, rootCmd)
	}
	return
}

func generateUsageTemplate(cmd, rootCmd *cobra.Command) string {
	var ret string

	if cmd.Short != "" {
		ret += HelpSectionPrinter("Usage")
		ret += pterm.Sprintfln("%s %s", pterm.Gray(">"), pterm.Magenta(cmd.Short))
		ret += "\n"
	}

	ret += pterm.Sprintfln("%s [global options] command [options] [arguments...]", rootCmd.Use)

	return ret
}

func generateDescriptionTemplate(description string) string {
	var ret string

	if description != "" {
		ret += HelpSectionPrinter("Description")
		ret += pterm.DefaultParagraph.Sprintln(description)
	}

	return ret
}

func generateCommandsTemplate(commands []*cobra.Command) string {
	var ret string

	if len(commands) > 0 {
		ret += HelpSectionPrinter("Commands")
		data := pterm.TableData{}
		for _, command := range commands {
			data = append(data, []string{command.Use + " " + strings.Join(command.Aliases, " "), command.Short})
		}
		result, _ := pterm.DefaultTable.WithData(data).Srender()
		ret += result + "\n"
	}

	return ret
}

func generateFlagsTemplate(flags *pflag.FlagSet) string {
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
