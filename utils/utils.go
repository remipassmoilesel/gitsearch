//go:generate mockgen -package mock -destination ../test/mock/mocks_Utils.go gitlab.com/remipassmoilesel/gitsearch/utils Utils
package utils

import (
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"strings"
)

type Utils interface {
	ContainsString(a []string, x string) bool
	OpenWebBrowser(url string) error
	Pager(content string) error
	PrintLn(content interface{})
}

func NewUtils() Utils {
	return &UtilsImpl{}
}

type UtilsImpl struct {
}

func (s *UtilsImpl) ContainsString(a []string, x string) bool {
	for _, n := range a {
		if x == n {
			return true
		}
	}
	return false
}

// See: https://gist.github.com/hyg/9c4afcd91fe24316cbf0
func (s *UtilsImpl) OpenWebBrowser(url string) error {
	var err error

	switch runtime.GOOS {
	case "linux":
		err = exec.Command("xdg-open", url).Start()
	case "windows":
		err = exec.Command("rundll32", "url.dll,FileProtocolHandler", url).Start()
	case "darwin":
		err = exec.Command("open", url).Start()
	default:
		err = fmt.Errorf("unsupported platform")
	}
	return err
}

func (s *UtilsImpl) Pager(content string) error {
	pager := "/usr/bin/less"
	if len(os.Getenv("PAGER")) > 0 {
		pager = os.Getenv("PAGER")
	}
	cmd := exec.Command(pager)
	cmd.Stdin = strings.NewReader(content)
	cmd.Stdout = os.Stdout
	return cmd.Run()
}

// For testing purposes
func (s *UtilsImpl) PrintLn(content interface{}) {
	fmt.Println(content)
}
