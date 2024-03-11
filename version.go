package main

import "runtime/debug"

var version string

func init() {
	if version == "" {
		if info, ok := debug.ReadBuildInfo(); ok {
			version = info.Main.Version
		} else {
			version = "(unknown)"
		}
	}
}
