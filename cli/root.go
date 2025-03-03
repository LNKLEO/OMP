package cli

import (
	"fmt"
	"os"
	"strings"

	"github.com/LNKLEO/OMP/build"
	"github.com/spf13/cobra"
)

var (
	configFlag   string
	shellName    string
	printVersion bool

	// for internal use only
	silent bool

	// deprecated
	initialize bool
)

var RootCmd = &cobra.Command{
	Use:   "OMP",
	Short: "OMP is a tool to render my prompt",
	Long: "OMP is a tool to render my prompt",
	Run: func(cmd *cobra.Command, _ []string) {
		if initialize {
			runInit(strings.ToLower(shellName))
			return
		}
		if printVersion {
			fmt.Println(build.Version)
			return
		}

		_ = cmd.Help()
	},
}

func Execute() {
	if err := RootCmd.Execute(); err != nil {
		// software error
		os.Exit(70)
	}
}

func init() {
	RootCmd.PersistentFlags().StringVarP(&configFlag, "config", "c", "", "config file path")
	RootCmd.PersistentFlags().BoolVar(&silent, "silent", false, "do not print anything")
	RootCmd.Flags().BoolVar(&printVersion, "version", false, "print the version number and exit")

	// Deprecated flags, should be kept to avoid breaking CLI integration.
	RootCmd.Flags().BoolVarP(&initialize, "init", "i", false, "init")
	RootCmd.Flags().StringVarP(&shellName, "shell", "s", "", "shell")

	// Hide flags that are deprecated or for internal use only.
	_ = RootCmd.PersistentFlags().MarkHidden("silent")

	// Disable completions
	RootCmd.CompletionOptions.DisableDefaultCmd = true
}
