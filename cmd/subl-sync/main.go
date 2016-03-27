package main

import (
	"os"
	"os/exec"
	"log"
	"fmt"
)

func main() {
	args := os.Args

	if len(args) < 2 {
		help()
	}

	operation := args[1]

	switch operation {
	case "-pull":
		pull(args[2:])
	case "-push":
		push(args[2:])
	case "-help":
		fallthrough
	default:
		help()
	}
}

func help() {
	fmt.Println("Sync the sublime text preferences and packages of this computer with a git repo.")
	fmt.Println("Set the path to your sublime preferences: export SUBLIME_TEXT_PACKAGES=path/to/my/pacakges")
	fmt.Println("Run with: subl-sync [-pull | -push] [-f]")
	os.Exit(0)
}

func push(options []string) {
	SUBLIME_TEXT_PACKAGES := os.Getenv("SUBLIME_TEXT_PACKAGES")

	// first stage the files
	stageCmd := exec.Command("git", "add", "-A")
	stageCmd.Dir = SUBLIME_TEXT_PACKAGES

	out, err := stageCmd.Output()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("%s", out)

	// then commit the files
	commitCmd := exec.Command("git", "commit", "-m", "Pushing up latest package and settings changes.")
	commitCmd.Dir = SUBLIME_TEXT_PACKAGES

	out, err = commitCmd.Output()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("%s", out)

	// then push the files
	pushCmd := exec.Command("git", "push", "origin", "master")
	pushCmd.Dir = SUBLIME_TEXT_PACKAGES

	out, err = pushCmd.Output()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("%s", out)

	os.Exit(0)
}

func pull(options []string) {
	SUBLIME_TEXT_PACKAGES := os.Getenv("SUBLIME_TEXT_PACKAGES")

	// first stash current changes
	stashCmd := exec.Command("git", "stash")
	stashCmd.Dir = SUBLIME_TEXT_PACKAGES

	out, err := stashCmd.Output()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("%s", out)

	// then git pull
	pullCmd := exec.Command("git", "pull", "origin", "master")
	pullCmd.Dir = SUBLIME_TEXT_PACKAGES

	out, err = pullCmd.Output()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("%s", out)

	// dont reapply stash if -f
	if len(options) > 0 && options[0] == "-f" {
		// NO-OP
	} else {
		stashCmd := exec.Command("git", "pull", "origin", "master")
		stashCmd.Dir = SUBLIME_TEXT_PACKAGES

		out, err = stashCmd.Output()
		if err != nil {
			log.Fatal(err)
		}
		fmt.Printf("%s", out)
	}

	os.Exit(0)
}

