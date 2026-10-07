//go:build desktoponly

package main

import _ "embed"

// The desktop launcher does not carry the CLI binary; the desktop app has its own backend.
var opencodeBin []byte

//go:embed assets/opencode.json
var defaultConfig []byte

//go:embed assets/opencode.strata.json
var strataConfig []byte

//go:embed assets/claude.strata.settings.json
var claudeConfig []byte

//go:embed assets/claude.desktop.strata.json
var claudeDesktopConfig []byte
