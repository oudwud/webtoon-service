package server

import (
	"image"
	"image/draw"
	"image/png"
	"io"
	"mime/multipart"
	"net/http"
	"strconv"

	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"

	"github.com/gin-gonic/gin"
)

func mergeImages(c *gin.Context) error {
	form, err := c.MultipartForm()
	if err != nil {
		return wrapError(http.StatusBadRequest, err, "fail to get multipart form")
	}
	defer form.RemoveAll()

	width, err := strconv.Atoi(c.PostForm("width"))
	if err != nil {
		return wrapError(http.StatusBadRequest, nil, "invalid width value: %v", c.PostForm("width"))
	}

	height, err := strconv.Atoi(c.PostForm("height"))
	if err != nil {
		return wrapError(http.StatusBadRequest, nil, "invalid height value: %v", c.PostForm("height"))
	}

	log.Debugf("width:%d, height:%d", width, height)

	files := form.File["files"]
	if files == nil {
		return wrapError(http.StatusBadRequest, nil, "no files in the multipart form")
	}

	img, err := merge(width, height, files)
	if err != nil {
		return wrapError(http.StatusInternalServerError, err, "fail to merge images")
	}

	clientGone := c.Stream(func(w io.Writer) bool {
		if err := png.Encode(w, img); err != nil {
			log.Error("fail to encode the merged image: ", err)
		}
		return false
	})
	if clientGone {
		log.Error("fail to write the merge image to the client. the client was gone.")
	}

	return nil
}

func merge(width, height int, files []*multipart.FileHeader) (image.Image, error) {
	dest := image.NewNRGBA(image.Rectangle{
		Min: image.Point{0, 0},
		Max: image.Point{width, height},
	})
	curDestHeight := 0

	var err error
	for _, file := range files {
		curDestHeight, err = putImage(dest, curDestHeight, file)
		if err != nil {
			return nil, errors.Wrap(err, "fail to put image")
		}
	}

	return dest, nil
}

func putImage(dest *image.NRGBA, curDestHeight int, file *multipart.FileHeader) (int, error) {
	destWidth := dest.Bounds().Max.X

	f, err := file.Open()
	if err != nil {
		return 0, errors.Wrapf(err, "fail to open the multipart file: %s", file.Filename)
	}
	defer f.Close()

	img, _, err := image.Decode(f)
	if err != nil {
		return 0, errors.Wrapf(err, "fail to decode the image: %s", file.Filename)
	}

	startingPoint := image.Point{(destWidth - img.Bounds().Dx()) / 2, curDestHeight}

	rect := image.Rectangle{
		Min: startingPoint,
		Max: startingPoint.Add(img.Bounds().Size()),
	}

	draw.Draw(dest, rect, img, image.Point{0, 0}, draw.Src)

	return rect.Max.Y, nil
}
