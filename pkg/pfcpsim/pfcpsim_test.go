// SPDX-License-Identifier: Apache-2.0
// Copyright 2022-present Open Networking Foundation

package pfcpsim

import (
	"sync"
	"testing"
	"time"
)

func TestNewPFCPClientUsesStableRecoveryTimestamp(t *testing.T) {
	before := time.Now()
	client := NewPFCPClient("127.0.0.1")
	after := time.Now()

	if client.recoveryTimestamp.Before(before) || client.recoveryTimestamp.After(after) {
		t.Fatalf("recovery timestamp %v was not initialized at client startup", client.recoveryTimestamp)
	}

	original := client.recoveryTimestamp
	time.Sleep(time.Millisecond)
	if client.recoveryTimestamp != original {
		t.Fatalf("recovery timestamp changed from %v to %v", original, client.recoveryTimestamp)
	}
}

func TestHeartbeatWorkerCanOnlyStartOnce(t *testing.T) {
	client := NewPFCPClient("127.0.0.1")
	started := 0
	var lock sync.Mutex
	start := func() {
		client.heartbeatOnce.Do(func() {
			lock.Lock()
			started++
			lock.Unlock()
		})
	}

	var workers sync.WaitGroup
	for i := 0; i < 10; i++ {
		workers.Add(1)
		go func() {
			defer workers.Done()
			start()
		}()
	}
	workers.Wait()

	if started != 1 {
		t.Fatalf("expected one heartbeat worker, got %d", started)
	}
}

func TestResetSessions(t *testing.T) {
	InsertSession(1, &PFCPSession{})
	InsertSession(11, &PFCPSession{})

	ResetSessions()

	if count := GetActiveSessionNum(); count != 0 {
		t.Fatalf("expected no active sessions after reset, got %d", count)
	}
}
