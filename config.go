package comicknife

import (
	"image/png"

	"github.com/disintegration/imaging"
)

type ImageConfig struct {
	Direction           string
	Quality             int
	Rotate              int
	OutputFormat        string
	PNGCompressionLevel int
	JPEGQuality         int
	FoceCrop            bool
}

func NewDefaultConfig() *ImageConfig {
	return &ImageConfig{
		Direction: "RL",
	}
}

func (c *ImageConfig) GetLRImageName() (string, string) {
	switch c.Direction {
	case "RL":
		return "_1", "_0"
	case "LR":
		return "_0", "_1"
	default:
		panic("not support direction")
	}
}

func (c *ImageConfig) ImgOpts() []imaging.EncodeOption {
	opts := make([]imaging.EncodeOption, 0)
	opts = append(opts,
		imaging.PNGCompressionLevel(png.CompressionLevel(c.PNGCompressionLevel)),
		imaging.JPEGQuality(c.JPEGQuality))
	return opts
}
