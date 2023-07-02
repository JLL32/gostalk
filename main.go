package main

import (
	"fmt"
	"os"
	"os/exec"
)

func main() {
  cmd := exec.Command("go", "run", os.Args[1])

  cmd.Stdout = os.Stdout
  cmd.Stderr = os.Stderr

  err := cmd.Run()

  if err != nil {
      fmt.Fprintln(os.Stderr, fmt.Sprint(err))
      return
  }
}
