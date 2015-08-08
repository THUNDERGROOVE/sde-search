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
		Version = string(ver[:len(ver)-1]) // Last byte is garbage.
		os.Chdir(os.Getenv("OPENSHIFT_DATA_DIR"))
		Branch = "master" // No other branches will be deployed.  I don't even feel bad about this
	}
}
