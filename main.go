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
	_ = cacheDir

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
			} else {
				modelName = filepath.Base(modelPath)
			}
			break
		}
	}
	if modelPath == "" {
		slog.Warn("no model")
		return
	}
	pt("model: %s\n", modelName)

	text := strings.Join(os.Args[1:], " ")

	cmd := exec.Command(
		llamaCPPPath,
		"--model", modelPath,
		"--prompt", `[INST]
You master every language in the world.
And you are really good at teaching languages.
Please explain the following text in simple, easy to understand English: `+text+`.
Then explain the text in Chinese.
After explaining, please give at least 5 examples of how to use it in English.
[/INST]
`,
		//"--prompt-cache", filepath.Join(cacheDir, "llm-translate-cache."+modelName),
		"--ctx-size", "4096",
		"--color",
		"--mlock",
		"--log-disable",
		"-ngl", "20",
	)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	ce(cmd.Run())
}
