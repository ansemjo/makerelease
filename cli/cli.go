package main

import (
	"fmt"
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
	Run: func(cmd *cobra.Command, args []string) {
		err := makeRelease(infile, outdir)
		if err != nil {
			if infile != nil {
				infile.Close()
			}
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
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
