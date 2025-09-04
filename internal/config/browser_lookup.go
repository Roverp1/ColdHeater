package config

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"

	"github.com/go-rod/rod/lib/launcher"
	"gopkg.in/yaml.v3"
)

type BrowserPath struct {
	path string `yaml:"browser_path"`
}

func GetBrowserBin() (string, error) {
	pathFile, err := os.ReadFile("configs/browser_path.yaml")
	if err != nil {
		if os.IsNotExist(err) {
			fmt.Println("It seems your browser path wasn't configurated before")
			return lookBrowser()
		} else {
			return "", fmt.Errorf("failed to read a browser path file:\n%v", err)
		}
	}

	var path BrowserPath
	err = yaml.Unmarshal(pathFile, &path)
	if err != nil {
		return "", fmt.Errorf("failed to parse browser path file:\n%v", err)
	}
	return path.path, nil
}

func lookBrowser() (browserBin string, err error) {
	
	browserBin, exists := lookPath()
	if exists {
		fmt.Printf("Automaticly found \"Brave\" browser:\n%s\n", browserBin)
		err = initBrowserPathFile(&browserBin)
		if err != nil {
			err = fmt.Errorf("browser path file was not initiated:\n%v", err)
		}
		return	
	} 

	browserBin, exists = launcher.LookPath()
	if exists {
		fmt.Printf("Automaticly found browser:\n%s\n", browserBin)
		fmt.Println("Attention: it's recommended to use \"Brave\" browser")
		err = initBrowserPathFile(&browserBin)
		if err != nil {
			err = fmt.Errorf("browser path file was not initiated:\n%v", err)
		}
		return
	}

	fmt.Println("Could not locate any installed browsers.")
	browserBin = initOwnPath()
	if browserBin != "aborted" {
		fmt.Printf("Entered path:\n%s\n", browserBin)		
		err = initBrowserPathFile(&browserBin)
		if err != nil {
			err = fmt.Errorf("browser path file was not initiated:\n%v", err)
		}
		return
	}
	fmt.Println("Aborted.")
	err = fmt.Errorf("aborted by user")
	return
}

func initBrowserPathFile(browserBin *string) (err error) {
	var input int
	for {
		options := [...]string{"Use and save as default", "Use once", "Use other browser (type own path)", "Abort"}

		for i, option := range options{
			fmt.Printf("%d. %s\n", i+1, option)
		}

		fmt.Scanln(&input)
		switch input{
		case 1:
			var path BrowserPath
			var data []byte
			path.path = *browserBin
			
			data, err = yaml.Marshal(&path)
			if err != nil {
				err = fmt.Errorf("failed to convert a browser path file as YAML:\n%v", err)
				return
			}

			err = os.WriteFile("configs/browser_path.yaml", data, 0644)
			if err != nil {
				err = fmt.Errorf("failed to save a browser path file:\n%v", err)
				return
			}

			return

		case 2:
			return
		case 3:
			i := initOwnPath()
			if i != "aborted" {
				fmt.Printf("Entered path:\n%s\n", i)
				*browserBin = i
			} else {
				fmt.Printf("Aborted. Current path:\n%s\n", *browserBin)
			}
		case 4:
			err = fmt.Errorf("aborted by user")
			return
		default:
			fmt.Println("incorrect input")
		}
	}
}

func initOwnPath()(browserBin string){
	var input string
	for {
		fmt.Println("Please, type the browser path to be used.")
		fmt.Println("Enter \"abort\" to abort.")
		fmt.Scanln(&input)
		switch input {
		case "abort":
			browserBin = "aborted"
			return
		default:
			var err error
			browserBin, err = exec.LookPath(input)
			if err == nil {
				return
			}
			fmt.Printf("Incorrect path:\n%v\n", err)
		}
	}
}

//scrapped from go-rod-lib
func lookPath() (found string, has bool) {
	list := map[string][]string{
		"darwin": {
			"-", 									//todo: find out where is brave bin on macos
		},
		"linux": {
			"brave-browser",
			"/usr/bin/brave-browser",
		},
		"openbsd": {
			"brave-browser",
		},
		"windows": append([]string{"chrome", "edge"}, expandWindowsExePaths(
			`BraveSoftware\Brave-Browser\Application\brave.exe`,
		)...),
	}[runtime.GOOS]

	for _, path := range list {
		var err error
		found, err = exec.LookPath(path)
		has = err == nil
		if has {
			break
		}
	}

	return
}

// _,,_,,_,,_
func expandWindowsExePaths(list ...string) []string {
	newList := []string{}
	for _, p := range list {
		newList = append(
			newList,
			filepath.Join(os.Getenv("ProgramFiles"), p),
			filepath.Join(os.Getenv("ProgramFiles(x86)"), p),
			filepath.Join(os.Getenv("LocalAppData"), p),
		)
	}

	return newList
}
