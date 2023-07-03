package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"syscall"
	"time"
)

func main() {

	// get all go source files recursively
	// capture each file's stat in a map with the absolute path as the key
	// iterate over the files and compare each stat with its older version in the map
	// if they don't match rerun the command

	cmd := startServer()
	defer killServer(cmd)

	initialStats := make(map[string]os.FileInfo, 0)

	for {
		changed := false
		err := filepath.Walk(".",
			func(path string, info os.FileInfo, err error) error {
				if err != nil {
					return err
				}

				if filepath.Ext(path) != ".go" {
					return nil
				}

				actualStat, err := os.Stat(path)
				if err != nil {
					return err
				}

				initialStat, ok := initialStats[path]
				if !ok {
					initialStats[path] = actualStat
					initialStat = initialStats[path]
				}

				if actualStat.Size() != initialStat.Size() || actualStat.ModTime() != initialStat.ModTime() {
					changed = true
				}

				initialStats[path] = actualStat
				return nil
			})

		if changed {
			fmt.Println("killing the previous instance")
			killServer(cmd)

			fmt.Println("rerunning the server")
			cmd = startServer()
		}

		if err != nil {
			log.Fatal(err)
		}

		time.Sleep(time.Second * 1)
	}
}

func startServer() *exec.Cmd {
	cmd := exec.Command("go", "run", os.Args[1])
	cmd.SysProcAttr = &syscall.SysProcAttr{Setsid: true}
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	err := cmd.Start()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(cmd.Process.Pid)
	return cmd
}

func killServer(cmd *exec.Cmd) {
	if err := syscall.Kill(-cmd.Process.Pid, syscall.SIGKILL); err != nil {
		log.Println("failed to kill: ", err)
	}
}
