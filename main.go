/******************************************************************************
	main.go
	entrypoint for the go process
	all other functions are	initiated from this main function.
******************************************************************************/
package main

// main - main process starting point
func main() {
	if cmdLnEmailAddress != "" && cmdLnEmailFrom != "" && cmdLnEmailHost != "" {
		stockMonitor()
	} else {
		stockCurrent()
	}
}
