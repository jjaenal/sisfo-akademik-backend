package entity

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestStudent_Validate(t *testing.T) {
	tests := []struct {
		name     string
		student  Student
		want     map[string]string
		hasError bool
	}{
		{
			name: "valid student",
			student: Student{
				Name:     "John Doe",
				TenantID: "tenant-1",
				Status:   "active",
			},
			want:     map[string]string{},
			hasError: false,
		},
		{
			name: "missing name",
			student: Student{
				TenantID: "tenant-1",
			},
			want: map[string]string{
				"name": "Name is required",
			},
			hasError: true,
		},
		{
			name: "missing tenant_id",
			student: Student{
				Name: "John Doe",
			},
			want: map[string]string{
				"tenant_id": "Tenant ID is required",
			},
			hasError: true,
		},
		{
			name: "default status",
			student: Student{
				Name:     "John Doe",
				TenantID: "tenant-1",
			},
			want:     map[string]string{},
			hasError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.student.Validate()
			if tt.hasError {
				assert.Greater(t, len(got), 0)
				for k, v := range tt.want {
					assert.Equal(t, v, got[k])
				}
			} else {
				assert.Empty(t, got)
			}
			
			if tt.name == "default status" {
				assert.Equal(t, "active", tt.student.Status)
			}
		})
	}
}
