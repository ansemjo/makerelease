package main

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"
)

// from go/build/syslist.go, 2018-08-18
const goosList = "android darwin dragonfly freebsd js linux nacl netbsd openbsd plan9 solaris windows zos"
const goarchList = "386 amd64 amd64p32 arm armbe arm64 arm64be ppc64 ppc64le mips mipsle mips64 mips64le mips64p32 mips64p32le ppc riscv riscv64 s390 s390x sparc sparc64 wasm"

// alternative build target list
var (
	targets        []string
	targetsFlag    = []string{"targets", "T", "alternative build target(s)"}
	addTargetsFlag = func(cmd *cobra.Command) {
		cmd.Flags().StringArrayVarP(&targets, targetsFlag[0], targetsFlag[1], []string{}, targetsFlag[2])
	}
	checkTargetFlag = func(cmd *cobra.Command) (err error) {
		if cmd.Flag(targetsFlag[0]).Changed {

			// convert consts to slices
			osList := strings.Split(goosList, " ")
			archList := strings.Split(goarchList, " ")

			for _, target := range targets {
				t := strings.Split(target, "/")
				if len(t) != 2 {
					return fmt.Errorf("could not parse target: %s", target)
				}
				if !contains(osList, t[0]) {
					return fmt.Errorf("invalid os: %s", target)
				}
				if !contains(archList, t[1]) {
					return fmt.Errorf("invalid arch: %s", target)
				}
			}

		}
		return
	}
)

// simple check wether a slice contains a string
func contains(slice []string, str string) bool {
	for _, s := range slice {
		if s == str {
			return true
		}
	}
	return false
}
