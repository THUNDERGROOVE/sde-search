package main

import (
	"log"
	"os"
	"os/exec"
	"path/filepath"
)

func init() {
	if Version == "" || Branch == "" {
		log.Printf("Version info isn't set.  Must be on openshift.  Let's try to get it from the environment")
		os.Chdir(filepath.Join(os.Getenv("HOME"), "git", "sdesearch.git"))
		cmd := exec.Command("git", "rev-parse", "--short", "HEAD")
		ver, _ := cmd.CombinedOutput()
		Version = string(ver)
	}
}
