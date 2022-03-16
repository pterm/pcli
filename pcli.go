package pcli

import (
	"errors"
	"io/ioutil"
	"net/http"
	"runtime"
	"strings"
	"time"

	"github.com/pterm/pterm"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"github.com/tidwall/gjson"
)

var rootCmd *cobra.Command

// DisableUpdateChecking automatically checks if a new version of your application is pushed, and notifies the user.
var DisableUpdateChecking = false

// SetRootCmd sets your rootCmd.
func SetRootCmd(cmd *cobra.Command) {
	rootCmd = cmd
}

type meta struct {
	Username string
	Reponame string
}

var AppInfo = meta{}

// Setup replaces some cobra functions with pcli functions for nicer output.
func Setup() {
	rootCmd.AddCommand(GetCiCommand())
	rootCmd.SetFlagErrorFunc(FlagErrorFunc())
	rootCmd.SetGlobalNormalizationFunc(GlobalNormalizationFunc())
	rootCmd.SetHelpFunc(HelpFunc())
	rootCmd.SetHelpTemplate(HelpTemplate())
	rootCmd.SetUsageFunc(UsageFunc())
	rootCmd.SetUsageTemplate(UsageTemplate())
	rootCmd.SetVersionTemplate(VersionTemplate())
	rootCmd.SetOut(PcliOut())
	rootCmd.SetErr(Err())
}

// SetRepo must be set to your repository path (eg. pterm/cli-template).
func SetRepo(repo string) error {
	parts := strings.Split(repo, "/")
	if len(parts) != 2 {
		return errors.New("repo must be set in this pattern: username/reponame, eg.: pterm/cli-template")
	}
	AppInfo.Username = parts[0]
	AppInfo.Reponame = parts[1]

	return nil
}

func getRepoPath() string {
	return AppInfo.Username + "/" + AppInfo.Reponame
}

// CheckForUpdates checks if a new version of your application is pushed, and notifies the user, if DisableUpdateChecking is true.
func CheckForUpdates() error {
	if !DisableUpdateChecking && AppInfo.Username != "" && AppInfo.Reponame != "" {
		resp, err := http.Get(pterm.Sprintf("https://api.github.com/repos/%s/releases/latest", getRepoPath()))
		if err != nil {
			return err
		}
		defer resp.Body.Close()

		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return err
		}

		tagName := gjson.Get(string(body), "tag_name").String()

		if rootCmd.Version != tagName && tagName != "" {
			format := "A new version of %s is availble (%s)!\n"
			format += "You can install the new version with: "

			switch runtime.GOOS {
			case "windows":
				format += pterm.Magenta(pterm.Sprintf(`iwr instl.sh/%s/windows | iex`, getRepoPath()))
			case "darwin":
				format += pterm.Magenta(pterm.Sprintf(`curl -sSL instl.sh/%s/macos | bash`, getRepoPath()))
			default:
				format += pterm.Magenta(pterm.Sprintf(`curl -sSL instl.sh/%s/linux | bash`, getRepoPath()))
			}
			pterm.Info.Printfln(format, rootCmd.Name(), pterm.Magenta(tagName))
		}
	}

	return nil
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
	if cmd.CommandPath() != rootCmd.CommandPath() {
		md += pterm.Sprintfln("# ... %s", strings.TrimSpace(strings.TrimPrefix(cmd.CommandPath(), rootCmd.Use)))
		md += pterm.Sprintfln("`%s`", cmd.CommandPath())
	} else {
		md += pterm.Sprintfln("# %s", cmd.CommandPath())
	}
	md += generateUsageTemplate(cmd)

	if cmd.Long != "" {
		md += pterm.Sprintfln("\n## Description\n\n```\n%s\n```", cmd.Long)
	}

	if cmd.Example != "" {
		md += pterm.Sprintfln("## Examples\n\n```bash\n%s\n```", cmd.Example)
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

	if len(cmd.Commands()) > 0 {
		md += HelpSectionPrinter("Commands")
		var data [][]string
		for _, command := range cmd.Commands() {
			if command.Hidden {
				continue
			}
			data = append(data, []string{command.CommandPath(), command.Short})
		}
		md += "|Command|Usage|\n"
		md += "|-------|-----|\n"
		for _, d := range data {
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

// GenerateMarkdownDoc walks trough every subcommand of rootCmd and creates a documentation written in Markdown for it.
func GenerateMarkdownDoc(command *cobra.Command) (markdown MarkdownDocument) {
	if !command.Hidden {
		return MarkdownDocument{
			Name:     command.CommandPath(),
			Markdown: generateMarkdown(command),
			Command:  command,
			Filename: strings.ReplaceAll(command.CommandPath(), " ", "_"),
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

	ret += pterm.Sprintfln("%s", pterm.LightMagenta(getParentString(cmd)))

	return ret
}

func getParentString(cmd *cobra.Command) (ret string) {
	c := cmd.Parent()

	ret = cmd.Use

	for c != nil {
		ret = c.Name() + " " + ret
		c = c.Parent()
	}

	return
}

func generateDescriptionTemplate(description string) string {
	var ret string

	if description != "" {
		ret += HelpSectionPrinter("Description")
		ret += pterm.Sprintln(description)
	}

	return ret
}

func generateExamplesTemplate(cmd *cobra.Command) string {
	var ret string

	if cmd.Example != "" {
		ret += HelpSectionPrinter("Examples")
		ret += cmd.Example + "\n"
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
