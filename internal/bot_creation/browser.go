package botcreation

import (
	"fmt"
	"time"

	"github.com/go-rod/rod"
	"github.com/go-rod/rod/lib/launcher"
	"github.com/go-rod/stealth"
)

type StealthBrowser struct {
	Browser *rod.Browser
}

func NewUserModBrowser(delay time.Duration) (*rod.Browser, error) {
	// var browserBin string = "/usr/bin/chromium"
	//
	// launcherURL, err := launcher.NewUserMode().Bin(browserBin).UserDataDir("tmp/rod-profile").Launch()
	// if err != nil {
	// 	return nil, fmt.Errorf("Failed to launch '%s' browser: %w", browserBin, err)
	// }
	browserPort := "ws://browser-template:9221"

	browser := rod.New().ControlURL(browserPort).SlowMotion(delay).NoDefaultDevice()
	err := browser.Connect()
	if err != nil {
		return nil, fmt.Errorf("Failed to connect to the browser on port '%s': %w", browserPort, err)
	}

	return browser, nil
}

func NewStealthBrowser(headless bool, delay time.Duration) *StealthBrowser {
	var launcherURL string

	launcherURL = launcher.New().
		Headless(headless).
		MustLaunch()

	browser := rod.New().ControlURL(launcherURL).SlowMotion(delay).NoDefaultDevice().MustConnect()

	return &StealthBrowser{
		Browser: browser,
	}
}

func (sb *StealthBrowser) NewStealthPage() *rod.Page {
	page := stealth.MustPage(sb.Browser)

	return page
}
