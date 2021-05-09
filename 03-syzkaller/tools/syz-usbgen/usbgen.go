// Copyright 2019 syzkaller project authors. All rights reserved.
// Use of this source code is governed by Apache 2 LICENSE that can be found in the LICENSE file.

package main

import (
	"encoding/hex"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"regexp"
	"sort"

	"github.com/ZHYfeng/2018-Dependency/03-syzkaller/pkg/osutil"
)

func main() {
	flag.Parse()
	args := flag.Args()
	if len(args) != 2 {
		usage()
	}

	syslog, err := ioutil.ReadFile(args[0])
	if err != nil {
		failf("failed to read file %v: %v", args[0], err)
	}

	r := regexp.MustCompile(`usbID: [0-9a-f]{48}`)
	matches := r.FindAll(syslog, -1)
	uniqueMatches := make(map[string]bool)
	for _, match := range matches {
		uniqueMatches[string(match)] = true
	}
	sortedMatches := make([]string, 0)
	for match := range uniqueMatches {
		match = match[len("usbID: "):]
		match = match[:34]
		sortedMatches = append(sortedMatches, match)
	}
	sort.Strings(sortedMatches)

	usbIDs := make([]byte, 0)
	usbIDs = append(usbIDs, []byte("// AUTOGENERATED FILE\n")...)
	usbIDs = append(usbIDs, []byte("// See docs/linux/external_fuzzing_usb.md\n")...)
	usbIDs = append(usbIDs, []byte("\n")...)
	usbIDs = append(usbIDs, []byte("package linux\n")...)
	usbIDs = append(usbIDs, []byte("\n")...)
	usbIDs = append(usbIDs, []byte("var usbIDs = ")...)
	for i, match := range sortedMatches {
		decodedMatch, err := hex.DecodeString(match)
		if err != nil {
			failf("failed to decode hes string %v: %v", match, err)
		}
		prefix := "\t"
		suffix := " +"
		if i == 0 {
			prefix = ""
		}
		if i == len(sortedMatches)-1 {
			suffix = ""
		}
		usbID := fmt.Sprintf("%v%#v%v\n", prefix, string(decodedMatch), suffix)
		usbIDs = append(usbIDs, []byte(usbID)...)
	}

	if err := osutil.WriteFile(args[1], usbIDs); err != nil {
		failf("failed to output file %v: %v", args[1], err)
	}

	fmt.Printf("%v ids written\n", len(sortedMatches))
}

func usage() {
	fmt.Fprintf(os.Stderr, "usage:\n")
	fmt.Fprintf(os.Stderr, "  syz-usbgen syslog.txt sys/linux/init_vusb_ids.go\n")
	os.Exit(1)
}

func failf(msg string, args ...interface{}) {
	fmt.Fprintf(os.Stderr, msg+"\n", args...)
	os.Exit(1)
}
