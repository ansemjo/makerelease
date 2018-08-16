package main

import (
	"os"

	"github.com/spf13/cobra"
)

// main cli command
var cmd = &cobra.Command{
	Use:   "makerelease",
	Short: "Make reproducible releases by building them in a container.",
	PreRunE: func(cmd *cobra.Command, args []string) (err error) {
		err = checkOutDirFlag(cmd)
		if err != nil {
			return
		}
		return checkInFileFlag(cmd)
	},
	RunE: func(cmd *cobra.Command, args []string) (err error) {
		return makeRelease(infile, outdir)
	},
}

func init() {
	cmd.Flags().SortFlags = false
	addOutdirFlag(cmd)
	addInfileFlag(cmd)
}

func main() {
	if err := cmd.Execute(); err != nil {
		os.Exit(1)
	}
}
