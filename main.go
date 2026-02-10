package main

import (
	"flag"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"slices"

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
	var (
		imgToPaaPath = flags.String("image-to-paa", `C:\Program Files (x86)\Steam\steamapps\common\DayZ Tools\Bin\ImageToPAA\ImageToPAA.exe`, "Path to the ImageToPAA executable")
		sourceRoot   = flags.String("source", "./source/", "Path to the source directory")
		outputRoot   = flags.String("output", "./build/", "Path to the output directory")
		yes          = flags.Bool("yes", false, "Automatically confirm all prompts (use with caution)")
		_            = flags.String("config", "", "config file (optional)")
	)

	err := ff.Parse(flags, os.Args[1:],
		ff.WithConfigFileFlag("config"),
		ff.WithConfigFileParser(ff.PlainParser),
	)
	if err != nil {
		flags.Usage()
		os.Exit(1)
	}

	fmt.Println("===================================================")
	fmt.Printf("ImageToPAA Path: %s\n", *imgToPaaPath)
	fmt.Printf("    Source Path: %s\n", *sourceRoot)
	fmt.Printf("    Output Path: %s\n", *outputRoot)
	fmt.Printf("   Auto-confirm: %t\n", *yes)
	fmt.Println("===================================================")

	// ensure our build output directory exists
	exists, err := BuildRootExists(*outputRoot)
	must(err)
	if !exists {
		fmt.Printf("Creating build output directory %q\n", *outputRoot)
		must(os.MkdirAll(*outputRoot, 0755))
	}

	addons, err := GetAddonsToBuild(*sourceRoot)
	must(err)

	toClean, err := GetAddonOutputDirsToClean(*outputRoot, addons)
	must(err)

	skip := []string{}

	for _, addon := range toClean {
		confirm, err := yesOrNo(*yes, fmt.Sprintf("⚠️ The build directory %q will be removed and recreated. Continue? [y/N] ", filepath.Join(*outputRoot, addon)))
		must(err)
		if confirm {
			must(os.RemoveAll(filepath.Join(*outputRoot, addon)))
		} else {
			skip = append(skip, addon)
		}
	}

	if len(addons)-len(skip) == 0 {
		fmt.Println("Nothing to build!")
		os.Exit(0)
	}

	for _, name := range addons {
		if slices.Contains(skip, name) {
			fmt.Printf("⏭️ Skipping: %q\n", name)
			continue
		}

		addonOutputPath := filepath.Join(*outputRoot, name)
		addonInputPath := filepath.Join(*sourceRoot, name)

		fmt.Printf("🛠️ Building: %s\n", name)
		fmt.Printf("   📂 Creating build output directory %q\n", addonOutputPath)
		must(os.MkdirAll(addonOutputPath, 0755))

		err = fs.WalkDir(os.DirFS(addonInputPath), ".", func(path string, d fs.DirEntry, err error) error {
			if err != nil {
				fmt.Fprintln(os.Stderr, err.Error())
				os.Exit(1)
			}

			infile := filepath.Join(addonInputPath, path)

			if shouldCopy(path) {
				fmt.Printf("   📄 Copying    : %q\n", path)
				must(copyFileWithPath(infile, filepath.Join(addonOutputPath, path)))
			} else if shouldConvert(path) {
				fmt.Printf("   🔁 Converting : %q\n", path)
				must(convertWithPath(
					infile,
					filepath.Join(addonOutputPath, swapExtension(path, ".paa")),
					*imgToPaaPath,
				))
			}

			return nil
		})
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error building addon %q: %s\n", name, err.Error())
		}
	}

	fmt.Println("🎉 Done!")
}
