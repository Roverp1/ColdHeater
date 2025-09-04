package config

import (
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
)

func GetBrowserBin() (browserBin string, err error) {
	
	return 
}

func lookPath() (found string, has bool) {
	list := map[string][]string{
		"darwin": {
			"-",
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
