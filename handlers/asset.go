package handlers

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"regexp"
	"strconv"
	"strings"

	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog/log"

	"github.com/alexferl/air/asset"
	"github.com/alexferl/air/util"
)

const (
	maxWidth   = 5000
	maxHeight  = 5000
	maxQuality = 100
)

func (h *Handler) Asset(c echo.Context) error {
	id := c.Param("id")
	format := c.QueryParam("format")
	crop := c.QueryParam("crop")

	if len(crop) > 0 {
		if _, ok := asset.StringToInterestingTypes[crop]; !ok {
			m := fmt.Sprintf("Unknown crop algorithm '%s'", crop)
			return c.JSON(http.StatusBadRequest, ErrorResponse{m})
		}
	}

	width, height, quality, err := parseParams(c.QueryParams())
	if err != nil {
		return c.JSON(http.StatusBadRequest, ErrorResponse{err.Error()})
	}

	path, err := util.GetFullPathFromSha256(id)
	if err != nil {
		return c.JSON(http.StatusBadRequest, ErrorResponse{"Invalid id"})
	}

	f, err := h.Storage.Get(path)
	if err != nil {
		return c.JSON(http.StatusNotFound, ErrorResponse{"File not found"})
	}
	defer f.Close()

	a, err := asset.New(f)
	if err != nil {
		log.Error().Msgf("Failed to create asset: %v", err)
		return c.JSON(http.StatusInternalServerError, ErrorResponse{"Error sending file"})
	}
	defer util.CleanupTempFile(a.File)

	var reader io.Reader
	if a.Type == "image" {
		out, err := util.CreateTempFile()
		if err != nil {
			return c.JSON(http.StatusInternalServerError, ErrorResponse{"Error creating temporary file"})
		}
		defer util.CleanupTempFile(out)

		rp := asset.NewResizeParams()
		rp.Width = width
		rp.Height = height
		rp.Quality = quality

		if len(format) == 0 {
			format = strings.Split(a.Ext, ".")[1]
		}

		if val, ok := asset.ImageTypes[format]; ok {
			rp.ImageType = val
		} else {
			m := fmt.Sprintf("Unknown image format '%s'", format)
			return c.JSON(http.StatusBadRequest, ErrorResponse{m})
		}

		b, err := a.Resize(rp)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, ErrorResponse{"Error resizing image"})
		}
		reader = bytes.NewReader(b)

	}

	c.Response().Header().Set("Cache-Control", "public, max-age=604800")
	return c.Stream(http.StatusOK, a.ContentType, reader)
}

func parseParams(params url.Values) (int, int, int, error) {
	w := params.Get("width")
	h := params.Get("height")
	s := params.Get("size")
	q := params.Get("quality")

	var width, height, quality int
	var err error

	if w != "" {
		width, err = strconv.Atoi(w)
		if err != nil {
			return 0, 0, 0, errors.New("width must be a number")
		}
	}
	if h != "" {
		height, err = strconv.Atoi(h)
		if err != nil {
			return 0, 0, 0, errors.New("height must be a number")
		}
	}

	if q != "" {
		quality, err = strconv.Atoi(q)
		if err != nil {
			return 0, 0, 0, errors.New("quality must be a number")
		}
	}

	if len(s) > 0 {
		r := regexp.MustCompile(`(\d+)[x](\d+)$`)
		size := r.FindStringSubmatch(s)
		if len(size) == 3 {
			width, _ = strconv.Atoi(size[1])
			height, _ = strconv.Atoi(size[2])
		} else {
			r := regexp.MustCompile(`(\d+)$`)
			m := r.FindString(s)
			if m == "" {
				return 0, 0, 0, errors.New("incorrect format for size")
			} else {
				width, _ = strconv.Atoi(s)
			}
		}
	}

	if width < 0 {
		return 0, 0, 0, errors.New("width cannot be less than 0")
	}
	if height < 0 {
		return 0, 0, 0, errors.New("height cannot be less than 0")
	}
	if quality < 0 {
		return 0, 0, 0, errors.New("quality cannot be less than 0")
	}
	if width > maxWidth {
		return 0, 0, 0, errors.New(fmt.Sprintf("width cannot be above %d", maxWidth))
	}
	if height > maxHeight {
		return 0, 0, 0, errors.New(fmt.Sprintf("height cannot be above %d", maxHeight))
	}
	if quality > maxQuality {
		return 0, 0, 0, errors.New(fmt.Sprintf("quality cannot be above %d", maxQuality))
	}

	return width, height, quality, nil
}
