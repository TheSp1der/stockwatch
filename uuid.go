package main

import (
	"github.com/gofrs/uuid"
	"time"
)

// getUUID returns a version 5 UUID based on the X500 namespace and
// the current local time.
func getUUID() string {
	ns, _ := uuid.FromString("6ba7b814-9dad-11d1-80b4-00c04fd430c8")
	
	return uuid.NewV5(ns, time.Now().Local().Format(timeFormat)).String()
}