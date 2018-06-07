package comicknife

import (
	"archive/zip"
	"compress/flate"
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/disintegration/imaging"
)

type ZipLoader struct {
	ImageConfig *ImageConfig
	FilePath    string
	OutputPath  string
	Ext         string
}

func NewZipLoader(path, outPath string, cfg *ImageConfig) (*ZipLoader, error) {
	ext := filepath.Ext(path)
	return &ZipLoader{
		ImageConfig: cfg,
		FilePath:    path,
		Ext:         ext,
	}, nil
}

func (z *ZipLoader) Process() error {
	imgs, err := loadImagesFromZipFile(z.FilePath, z.ImageConfig)
	if err != nil {
		return err
	}

	fmt.Printf("Load %d images\n", len(imgs))
	resultImgs, err := BatchCut(imgs, z.ImageConfig)
	if err != nil {
		return err
	}

	filename := filepath.Base(z.FilePath)
	return imagesToZip(resultImgs, filepath.Join(z.OutputPath, filename))
}

func imagesToZip(imgs []Image, path string) error {
	d, err := os.Create(path)
	if err != nil {
		return err
	}
	defer d.Close()
	w := zip.NewWriter(d)

	w.RegisterCompressor(zip.Deflate, func(out io.Writer) (io.WriteCloser, error) {
		return flate.NewWriter(out, flate.BestCompression)
	})
	defer w.Close()

	for _, img := range imgs {
		zipFile, err := w.Create(img.OutputFileName())
		if err != nil {
			return err
		}
		ba, err := img.Bytes()
		if err != nil {
			return err
		}
		_, err = zipFile.Write(ba)
		if err != nil {
			return err
		}
	}
	return nil
}

func loadImagesFromZipFile(path string, cfg *ImageConfig) ([]Image, error) {
	ret := make([]Image, 0)
	reader, err := zip.OpenReader(path)
	if err != nil {
		return nil, err
	}
	for _, file := range reader.File {
		if file.FileInfo().IsDir() {
			fmt.Printf("DIR: %s\n", file.Name)
			continue
		}
		rc, err := file.Open()
		if err != nil {
			return nil, err
		}
		defer rc.Close()
		filename := file.Name
		fmt.Printf("Load image %s\n", filename)
		src, err := imaging.Decode(rc)
		if err != nil {
			return nil, err
		}
		decodeFormatName := filename
		if cfg.OutputFormat != "" {
			decodeFormatName = "abcd." + cfg.OutputFormat
		}
		outputFormat, err := imaging.FormatFromFilename(decodeFormatName)
		if err != nil {
			return nil, err
		}
		ret = append(ret, Image{Image: src, FileName: filename, OutFormat: outputFormat, cfg: cfg})
	}

	return ret, nil
}
