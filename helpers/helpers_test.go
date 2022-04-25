/*
Copyright 2022 Equinix, Inc.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package helpers_test

import (
	"os"
	"path/filepath"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"github.com/metal-toolbox/auditevent/helpers"
)

func TestOpenAuditLogFileUntilSuccess(t *testing.T) {
	t.Parallel()

	var wg sync.WaitGroup
	wg.Add(1)

	tmpdir := t.TempDir()
	tmpfile := filepath.Join(tmpdir, "audit.log")

	go func() {
		defer wg.Done()
		time.Sleep(time.Second)
		fd, err := os.OpenFile(tmpfile, os.O_RDONLY|os.O_CREATE, 0o600)
		require.NoError(t, err)
		err = fd.Close()
		require.NoError(t, err)
	}()

	fd, err := helpers.OpenAuditLogFileUntilSuccess(tmpfile)
	require.NoError(t, err)
	require.NotNil(t, fd)

	err = fd.Close()
	require.NoError(t, err)

	// We wait so we don't leak file descriptors
	wg.Wait()

	err = os.Remove(tmpfile)
	require.NoError(t, err)
}

func TestOpenAuditLogFileError(t *testing.T) {
	t.Parallel()

	var wg sync.WaitGroup
	wg.Add(1)

	tmpdir := t.TempDir()
	tmpfile := filepath.Join(tmpdir, "audit.log")

	go func() {
		defer wg.Done()
		time.Sleep(time.Second)
		// This file is read only
		fd, err := os.OpenFile(tmpfile, os.O_RDONLY|os.O_CREATE, 0o500)
		require.NoError(t, err)
		err = fd.Close()
		require.NoError(t, err)
	}()

	fd, err := helpers.OpenAuditLogFileUntilSuccess(tmpfile)
	require.Error(t, err)
	require.Nil(t, fd)

	// We wait so we don't leak file descriptors
	wg.Wait()

	err = os.Remove(tmpfile)
	require.NoError(t, err)
}