// Copyright (c) 2023 Silverton Data, Inc.
// You may use, distribute, and modify this code under the terms of the Apache-2.0 license, a copy of
// which may be found at https://github.com/silverton-io/buz/blob/main/LICENSE

package envelope

import (
	"strconv"
	"strings"
)

type OutputLocation struct {
	Path        string
	DatabaseFqn string
	Namespace   string
}

// Build an OutputLocation using Envelope contents
func NewOutputLocationFromEnvelope(e *Envelope) OutputLocation {
	// Example output path: isValid=true/vendor=io.silverton/namespace=gettingStarted.example/version=1.1/year=2023/month=3/day=10/
	year, month, day := e.BuzTimestamp.Date()
	datePath := "/year=" + strconv.Itoa(year) + "/month=" + strconv.Itoa(int(month)) + "/day=" + strconv.Itoa(day) + "/"
	outputPath := "isValid=" + strconv.FormatBool(e.IsValid) + "/vendor=" + e.Vendor + "/namespace=" + e.Namespace + "/version=" + e.Version + datePath
	// Example output database fqn: io_silverton.gettingstarted_example_1
	outputDbFqn := strings.Replace(e.Vendor, ".", "_", -1) + "." + strings.Replace(e.Namespace, ".", "_", -1) + "_" + strings.Split(e.Version, ".")[0]
	// Example output namespace: io.silverton.gettingStarted.example.1.1
	outputNamespace := strings.Join([]string{e.Vendor, e.Namespace, e.Version}, ".")
	return OutputLocation{
		Path:        outputPath,
		DatabaseFqn: outputDbFqn,
		Namespace:   outputNamespace,
	}
}
