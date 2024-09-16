package services

import (
	"testing"
)

func TestUserMutex_GetMutex(t *testing.T) {
	tests := []struct {
		name       string
		setup      func(um *UserMutex)
		expectedOK bool
	}{
		{
			"existing_user",
			func(um *UserMutex) {
				um.SetMutex(1)
			},
			true,
		},
		{
			"non_existing_user",
			func(um *UserMutex) {},
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			um := NewUserMutex()
			tt.setup(um)
			_, ok := um.GetMutex(1)
			if ok != tt.expectedOK {
				t.Errorf("Expected %v, got %v\n",
					tt.expectedOK, ok)
			}
		})
	}
}

func TestUserMutex_SetMutex(t *testing.T) {
	tests := []struct {
		name  string
		setup func(um *UserMutex)
	}{
		{
			"repeated_set",
			func(um *UserMutex) {
				um.SetMutex(1)
			},
		},
		{
			"first_set",
			func(um *UserMutex) {},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			um := NewUserMutex()
			tt.setup(um)
			got := um.SetMutex(1)
			if _, ok := um.GetMutex(1); !ok || got == nil {
				t.Errorf("Expected user ID 1 to be set, but it wasn't")
			}
		})
	}
}

func TestUserMutex_DeleteMutex(t *testing.T) {
	tests := []struct {
		name       string
		setup      func(um *UserMutex)
		expectedOK bool
	}{
		{
			"existing_unlocked_user",
			func(um *UserMutex) {
				um.SetMutex(1)
			},
			true,
		},
		{
			"existing_locked_user",
			func(um *UserMutex) {
				mutex := um.SetMutex(1)
				mutex.Lock()
			},
			false,
		},
		{
			"non_existing_user",
			func(um *UserMutex) {},
			true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			um := NewUserMutex()
			tt.setup(um)
			err := um.DeleteMutex(1)
			if tt.expectedOK {
				if err != nil {
					t.Errorf("Expected no error, got %v\n", err)
					return
				}
				_, ok := um.GetMutex(1)
				if ok {
					t.Errorf("Expected user ID 1 to be deleted, but it wasn't")
				}
			} else {
				if err == nil {
					t.Errorf("Expected error, got none")
				}
			}
		})
	}
}
