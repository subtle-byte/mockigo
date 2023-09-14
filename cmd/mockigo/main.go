package main

import (
	"flag"
	"log"
	"os"
	"strings"

	"github.com/subtle-byte/mockigo/internal/generator"
	"github.com/subtle-byte/mockigo/internal/util"
)

func main() {
	var targets, outFile string
	var testPkg, gogen bool
	flag.StringVar(&targets, "targets", "", "Comma separated list of interfaces/function types to generate mocks for, all if empty")
	flag.BoolVar(&testPkg, "test-pkg", false, "Should the generated mocks to be placed in the package with _test suffix")
	flag.StringVar(&outFile, "out-file", "mocks_test.go", "Output file")
	flag.BoolVar(&gogen, "gogen", true, "Generate go:generate in the output file to ease the regeneration")
	flag.Parse()

	targetsSplitted := strings.Split(targets, ",")
	if len(targetsSplitted) == 1 && targetsSplitted[0] == "" {
		targetsSplitted = nil
	}
	gogenCmd := ""
	if gogen {
		sameDir, err := util.FileInDir(".", outFile)
		if err != nil {
			log.Fatalln("ERROR:", err)
		}
		if sameDir {
			gogenCmd = "mockigo " + strings.Join(os.Args[1:], " ")
		} else {
			log.Println("WARNING: go:generate is not added to output because generating to the different directory")
		}
	}
	err := generator.Generate(generator.Config{
		TargetPkgDirPath: ".",
		Targets: generator.Targets{
			Include:    len(targetsSplitted) == 0,
			Exceptions: util.SliceToSet(targetsSplitted),
		},
		OutPkgName: func(inspectedPkgName string) string {
			if testPkg {
				inspectedPkgName += "_test"
			}
			return inspectedPkgName
		},
		OutFilePath: outFile,
		OutPublic:   true,
		GoGenCmd:    gogenCmd,
	})
	if err != nil {
		log.Fatalln("ERROR:", err)
	}
}
