package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"syscall"
	"time"

	"github.com/getlantern/systray"
)

const diskIcon = "⛁"

type options struct {
	color    bool
	warn     float64
	critical float64
}

func main() {
	var opts options
	flag.BoolVar(&opts.color, "color", false, "show yellow and red low-space indicators")
	flag.Float64Var(&opts.warn, "warn", 20, "free-space percentage at or below which to show yellow")
	flag.Float64Var(&opts.critical, "err", 10, "free-space percentage at or below which to show red")
	flag.Parse()
	if flag.NArg() != 0 || !opts.validThresholds() {
		fmt.Fprintln(os.Stderr, "Use flags only, with thresholds satisfying 0 <= -err < -warn <= 100.")
		flag.Usage()
		os.Exit(2)
	}
	systray.Run(func() { onReady(opts) }, nil)
}

func (o options) validThresholds() bool {
	return o.critical >= 0 && o.critical < o.warn && o.warn <= 100
}

func (o options) icon(freePercent float64) string {
	switch {
	case o.color && freePercent <= o.critical:
		return "🔴"
	case o.color && freePercent <= o.warn:
		return "🟡"
	default:
		return diskIcon
	}
}

func readDiskSpace() (freeBytes uint64, freePercent float64, err error) {
	var stats syscall.Statfs_t
	if err := syscall.Statfs("/", &stats); err != nil {
		return 0, 0, err
	}
	if stats.Bsize == 0 || stats.Blocks == 0 || stats.Bavail > stats.Blocks {
		return 0, 0, fmt.Errorf("invalid capacity for startup filesystem")
	}
	return stats.Bavail * uint64(stats.Bsize), 100 * float64(stats.Bavail) / float64(stats.Blocks), nil
}

func updateFreeSpace(opts options) {
	freeBytes, freePercent, err := readDiskSpace()
	if err != nil {
		log.Printf("Reading free disk space: %v", err)
		systray.SetTitle(diskIcon + " —")
		systray.SetTooltip("Free disk space unavailable")
		return
	}
	freeGiB := float64(freeBytes) / (1 << 30)
	systray.SetTitle(fmt.Sprintf("%s %.2f GiB", opts.icon(freePercent), freeGiB))
	freeGB := float64(freeBytes) / 1e9
	systray.SetTooltip(fmt.Sprintf("%.0f GB, %.0f%% Free", freeGB, freePercent))
}

func onReady(opts options) {
	systray.SetTitle(diskIcon + " Init")
	mQuit := systray.AddMenuItem("Quit", "Quit the whole app")

	go func() {
		<-mQuit.ClickedCh
		systray.Quit()
	}()

	go func() {
		for {
			updateFreeSpace(opts)
			time.Sleep(10 * time.Second)
		}
	}()
}
