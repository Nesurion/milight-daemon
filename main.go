package main

import (
	"encoding/json"
	"errors"
	"flag"
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
	ginMode := flag.String("mode", gin.ReleaseMode, "Gin Mode (debug, release, test)")
	flag.Parse()
	setMode(*ginMode)
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
		id, err := parseGroup(c)
		if err != nil {
			c.AbortWithError(400, err)
			return
		}
		msg := gin.H{
			"command": "on",
			"group":   0,
		}
		if id == -1 {
			for _, g := range controller.Groups {
				err = g.On()
				if err != nil {
					err = errors.New("failed to send off")
					c.AbortWithError(500, err)
					return
				}
			}
		} else {
			err = controller.Groups[id].On()
			if err != nil {
				err = errors.New("failed to send off")
				c.AbortWithError(500, err)
				return
			}
			msg["group"] = controller.Groups[id].Id
		}
		c.JSON(200, msg)
	})

	router.POST("/off", func(c *gin.Context) {
		id, err := parseGroup(c)
		if err != nil {
			c.AbortWithError(400, err)
			return
		}
		msg := gin.H{
			"command": "off",
			"group":   0,
		}
		if id == -1 {
			for _, g := range controller.Groups {
				err = g.Off()
				if err != nil {
					err = errors.New("failed to send off")
					c.AbortWithError(500, err)
					return
				}
			}
		} else {
			err = controller.Groups[id].Off()
			if err != nil {
				err = errors.New("failed to send off")
				c.AbortWithError(500, err)
				return
			}
			msg["group"] = controller.Groups[id].Id
		}
		c.JSON(200, msg)
	})

	router.POST("/color", func(c *gin.Context) {
		id, err := parseGroup(c)
		if err != nil {
			c.AbortWithError(400, err)
			return
		}
		color, err := parseColorRGB(c)
		if err != nil {
			c.AbortWithError(400, err)
			return
		}
		msg := gin.H{
			"command": "color",
			"color":   color,
			"group":   0,
		}
		if id == -1 {
			for _, g := range controller.Groups {
				err := g.SendColor(color)
				if err != nil {
					err = errors.New("failed to send color")
					c.AbortWithError(500, err)
					return
				}
			}
		} else {
			err = controller.Groups[id].SendColor(color)
			if err != nil {
				err = errors.New("failed to send color")
				c.AbortWithError(500, err)
				return
			}
			msg["group"] = controller.Groups[id].Id
		}
		c.JSON(200, msg)
	})

	router.POST("/brightness", func(c *gin.Context) {
		var bl uint8
		id, err := parseGroup(c)
		if err != nil {
			c.AbortWithError(400, err)
			return
		}
		bl, err = parseBrightnessLevel(c)
		if err != nil {
			c.AbortWithError(400, err)
			return
		}
		msg := gin.H{
			"command": "brightness",
			"level":   bl,
			"group":   0,
		}
		if id == -1 {
			for _, g := range controller.Groups {
				err = g.SetBri(bl)
				if err != nil {
					err = errors.New("failed to set brightness")
					c.AbortWithError(500, err)
					return
				}
			}
		} else {
			err = controller.Groups[id].SetBri(bl)
			if err != nil {
				err = errors.New("failed to set brightness")
				c.AbortWithError(500, err)
				return
			}
			msg["group"] = controller.Groups[id].Id
		}
		c.JSON(200, msg)
	})

	router.POST("/hue", func(c *gin.Context) {
		id, err := parseGroup(c)
		if err != nil {
			c.AbortWithError(400, err)
			return
		}
		color, err := parseColorName(c)
		if err != nil {
			c.AbortWithError(400, err)
			return
		}
		msg := gin.H{
			"command": "hue",
			"color":   fmt.Sprintf("%x", color),
			"group":   0,
		}
		if id == -1 {
			for _, g := range controller.Groups {
				err = g.SetHue(color)
				if err != nil {
					err = errors.New("failed to set hue")
					c.AbortWithError(500, err)
					return
				}
			}
		} else {
			err = controller.Groups[id].SetHue(color)
			if err != nil {
				err = errors.New("failed to set hue")
				c.AbortWithError(500, err)
				return
			}
			msg["group"] = controller.Groups[id].Id
		}
		c.JSON(200, msg)
	})

	router.POST("/white", func(c *gin.Context) {
		id, err := parseGroup(c)
		if err != nil {
			c.AbortWithError(400, err)
			return
		}
		msg := gin.H{
			"command": "white",
			"group":   0,
		}
		if id == -1 {
			for _, g := range controller.Groups {
				err = g.White()
				if err != nil {
					err = errors.New("failed to set white")
					c.AbortWithError(500, err)
					return
				}
			}
		} else {
			err = controller.Groups[id].White()
			if err != nil {
				err = errors.New("failed to set white")
				c.AbortWithError(500, err)
				return
			}
			msg["group"] = controller.Groups[id].Id
		}
		c.JSON(200, msg)
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

func parseGroup(c *gin.Context) (int, error) {
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

func parseColorRGB(c *gin.Context) (colorful.Color, error) {
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
	fmt.Println(rgb)
	color := colorful.Color{
		rgb["r"] / 255.0,
		rgb["g"] / 255.0,
		rgb["b"] / 255.0,
	}
	return color, nil
}

func parseBrightnessLevel(c *gin.Context) (uint8, error) {
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

func parseColorName(c *gin.Context) (uint8, error) {
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

func setMode(mode string) {
	switch mode {
	case "release":
		gin.SetMode(gin.ReleaseMode)
	case "debug":
		gin.SetMode(gin.DebugMode)
	case "test":
		gin.SetMode(gin.TestMode)
	default:
		panic("mode unavailable. (debug, release, test)")
	}
}
