//
// Copyright © 2017 Ikey Doherty <ikey@solus-project.com>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
//

// Package libpreston provides the core functionality for the Preston CLI binary,
// and is the main entry point into all functionality. It provides various scanners
// and methods to verify the conformance of a given package.
package libpreston

import (
	"fmt"
	"os"
	"path/filepath"
)

// TreeFunc is a simple callback type which takes a full path as a string to
// check and read.
type TreeFunc func(path string) error

// TreeScanner is used to scan a source tree and do useful Things with said tree.
// Predominantly this is used for scanning and discovering licenses within the tree.
type TreeScanner struct {
	BaseDir      string   // Base directory to actually scan
	ignoredPaths []string // List of paths to always ignore
}

// NewTreeScanner will return a scanner for the given directory
func NewTreeScanner(basedir string) *TreeScanner {
	return &TreeScanner{
		BaseDir: basedir,
		ignoredPaths: []string{
			".git", // Really no sense digging inside these
		},
	}
}

// Scan will do the grunt work of actually _scanning_ the tree
func (t *TreeScanner) Scan() error {
	return filepath.Walk(t.BaseDir, t.walker)
}

// isIgnored is currently stupidly simple, but will eventually support
// patterns.
func (t *TreeScanner) isIgnored(path string) bool {
	for _, p := range t.ignoredPaths {
		if path == p {
			return true
		}
	}
	return false
}

// walker handles each item in the tree-walk, skipping "special" paths
func (t *TreeScanner) walker(path string, info os.FileInfo, err error) error {
	if info == nil {
		return nil
	}
	if t.isIgnored(path) {
		if info.IsDir() {
			return filepath.SkipDir
		}
		return nil
	}
	if info.IsDir() {
		return nil
	}
	fmt.Fprintf(os.Stderr, "Path: %v\n", path)
	return nil
}
