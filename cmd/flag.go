package main

import (
	"github.com/eternnoir/comicknife"
	"github.com/urfave/cli"
)

var Config *comicknife.ImageConfig = comicknife.NewDefaultConfig()
var (
	flagOutputPath     = ""
	flagRunnerNumber   = 10
	flagParallelNumber = 2
)

var flags []cli.Flag = []cli.Flag{
	cli.BoolFlag{
		Name:        "fc",
		Usage:       "Force crop event image height > width",
		Destination: &Config.FoceCrop,
	},
	cli.StringFlag{
		Name:        "d",
		Value:       "RL",
		Usage:       "Direction. eg. \"RL\" or \"LR\"",
		Destination: &Config.Direction,
	},

	cli.StringFlag{
		Name:        "o",
		Value:       "./split",
		Usage:       "Output folder",
		Destination: &flagOutputPath,
	},
	cli.StringFlag{
		Name:        "f",
		Value:       "",
		Usage:       "Output image format. eg. jpg, png",
		Destination: &Config.OutputFormat,
	},
	cli.IntFlag{
		Name:        "pc",
		Value:       0,
		Usage:       "PNGCompressionLevel. 0, -1, -2, -3",
		Destination: &Config.PNGCompressionLevel,
	},
	cli.IntFlag{
		Name:        "jq",
		Value:       100,
		Usage:       "JPEGQuality.",
		Destination: &Config.JPEGQuality,
	},
	cli.IntFlag{
		Name:        "rn",
		Value:       10,
		Usage:       "Runner number.",
		Destination: &flagRunnerNumber,
	},
	cli.IntFlag{
		Name:        "pn",
		Value:       2,
		Usage:       "Parallel number.",
		Destination: &flagParallelNumber,
	},
}
