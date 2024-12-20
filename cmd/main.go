package main

import "github.com/xbmlz/ipa"

func main() {
	app := ipa.New()

	app.Get("/ping", func(c *ipa.Context) error {
		c.Logger.Info("pong")
		c.GinContext.JSON(200, "pong")
		return nil
	})

	app.Run()
}
