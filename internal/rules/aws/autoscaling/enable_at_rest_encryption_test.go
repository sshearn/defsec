package autoscaling

import (
	"testing"

	"github.com/aquasecurity/defsec/internal/types"

	"github.com/aquasecurity/defsec/pkg/state"

	"github.com/aquasecurity/defsec/pkg/providers/aws/autoscaling"
	"github.com/aquasecurity/defsec/pkg/providers/aws/ec2"
	"github.com/aquasecurity/defsec/pkg/scan"

	"github.com/stretchr/testify/assert"
)

func TestCheckEnableAtRestEncryption(t *testing.T) {
	tests := []struct {
		name     string
		input    autoscaling.Autoscaling
		expected bool
	}{
		{
			name: "Autoscaling unencrypted root block device",
			input: autoscaling.Autoscaling{
				LaunchConfigurations: []autoscaling.LaunchConfiguration{
					{
						Metadata: types.NewTestMetadata(),
						RootBlockDevice: &ec2.BlockDevice{
							Metadata:  types.NewTestMetadata(),
							Encrypted: types.Bool(false, types.NewTestMetadata()),
						},
					},
				},
			},
			expected: true,
		},
		{
			name: "Autoscaling unencrypted EBS block device",
			input: autoscaling.Autoscaling{
				LaunchConfigurations: []autoscaling.LaunchConfiguration{
					{
						Metadata: types.NewTestMetadata(),
						EBSBlockDevices: []*ec2.BlockDevice{
							{
								Metadata:  types.NewTestMetadata(),
								Encrypted: types.Bool(false, types.NewTestMetadata()),
							},
						},
					},
				},
			},
			expected: true,
		},
		{
			name: "Autoscaling encrypted root and EBS block devices",
			input: autoscaling.Autoscaling{
				LaunchConfigurations: []autoscaling.LaunchConfiguration{
					{
						Metadata: types.NewTestMetadata(),
						RootBlockDevice: &ec2.BlockDevice{
							Metadata:  types.NewTestMetadata(),
							Encrypted: types.Bool(true, types.NewTestMetadata()),
						},
						EBSBlockDevices: []*ec2.BlockDevice{
							{
								Metadata:  types.NewTestMetadata(),
								Encrypted: types.Bool(true, types.NewTestMetadata()),
							},
						},
					},
				},
			},
			expected: false,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			var testState state.State
			testState.AWS.Autoscaling = test.input
			results := CheckEnableAtRestEncryption.Evaluate(&testState)
			var found bool
			for _, result := range results {
				if result.Status() == scan.StatusFailed && result.Rule().LongID() == CheckEnableAtRestEncryption.Rule().LongID() {
					found = true
				}
			}
			if test.expected {
				assert.True(t, found, "Rule should have been found")
			} else {
				assert.False(t, found, "Rule should not have been found")
			}
		})
	}
}
