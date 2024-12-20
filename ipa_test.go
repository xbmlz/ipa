package ipa

import "testing"

func TestLoadConfig(t *testing.T) {
	app := App{}

	app.loadConfig()

	if app.Config == nil {
		t.Errorf("config was not read")
	}
}

func TestAddRoute(t *testing.T) {

}
