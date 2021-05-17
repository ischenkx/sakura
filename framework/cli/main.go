package main

import (
	"flag"
	"fmt"
	"github.com/ischenkx/notify/framework/generator"
	"github.com/ischenkx/notify/framework/parser"
	"github.com/fatih/color"
	"os"
	"path/filepath"
	"time"
)

var (
	outputPath = flag.String("out", "./cmd/main.go", "output file path")
	inputPath = flag.String("src", "./", "source folder for code generation")
	buildFlag = flag.Bool("build", false, "building source code")
)

func main() {
	flag.Parse()

	if *buildFlag {
		color.Blue("Building...")

		var absErr error

		*inputPath, absErr = filepath.Abs(*inputPath)

		if absErr != nil {
			color.Red("Failed to get absolute path (%s): \n\t%s", *inputPath, absErr)
			return
		}

		*outputPath, absErr = filepath.Abs(*outputPath)

		if absErr != nil {
			color.Red("Failed to get absolute path (%s): \n\t%s", *inputPath, absErr)
		}

		now := time.Now()
		info, err := (&parser.Parser{Folder: *inputPath}).CollectInfo()
		color.Blue("Info collection time: %s", time.Since(now))

		now = time.Now()
		data, e := generator.Generate(info)
		color.Blue("Generation time: %s", time.Since(now))

		if e != nil {
			err.Error(e.Error())
		}


		if len(err.Errors()) == 0 {
			color.Green("Errors: 0")
		} else {
			color.Red("Errors: %d", len(err.Errors()))
			for _, s := range err.Errors() {
				color.Blue(s)
			}
		}

		if len(err.Warnings()) == 0 {
			color.Green("Warnings: 0")
		} else {
			color.Yellow("Warnings: %d", len(err.Warnings()))
			for _, s := range err.Warnings() {
				color.Blue(s)
			}
		}

		if err.Ok() {
			color.Green("Successfully generated files")
			if e := os.WriteFile(*outputPath, []byte(data), 0777); e != nil {
				color.Red("Failed to write files: %s", e)
			} else {
				color.Green("Successfully wrote files")
			}
		} else {
			color.Red("Failed to generate source files")
		}
	} else {
		fmt.Println("Notify is a framework for building real-time applications")
	}
}