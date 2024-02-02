package main

import (
	"log/slog"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

func main() {
	var llamaCPPPath string
	for _, cmd := range []string{
		"llama.cpp",
		"llama",
	} {
		path, err := exec.LookPath(cmd)
		if err == nil {
			llamaCPPPath = path
			break
		}
	}
	if llamaCPPPath == "" {
		slog.Warn("llama.cpp not found")
		return
	}

	homeDir, err := os.UserHomeDir()
	ce(err)
	cacheDir, err := os.UserCacheDir()
	ce(err)

	var modelPath string
	for _, path := range []string{
		filepath.Join(homeDir, ".llama.model"),
		filepath.Join(homeDir, "llama.model"),
	} {
		_, err := os.Stat(path)
		if err == nil {
			modelPath = path
			break
		}
	}
	if modelPath == "" {
		slog.Warn("no model")
		return
	}

	text := strings.Join(os.Args[1:], " ")

	cmd := exec.Command(
		llamaCPPPath,
		"--model", modelPath,
		"--prompt", `
Please translate this text to simple and basic English: [`+text+`].
If the text includes less than 10 words, give at least 5 examples of how to use it.
`,
		"--prompt-cache", filepath.Join(cacheDir, "llm-translate-cache"),
		"--ctx-size", "0",
		"--color",
		"--mlock",
		"--log-disable",
	)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	ce(cmd.Run())
}
