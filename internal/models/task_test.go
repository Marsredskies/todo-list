package models

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestValidator(t *testing.T) {
	testCases := []struct {
		name        string
		input       Task
		expectedErr error
	}{
		{
			name:        "test ok",
			input:       Task{0, "name", "somebody", "to do", "sometext"},
			expectedErr: nil,
		},

		{
			name:        "test invalid task status",
			input:       Task{0, "name", "somebody", "dingus", "sometext"},
			expectedErr: errInvalidTaskStatus,
		},

		{
			name:        "test empty task name",
			input:       Task{0, "", "somebody", "done", "sometext"},
			expectedErr: errEmptyTaskName,
		},

		{
			name:        "test empty task description",
			input:       Task{0, "name", "somebody", "in progress", ""},
			expectedErr: errEmptyTaskDescription,
		},
		{
			name:        "test empty task description",
			input:       Task{0, "name", "", "in progress", "sometext"},
			expectedErr: errNoAssignee,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			err := tc.input.Validate()
			require.Equal(t, tc.expectedErr, err)
		})
	}
}
