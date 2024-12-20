package main

import "github.com/xbmlz/ipa"

type User struct {
	ID       int    `json:"id"`
	Username string `json:"username"`
	Password string `json:"password"`
}

func main() {
	app := ipa.New()

	app.Get("/ping", func(c *ipa.Context) error {
		// log
		c.Log.Info("pong")

		// TODO db
		// user, err := c.DB.Find(&User{})

		// TODO redis
		// res, err := c.Redis.Get("key")

		c.JSON(200, "pong")
		return nil
	})

	app.Run()
}
