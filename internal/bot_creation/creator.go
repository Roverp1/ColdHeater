package botcreation

import (
	"fmt"
	"time"

	"github.com/go-rod/rod/lib/proto"
)

func CreateGmailBot() error {
	browser := NewStealthBrowser(false, time.Second)

	page := browser.NewStealthPage()
	err := page.MustNavigate("https://accounts.google.com/signin").WaitLoad()
	if err != nil {
		return fmt.Errorf("Failed to navigate to the account cretion page")
	}

	createAccElement, err := page.ElementR("span", "/^Create account$/")
	if err != nil {
		return fmt.Errorf("Failed to select 'Create Account' element:\n%w", err)
	}

	err = createAccElement.Click(proto.InputMouseButtonLeft, 1)
	if err != nil {
		return fmt.Errorf("Failed to click 'Create account element'")
	}

	forPresonalUseElement, err := page.ElementR("li", "/^For my personal use$/")
	if err != nil {
		return fmt.Errorf("Failed to select 'For my personal use' element:\n%w", err)
	}

	err = forPresonalUseElement.Click(proto.InputMouseButtonLeft, 1)
	if err != nil {
		return fmt.Errorf("Failed to click 'For personal use element'")
	}

	time.Sleep(time.Hour)

	return nil
}
