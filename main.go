package main

import (
	"fmt"
	"os/exec"
	"regexp"
	"time"

	"github.com/getlantern/systray"
)

func main() {
	systray.Run(onReady, onExit)
}

func formatBytesToGb(bytes int64) string {
	return fmt.Sprintf("%.2f", float64(bytes)/1024/1024/1024)
}

var diskEmoji = "⛁"

func updateFreeSpace() {
	process := exec.Command("diskutil", "info", "/")
	process.Wait()
	output, err := process.CombinedOutput()

	if err != nil {
		fmt.Printf("error running diskutil: %v\nOutput: %s", err, string(output))
		return
	}

	re := regexp.MustCompile(`Container Free Space:\s+[\d.]+\s+GB\s+\((\d+)\s+Bytes\)`)
	matches := re.FindStringSubmatch(string(output))
	if len(matches) < 2 {
		fmt.Println("could not find free space in diskutil output")
		return
	}

	var freeBytes int64
	_, err = fmt.Sscanf(matches[1], "%d", &freeBytes)
	if err != nil {
		fmt.Printf("error parsing free space: %v", err)
		return
	}

	freeGB := formatBytesToGb(freeBytes)
	systray.SetTitle(diskEmoji + " " + freeGB + " GB")
	tooltip := fmt.Sprintf("Free disk space is %s GB (%d Bytes)", freeGB, freeBytes)
	systray.SetTooltip(tooltip)

}

func onReady() {

	systray.SetTitle(diskEmoji + " Init")
	mQuit := systray.AddMenuItem("Quit", "Quit the whole app")

	go func() {
		<-mQuit.ClickedCh
		systray.Quit()
	}()

	go func() {
		ticker := time.NewTicker(10 * time.Second)
		defer ticker.Stop()

		for range ticker.C {
			updateFreeSpace()
		}
	}()

	updateFreeSpace()
}

func onExit() {
	// clean up here
}
