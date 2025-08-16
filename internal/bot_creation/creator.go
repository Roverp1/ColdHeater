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

func CreateGmailBot() error {
	browser := NewStealthBrowser(false, time.Second/2)

	bot := generateBotData()
	fmt.Println("birthDate:", bot.BirthDate.Format("2006-01-02"))

	page := browser.NewStealthPage()
	err := page.MustNavigate("https://accounts.google.com/signin").MustWindowFullscreen().WaitLoad()
	if err != nil {
		return fmt.Errorf("Failed to navigate to the signin page")
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

	err = page.WaitLoad()
	if err != nil {
		return fmt.Errorf("Failed to navigate to account creation page: %w", err)
	}

	err = rod.Try(func() {
		page.MustElement("#firstName").MustInput(bot.FirstName)
		page.MustElement("#lastName").MustInput(bot.LastName)

		page.MustElementR("span", "/^Next$/").MustClick()
		page.MustWaitLoad()
	})
	if err != nil {
		return fmt.Errorf("Failed to pass 'Enter your name' page: %w", err)
	}

	err = rod.Try(func() {
		page.MustElement("#month").MustClick().
			MustElementR("li", fmt.Sprintf("/^%s$/", bot.BirthDate.Month().String())).MustClick()
		page.MustElement("#day").MustInput(fmt.Sprintf("%d", bot.BirthDate.Day()))
		page.MustElement("#year").MustInput(fmt.Sprintf("%d", bot.BirthDate.Year()))
		page.MustElementR("div", "/^Gender$/").MustClick().
			MustElementR("li", fmt.Sprintf("/^%s$/i", bot.Gender)).MustClick()

		page.MustElementR("span", "/^Next$/").MustClick()
		page.MustWaitLoad()
	})
	if err != nil {
		return fmt.Errorf("Failed to pass 'Basic information' page: %w", err)
	}

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
