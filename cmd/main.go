package main

import (
	"errors"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"sync"

	"github.com/eternnoir/comicknife"
	"github.com/urfave/cli"
)

var pnchain chan struct{}

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
	if flagParallelNumber < 1 {
		flagParallelNumber = 2
	}
	if flagRunnerNumber < 1 {
		flagRunnerNumber = 10
	}

	pnchain = make(chan struct{}, flagParallelNumber)

	comicknife.LimitChain = make(chan struct{}, flagRunnerNumber)
	var wg sync.WaitGroup
	for _, filepath := range c.Args() {
		wg.Add(1)
		pnchain <- struct{}{}
		go func(f string) {
			defer func() {
				<-pnchain
			}()
			if err := runOneFile(&wg, f); err != nil {
				fmt.Printf("[ERROR] %s\n", err.Error())
			}
		}(filepath)
	}
	wg.Wait()
	fmt.Println("All done")
	return nil
}

func runOneFile(wg *sync.WaitGroup, path string) error {
	defer wg.Done()
	fmt.Printf("Start to get loader for %s\n. Output to :%s\n", path, flagOutputPath)

	loader, err := BuildLoader(path, flagOutputPath)
	if err != nil {
		return err
	}
	if err := loader.Process(); err != nil {
		return err
	}
	return nil
}

func BuildLoader(path, outPath string) (comicknife.Loader, error) {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return nil, err
	}

	ext := filepath.Ext(path)
	fmt.Println(ext)
	switch ext {
	case ".cbz", ".zip":
		fmt.Printf("Ext is %s use zip loader.\n", ext)
		return comicknife.NewZipLoader(path, outPath, Config)
	default:
		return nil, errors.New(fmt.Sprintf("%s not support format", ext))
	}
	return nil, errors.New(fmt.Sprintf("%s not support format", ext))
}
