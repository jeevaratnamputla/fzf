//go:build windows

package fzf

import (
	"io"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"sync/atomic"
)

var shPath atomic.Value

func sh(bash bool) (string, error) {
	if cached := shPath.Load(); cached != nil {
		return cached.(string), nil
	}

	name := "sh"
	if bash {
		name = "bash"
	}
	cmd := exec.Command("cygpath", "-w", "/usr/bin/"+name)
	bytes, err := cmd.Output()
	if err != nil {
		return "", err
	}

	sh := strings.TrimSpace(string(bytes))
	shPath.Store(sh)
	return sh, nil
}

func mkfifo(path string, mode uint32) (string, error) {
	m := strconv.FormatUint(uint64(mode), 8)
	cmd := exec.Command("mkfifo", "-m", m, path)
	if err := cmd.Run(); err != nil {
		return path, err
	}
	return path + ".lnk", nil
}

func withOutputPipe(output string, task func(io.ReadCloser)) error {
	f, err := os.Open(output)
	if err != nil {
		return err
	}
	defer f.Close()
	task(f)
	return nil
}

func withInputPipe(input string, task func(io.WriteCloser)) error {
	f, err := os.OpenFile(input, os.O_WRONLY, 0)
	if err != nil {
		return err
	}
	defer f.Close()
	task(f)
	return nil
}
