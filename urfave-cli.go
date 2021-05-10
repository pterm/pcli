package pcli

import (
	"io"
	"reflect"
	"strings"

	"github.com/pterm/pterm"
	"github.com/urfave/cli/v2"
)

// UrfaveCliVersionPrinter is a drop in replacement for urfave/cli `VersionPrinter` function.
//
// Usage:
//     cli.VersionPrinter = pcli.UrfaveCliVersionPrinter
func UrfaveCliVersionPrinter(c *cli.Context) {
	pterm.Print(GenerateVersionString(c.App.Name, c.App.Version))
}

// UrfaveCliHelpPrinterCustom is a drop in replacement for urfave/cli `HelpPrinterCustom` function.
//
// Usage:
//     cli.HelpPrinterCustom = pcli.UrfaveCliHelpPrinterCustom(app)
func UrfaveCliHelpPrinterCustom(app *cli.App) func(w io.Writer, templ string, data interface{}, customFunc map[string]interface{}) {
	return func(w io.Writer, templ string, data interface{}, customFunc map[string]interface{}) {
		if command, ok := data.(*cli.Command); ok {
			_, _ = w.Write([]byte(urfaveCliGenerateCommandHelpTemplate(command, app)))
		} else if a, ok := data.(*cli.App); ok {
			if reflect.DeepEqual(a, app) {
				_, _ = w.Write([]byte(urfaveCliGenerateAppHelpTemplate(app, app)))
			} else {
				_, _ = w.Write([]byte(urfaveCliGenerateSubcommandHelpTemplate(a, app)))
			}
		}
	}
}

func urfaveCliGenerateAppHelpTemplate(app *cli.App, mainApp *cli.App) string {
	var ret string

	ret += GenerateTitleString(mainApp.HelpName, mainApp.Version)
	ret += urvafeCliGenerateAuthorsTemplate(app.Authors)
	ret += urvafeCliGenerateUsageTemplate(app.Usage, app.UsageText, mainApp)
	ret += urvafeCliGenerateDescriptionTemplate(app.Description)
	ret += urvafeCliGenerateCommandsTemplate(app.VisibleCommands())
	ret += urvafeCliGenerateFlagsTemplate(app.VisibleFlags())
	ret += GenerateCopyrightString(app.Copyright)

	return ret
}

func urfaveCliGenerateCommandHelpTemplate(command *cli.Command, mainApp *cli.App) string {
	var ret string

	ret += GenerateTitleString(command.HelpName, mainApp.Version)
	ret += urvafeCliGenerateUsageTemplate(command.Usage, command.UsageText, mainApp)
	ret += urvafeCliGenerateDescriptionTemplate(command.Description)
	ret += urvafeCliGenerateFlagsTemplate(command.VisibleFlags())
	ret += GenerateCopyrightString(mainApp.Copyright)

	return ret
}

func urfaveCliGenerateSubcommandHelpTemplate(cmd *cli.App, mainApp *cli.App) string {
	var ret string

	ret += GenerateTitleString(cmd.HelpName, mainApp.Version)
	ret += urvafeCliGenerateUsageTemplate(cmd.Usage, cmd.UsageText, mainApp)
	ret += urvafeCliGenerateDescriptionTemplate(cmd.Description)
	ret += urvafeCliGenerateAuthorsTemplate(cmd.Authors)
	ret += urvafeCliGenerateCommandsTemplate(cmd.VisibleCommands())
	ret += urvafeCliGenerateFlagsTemplate(cmd.VisibleFlags())
	ret += GenerateCopyrightString(mainApp.Copyright)

	return ret
}

func urvafeCliGenerateUsageTemplate(usage, usageText string, mainApp *cli.App) string {
	var ret string

	if usage != "" {
		ret += HelpSectionPrinter("Usage")
		ret += pterm.Sprintfln("%s %s", pterm.Gray(">"), pterm.Magenta(usage))
		ret += "\n"
	}
	if usageText != "" {
		ret += usageText
	} else {
		ret += pterm.Sprintfln("%s [global options] command [options] [arguments...]", mainApp.Name)
	}

	return ret
}

func urvafeCliGenerateDescriptionTemplate(description string) string {
	var ret string

	if description != "" {
		ret += HelpSectionPrinter("Description")
		ret += pterm.DefaultParagraph.Sprintln(description)
	}

	return ret
}

func urvafeCliGenerateAuthorsTemplate(authors []*cli.Author) string {
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

func urvafeCliGenerateCommandsTemplate(commands []*cli.Command) string {
	var ret string

	if len(commands) > 0 {
		ret += HelpSectionPrinter("Commands")
		data := pterm.TableData{}
		for _, command := range commands {
			data = append(data, []string{command.Name + " " + strings.Join(command.Aliases, " "), command.Usage})
		}
		result, _ := pterm.DefaultTable.WithData(data).Srender()
		ret += result + "\n"
	}

	return ret
}

func urvafeCliGenerateFlagsTemplate(flags []cli.Flag) string {
	var ret string

	if len(flags) > 0 {
		ret += HelpSectionPrinter("Flags")
		data := pterm.TableData{}
		for _, flag := range flags {
			parts := strings.Split(flag.String(), "\t")
			flagNames := flag.Names()
			for i, name := range flagNames {
				if len(name) > 1 {
					flagNames[i] = "--" + name
				} else {
					flagNames[i] = "-" + name
				}
			}
			names := strings.Join(flagNames, ", ")
			data = append(data, []string{names, parts[1]})
		}
		result, _ := pterm.DefaultTable.WithData(data).Srender()
		ret += result + "\n\n"
	}

	return ret
}
