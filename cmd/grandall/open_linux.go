//+build: linux

package main

import (
	"fmt"
	"os/exec"
)

// BUG: OpenURL on linux only knows how to use xdg-open.
func OpenURL(url string) error {
	err := exec.Command("xdg-open", url).Run()
	if err != nil {
		return fmt.Errorf("xdg-open (%v)", err)
	}
	return nil
}
