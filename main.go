package main

import (
	"flag"
	"fmt"
	"io/fs"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

func main() {
	imgToPaaPath := flag.String("converter", `C:\Program Files (x86)\Steam\steamapps\common\DayZ Tools\Bin\ImageToPAA\ImageToPAA.exe`, "Path to the ImageToPAA executable")
	sourceRoot := flag.String("source", "./source/", "Path to the source directory")
	outputRoot := flag.String("output", "./", "Path to the output directory")
	flag.Parse()

	fileSystem := os.DirFS(*sourceRoot)

	err := fs.WalkDir(fileSystem, ".", func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			log.Fatal(err)
		}
		if strings.HasSuffix(path, ".png") {
			infile := filepath.Join(*sourceRoot, path)
			outfile := filepath.Join(*outputRoot, path)
			fmt.Println("Converting", path)
			cmd := exec.Command(*imgToPaaPath, infile, outfile)
			out, err := cmd.CombinedOutput()
			if err != nil {
				fmt.Println("error:", err)
				fmt.Println(string(out))
			}
		}
		return nil
	})

	if err != nil {
		fmt.Println("error:", err)
	}
}
