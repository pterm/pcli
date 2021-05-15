package pcli

import (
	"path/filepath"
	"runtime"
	"strings"
	"time"

	"github.com/pterm/pterm"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

var rootCmd *cobra.Command
var rootPath string

// SetRootCmd sets your rootCmd.
func SetRootCmd(cmd *cobra.Command) {
	rootCmd = cmd
	_, scriptPath, _, _ := runtime.Caller(1)
	rootPath = filepath.Join(scriptPath, "../../")
}

// GetCiCommand returns a custom crafted CI command. This must be used when using https://github.com/pterm/cli-template.
func GetCiCommand() *cobra.Command {
	return ptermCICmd
}

func generateMarkdown(cmd *cobra.Command) (md string) {
	md += generateMarkdownTree(cmd)
	md += "\n\n---\n"
	md += "> **Documentation automatically generated with [PTerm](https://github.com/pterm/cli-template) on " + time.Now().Format("02 January 2006") + "**\n"

	return
}

// generateMarkdownTree generates a help document written in markdown for a command.
func generateMarkdownTree(cmd *cobra.Command) (md string) {
	if cmd.Hidden {
		return
	}
	pterm.DisableColor()
	md += pterm.Sprintfln("# %s", cmd.CommandPath())
	md += generateUsageTemplate(cmd)
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
	}

	for _, c := range cmd.Commands() {
		md += generateMarkdownTree(c)
	}

	pterm.EnableColor()

	return
}

// MarkdownDocument contains the command and it's markdown documentation.
type MarkdownDocument struct {
	Name     string
	Markdown string
	Command  *cobra.Command
	Filename string
}

// GenerateMarkdownDocs walks trough every subcommand of rootCmd and creates a documentation written in Markdown for it.
func GenerateMarkdownDocs(command *cobra.Command) (markdown []MarkdownDocument) {
	if !command.Hidden {
		markdown = append(markdown, MarkdownDocument{
			Name:     command.Name(),
			Markdown: generateMarkdown(command),
			Command:  command,
			Filename: strings.ReplaceAll(command.Name(), " ", "_"),
		})
	}
	for _, cmd := range command.Commands() {
		if !cmd.Hidden {
			markdown = append(markdown, MarkdownDocument{
				Name:     cmd.Name(),
				Markdown: generateMarkdown(cmd),
				Command:  cmd,
				Filename: strings.ReplaceAll(cmd.Name(), " ", "_"),
			})
			GenerateMarkdownDocs(cmd)
		}
	}
	return
}

func generateUsageTemplate(cmd *cobra.Command) string {
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
