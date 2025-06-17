package main

import (
	"bytes"
	"fmt"
	"os/exec"
	"strconv"
	"strings"
	"time"

	"github.com/getlantern/systray"
)

func main() {
	systray.Run(onReady, onExit)
}

func onReady() {
	systray.SetTitle("🔋EmcodeBattery")
	systray.SetTooltip("Battery Monitor")
	mQuit := systray.AddMenuItem("Quit", "Quit the app")

	go func() {
		for {
			checkBattery()
			time.Sleep(3 * time.Minute)
			//time.Sleep(20 * time.Second)
		}
	}()

	go func() {
		<-mQuit.ClickedCh
		systray.Quit()
	}()
}

func onExit() {
	// Cleanup if needed
}

func checkBattery() {
	lowBatteryVal, fullBatteryThresh := 35, 96

	out, err := exec.Command("pmset", "-g", "batt").Output()
	if err != nil {
		fmt.Println("Failed to get battery info:", err)
		return
	}

	batteryStr := string(out)
	percentStr := extractBatteryPercent(batteryStr)
	if percentStr == "" {
		fmt.Println("Could not parse battery percentage.")
		return
	}

	percent, err := strconv.Atoi(percentStr)
	if err != nil {
		fmt.Println("Invalid battery percentage:", err)
		return
	}

	if percent <= lowBatteryVal {
		message := fmt.Sprintf("Battery is at %d%% \nPlease plug in power...", percent)
		showNotification("Low Battery Alert", message)
		// Play sound (custom file)

	}
	if percent >= fullBatteryThresh {
		message := fmt.Sprintf("Battery is at %d%% \nUnplug Power Source...", percent)
		showNotification("Full Battery Charge", message)
		// Play sound (custom file)

	}
}

func extractBatteryPercent(input string) string {
	lines := strings.Split(input, "\n")
	for _, line := range lines {
		if strings.Contains(line, "%") {
			parts := strings.Fields(line)
			for _, p := range parts {
				if strings.HasSuffix(p, "%;") || strings.HasSuffix(p, "%") {
					p = strings.TrimRight(p, "%;")
					return p
				}
			}
		}
	}
	return ""
}

func showNotification(title, message string) {
	script := fmt.Sprintf(`display notification "%s" with title "%s"`, message, title)
	cmd := exec.Command("osascript", "-e", script)
	var stderr bytes.Buffer
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		fmt.Println("Failed to show notification:", stderr.String())
	}
	times := 3
	for i := 0; i < times; i++ {
		exec.Command("afplay", "/System/Library/Sounds/Glass.aiff").Start()
		time.Sleep(2 * time.Second)
	}

}
