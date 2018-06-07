package main

import (
	"errors"
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/eternnoir/comicknife"
	"github.com/urfave/cli"
)

func main() {
	app := cli.NewApp()
	app.Name = "ComicKnife"
	app.Usage = "Split your comic"
	app.Action = run
	app.Flags = flags
	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}

func run(c *cli.Context) error {
	path := c.Args().Get(0)
	outPath := c.Args().Get(1)
	if outPath == "" {
		outPath = "./"
	}
	fmt.Printf("Start to get loader for %s\n. Output to :%s", path, outPath)

	loader, err := BuildLoader(path, outPath)
	if err != nil {
		return err
	}
	if err := loader.Process(); err != nil {
		return err
	}
	fmt.Println("All done")
	return nil
}

func BuildLoader(path, outPath string) (comicknife.Loader, error) {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return nil, err
	}

	ext := filepath.Ext(path)
	switch ext {
	case ".zip":
	case ".cbz":
		fmt.Printf("Ext is %s use zip loader.\n", ext)
		return comicknife.NewZipLoader(path, outPath, Config)
	default:
		return nil, errors.New(fmt.Sprintf("%s not support format", ext))
	}
	return nil, errors.New(fmt.Sprintf("%s not support format", ext))
}
