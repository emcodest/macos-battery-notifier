// this needs a laucher agent to run as background
package main

import (
	"bytes"
	"fmt"
	"os/exec"
	"strconv"
	"strings"
)


func main() {
	lowBatteryVal, fullBatteryThresh := 32, 100
	// Step 1: Get battery info using pmset
	out, err := exec.Command("pmset", "-g", "batt").Output()
	if err != nil {
		fmt.Println("Failed to get battery info:", err)
		return
	}

	// Step 2: Extract battery percentage from output
	batteryStr := string(out)
	percentStr := extractBatteryPercent(batteryStr)
	fmt.Println("## battery %", percentStr)
	if percentStr == "" {
		fmt.Println("Could not parse battery percentage.")
		return
	}

	percent, err := strconv.Atoi(percentStr)
	if err != nil {
		fmt.Println("Invalid battery percentage:", err)
		return
	}

	// Step 3: If battery < 30%, show a macOS notification
	if percent <= lowBatteryVal {
		message := fmt.Sprintf("Battery is at %d%% \nPlease plug in power...", percent)
		showNotification("Low Battery Alert", message)
	}
	if percent >= fullBatteryThresh {
		message := fmt.Sprintf("Battery is at %d%% \nUnplug Power Source...", percent)
		showNotification("Full Battery Charge", message)
	}
}

func extractBatteryPercent(input string) string {
	// Looks for pattern like "58%;"
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
}
