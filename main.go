package main

import (
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strconv"

	"config"
	"milight"

	"github.com/evq/go-limitless"
	"github.com/gin-gonic/gin"
	"github.com/lucasb-eyer/go-colorful"
)

func main() {
	ginMode := flag.String("mode", gin.ReleaseMode, "Gin Mode (debug, release, test)")
	flag.Parse()
	SetMode(*ginMode)
	router := gin.Default()

	c, err := ginconfig.ParseConfig("milight-daemon.conf")
	if err != nil {
		panic("failed to parse config file")
	}
	host := fmt.Sprintf("0.0.0.0:%d", c.Port)

	// create limitless controller
	controller := limitless.LimitlessController{
		Host: c.Bridge,
	}
	groups := milight.Groups(&controller)
	controller.Groups = groups

	router.POST("/on", func(c *gin.Context) {
		id, err := milight.ParseGroup(c)
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
		id, err := milight.ParseGroup(c)
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

	router.POST("/rgb", func(c *gin.Context) {
		id, err := milight.ParseGroup(c)
		if err != nil {
			c.AbortWithError(400, err)
			return
		}
		rgb, err := milight.ParseRGB(c)
		if err != nil {
			c.AbortWithError(400, err)
			return
		}
		msg := gin.H{
			"command": "rgb",
			"rgb":     rgb,
			"group":   0,
		}
		if id == -1 {
			for _, g := range controller.Groups {
				err := g.SendColor(rgb)
				if err != nil {
					err = errors.New("failed to send color")
					c.AbortWithError(500, err)
					return
				}
			}
		} else {
			err = controller.Groups[id].SendColor(rgb)
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
		id, err := milight.ParseGroup(c)
		if err != nil {
			c.AbortWithError(400, err)
			return
		}
		bl, err = milight.ParseBrightnessLevel(c)
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

	router.POST("/color", func(c *gin.Context) {
		id, err := milight.ParseGroup(c)
		if err != nil {
			c.AbortWithError(400, err)
			return
		}
		color, err := milight.ParseColor(c)
		if err != nil {
			c.AbortWithError(400, err)
			return
		}
		msg := gin.H{
			"command": "color",
			"color":   fmt.Sprintf("%x", color),
			"group":   0,
		}
		if id == -1 {
			for _, g := range controller.Groups {
				err = g.SetHue(color)
				if err != nil {
					err = errors.New("failed to set color")
					c.AbortWithError(500, err)
					return
				}
			}
		} else {
			err = controller.Groups[id].SetHue(color)
			if err != nil {
				err = errors.New("failed to set color")
				c.AbortWithError(500, err)
				return
			}
			msg["group"] = controller.Groups[id].Id
		}
		c.JSON(200, msg)
	})

	router.POST("/white", func(c *gin.Context) {
		id, err := milight.ParseGroup(c)
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

func SetMode(mode string) {
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
