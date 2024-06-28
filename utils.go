package main

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"strings"
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

func runCmdWithConfirm(prompt string, cmd string, silent bool) {
	fmt.Printf("%s:\n \033[32m %s \033[0m\n", prompt, cmd)
	if !silent {
		yn := StringPrompt("Whether to execute the above command?(yes/no/y/n, default:no):")
		yn = strings.ToLower(yn)
		if yn == "y" || yn == "yes" {
			fmt.Printf("Sudo ")
			out, err := RunCommand("sh", "-c", cmd)
			if err != nil {
				fmt.Printf("failed to execute cmd(%s), err: %v, stdout: %s\n\n", cmd, err, out)
				// os.Exit(0)
			} else {
				fmt.Printf("execution succeed!\n")
			}
		}
	}
	fmt.Printf("\n")
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

func StringPrompt(label string) string {
	var s string
	r := bufio.NewReader(os.Stdin)
	for {
		fmt.Fprint(os.Stderr, label+" ")
		s, _ = r.ReadString('\n')
		if s != "" {
			break
		}
	}
	return strings.TrimSpace(s)
}

func isMacOsx() bool{
    return runtime.GOOS == "darwin"
}
func isWindowsOs() bool{
    return runtime.GOOS == "windows"
}
