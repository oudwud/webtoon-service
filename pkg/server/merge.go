package server

import (
	"image"
	"image/color"
	"image/png"
	"io"
	"mime/multipart"
	"net/http"
	"strconv"

	"github.com/disintegration/imaging"
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
		err := imaging.Encode(w, img, imaging.PNG, imaging.PNGCompressionLevel(png.DefaultCompression))
		if err != nil {
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
	dest := imaging.New(width, height, color.NRGBA{0, 0, 0, 0})
	curDestHeight := 0

	var err error
	for _, file := range files {
		dest, curDestHeight, err = putImage(dest, curDestHeight, file)
		if err != nil {
			return nil, errors.Wrap(err, "fail to put image")
		}
	}

	return dest, nil
}

func putImage(dest *image.NRGBA, curDestHeight int, file *multipart.FileHeader) (*image.NRGBA, int, error) {
	destWidth := dest.Bounds().Max.X

	f, err := file.Open()
	if err != nil {
		return nil, 0, errors.Wrapf(err, "fail to open the multipart file: %s", file.Filename)
	}
	defer f.Close()

	img, err := imaging.Decode(f)
	if err != nil {
		return nil, 0, errors.Wrapf(err, "fail to decode the image: %s", file.Filename)
	}

	x := (destWidth - img.Bounds().Max.X) / 2
	dest = imaging.Paste(dest, img, image.Pt(x, curDestHeight))
	return dest, curDestHeight + img.Bounds().Max.Y, nil
}
