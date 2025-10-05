package botcreation

import (
	"coldheater/internal/database"
	"fmt"
	"math/rand"
	"time"

	"github.com/brianvoe/gofakeit/v7"
	"github.com/go-rod/rod"
	"github.com/go-rod/rod/lib/proto"
)

func CreateGmailBot(verificationAcc database.VerificationAccount) (bot *database.Bot, wasVerificationAccUsed bool, err error) {
	userBrowser, err := NewUserModBrowser(time.Second / 4)
	if err != nil {
		return nil, wasVerificationAccUsed, fmt.Errorf("Failed to create new user mode browser:\n%w", err)
	}

	browser, err := userBrowser.Incognito()
	if err != nil {
		return nil, wasVerificationAccUsed, fmt.Errorf("Failed to create incognito browser:\n%w", err)
	}

	bot = generateBotData()

	gmailSignUpPage, err := browser.Page(proto.TargetCreateTarget{})
	if err != nil {
		return nil, wasVerificationAccUsed, fmt.Errorf("Failed to create new browser page:\n%w", err)
	}

	err = rod.Try(func() {
		gmailSignUpPage.MustNavigate("https://accounts.google.com/signin").MustWindowFullscreen().MustWaitLoad()
	})
	if err != nil {
		return nil, wasVerificationAccUsed, fmt.Errorf("Failed to navigate to the signin page")
	}

	createAccElement, err := gmailSignUpPage.ElementR("span", "/^Create account$/")
	if err != nil {
		return nil, wasVerificationAccUsed, fmt.Errorf("Failed to select 'Create Account' element:\n%w", err)
	}

	err = createAccElement.Click(proto.InputMouseButtonLeft, 1)
	if err != nil {
		return nil, wasVerificationAccUsed, fmt.Errorf("Failed to click 'Create account element':\n%w", err)
	}

	forPresonalUseElement, err := gmailSignUpPage.ElementR("li", "/^For my personal use$/")
	if err != nil {
		return nil, wasVerificationAccUsed, fmt.Errorf("Failed to select 'For my personal use' element:\n%w", err)
	}

	err = forPresonalUseElement.Click(proto.InputMouseButtonLeft, 1)
	if err != nil {
		return nil, wasVerificationAccUsed, fmt.Errorf("Failed to click 'For personal use element'")
	}

	err = gmailSignUpPage.WaitLoad()
	if err != nil {
		return nil, wasVerificationAccUsed, fmt.Errorf("Failed to navigate to account creation page: %w", err)
	}

	err = rod.Try(func() {
		gmailSignUpPage.MustElement("#firstName").MustInput(bot.FirstName)
		gmailSignUpPage.MustElement("#lastName").MustInput(bot.LastName)

		gmailSignUpPage.MustElementR("span", "/^Next$/").MustClick()
		gmailSignUpPage.MustWaitLoad()
	})
	if err != nil {
		return nil, wasVerificationAccUsed, fmt.Errorf("Failed to pass 'Enter your name' page: %w", err)
	}

	err = rod.Try(func() {
		gmailSignUpPage.MustElement("#month").MustClick().
			MustElementR("li", fmt.Sprintf("/^%s$/", bot.BirthDate.Month().String())).MustClick()
		gmailSignUpPage.MustElement("#day").MustInput(fmt.Sprintf("%d", bot.BirthDate.Day()))
		gmailSignUpPage.MustElement("#year").MustInput(fmt.Sprintf("%d", bot.BirthDate.Year()))
		gmailSignUpPage.MustElementR("div", "/^Gender$/").MustClick().
			MustElementR("li", fmt.Sprintf("/^%s$/i", bot.Gender)).MustClick()

		gmailSignUpPage.MustElementR("span", "/^Next$/").MustClick()
		gmailSignUpPage.MustWaitLoad()
	})
	if err != nil {
		return nil, wasVerificationAccUsed, fmt.Errorf("Failed to pass 'Basic information' page: %w", err)
	}

	err = rod.Try(func() {
		gmailSignUpPage.MustElement("#emailPhone").MustInput(verificationAcc.Email)
		gmailSignUpPage.MustElementR("span", "/^Next$/").MustClick()
		gmailSignUpPage.MustWaitLoad()
	})
	if err != nil {
		return nil, wasVerificationAccUsed, fmt.Errorf("Failed to pass 'Add phone or email' page: %w", err)
	}

	_, err = gmailSignUpPage.Race().
		ElementR("span", "/^Verify your email address$/").Handle(func(e *rod.Element) error {
		_, err := GetVerificationCode(browser, verificationAcc)
		if err != nil {
			return fmt.Errorf("Failed to get verification code from account %s: %w", verificationAcc.Email, err)
		}

		wasVerificationAccUsed = true

		return nil
	}).
		ElementR("div", "/That username is taken. Try another./").Handle(func(e *rod.Element) error {
		err := rod.Try(func() {
			// more specific selector?
			gmailSignUpPage.MustElementR("button", "/an email address or phone number/").MustClick()
			gmailSignUpPage.MustWaitLoad()
		})
		if err != nil {
			return fmt.Errorf("Failed to pass 'use email or phone number page:\n%w", err)
		}

		return nil
	}).
		Do()
	if err != nil {
		return nil, wasVerificationAccUsed, fmt.Errorf("Race failed:\n%w", err)
	}

	var availableEmail string

	err = rod.Try(func() {
		const emailRegex string = `/^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$/`

		emailElement := gmailSignUpPage.MustElementR("div", emailRegex)

		availableEmail = emailElement.MustText()
		emailElement.MustClick()

		gmailSignUpPage.MustElementR("span", "/^Next$/").MustClick().MustWaitLoad()
	})
	if err != nil {
		return nil, wasVerificationAccUsed, fmt.Errorf("Failed to get available email:\n%w", err)
	}
	bot.Email = availableEmail

	err = rod.Try(func() {
		gmailSignUpPage.MustElement(`input[aria-label="Password"]`).MustInput(bot.Password)
		gmailSignUpPage.MustElement(`input[aria-label="Confirm"]`).MustInput(bot.Password)
		gmailSignUpPage.MustElementR("span", "/^Next$/").MustClick().MustWaitLoad()
	})
	if err != nil {
		return nil, wasVerificationAccUsed, fmt.Errorf("Failed to pass 'create password' page:\n%w", err)
	}

	time.Sleep(time.Minute)

	// temp nil, in future return bot
	return nil, wasVerificationAccUsed, nil
}

func GetVerificationCode(browser *rod.Browser, verificationAcc database.VerificationAccount) (verificationCode *string, err error) {
	const gmailNextElement string = "/^Next$/"

	verificationCodePage, err := browser.Page(proto.TargetCreateTarget{})
	if err != nil {
		return nil, fmt.Errorf("Failed to create verificationCodePage:\n%w", err)
	}

	err = rod.Try(func() {
		verificationCodePage.MustNavigate("https://accounts.google.com/signin").MustWaitLoad()
	})
	if err != nil {
		return nil, fmt.Errorf("Failed to navigate to verificationCodePage:\n%w", err)
	}

	err = rod.Try(func() {
		verificationCodePage.MustElement("#identifierId").MustInput(verificationAcc.Email)
		verificationCodePage.MustElementR("span", gmailNextElement).MustClick()
		verificationCodePage.MustWaitLoad()
	})
	if err != nil {
		return nil, fmt.Errorf("Failed to pass 'Sign in' page: %w", err)
	}

	err = rod.Try(func() {
		verificationCodePage.MustElement("input[aria-label=\"Enter your password\"]").MustInput(verificationAcc.Password)
		verificationCodePage.MustElementR("span", gmailNextElement).MustClick()
		verificationCodePage.MustWaitLoad()
	})
	if err != nil {
		return nil, fmt.Errorf("Failed to pass 'input password' page:\n%w", err)
	}

	var inboxPage *rod.Page

	err = rod.Try(func() {
		event := proto.TargetTargetCreated{}
		wait := browser.WaitEvent(&event)

		// might need beter handling, if 'products' window wont always be always open after login
		verificationCodePage.MustElementR("span", "/^Gmail$/").MustClick()

		wait()
		inboxPage = browser.MustPageFromTargetID(event.TargetInfo.TargetID)
	})
	if err != nil {
		return nil, fmt.Errorf("Failed to open inbox page:\n%w", err)
	}

	err = rod.Try(func() {
		_ = inboxPage.TargetID
	})

	return verificationCode, nil
}

func generateBotData() *database.Bot {
	person := gofakeit.Person()
	username := gofakeit.Username()
	password := gofakeit.Password(true, true, true, true, false, 20)
	birthDate := randomDate(1950, 2000)

	return &database.Bot{
		FirstName: person.FirstName,
		LastName:  person.LastName,
		Username:  username,
		Password:  password,
		BirthDate: birthDate,
		Gender:    person.Gender,
	}
}

func randomDate(startYear, endYear int) time.Time {
	// Random year in range (inclusive)
	year := rand.Intn(endYear-startYear+1) + startYear

	// Random month (1-12)
	month := rand.Intn(12) + 1

	// Random day (1-28 to avoid month/leap year issues)
	day := rand.Intn(28) + 1

	return time.Date(year, time.Month(month), day, 0, 0, 0, 0, time.UTC)
}
