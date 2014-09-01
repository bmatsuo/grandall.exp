//+build: linux

package main

import (
	"fmt"
	"os/exec"
)

// BUG: OpenURL on linux only knows how to use xdg-open.
func OpenURL(url string) error {
	p, err := exec.Command("xdg-open", url).CombinedOutput()
	if err != nil {
		return fmt.Errorf("%q (%v)", p, err)
	}
	return nil
}
