package asset

import (
	"bufio"
	"bytes"
	"crypto/sha256"
	"errors"
	"fmt"
	"io"
	"mime"
	"net/http"
	"os"
	"strings"

	"github.com/davidbyttow/govips/v2/vips"
	"github.com/rs/zerolog/log"

	"github.com/alexferl/air/util"
)

const (
	IMAGE                  = "image"
	JPEG                   = "image/jpeg"
	PNG                    = "image/png"
	WEBP                   = "image/webp"
	DefaultImageType       = vips.ImageTypeJPEG
	DefaultInterestingType = vips.InterestingNone
)

var imageTypes = []string{JPEG, PNG, WEBP}

type Asset struct {
	File        *os.File
	ContentType string
	// Preferred extension for files that have more than one
	Ext        string
	Extensions []string
	Name       string
	PathPrefix string
	Path       string
	Sha256     string
	Type       string

	buf bytes.Buffer
}

func New(r io.Reader) (*Asset, error) {
	f, err := util.CreateTempFile()
	if err != nil {
		return nil, err
	}

	var buf bytes.Buffer
	if _, err := io.Copy(&buf, r); err != nil {
		log.Error().Msgf("Failed to copy buffer: %v", err)
		return nil, err
	}

	_, err = f.Write(buf.Bytes())
	if err != nil {
		log.Error().Msgf("Failed to write buffer to temp file: %v", err)
		return nil, err
	}

	// rewind file
	f.Seek(0, 0)

	a := &Asset{
		File: f,
	}

	err = a.load()
	if err != nil {
		return nil, err
	}

	return a, nil
}

type ResizeParams struct {
	Width       int
	Height      int
	Quality     int
	ImageType   vips.ImageType
	Interesting vips.Interesting
}

func NewResizeParams() *ResizeParams {
	return &ResizeParams{
		ImageType:   DefaultImageType,
		Interesting: DefaultInterestingType,
	}
}

func (a *Asset) Resize(rp *ResizeParams) ([]byte, error) {
	defer a.rewind()

	if a.Type != IMAGE {
		return nil, errors.New("file type doesn't support resizing")
	}

	image, err := vips.LoadImageFromBuffer(a.buf.Bytes(), vips.NewImportParams())
	if err != nil {
		log.Error().Msgf("Failed to load image: %v", err)
		return nil, err
	}

	var force bool
	if rp.Width > 0 && rp.Height > 0 {
		force = true
	}

	if rp.Width == 0 {
		rp.Width = image.Width()
	}
	if rp.Height == 0 {
		rp.Height = image.Height()
	}

	log.Debug().Msgf(
		"Resize called with: width: %d height: %d quality: %d imageType: %s interesting: %s",
		rp.Width,
		rp.Height,
		rp.Quality,
		strings.Split(rp.ImageType.FileExt(), ".")[1],
		InterestingTypesToString[rp.Interesting],
	)

	if !force {
		err = image.Thumbnail(rp.Width, rp.Height, rp.Interesting)
	} else {
		err = image.ThumbnailWithSize(rp.Width, rp.Height, rp.Interesting, vips.SizeForce)
	}
	if err != nil {
		log.Error().Msgf("Failed to create thumbnail: %v", err)
		return nil, err
	}

	err = image.RemoveMetadata()
	if err != nil {
		log.Error().Msgf("Failed to remove metadata: %v", err)
		return nil, err
	}

	if rp.Quality > 0 {
	}

	var b []byte
	var exportErr error
	switch rp.ImageType {
	case vips.ImageTypeJPEG:
		p := vips.NewJpegExportParams()
		if rp.Quality > 0 {
			p.Quality = rp.Quality
		}
		b, _, exportErr = image.ExportJpeg(p)
	case vips.ImageTypePNG:
		b, _, exportErr = image.ExportPng(vips.NewPngExportParams())
	case vips.ImageTypeWEBP:
		p := vips.NewWebpExportParams()
		if rp.Quality > 0 {
			p.Quality = rp.Quality
		}
		b, _, exportErr = image.ExportWebp(p)
	default:
		return nil, errors.New("unknown image type")
	}

	if exportErr != nil {
		log.Error().Msgf("Failed to export image with type %s: %v", vips.ImageTypes[rp.ImageType], exportErr)
		return nil, exportErr
	}

	a.setMimeTypeFromExt(rp.ImageType.FileExt())
	a.setExtensionsFromMimeType(a.ContentType)

	return b, nil
}

func (a *Asset) load() error {
	defer a.rewind()

	r := bufio.NewReader(a.File)
	var buf bytes.Buffer
	tee := io.TeeReader(r, &buf)

	h := sha256.New()
	if _, err := io.Copy(h, tee); err != nil {
		log.Error().Msgf("Failed to hash file")
		return err
	}

	a.detectContentType(buf.Bytes())
	a.Sha256 = fmt.Sprintf("%x", h.Sum(nil))
	a.Name = a.Sha256
	a.Path, _ = util.GetFullPathFromSha256(a.Sha256)
	prefix := strings.Split(a.Path, "/")[0:3]
	a.PathPrefix = strings.Join(prefix, "/")
	a.buf = buf

	return nil
}

func (a *Asset) detectContentType(buf []byte) {
	a.ContentType = http.DetectContentType(buf)
	a.setExtensionsFromMimeType(a.ContentType)
	for _, t := range imageTypes {
		if t == a.ContentType {
			a.Type = strings.Split(a.ContentType, "/")[0]
		}
	}
}

func (a *Asset) setExtensionsFromMimeType(typ string) {
	exts, _ := mime.ExtensionsByType(typ)
	if len(exts) == 1 {
		a.Ext = exts[0]
	}

	if len(exts) > 1 {
		a.Extensions = exts
		switch typ {
		case JPEG:
			a.Ext = ".jpeg"
		default:
			a.Ext = exts[0]
		}
	}
}

func (a *Asset) setMimeTypeFromExt(ext string) {
	a.ContentType = mime.TypeByExtension(ext)
}

func (a *Asset) rewind() {
	_, err := a.File.Seek(0, 0)
	if err != nil {
		log.Error().Msgf("Failed to seek file %s: %v", a.File.Name(), err)
	}
}
