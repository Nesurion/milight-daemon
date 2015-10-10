package milight

import (
	"errors"
	"strconv"

	"github.com/evq/go-limitless"
	"github.com/gin-gonic/gin"
	"github.com/lucasb-eyer/go-colorful"
)

const (
	BRIGHTNESS_RATIO  = 4
	BRIGHTNESS_OFFSET = 2
	BRIGHTNESS_MIN    = 1
	BRIGHTNESS_MAX    = 100
)

func Groups(c *limitless.LimitlessController) []limitless.LimitlessGroup {
	g := make([]limitless.LimitlessGroup, 4, 4)
	for i := 0; i < 4; i++ {
		g[i] = limitless.LimitlessGroup{
			Id:         i + 1,
			Controller: c,
		}
	}
	return g
}

func ParseGroup(c *gin.Context) (int, error) {
	group := c.Query("group")
	id, err := strconv.Atoi(group)
	if err != nil {
		err = errors.New("failed to parse group")
		return -1, err
	}
	if id < 0 || id > 4 {
		err = errors.New("invalid id. must be <= 0 or >= 4")
		return -1, err
	}
	// use id as index for groups
	id = id - 1
	return id, nil
}

func ParseRGB(c *gin.Context) (colorful.Color, error) {
	rgb := map[string]float64{
		"r": 0,
		"g": 0,
		"b": 0,
	}
	var err error
	for k, v := range rgb {
		if v, err = strconv.ParseFloat(c.Query(k), 64); err != nil {
			err = errors.New("failed to parse color")
			return colorful.Color{}, err
		}
		if v < 0 || v > 255 {
			err = errors.New("invalid color value. must be <= 0 or >= 255")
			return colorful.Color{}, err
		}
		rgb[k] = v
	}
	color := colorful.Color{
		rgb["r"] / 255.0,
		rgb["g"] / 255.0,
		rgb["b"] / 255.0,
	}
	return color, nil
}

func ParseBrightnessLevel(c *gin.Context) (uint8, error) {
	level := c.Query("level")
	b64, err := strconv.ParseUint(level, 10, 8)
	if err != nil {
		err = errors.New("failed to parse brightness level")
		return 0, err
	}
	b := uint8(b64)
	if b < BRIGHTNESS_MIN || b > BRIGHTNESS_MAX {
		err = errors.New("invalid brightness level. Must be between 1-100")
		return 0, err
	}
	b = b/BRIGHTNESS_RATIO + BRIGHTNESS_OFFSET
	return b, nil
}

func ParseColor(c *gin.Context) (uint8, error) {
	Colors := map[string]uint8{
		"violet":        0x00,
		"blue":          0x10,
		"baby_blue":     0x20,
		"aqua":          0x30,
		"mint":          0x40,
		"seafoam_green": 0x50,
		"green":         0x60,
		"lime_green":    0x70,
		"yellow":        0x80,
		"yellow_orange": 0x90,
		"orange":        0xA0,
		"red":           0xB0,
		"pink":          0xC0,
		"fusia":         0xD0,
		"lilac":         0xE0,
		"lavendar":      0xF0,
	}

	color := c.Query("color")
	colorHex, ok := Colors[color]
	if !ok {
		err := errors.New("invalid color name")
		return 0, err
	}
	return colorHex, nil
}
