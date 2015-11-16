package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/nesurion/milight-daemon/milight"
)

const (
	VERSION = "0.1.0"
)

func main() {
	if len(os.Args) > 1 {
		// if argument version was given, print VERSION and exit
		if os.Args[1] == "version" {
			fmt.Println(VERSION)
			os.Exit(0)
		}
	}

	ginMode := flag.String("mode", gin.ReleaseMode, "Gin Mode (debug, release, test)")
	flag.Parse()
	SetMode(*ginMode)
	router := gin.Default()
	router.Use(CORSMiddleware())

	c, err := milight.NewConfig("milight-daemon.conf")
	if err != nil {
		panic("failed to parse config file")
	}
	host := fmt.Sprintf("0.0.0.0:%d", c.Port)
	fmt.Println("=== Milight Daemon ===")
	fmt.Printf("Version %s\n", VERSION)
	fmt.Printf("Running on %s\n", host)

	// create limitless controller
	mc, err := milight.NewClient(c)
	if err != nil {
		panic(err)
	}

	router.POST("/on", func(c *gin.Context) {
		id, err := ParseGroup(c)
		if err != nil {
			c.AbortWithError(400, err)
			return
		}
		err = mc.On(id)
		if err != nil {
			c.AbortWithError(500, err)
			return
		}
		msg := gin.H{
			"command": "on",
			"group":   id,
		}
		c.JSON(200, msg)
	})

	router.POST("/off", func(c *gin.Context) {
		id, err := ParseGroup(c)
		if err != nil {
			c.AbortWithError(400, err)
			return
		}
		err = mc.Off(id)
		if err != nil {
			c.AbortWithError(500, err)
			return
		}
		msg := gin.H{
			"command": "off",
			"group":   id,
		}
		c.JSON(200, msg)
	})

	router.POST("/rgb", func(c *gin.Context) {
		id, err := ParseGroup(c)
		if err != nil {
			c.AbortWithError(400, err)
			return
		}
		rgb, err := ParseRGB(c)
		if err != nil {
			c.AbortWithError(400, err)
			return
		}
		err = mc.Rgb(id, rgb)
		if err != nil {
			c.AbortWithError(500, err)
			return
		}
		msg := gin.H{
			"command": "rgb",
			"rgb":     rgb,
			"group":   id,
		}
		c.JSON(200, msg)
	})

	router.POST("/brightness", func(c *gin.Context) {
		id, err := ParseGroup(c)
		if err != nil {
			c.AbortWithError(400, err)
			return
		}
		bl, err := ParseBrightnessLevel(c)
		if err != nil {
			c.AbortWithError(400, err)
			return
		}
		err = mc.Brightness(id, bl)
		if err != nil {
			c.AbortWithError(500, err)
			return
		}
		msg := gin.H{
			"command": "brightness",
			"level":   bl,
			"group":   id,
		}
		c.JSON(200, msg)
	})

	router.POST("/color", func(c *gin.Context) {
		id, err := ParseGroup(c)
		if err != nil {
			c.AbortWithError(400, err)
			return
		}
		color, err := ParseColor(c)
		if err != nil {
			c.AbortWithError(400, err)
			return
		}
		err = mc.Color(id, color)
		if err != nil {
			c.AbortWithError(500, err)
			return
		}
		msg := gin.H{
			"command": "color",
			"color":   fmt.Sprintf("%x", color),
			"group":   id,
		}
		c.JSON(200, msg)
	})

	router.POST("/white", func(c *gin.Context) {
		id, err := ParseGroup(c)
		if err != nil {
			c.AbortWithError(400, err)
			return
		}
		err = mc.White(id)
		if err != nil {
			c.AbortWithError(500, err)
			return
		}
		msg := gin.H{
			"command": "white",
			"group":   id,
		}
		c.JSON(200, msg)
	})

	router.POST("/night", func(c *gin.Context) {
		id, err := ParseGroup(c)
		if err != nil {
			c.AbortWithError(400, err)
			return
		}
		err = mc.Night(id)
		if err != nil {
			c.AbortWithError(500, err)
			return
		}
		msg := gin.H{
			"command": "night",
			"group":   id,
		}
		c.JSON(200, msg)
	})

	router.POST("/disco", func(c *gin.Context) {
		id, err := ParseGroup(c)
		if err != nil {
			c.AbortWithError(400, err)
			return
		}
		speed, err := parseSpeed(c)
		if err != nil {
			c.AbortWithError(400, err)
			return
		}
		err = mc.Disco(id, speed)
		if err != nil {
			c.AbortWithError(500, err)
			return
		}
		msg := gin.H{
			"command": "disco",
			"speed":   speed,
			"group":   id,
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

func CORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Max-Age", "86400")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE, UPDATE")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Origin, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
		c.Writer.Header().Set("Access-Control-Expose-Headers", "Content-Length")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")

		if c.Request.Method == "OPTIONS" {
			fmt.Println("OPTIONS")
			c.AbortWithStatus(200)
		} else {
			c.Next()
		}
	}
}
