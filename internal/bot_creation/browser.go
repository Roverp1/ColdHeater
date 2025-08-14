package botcreation

import (
	"github.com/go-rod/rod"
	"github.com/go-rod/rod/lib/launcher"
	"github.com/go-rod/stealth"
)

type StealthBrowser struct {
	Browser *rod.Browser
}

func NewStealthBrowser(headless bool) *StealthBrowser {
	l := launcher.New().
		Headless(headless).
		MustLaunch()

	browser := rod.New().ControlURL(l).NoDefaultDevice().MustConnect()

	return &StealthBrowser{
		Browser: browser,
	}
}

func (sb *StealthBrowser) NewStealthPage() *rod.Page {
	page := stealth.MustPage(sb.Browser)

	return page
}
