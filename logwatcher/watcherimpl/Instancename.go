package watcherimpl

import (
	"bytes"
	"os/exec"
	"strings"
)

func cmdLineOutBeautify(message string) string {
	message = strings.Replace(message, "\n", " ", -1)
	message = strings.Replace(message, "\r", " ", -1)
	return message
}

func GetThisInstanceName() string {
	cmd := exec.Command("bash", "-c", "hostname")
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		return ""
	}
	return cmdLineOutBeautify(string(stdout.Bytes()))
}
