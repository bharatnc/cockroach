// Copyright 2020 The Cockroach Authors.
//
// Use of this software is governed by the Business Source License
// included in the file licenses/BSL.txt.
//
// As of the Change Date specified in that file, in accordance with
// the Business Source License, use of this software will be governed
// by the Apache License, Version 2.0, included in the file
// licenses/APL.txt.

package colcontainerutils

import (
	"testing"

	"github.com/cockroachdb/cockroach/pkg/base"
	"github.com/cockroachdb/cockroach/pkg/sql/colcontainer"
	"github.com/cockroachdb/cockroach/pkg/storage/engine"
	"github.com/cockroachdb/cockroach/pkg/storage/engine/fs"
	"github.com/cockroachdb/cockroach/pkg/testutils"
)

const inMemDirName = "testing"

// NewTestingDiskQueueCfg returns a DiskQueueCfg and a non-nil cleanup function.
func NewTestingDiskQueueCfg(t testing.TB, inMem bool) (colcontainer.DiskQueueCfg, func()) {
	t.Helper()

	var (
		cfg       colcontainer.DiskQueueCfg
		cleanup   func()
		testingFS fs.FS
		path      string
	)

	if inMem {
		ngn := engine.NewDefaultInMem()
		testingFS = ngn.(fs.FS)
		if err := testingFS.CreateDir(inMemDirName); err != nil {
			t.Fatal(err)
		}
		path = inMemDirName
		cleanup = ngn.Close
	} else {
		ngn, err := engine.NewDefaultEngine(0 /* cacheSize */, base.StorageConfig{})
		if err != nil {
			t.Fatal(err)
		}
		testingFS = ngn.(fs.FS)
		tempPath, dirCleanup := testutils.TempDir(t)
		path = tempPath
		cleanup = func() {
			ngn.Close()
			dirCleanup()
		}
	}
	cfg.FS = testingFS
	cfg.Path = path

	if err := cfg.EnsureDefaults(); err != nil {
		t.Fatal(err)
	}

	return cfg, cleanup
}
