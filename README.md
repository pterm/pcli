<h1 align="center">PCLI ✨</h1>

PCLI features drop-in functions to integrate PTerm into popular CLI frameworks for Go.

## spf13/cobra

|Feature|Replace With|
|-------|------------|
|Version Output|`rootCmd.SetVersionTemplate(pcli.GenerateVersionString(rootCmd.Name(), rootCmd.Version))`|
|Help Function|`rootCmd.SetHelpFunc(pcli.Spf13CobraHelpFunc(rootCmd))`|
|Error Function|`rootCmd.SetFlagErrorFunc(pcli.Spf13CobraFlagErrorFunc(rootCmd))`|

## urfave/cli

|Feature|Replace With|
|-------|------------|
|Version Output|`cli.VersionPrinter = pcli.UrfaveCliVersionPrinter`|
|Help Function|`cli.HelpPrinterCustom = pcli.UrfaveCliHelpPrinterCustom(app)`|
