package user

import (
	"testing"

	"server/model"
)

func TestValidateBatchUserStatusTargets(t *testing.T) {
	tests := []struct {
		name       string
		ids        []uint
		status     int
		operatorID uint
		users      []model.SysUser
		wantErr    bool
	}{
		{
			name:       "rejects empty ids",
			ids:        nil,
			status:     0,
			operatorID: 2,
			wantErr:    true,
		},
		{
			name:       "rejects invalid status",
			ids:        []uint{2},
			status:     2,
			operatorID: 1,
			users:      []model.SysUser{{BaseModel: model.BaseModel{ID: 2}, Username: "tester"}},
			wantErr:    true,
		},
		{
			name:       "rejects self disable",
			ids:        []uint{2, 3},
			status:     0,
			operatorID: 2,
			users: []model.SysUser{
				{BaseModel: model.BaseModel{ID: 2}, Username: "operator"},
				{BaseModel: model.BaseModel{ID: 3}, Username: "member"},
			},
			wantErr: true,
		},
		{
			name:       "rejects protected admin",
			ids:        []uint{3},
			status:     0,
			operatorID: 2,
			users: []model.SysUser{
				{
					BaseModel: model.BaseModel{ID: 3},
					Username:  "admin",
					Roles:     []model.SysRole{{Code: "admin"}},
				},
			},
			wantErr: true,
		},
		{
			name:       "rejects self enable",
			ids:        []uint{2},
			status:     1,
			operatorID: 2,
			users: []model.SysUser{
				{BaseModel: model.BaseModel{ID: 2}, Username: "operator"},
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validateBatchUserStatusTargets(tt.ids, tt.status, tt.operatorID, tt.users)
			if (err != nil) != tt.wantErr {
				t.Fatalf("validateBatchUserStatusTargets() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
