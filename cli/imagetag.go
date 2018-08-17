package main

import (
	"github.com/spf13/cobra"
)

// the tagged docker image to use
var (
	tag        string
	tagFlag    = []string{"tag", "t", "ansemjo/makerelease", "docker tag to use"}
	addTagFlag = func(cmd *cobra.Command) {
		cmd.Flags().StringVarP(&tag, tagFlag[0], tagFlag[1], tagFlag[2], tagFlag[3])
	}
)
