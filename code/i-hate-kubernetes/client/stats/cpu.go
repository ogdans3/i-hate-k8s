package stats

import (
	"fmt"
	"os"
	"strconv"
	"strings"
)

func getProcessCPUStats(pid int) (utime, stime uint64, err error) {
	// Read the /proc/[pid]/stat file
	statFile := fmt.Sprintf("/proc/%d/stat", pid)
	data, err := os.ReadFile(statFile)
	if err != nil {
		return 0, 0, err
	}

	// Parse the content of /proc/[pid]/stat
	parts := strings.Fields(string(data))
	if len(parts) < 14 {
		return 0, 0, fmt.Errorf("invalid stat format")
	}

	// The 14th and 15th fields are user and system CPU time (utime, stime)
	utime, err = strconv.ParseUint(parts[13], 10, 64)
	if err != nil {
		return 0, 0, err
	}

	stime, err = strconv.ParseUint(parts[14], 10, 64)
	if err != nil {
		return 0, 0, err
	}

	return utime, stime, nil
}

var initialUtime, initialStime uint64

func GetCpu(pid int) *float64 {
	if initialUtime == 0 {
		initialUtime, initialStime, _ = getProcessCPUStats(pid)
	}
	currentUtime, currentStime, err := getProcessCPUStats(pid)

	if err != nil {
		return nil
	}
	// Calculate the difference in user and system CPU times
	utimeDelta := float64(currentUtime - initialUtime)
	stimeDelta := float64(currentStime - initialStime)
	totalTime := utimeDelta + stimeDelta

	// Get the number of CPU cores
	numCores, err := os.ReadFile("/proc/cpuinfo")
	if err != nil {
		fmt.Println("Error reading CPU info:", err)
		return nil
	}
	coreCount := strings.Count(string(numCores), "processor")

	// Read the total system uptime from /proc/uptime
	uptimeData, err := os.ReadFile("/proc/uptime")
	if err != nil {
		fmt.Println("Error reading system uptime:", err)
		return nil
	}
	uptimeFields := strings.Fields(string(uptimeData))
	systemUptime, err := strconv.ParseFloat(uptimeFields[0], 64)
	if err != nil {
		fmt.Println("Error parsing system uptime:", err)
		return nil
	}

	// Calculate the CPU usage as a percentage
	//TODO: Why do i multiply by 100_000 here? Some error in the code somewhere
	cpuUsage := (totalTime / (systemUptime * float64(coreCount))) * 100 * 100_000

	//fmt.Printf("CPU Usage of process %d: %.2f%% %.10f\n\n", pid, cpuUsage, cpuUsage)

	// Update the previous time values
	initialUtime, initialStime = currentUtime, currentStime
	return &cpuUsage

}
