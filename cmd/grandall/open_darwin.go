//+build: darwin

package main

import (
	"fmt"
	"log"
	"os/exec"
)

func OpenURL(url string) error {
	log.Printf("open %q", url)
	p, err := exec.Command("open", url).CombinedOutput()
	if err != nil {
		return fmt.Errorf("%q (%v)", p, err)
	}
	return nil
}
