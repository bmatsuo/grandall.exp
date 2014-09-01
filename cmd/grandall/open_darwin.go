//+build: darwin

package main

import (
	"fmt"
	"log"
	"os/exec"
)

func OpenURL(url string) error {
	log.Printf("open %q", url)
	err := exec.Command("open", url).Run()
	if err != nil {
		return fmt.Errorf("open (%v)", err)
	}
	return nil
}
