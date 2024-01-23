package capabilities

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink-common/pkg/values"
)

type smockCapability struct {
	Validatable
	CapabilityInfo
}

func (m *smockCapability) Start(ctx context.Context, config values.Map) (values.Value, error) {
	return nil, nil
}

func (m *smockCapability) Stop(ctx context.Context) error {
	return nil
}

func (m *smockCapability) Execute(ctx context.Context, inputs values.Map) (values.Value, error) {
	return nil, nil
}

type cmockCapability struct {
	Validatable
	CapabilityInfo
}

func (m *cmockCapability) Start(ctx context.Context, config values.Map) (values.Value, error) {
	return nil, nil
}

func (m *cmockCapability) Stop(ctx context.Context) error {
	return nil
}

func (m *cmockCapability) Execute(ctx context.Context, callback chan values.Map, inputs values.Map) (values.Value, error) {
	return nil, nil
}

func TestRegistry(t *testing.T) {
	ctx := context.Background()

	r := NewRegistry()

	id := "capability-1"
	ci, err := NewCapabilityInfo(
		id,
		CapabilityTypeAction,
		"capability-1-description",
		"v1.0.0",
	)
	require.NoError(t, err)

	c := &smockCapability{CapabilityInfo: ci}
	err = r.Add(ctx, c)
	require.NoError(t, err)

	gc, err := r.Get(ctx, id)
	require.NoError(t, err)

	assert.Equal(t, c, gc)

	cs := r.List(ctx)
	assert.Len(t, cs, 1)
	assert.Equal(t, c, cs[0])
}

func TestRegistry_NoDuplicateIDs(t *testing.T) {
	ctx := context.Background()
	r := NewRegistry()

	id := "capability-1"
	ci, err := NewCapabilityInfo(
		id,
		CapabilityTypeAction,
		"capability-1-description",
		"v1.0.0",
	)
	require.NoError(t, err)

	c := &smockCapability{CapabilityInfo: ci}
	err = r.Add(ctx, c)
	require.NoError(t, err)

	ci, err = NewCapabilityInfo(
		id,
		CapabilityTypeReport,
		"capability-2-description",
		"v1.0.0",
	)
	require.NoError(t, err)
	c2 := &smockCapability{CapabilityInfo: ci}

	err = r.Add(ctx, c2)
	assert.ErrorContains(t, err, "capability with id: capability-1 already exists")
}

func TestRegistry_ChecksExecutionAPIByType(t *testing.T) {
	tcs := []struct {
		name          string
		newCapability func() Capability
		errContains   string
	}{
		{
			name: "trigger, sync",
			newCapability: func() Capability {
				id := uuid.New().String()
				ci, err := NewCapabilityInfo(
					id,
					CapabilityTypeTrigger,
					"capability-1-description",
					"v1.0.0",
				)
				require.NoError(t, err)

				return &smockCapability{CapabilityInfo: ci}
			},
			errContains: "does not satisfy AsynchronousCapability interface",
		},
		{
			name: "reports, sync",
			newCapability: func() Capability {
				id := uuid.New().String()
				ci, err := NewCapabilityInfo(
					id,
					CapabilityTypeReport,
					"capability-1-description",
					"v1.0.0",
				)
				require.NoError(t, err)

				return &smockCapability{CapabilityInfo: ci}
			},
			errContains: "does not satisfy AsynchronousCapability interface",
		},
		{
			name: "action, sync",
			newCapability: func() Capability {
				id := uuid.New().String()
				ci, err := NewCapabilityInfo(
					id,
					CapabilityTypeAction,
					"capability-1-description",
					"v1.0.0",
				)
				require.NoError(t, err)

				return &smockCapability{CapabilityInfo: ci}
			},
		},
		{
			name: "target, sync",
			newCapability: func() Capability {
				id := uuid.New().String()
				ci, err := NewCapabilityInfo(
					id,
					CapabilityTypeTarget,
					"capability-1-description",
					"v1.0.0",
				)
				require.NoError(t, err)

				return &smockCapability{CapabilityInfo: ci}
			},
		},
		{
			name: "trigger, async",
			newCapability: func() Capability {
				id := uuid.New().String()
				ci, err := NewCapabilityInfo(
					id,
					CapabilityTypeTrigger,
					"capability-1-description",
					"v1.0.0",
				)
				require.NoError(t, err)

				return &cmockCapability{CapabilityInfo: ci}
			},
		},
		{
			name: "reports, async",
			newCapability: func() Capability {
				id := uuid.New().String()
				ci, err := NewCapabilityInfo(
					id,
					CapabilityTypeReport,
					"capability-1-description",
					"v1.0.0",
				)
				require.NoError(t, err)

				return &cmockCapability{CapabilityInfo: ci}
			},
		},
		{
			name: "action, async",
			newCapability: func() Capability {
				id := uuid.New().String()
				ci, err := NewCapabilityInfo(
					id,
					CapabilityTypeAction,
					"capability-1-description",
					"v1.0.0",
				)
				require.NoError(t, err)

				return &cmockCapability{CapabilityInfo: ci}
			},
			errContains: "does not satisfy SynchronousCapability interface",
		},
		{
			name: "target, async",
			newCapability: func() Capability {
				id := uuid.New().String()
				ci, err := NewCapabilityInfo(
					id,
					CapabilityTypeTarget,
					"capability-1-description",
					"v1.0.0",
				)
				require.NoError(t, err)

				return &cmockCapability{CapabilityInfo: ci}
			},
			errContains: "does not satisfy SynchronousCapability interface",
		},
	}

	ctx := context.Background()
	reg := NewRegistry()
	for _, tc := range tcs {
		c := tc.newCapability()
		err := reg.Add(ctx, c)
		if tc.errContains == "" {
			require.NoError(t, err)

			info := c.Info()
			id := info.Id
			switch info.CapabilityType {
			case CapabilityTypeAction:
				_, err := reg.GetAction(ctx, id)
				require.NoError(t, err)
			case CapabilityTypeTarget:
				_, err := reg.GetTarget(ctx, id)
				require.NoError(t, err)
			case CapabilityTypeTrigger:
				_, err := reg.GetTrigger(ctx, id)
				require.NoError(t, err)
			case CapabilityTypeReport:
				_, err := reg.GetReport(ctx, id)
				require.NoError(t, err)
			}
		} else {
			assert.ErrorContains(t, err, tc.errContains)
		}
	}
}