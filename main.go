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

	var modelPath, modelName string
	for _, path := range []string{
		filepath.Join(homeDir, ".llama.model"),
		filepath.Join(homeDir, "llama.model"),
	} {
		info, err := os.Lstat(path)
		if err == nil {
			modelPath = path
			if info.Mode()&os.ModeSymlink > 0 {
				dest, err := os.Readlink(path)
				if err != nil {
					continue
				}
				modelName = filepath.Base(dest)
			}
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
Please translate this text to simple, concise English: [`+text+`].
If the translated text includes less than 10 words, give at least 5 examples of how to use it in English.
`,
		"--prompt-cache", filepath.Join(cacheDir, "llm-translate-cache."+modelName),
		"--ctx-size", "0",
		"--color",
		"--mlock",
		"--log-disable",
	)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	ce(cmd.Run())
}
