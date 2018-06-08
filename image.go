package comicknife

import (
	"bufio"
	"bytes"
	"fmt"
	"image"
	"path/filepath"
	"strings"
	"sync"

	"github.com/disintegration/imaging"
)

var LimitChain = make(chan struct{}, 10)

type Image struct {
	image.Image
	imgConfig image.Config
	FileName  string
	OutFormat imaging.Format
	cfg       *ImageConfig
}

func (i *Image) Cut() ([]Image, error) {
	ic, err := i.GetConfig()
	if err != nil {
		return nil, err
	}
	w := ic.Width
	h := ic.Height
	if (ic.Height > ic.Width) && !i.cfg.FoceCrop {
		return []Image{*i}, nil
	}
	var extension = filepath.Ext(i.FileName)
	var name = i.FileName[0 : len(i.FileName)-len(extension)]
	leftImg := imaging.CropAnchor(i.Image, w/2, h, imaging.TopLeft)
	rightImg := imaging.CropAnchor(i.Image, w/2, h, imaging.TopRight)

	lfile, rfile := i.cfg.GetLRImageName()

	leftImage := Image{Image: leftImg, FileName: name + lfile, OutFormat: i.OutFormat, cfg: i.cfg}
	rightImage := Image{Image: rightImg, FileName: name + rfile, OutFormat: i.OutFormat, cfg: i.cfg}
	ret := make([]Image, 0)
	ret = append(ret, leftImage, rightImage)
	return ret, nil
}

func (i *Image) GetConfig() (image.Config, error) {
	ba, err := i.Bytes()
	if err != nil {
		return image.Config{}, err
	}
	r := bytes.NewReader(ba)
	config, _, err := image.DecodeConfig(r)
	return config, err
}

func (i *Image) OutputFileName() string {
	return i.FileName + "." + strings.ToLower(i.OutFormat.String())
}

func (i *Image) Bytes() ([]byte, error) {
	var b bytes.Buffer
	if err := imaging.Encode(bufio.NewWriter(&b), i.Image, i.OutFormat, i.cfg.ImgOpts()...); err != nil {
		return nil, err
	}
	return b.Bytes(), nil
}

type FolderImages struct {
	Images  []Image
	DirName string
}

func BatchCut(imgs []Image, cfg *ImageConfig) ([]Image, error) {
	ret := make([]Image, 0)
	mux := &sync.Mutex{}

	var wg sync.WaitGroup
	for _, img := range imgs {
		wg.Add(1)
		LimitChain <- struct{}{}
		go func(limg Image) {
			defer wg.Done()
			defer func() {
				<-LimitChain
			}()
			lrImgs, err := limg.Cut()
			if err != nil {
				panic(err)
			}
			mux.Lock()
			ret = append(ret, lrImgs...)
			mux.Unlock()
		}(img)
	}
	fmt.Println("Wait for processing")
	wg.Wait()
	return ret, nil
}
