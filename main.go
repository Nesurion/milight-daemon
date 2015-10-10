package main

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
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

type Config struct {
	Port   int    `json:"port"`
	Bridge string `json:"bridge"`
}

func main() {
	router := gin.Default()

	c, err := parseConfig("milight-daemon.conf")
	if err != nil {
		panic("failed to parse config file")
	}
	host := fmt.Sprintf("0.0.0.0:%d", c.Port)

	// create limitless controller
	controller := limitless.LimitlessController{
		Host: c.Bridge,
	}
	groups := groups(&controller)
	controller.Groups = groups

	router.POST("/on", func(c *gin.Context) {
		id := parseGroup(c)
		if id == 0 {
			for _, g := range controller.Groups {
				g.On()
			}
			return
		}
		controller.Groups[id].On()
	})

	router.POST("/off", func(c *gin.Context) {
		id := parseGroup(c)
		if id == 0 {
			for _, g := range controller.Groups {
				g.Off()
			}
			return
		}
		controller.Groups[id].Off()
	})

	router.POST("/color", func(c *gin.Context) {
		id := parseGroup(c)
		color := parseColorRGB(c)
		if id == 0 {
			for _, g := range controller.Groups {
				err := g.SendColor(color)
				if err != nil {
					c.String(500, "failed to send color")
				}
			}
		}
		err := controller.Groups[id].SendColor(color)
		if err != nil {
			c.String(500, "failed to send color")
		}
	})

	router.POST("/brightness", func(c *gin.Context) {
		id := parseGroup(c)
		bl := parseBrightnessLevel(c)
		if id == 0 {
			for _, g := range controller.Groups {
				g.SetBri(bl)
			}
			return
		}
		err := controller.Groups[id].SetBri(bl)
		if err != nil {
			c.String(500, "failed to send color")
		}
	})

	router.POST("/hue", func(c *gin.Context) {
		id := parseGroup(c)
		color := parseColorName(c)
		if id == 0 {
			for _, g := range controller.Groups {
				g.SetHue(color)
			}
			return
		}
		err := controller.Groups[id].SetHue(color)
		if err != nil {
			c.String(500, "failed to send color")
		}
	})

	router.POST("/white", func(c *gin.Context) {
		id := parseGroup(c)
		if id == 0 {
			for _, g := range controller.Groups {
				g.White()
			}
			return
		}
		controller.Groups[id].White()
	})

	router.Run(host)
}

func parseConfig(configPath string) (Config, error) {
	c := Config{}
	absConfigPath, err := filepath.Abs(configPath)
	if err != nil {
		return c, err
	}
	file, _ := os.Open(absConfigPath)
	decoder := json.NewDecoder(file)
	err = decoder.Decode(&c)
	if err != nil {
		return c, err
	}
	return c, nil
}

func groups(c *limitless.LimitlessController) []limitless.LimitlessGroup {
	g := make([]limitless.LimitlessGroup, 4, 4)
	for i := 0; i < 4; i++ {
		g[i] = limitless.LimitlessGroup{
			Id:         i + 1,
			Controller: c,
		}
	}
	return g
}

func parseGroup(c *gin.Context) int {
	group := c.Query("group")
	var id int
	var err error
	if id, err = strconv.Atoi(group); err != nil {
		// return http error instead of logging
		c.String(500, "failed to parse group")
	}
	if id < 0 || id > 4 {
		c.String(400, "invalid id. must be <= 0 or >= 4")
	}
	// use id as index for groups
	id = id + 1
	return id
}

func parseColorRGB(c *gin.Context) colorful.Color {
	rgb := map[string]float64{
		"r": 0,
		"g": 0,
		"b": 0,
	}
	var err error
	for k, v := range rgb {
		if v, err = strconv.ParseFloat(c.Query(k), 64); err != nil {
			c.String(500, "failed to parse color")
		}
		if v < 0 || v > 255 {
			c.String(400, "invalid color value %s %d. must be <= 0 or >= 255")
		}
	}
	color := colorful.Color{
		rgb["r"] / 255.0,
		rgb["g"] / 255.0,
		rgb["b"] / 255.0,
	}
	return color
}

func parseBrightnessLevel(c *gin.Context) uint8 {
	level := c.Query("level")
	var err error
	var b64 uint64
	if b64, err = strconv.ParseUint(level, 10, 8); err != nil {
		c.String(500, "failed to parse brightness level")
	}
	b := uint8(b64)
	if b < BRIGHTNESS_MIN || b > BRIGHTNESS_MAX {
		c.String(400, "invalid brightness level. Must be between 1-100")
	}
	b = b/BRIGHTNESS_RATIO + BRIGHTNESS_OFFSET
	return b
}

func parseColorName(c *gin.Context) uint8 {
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
		c.String(400, "invalid color name %s", color)
	}
	return colorHex
}
