package handlers

import (
	"bufio"
	"io"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog/log"
	"github.com/spf13/viper"

	"github.com/alexferl/air/asset"
	"github.com/alexferl/air/util"
)

func (h *Handler) Upload(c echo.Context) error {
	var maxFileSize = viper.GetInt64("max-file-size") << 20

	c.Request().Body = http.MaxBytesReader(c.Response(), c.Request().Body, maxFileSize+1024)
	reader, err := c.Request().MultipartReader()
	if err != nil {
		if err.Error() == "request Content-Type isn't multipart/form-data" {
			return c.JSON(http.StatusUnsupportedMediaType, ErrorResponse{err.Error()})
		} else if err.Error() == "no multipart boundary param in Content-Type" {
			return c.JSON(http.StatusUnprocessableEntity, ErrorResponse{err.Error()})
		} else {
			log.Error().Msgf("Failed to process multipart form: %v", err)
			return c.JSON(http.StatusInternalServerError, ErrorResponse{"Error processing form"})
		}
	}

	p, err := reader.NextPart()
	if err != nil {
		log.Error().Msgf("Failed to read multipart form: %v", err)
		return c.JSON(http.StatusInternalServerError, ErrorResponse{"Error reading form"})
	}

	if p.FormName() != "file" {
		return c.JSON(http.StatusBadRequest, ErrorResponse{"'file' field is required"})
	}

	buf := bufio.NewReader(p)
	f, err := util.CreateTempFile()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, ErrorResponse{"Error creating temporary file"})
	}
	defer util.CleanupTempFile(f)

	lmt := io.LimitReader(buf, maxFileSize+1)
	written, err := io.Copy(f, lmt)
	if err != nil && err != io.EOF {
		log.Error().Msgf("Failed to write to temp file: %v", err)
		return c.JSON(http.StatusInternalServerError, ErrorResponse{"Error writing temporary file"})
	}

	if written > maxFileSize {
		return c.JSON(http.StatusRequestEntityTooLarge, ErrorResponse{"File is too large"})
	}

	f.Seek(0, 0) // rewind file

	a, err := asset.New(f)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, ErrorResponse{"Error saving file"})
	}
	defer util.CleanupTempFile(a.File)

	err = h.Storage.Put(a)
	if err != nil {
		log.Error().Msgf("Failed to save file: %v", err)
		return c.JSON(http.StatusInternalServerError, ErrorResponse{"Error saving file to storage"})
	}

	c.Response().Header().Set("Location", "/"+a.Name)
	return c.JSON(http.StatusCreated, map[string]string{"id": a.Name})
}
