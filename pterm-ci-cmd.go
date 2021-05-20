package pcli

import (
	"io/fs"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"text/template"
	"time"

	"github.com/pterm/pterm"
	"github.com/spf13/cobra"
)

// ptermCICmd represents the ptermCi command
// ! Do not delete this file. It it used inside the CI system.
var ptermCICmd = &cobra.Command{
	Use:   "ptermCI",
	Short: "Run internal CI-System to update documentation.",
	Long: `This command is used in the CI-System to generate new documentation of the CLI tool.
It should not be used outside the development of this tool.`,
	Run: func(cmd *cobra.Command, args []string) {
		pterm.Info.Printfln("Running PtermCI for %s", rootCmd.Name())
		started := time.Now()
		originURL := detectOriginURL()

		if !strings.Contains(originURL, "/cli-template") {
			if _, err := os.Stat(getPathTo("./setup/main.go")); err == nil {
				pterm.DefaultSection.Println("Running template setup")
				cmd := exec.Command("go", "run", getPathTo("./setup/main.go"))
				cmd.Stdout = os.Stdout
				cmd.Stderr = os.Stderr
				pterm.Fatal.PrintOnError(cmd.Run())
				pterm.DefaultSection.Println("Deleting setup script")
				pterm.Fatal.PrintOnError(os.RemoveAll(getPathTo("./setup")))
			}
		}

		pterm.DefaultSection.Println("Generating markdown documentation")

		markdownDoc := GenerateMarkdownDoc(rootCmd)
		pterm.Fatal.PrintOnError(ioutil.WriteFile(getPathTo("/docs/docs.md"), []byte(markdownDoc.Markdown), 0777))

		project := struct {
			ProjectPath string
			Name        string
			RepoName    string
			UserName    string
			URL         string
			Short       string
			Long        string

			GitHubPagesURL string

			InstallCommandWindows string
			InstallCommandLinux   string
			InstallCommandMacOS   string
		}{}

		projectParts := strings.Split(strings.TrimPrefix(originURL, "https://github.com/"), "/")

		project.UserName = projectParts[0]
		project.RepoName = projectParts[1]
		project.ProjectPath = pterm.Sprintf("%s/%s", project.UserName, project.RepoName)
		project.Name = rootCmd.Name()
		project.URL = pterm.Sprintf("https://github.com/%s", project.ProjectPath)
		project.Short = rootCmd.Short
		project.Long = rootCmd.Long

		project.InstallCommandWindows = pterm.Sprintf(`iwr -useb instl.sh/%s/windows | iex`, project.ProjectPath)
		project.InstallCommandLinux = pterm.Sprintf(`curl -fsSL instl.sh/%s/linux | bash`, project.ProjectPath)
		project.InstallCommandMacOS = pterm.Sprintf(`curl -fsSL instl.sh/%s/macos | bash`, project.ProjectPath)

		project.GitHubPagesURL = pterm.Sprintf("https://%s.github.io/%s", project.UserName, project.RepoName)

		pterm.DefaultSection.Println("Processing '*.template.[md|html|js|css]' files")

		walkOverExt("", ".template.md,.template.html,.template.js,.template.css", func(path string) {
			contentBytes, err := ioutil.ReadFile(path)
			content := string(contentBytes)
			tmpl, err := template.New(filepath.Base(path)).Parse(content)
			pterm.Fatal.PrintOnError(err)
			file, err := os.OpenFile(strings.ReplaceAll(path, ".template", ""), os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0777)
			pterm.Fatal.PrintOnError(err)
			pterm.Fatal.PrintOnError(tmpl.Execute(file, project))
			pterm.Fatal.PrintOnError(file.Close())
		})

		pterm.DefaultSection.Println("Copying README.md to docs/REAMDE.md")

		input, err := ioutil.ReadFile(getPathTo("/README.md"))
		pterm.Fatal.PrintOnError(err)
		pterm.Fatal.PrintOnError(ioutil.WriteFile(getPathTo("/docs/README.md"), input, 0777))

		pterm.Success.Printfln("The PTerm-CI System took %v to complete.", time.Since(started))
	},
	Hidden: true,
}

func detectOriginURL() (url string) {
	out, err := exec.Command("git", "remote", "-v").Output()
	pterm.Fatal.PrintOnError(err)
	pterm.Debug.Printfln("Git output:\n%s", string(out))

	output := string(out)

	for _, s := range strings.Split(output, "\n") {
		s = strings.TrimSpace(strings.TrimLeft(s, "origin"))
		if strings.HasPrefix(s, "https://github.com/") && strings.Contains(s, "push") {
			pterm.Debug.Printfln("Detected GitHub Repo: %s", s)
			url = strings.TrimSpace(strings.TrimRight(s, "(push)"))

			return
		}
	}

	return
}

func walkOverExt(path, exts string, f func(path string)) {
	_ = filepath.Walk(getPathTo(path), func(path string, info fs.FileInfo, err error) error {
		for _, ext := range strings.Split(exts, ",") {
			if strings.HasSuffix(path, ext) {
				f(path)
			}
		}
		return nil
	})
}

func getPathTo(file string) string {
	// dir, _ := os.Getwd()
	return filepath.Join("./", file)
}
