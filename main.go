package main

// main process starting point.
func main() {
	if cmdLnEmailAddress != "" && cmdLnEmailFrom != "" && cmdLnEmailHost != "" {
		stockMonitor()
	} else {
		stockCurrent()
	}
}
