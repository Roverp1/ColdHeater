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

func CreateGmailBot(verificationAcc database.VerificationAccount) error {
	browser := NewStealthBrowser(false, time.Second/4)

	bot := generateBotData()

	gmailSignUpPage := browser.NewStealthPage()
	err := gmailSignUpPage.MustNavigate("https://accounts.google.com/signin").MustWindowFullscreen().WaitLoad()
	if err != nil {
		return fmt.Errorf("Failed to navigate to the signin page")
	}

	createAccElement, err := gmailSignUpPage.ElementR("span", "/^Create account$/")
	if err != nil {
		return fmt.Errorf("Failed to select 'Create Account' element:\n%w", err)
	}

	err = createAccElement.Click(proto.InputMouseButtonLeft, 1)
	if err != nil {
		return fmt.Errorf("Failed to click 'Create account element'")
	}

	forPresonalUseElement, err := gmailSignUpPage.ElementR("li", "/^For my personal use$/")
	if err != nil {
		return fmt.Errorf("Failed to select 'For my personal use' element:\n%w", err)
	}

	err = forPresonalUseElement.Click(proto.InputMouseButtonLeft, 1)
	if err != nil {
		return fmt.Errorf("Failed to click 'For personal use element'")
	}

	err = gmailSignUpPage.WaitLoad()
	if err != nil {
		return fmt.Errorf("Failed to navigate to account creation page: %w", err)
	}

	err = rod.Try(func() {
		gmailSignUpPage.MustElement("#firstName").MustInput(bot.FirstName)
		gmailSignUpPage.MustElement("#lastName").MustInput(bot.LastName)

		gmailSignUpPage.MustElementR("span", "/^Next$/").MustClick()
		gmailSignUpPage.MustWaitLoad()
	})
	if err != nil {
		return fmt.Errorf("Failed to pass 'Enter your name' page: %w", err)
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
		return fmt.Errorf("Failed to pass 'Basic information' page: %w", err)
	}

	err = rod.Try(func() {
		gmailSignUpPage.MustElement("#emailPhone").MustInput(verificationAcc.Email)
		gmailSignUpPage.MustElementR("span", "/^Next$/").MustClick()
		gmailSignUpPage.MustWaitLoad()
	})
	if err != nil {
		return fmt.Errorf("Failed to pass 'Add phone or email' page: %w", err)
	}

	verificationEmailPage := browser.NewStealthPage()
	verificationEmailPage.Navigate("https://account.proton.me/login")

	time.Sleep(time.Hour)

	return nil
}

func generateBotData() database.Bot {
	person := gofakeit.Person()
	username := gofakeit.Username()
	password := gofakeit.Password(true, true, true, true, false, 20)
	birthDate := randomDate(1950, 2000)

	return database.Bot{
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
