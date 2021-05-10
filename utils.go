package pcli

import "github.com/pterm/pterm"

// GenerateCopyrightString generates a pretty looking copyright string.
func GenerateCopyrightString(copyright string) string {
	var ret string

	if copyright != "" {
		ret += pterm.Gray(copyright) + "\n"
	}

	return ret
}

// GenerateVersionString generates a pretty looking version string.
func GenerateVersionString(appName, version string) string {
	var ret string

	ret += pterm.Info.Sprintfln("%s is on version: %s", appName, pterm.Magenta(version))

	return ret
}

// GenerateTitleString generates a pretty looking title string.
func GenerateTitleString(helpName string, version string) string {
	var ret string

	ret += pterm.Sprintf("\n# %s | %s\n", pterm.Cyan(helpName), pterm.Green(version))

	return ret
}
