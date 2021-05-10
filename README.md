<h1 align="center">PCLI ✨</h1>

PCLI features drop-in functions to integrate PTerm into popular CLI frameworks for Go.

## spf13/cobra

The drop-in functions must be used in your `root.go` file inside the `init` function.

|Feature|Replace With|
|-------|------------|
|Version Output|`rootCmd.SetVersionTemplate(pcli.GenerateVersionString(rootCmd.Name(), rootCmd.Version))`|
|Help Function|`rootCmd.SetHelpFunc(pcli.Spf13CobraHelpFunc(rootCmd))`|
|Error Function|`rootCmd.SetFlagErrorFunc(pcli.Spf13CobraFlagErrorFunc(rootCmd))`|

## urfave/cli

The drop-in functions must be used in your `main` function, before `app.Run(os.Args)`.

|Feature|Replace With|
|-------|------------|
|Version Output|`cli.VersionPrinter = pcli.UrfaveCliVersionPrinter`|
|Help Function|`cli.HelpPrinterCustom = pcli.UrfaveCliHelpPrinterCustom(app)`|
