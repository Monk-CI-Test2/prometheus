// Copyright The Prometheus Authors
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package fileutil

import (
	"io"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/prometheus/prometheus/util/testutil"
)

func TestPreallocExtendTrunc(t *testing.T) {
	dir := testutil.NewTemporaryDirectory("test_prealloc_extend_trunc", t)
	defer dir.Close()

	filePath := filepath.Join(dir.Path(), "testfile")

	// 1. Create a file and write some data (10 bytes).
	f, err := os.OpenFile(filePath, os.O_CREATE|os.O_RDWR, 0o666)
	require.NoError(t, err)
	defer f.Close()

	data := []byte("0123456789")
	_, err = f.Write(data)
	require.NoError(t, err)

	// Verify initial size is 10 bytes.
	fi, err := f.Stat()
	require.NoError(t, err)
	require.Equal(t, int64(10), fi.Size())

	// 2. Call preallocExtendTrunc with size larger than file size (25 bytes).
	// This should extend the file to 25 bytes.
	err = preallocExtendTrunc(f, 25)
	require.NoError(t, err)

	fi, err = f.Stat()
	require.NoError(t, err)
	require.Equal(t, int64(25), fi.Size())

	// 3. Call preallocExtendTrunc with size smaller than file size (5 bytes).
	// This should NOT truncate the file. The file should remain 25 bytes.
	err = preallocExtendTrunc(f, 5)
	require.NoError(t, err)

	fi, err = f.Stat()
	require.NoError(t, err)
	require.Equal(t, int64(25), fi.Size())

	// 4. Verify that the seek offset is preserved.
	// Let's set offset to 4, call preallocExtendTrunc, and verify it is still 4.
	_, err = f.Seek(4, io.SeekStart)
	require.NoError(t, err)

	err = preallocExtendTrunc(f, 30)
	require.NoError(t, err)

	currOffset, err := f.Seek(0, io.SeekCurrent)
	require.NoError(t, err)
	require.Equal(t, int64(4), currOffset)

	fi, err = f.Stat()
	require.NoError(t, err)
	require.Equal(t, int64(30), fi.Size())
}

func TestPreallocateExtend(t *testing.T) {
	dir := testutil.NewTemporaryDirectory("test_preallocate_extend", t)
	defer dir.Close()

	filePath := filepath.Join(dir.Path(), "testfile")

	f, err := os.OpenFile(filePath, os.O_CREATE|os.O_RDWR, 0o666)
	require.NoError(t, err)
	defer f.Close()

	// Initial size 0.
	fi, err := f.Stat()
	require.NoError(t, err)
	require.Equal(t, int64(0), fi.Size())

	// Preallocate with size 50 and extendFile = true.
	err = Preallocate(f, 50, true)
	require.NoError(t, err)

	fi, err = f.Stat()
	require.NoError(t, err)
	require.Equal(t, int64(50), fi.Size())

	// Preallocate with size 20 (smaller) and extendFile = true.
	// It should NOT truncate the file.
	err = Preallocate(f, 20, true)
	require.NoError(t, err)

	fi, err = f.Stat()
	require.NoError(t, err)
	require.Equal(t, int64(50), fi.Size())
}
