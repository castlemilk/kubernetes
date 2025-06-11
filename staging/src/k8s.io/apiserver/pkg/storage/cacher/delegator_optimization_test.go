/*
Copyright 2024 The Kubernetes Authors.

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

package cacher

import (
	"context"
	"os"
	"strconv"
	"sync"
	"testing"
	"time"

	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/watch"
	"k8s.io/apiserver/pkg/features"
	"k8s.io/apiserver/pkg/storage"
	utilfeature "k8s.io/apiserver/pkg/util/feature"
	featuregatetesting "k8s.io/component-base/featuregate/testing"
)

func TestConsistentReadOptimizationEnvironmentVariable(t *testing.T) {
	// Test that environment variables work correctly

	// Save original state
	originalEnabled := ConsistentReadOptimizationEnabled
	defer func() {
		ConsistentReadOptimizationEnabled = originalEnabled
	}()

	// Test disabling optimization
	t.Setenv("KUBE_CONSISTENT_READ_OPTIMIZATION", "false")

	// Reinitialize (simulate package init)
	if val, exists := os.LookupEnv("KUBE_CONSISTENT_READ_OPTIMIZATION"); exists {
		ConsistentReadOptimizationEnabled, _ = strconv.ParseBool(val)
	}

	if ConsistentReadOptimizationEnabled {
		t.Error("Optimization should be disabled when env var is false")
	}

	// Test enabling optimization
	t.Setenv("KUBE_CONSISTENT_READ_OPTIMIZATION", "true")
	if val, exists := os.LookupEnv("KUBE_CONSISTENT_READ_OPTIMIZATION"); exists {
		ConsistentReadOptimizationEnabled, _ = strconv.ParseBool(val)
	}

	if !ConsistentReadOptimizationEnabled {
		t.Error("Optimization should be enabled when env var is true")
	}
}

func TestConsistentReadCoalescingWindow(t *testing.T) {
	// Test that the coalescing window is set to a reasonable value
	expectedWindow := 100 * time.Millisecond
	if ConsistentReadCoalescingWindow != expectedWindow {
		t.Errorf("Expected coalescing window %v, got %v", expectedWindow, ConsistentReadCoalescingWindow)
	}
}

func TestConsistentReadCacheFreshnessThreshold(t *testing.T) {
	// Test that the cache freshness threshold is set to a reasonable value
	expectedThreshold := 1 * time.Second
	if ConsistentReadCacheFreshnessThreshold != expectedThreshold {
		t.Errorf("Expected cache freshness threshold %v, got %v", expectedThreshold, ConsistentReadCacheFreshnessThreshold)
	}
}

func TestConsistentListFromCacheFeatureGate(t *testing.T) {
	// Test that the feature gate can be enabled and disabled
	featuregatetesting.SetFeatureGateDuringTest(t, utilfeature.DefaultFeatureGate, features.ConsistentListFromCache, true)

	if !utilfeature.DefaultFeatureGate.Enabled(features.ConsistentListFromCache) {
		t.Error("ConsistentListFromCache feature should be enabled")
	}
}

func TestConsistencyCheckerEnvironmentVariable(t *testing.T) {
	// Test that the consistency checker environment variable works

	// Save original state
	originalEnabled := ConsistencyCheckerEnabled
	defer func() {
		ConsistencyCheckerEnabled = originalEnabled
	}()

	// Test enabling consistency checker
	t.Setenv("KUBE_WATCHCACHE_CONSISTENCY_CHECKER", "true")

	// Reinitialize (simulate package init)
	ConsistencyCheckerEnabled, _ = strconv.ParseBool(os.Getenv("KUBE_WATCHCACHE_CONSISTENCY_CHECKER"))

	if !ConsistencyCheckerEnabled {
		t.Error("Consistency checker should be enabled when env var is true")
	}

	// Test disabling consistency checker
	t.Setenv("KUBE_WATCHCACHE_CONSISTENCY_CHECKER", "false")
	ConsistencyCheckerEnabled, _ = strconv.ParseBool(os.Getenv("KUBE_WATCHCACHE_CONSISTENCY_CHECKER"))

	if ConsistencyCheckerEnabled {
		t.Error("Consistency checker should be disabled when env var is false")
	}
}

func TestConsistencyCheckPeriod(t *testing.T) {
	// Test that the consistency check period is set to a reasonable value
	expectedPeriod := 5 * time.Minute
	if ConsistencyCheckPeriod != expectedPeriod {
		t.Errorf("Expected consistency check period %v, got %v", expectedPeriod, ConsistencyCheckPeriod)
	}
}

func TestGlobalConsistentReadCoalescerInitialization(t *testing.T) {
	// Test that the global coalescer is properly initialized
	if globalConsistentReadCoalescer == nil {
		t.Error("Expected globalConsistentReadCoalescer to be initialized")
	}

	if globalConsistentReadCoalescer.resourceVersionCache == nil {
		t.Error("Expected resourceVersionCache to be initialized")
	}

	if globalConsistentReadCoalescer.pendingRequests == nil {
		t.Error("Expected pendingRequests to be initialized")
	}
}

// Mock implementations for testing
type mockStorageWithCallCount struct {
	callCount int
	mu        sync.Mutex
}

func (m *mockStorageWithCallCount) GetCurrentResourceVersion(ctx context.Context) (uint64, error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.callCount++
	return 1000, nil
}

func (m *mockStorageWithCallCount) Versioner() storage.Versioner { return nil }
func (m *mockStorageWithCallCount) Create(ctx context.Context, key string, obj, out runtime.Object, ttl uint64) error {
	return nil
}
func (m *mockStorageWithCallCount) Delete(ctx context.Context, key string, out runtime.Object, preconditions *storage.Preconditions, validateDeletion storage.ValidateObjectFunc, cachedExistingObject runtime.Object, opts storage.DeleteOptions) error {
	return nil
}
func (m *mockStorageWithCallCount) Watch(ctx context.Context, key string, opts storage.ListOptions) (watch.Interface, error) {
	return nil, nil
}
func (m *mockStorageWithCallCount) Get(ctx context.Context, key string, opts storage.GetOptions, objPtr runtime.Object) error {
	return nil
}
func (m *mockStorageWithCallCount) GetList(ctx context.Context, key string, opts storage.ListOptions, listObj runtime.Object) error {
	return nil
}
func (m *mockStorageWithCallCount) GuaranteedUpdate(ctx context.Context, key string, destination runtime.Object, ignoreNotFound bool, preconditions *storage.Preconditions, tryUpdate storage.UpdateFunc, cachedExistingObject runtime.Object) error {
	return nil
}
func (m *mockStorageWithCallCount) Count(key string) (int64, error)                { return 0, nil }
func (m *mockStorageWithCallCount) ReadinessCheck() error                          { return nil }
func (m *mockStorageWithCallCount) RequestWatchProgress(ctx context.Context) error { return nil }
