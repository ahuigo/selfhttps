package main

import (
	"bytes"
	"errors"
	"os"
	"os/exec"
)

func RunCommand(name string, args ...string) (out string, err error) {
	var stderr bytes.Buffer
	var stdout bytes.Buffer
	cmd := exec.Command(name, args...)
	cmd.Stderr = &stderr
	cmd.Stdout = &stdout
	if err := cmd.Run(); err != nil {
		out = stdout.String()
		if stderr.Len() > 0 {
			return out, errors.New(stderr.String())
		}
		return out, err
	}
	return
}

func CommandExists(cmd string) bool {
	_, err := exec.LookPath(cmd)
	return err == nil
}

func IsFileExists(file string) bool {
	_, err := os.Stat(file)
	if os.IsNotExist(err) {
		return false
	} else if err != nil {
		panic(err)
	}
	return err == nil
}
