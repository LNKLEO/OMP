package cli

import (
	"fmt"
	"strings"

	"github.com/LNKLEO/OMP/font"
	"github.com/LNKLEO/OMP/runtime"
	"github.com/LNKLEO/OMP/terminal"

	"github.com/spf13/cobra"
)

var (
	zipFolder string

	fontCmd = &cobra.Command{
		Use:   "font [install|configure]",
		Short: "Manage fonts",
		Long: `Manage fonts.

This command is used to install fonts and configure the font in your terminal.

  - install: OMP font install 3270`,
		ValidArgs: []string{
			"install",
			"configure",
		},
		Run: func(cmd *cobra.Command, args []string) {
			if len(args) == 0 {
				_ = cmd.Help()
				return
			}
			switch args[0] {
			case "install":
				var fontName string
				if len(args) > 1 {
					fontName = args[1]
				}

				flags := &runtime.Flags{
					SaveCache: true,
				}

				env := &runtime.Terminal{}
				env.Init(flags)
				defer env.Close()

				terminal.Init(env.Shell())

				if !strings.HasPrefix(zipFolder, "/") {
					zipFolder += "/"
				}

				font.Run(fontName, env.Cache(), env.Root(), zipFolder)

				return
			case "configure":
				fmt.Println("not implemented")
			default:
				_ = cmd.Help()
			}
		},
	}
)

func init() {
	fontCmd.Flags().StringVar(&zipFolder, "zip-folder", "", "the folder inside the zip file to install fonts from")
	RootCmd.AddCommand(fontCmd)
}
