package cli

import (
	"fmt"

	"github.com/LNKLEO/OMP/cache"

	"github.com/spf13/cobra"
)

// getCmd represents the get command
var getCache = &cobra.Command{
	Use:   "cache [path|clear|edit]",
	Short: "Interact with the OMP cache",
	Long: `Interact with the OMP cache.

You can do the following:

- path: list cache path
- clear: remove all cache values
- edit: edit cache values`,
	ValidArgs: []string{
		"path",
		"clear",
	},
	Args: NoArgsOrOneValidArg,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 0 {
			_ = cmd.Help()
			return
		}

		switch args[0] {
		case "path":
			fmt.Println(cache.Path())
		case "clear":
			deletedFiles, err := cache.Clear(cache.Path(), true)
			if err != nil {
				fmt.Println(err)
				return
			}

			for _, file := range deletedFiles {
				fmt.Println("removed cache file:", file)
			}
		}
	},
}

func init() {
	RootCmd.AddCommand(getCache)
}
