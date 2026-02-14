package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"

	"github.com/peterbourgon/ff/v3"
)

func must(err error) {
	if err != nil {
		fmt.Fprintf(os.Stderr, "⛔ %v", err)
		os.Exit(1)
	}
}

func main() {
	flags := flag.NewFlagSet("mod-build", flag.ExitOnError)
	flags.Usage = func() {
		fmt.Fprintf(os.Stderr, "usage: %s [options] <source-directory>\n", filepath.Base(os.Args[0]))
		flags.PrintDefaults()
	}
	var (
		imgToPaaPath = flags.String("image-to-paa", `C:\Program Files (x86)\Steam\steamapps\common\DayZ Tools\Bin\ImageToPAA\ImageToPAA.exe`, "Path to the ImageToPAA executable")
		outputRoot   = flags.String("output", `P:\`, "Path to the output directory root (where built addons will be placed)")
		yes          = flags.Bool("yes", false, "Automatically confirm all prompts (use with caution)")
		clean        = flags.Bool("clean", false, "Clean output directory before building (deletes files which are not present in the source)")
		_            = flags.String("config", "", "config file (optional)")
	)

	err := ff.Parse(flags, os.Args[1:],
		ff.WithConfigFileFlag("config"),
		ff.WithConfigFileParser(ff.PlainParser),
	)
	if err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		fmt.Fprintln(os.Stderr, "")
		flags.Usage()
		os.Exit(1)
	}

	sourceDir := flags.Arg(0)
	if sourceDir == "" {
		fmt.Fprintln(os.Stderr, "error: source directory is required")
		fmt.Fprintln(os.Stderr, "")
		flags.Usage()
		os.Exit(1)
	}

	addonName := filepath.Base(sourceDir) // TODO: Or from $PBOPREFIX@.txt?
	outputDirectory := filepath.Join(*outputRoot, addonName)

	fmt.Println("===================================================")
	fmt.Printf("ImageToPAA Path: %s\n", *imgToPaaPath)
	fmt.Printf("    Source Path: %s\n", sourceDir)
	fmt.Printf("    Output Root: %s\n", *outputRoot)
	fmt.Printf("   Auto-confirm: %t\n", *yes)
	fmt.Printf("          Clean: %t\n", *clean)
	fmt.Println("---------------------------------------------------")
	fmt.Printf("      Addon Name: %s\n", addonName)
	fmt.Printf("Output Directory: %s\n", outputDirectory)
	fmt.Println("===================================================")

	source := NewSource(sourceDir)
	must(source.EnsureValid())

	output := NewOutput(outputDirectory)
	must(output.EnsureExists())

	confirm, err := yesOrNo(*yes, fmt.Sprintf("⚠️ The contents of %q will be removed or replaced. Continue? [y/N] ", outputDirectory))
	must(err)
	if !confirm {
		os.Exit(0)
	}

	task, err := source.Prepare()
	must(err)

	outputManifest, err := output.LoadManifest()
	must(err)

	if *clean {
		toClean, err := output.PathsToClean(task)
		must(err)

		for _, path := range toClean {
			fmt.Printf("🧹 Deleting   : %q\n", path)
			must(output.Remove(path))
		}
	}

	for _, path := range task.Copy {
		if outputManifest[path].SourceHash == task.Manifest[path].SourceHash {
			hash, err := output.Hash(path)
			must(err)
			if hash == outputManifest[path].SourceHash {
				fmt.Printf("⏭️ Unchanged  : %q\n", path)
				continue
			}
		}
		fmt.Printf("📄 Copying    : %q\n", path)
		must(output.Copy(source.RealPath(path), path))
	}

	for _, path := range task.Convert {
		if outputManifest[path].SourceHash == task.Manifest[path].SourceHash {
			hash, err := output.Hash(outputManifest[path].OutputPath)
			must(err)
			if hash == outputManifest[path].OutputHash {
				entry := task.Manifest[path]
				entry.OutputPath = outputManifest[path].OutputPath
				entry.OutputHash = outputManifest[path].OutputHash
				task.Manifest[path] = entry
				fmt.Printf("⏭️ Unchanged  : %q\n", path)
				continue
			}
		}
		fmt.Printf("🔁 Converting : %q\n", path)
		outputPath, outputHash, err := output.Convert(source.RealPath(path), path, *imgToPaaPath)
		must(err)
		entry := task.Manifest[path]
		entry.OutputPath = outputPath
		entry.OutputHash = outputHash
		task.Manifest[path] = entry
	}

	err = output.WriteManifest(task.Manifest)
	if err != nil {
		fmt.Fprintf(os.Stderr, "⚠️ Failed to write manifest file: %v\n", err)
	}

	fmt.Println("🎉 Done!")
}
