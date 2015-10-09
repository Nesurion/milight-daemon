package main

import (
	"encoding/json"
	"fmt"
	"log"
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
	router := gin.Default()
	c, err := parseConfig("milight-daemon.conf")
	if err != nil {
		panic("failed to parse config file")
	}
	host := fmt.Sprintf("0.0.0.0:%d", c.Port)
	// create controller
	controller := limitless.LimitlessController{
		Host: c.Bridge,
	}
	groups := groups(&controller)
	controller.Groups = groups

	fmt.Println(controller)

	router.POST("/on", func(c *gin.Context) {
		group := c.Query("group")
		var id int
		var err error
		if id, err = strconv.Atoi(group); err != nil {
			// return http error instead of logging
			log.Print("failed to parse group id")
		}
		switch id {
		case 1:
			controller.Groups[0].On()
		case 2:
			controller.Groups[1].On()
		case 3:
			controller.Groups[2].On()
		case 4:
			controller.Groups[3].On()
		default:
			for _, g := range controller.Groups {
				g.On()
			}
		}
	})
	router.POST("/off", func(c *gin.Context) {
		// this is ugly as hell
		//TODO: find a way to export all the parsing and stuff into a function
		group := c.Query("group")
		var id int
		var err error
		if id, err = strconv.Atoi(group); err != nil {
			log.Print("failed to parse group id")
			return
		}
		switch id {
		case 1:
			controller.Groups[0].Off()
		case 2:
			controller.Groups[1].Off()
		case 3:
			controller.Groups[2].Off()
		case 4:
			controller.Groups[3].Off()
		default:
			for _, g := range controller.Groups {
				g.Off()
			}
		}
	})
	router.POST("/color", func(c *gin.Context) {
		group := c.Query("group")
		r := c.Query("r")
		g := c.Query("g")
		b := c.Query("b")
		var id int
		var red, green, blue float64
		var err error
		if id, err = strconv.Atoi(group); err != nil {
			log.Print("failed to parse group id")
			return
		}
		if red, err = strconv.ParseFloat(r, 64); err != nil {
			log.Print("failed to parse color")
			return
		}
		if green, err = strconv.ParseFloat(g, 64); err != nil {
			log.Print("failed to parse color")
			return
		}
		if blue, err = strconv.ParseFloat(b, 64); err != nil {
			log.Print("failed to parse color")
			return
		}
		color := colorful.Color{red / 255.0, green / 255.0, blue / 255.0}
		switch id {
		case 1:
			err = controller.Groups[0].SendColor(color)
			if err != nil {
				log.Print("failed to send color")
			}
		case 2:
			err = controller.Groups[0].SendColor(color)
			if err != nil {
				log.Print("failed to send color")
			}
		case 3:
			err = controller.Groups[0].SendColor(color)
			if err != nil {
				log.Print("failed to send color")
			}
		case 4:
			err = controller.Groups[0].SendColor(color)
			if err != nil {
				log.Print("failed to send color")
			}
		default:
			for _, g := range controller.Groups {
				err = g.SendColor(color)
				if err != nil {
					log.Print("failed to send color")
				}
			}
		}
	})
	router.POST("/brightness", func(c *gin.Context) {
		group := c.Query("group")
		level := c.Query("level")
		var id int
		var b64 uint64
		var err error
		if id, err = strconv.Atoi(group); err != nil {
			log.Print("failed to parse group id")
			return
		}
		if b64, err = strconv.ParseUint(level, 10, 8); err != nil {
			log.Print("failed to parse brightness level")
			return
		}
		b := uint8(b64)
		if b < BRIGHTNESS_MIN || b > BRIGHTNESS_MAX {
			log.Print("failed to set brightness level. Must be between 1-100")
			return
		}
		b = b/BRIGHTNESS_RATIO + BRIGHTNESS_OFFSET
		switch id {
		case 1:
			controller.Groups[0].SetBri(b)
		case 2:
			controller.Groups[1].SetBri(b)
		case 3:
			controller.Groups[2].SetBri(b)
		case 4:
			controller.Groups[3].SetBri(b)
		default:
			for _, g := range controller.Groups {
				g.SetBri(b)
			}
		}
	})
	router.POST("/hue", func(c *gin.Context) {
		group := c.Query("group")
		color := c.Query("color")
		var id int
		var err error
		if id, err = strconv.Atoi(group); err != nil {
			log.Print("failed to parse group id")
			return
		}
		colorHex, ok := Colors[color]
		if !ok {
			log.Printf("color %s not available\n", color)
			return
		}
		switch id {
		case 1:
			controller.Groups[0].SetHue(colorHex)
		case 2:
			controller.Groups[1].SetHue(colorHex)
		case 3:
			controller.Groups[2].SetHue(colorHex)
		case 4:
			controller.Groups[3].SetHue(colorHex)
		default:
			for _, g := range controller.Groups {
				g.SetHue(colorHex)
			}
		}
	})
	router.POST("/white", func(c *gin.Context) {
		group := c.Query("group")
		var id int
		var err error
		if id, err = strconv.Atoi(group); err != nil {
			log.Print("failed to parse group id")
		}
		switch id {
		case 1:
			controller.Groups[0].White()
		case 2:
			controller.Groups[1].White()
		case 3:
			controller.Groups[2].White()
		case 4:
			controller.Groups[3].White()
		default:
			for _, g := range controller.Groups {
				g.White()
			}
		}
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
