package asset

import (
	"bytes"
	"io/ioutil"
	"os"
	"testing"

	"github.com/davidbyttow/govips/v2/vips"
	"github.com/stretchr/testify/assert"
)

func TestAsset(t *testing.T) {
	vips.Startup(nil)
	defer vips.Shutdown()

	images := []struct {
		orig        string
		width       int
		height      int
		transformed string
		origExt     string
		transExt    string
		origType    string
		transType   string
		transFmt    vips.ImageType
	}{
		{"cat.png", 640, 0, "cat_640.png", ".png", ".png", PNG, PNG, vips.ImageTypePNG},
		{"cat.png", 0, 0, "cat.webp", ".png", ".webp", PNG, WEBP, vips.ImageTypeWEBP},
		{"cat.png", 0, 0, "cat.jpg", ".png", ".jpeg", PNG, JPEG, vips.ImageTypeJPEG},
		{"tiger.jpg", 640, 0, "tiger_640.jpg", ".jpeg", ".jpeg", JPEG, JPEG, vips.ImageTypeJPEG},
		{"tiger.jpg", 640, 0, "tiger_640.webp", ".jpeg", ".webp", JPEG, WEBP, vips.ImageTypeWEBP},
		{"tiger.jpg", 0, 0, "tiger.png", ".jpeg", ".png", JPEG, PNG, vips.ImageTypePNG},
	}

	for _, img := range images {
		f, err := os.Open("../fixtures/" + img.orig)
		assert.NoError(t, err)

		a, err := New(f)
		assert.NoError(t, err)

		assert.Equal(t, img.origExt, a.Ext)
		assert.Equal(t, img.origType, a.ContentType)
		assert.Equal(t, "image", a.Type)

		b, err := a.Resize(&ResizeParams{Width: img.width, Height: img.height, ImageType: img.transFmt})
		assert.NoError(t, err)

		c, err := ioutil.ReadFile("../fixtures/" + img.transformed)
		assert.NoError(t, err)

		assert.Equal(t, img.transType, a.ContentType)
		assert.Equal(t, img.transExt, a.Ext)
		assert.Equal(t, 0, bytes.Compare(c, b))

		f.Close()
		a.File.Close()
		os.Remove(a.File.Name())
	}
}
