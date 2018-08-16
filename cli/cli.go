package main

import "github.com/spf13/cobra"

var cmd = &cobra.Command{
	Use:   "makerelease",
	Short: "make reproducible releases",
	PreRunE: func(cmd *cobra.Command, args []string) (err error) {
		return checkFileArgFlags(cmd)
	},
	Run: func(cmd *cobra.Command, args []string) {
		makeRelease(infile, outdirArg)
	},
}

func init() {
	cmd.Flags().SortFlags = false
	addFileArgFlags(cmd)
}
