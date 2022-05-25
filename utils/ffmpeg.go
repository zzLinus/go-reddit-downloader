package utils

import (
	"log"
	"os"
	"os/exec"
)

func MergeAudioVideo(cpath []string, dpath string) error {
	flags := []string{
		"-y",
	}

	for _, p := range cpath {
		flags = append(flags, "-i", p)
	}

	flags = append(flags, "-c:v", "copy", "-c:a", "copy", dpath)
	cmd := exec.Command("ffmpeg", flags...)
	err := cmd.Run()
	if err != nil {
		log.Fatalf("can't execute command:%s", cmd.String())
		return err
	}
	for _, p := range cpath {
		os.Remove(p)
	}

	return nil
}
