package service

import (
	"errors"
	"security/internal/core/domains"
	"security/internal/features/service/mocks"
	"testing"

	"go.uber.org/mock/gomock"
)

func TestGetDanger(t *testing.T) {
	tests := []struct {
		name        string
		md5         string
		wantErr     bool
		waitErr     error
		prepareMock func(m *mocks.MockDangerRepository)
	}{
		{
			name:    "success test",
			md5:     "d41d8cd98f00b204e9800998ecf8427e",
			wantErr: false,
			waitErr: nil,
			prepareMock: func(m *mocks.MockDangerRepository) {
				m.EXPECT().GetDanger(gomock.Any(), gomock.Any()).Return(domains.Danger{}, nil).Times(1)
			},
		},

		{
			name:    "fail test because Create Danger return err",
			md5:     "d41d8cd98f00b204e9800998ecf8427e",
			wantErr: true,
			waitErr: errorFromRepository,
			prepareMock: func(m *mocks.MockDangerRepository) {
				m.EXPECT().GetDanger(gomock.Any(), gomock.Any()).Return(domains.Danger{}, errorFromRepository).Times(1)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := t.Context()

			ctr := gomock.NewController(t)
			defer ctr.Finish()

			mock := mocks.NewMockDangerRepository(ctr)
			tt.prepareMock(mock)

			service := NewDangerService(mock)

			_, err := service.GetDanger(ctx, tt.md5)
			if (err != nil) != tt.wantErr {
				t.Errorf("want err=%v, but returned=%v:", tt.wantErr, err)
			}

			if !errors.Is(err, tt.waitErr) {
				t.Errorf("wait err=%v, but returned=%v:", tt.waitErr, err)
			}
		})
	}
}
