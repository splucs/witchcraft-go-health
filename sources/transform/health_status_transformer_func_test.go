// Copyright (c) 2020 Palantir Technologies. All rights reserved.
//
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

package transform

import (
	"context"
	"testing"

	werror "github.com/palantir/witchcraft-go-error"
	"github.com/palantir/witchcraft-go-health/conjure/witchcraft/api/health"
	"github.com/palantir/witchcraft-go-health/sources/store"
	"github.com/stretchr/testify/assert"
)

func TestSource(t *testing.T) {
	expected := health.HealthStatus{
		Checks: map[health.CheckType]health.HealthCheckResult{
			"a": {},
		},
	}
	mapper := func(in health.HealthStatus) health.HealthStatus {
		return expected
	}
	keyed := store.NewKeyedErrorHealthCheckSource("foo", "bar")
	keyed.Submit("foo", werror.Error("err"))
	source := NewSource(keyed, mapper)
	status := source.HealthStatus(context.Background())
	assert.Equal(t, expected, status)
}

func TestSourceNilChecks(t *testing.T) {
	mapper := func(in health.HealthStatus) health.HealthStatus {
		return health.HealthStatus{}
	}
	keyed := store.NewKeyedErrorHealthCheckSource("foo", "bar")
	source := NewSource(keyed, nil)
	assert.Equal(t, source.HealthStatus(context.Background()), health.HealthStatus{
		Checks: map[health.CheckType]health.HealthCheckResult{
			"foo": {
				Type:    "foo",
				State:   health.New_HealthState(health.HealthState_HEALTHY),
				Message: toString("bar"),
			},
		},
	})
	source = NewSource(nil, mapper)
	assert.Equal(t, source.HealthStatus(context.Background()), health.HealthStatus{})
}

func toString(s string) *string {
	return &s
}
