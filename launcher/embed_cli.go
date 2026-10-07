//go:build !desktoponly

package main

import _ "embed"

//go:embed assets/opencode.exe
var opencodeBin []byte

//go:embed assets/opencode.json
var defaultConfig []byte

//go:embed assets/opencode.strata.json
var strataConfig []byte

//go:embed assets/claude.strata.settings.json
var claudeConfig []byte

//go:embed assets/claude.desktop.strata.json
var claudeDesktopConfig []byte
