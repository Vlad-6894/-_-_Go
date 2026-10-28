package service

import (
	"errors"
	"security/internal/core/domains"
	core_errors "security/internal/core/errors"
	"security/internal/features/service/mocks"
	"testing"

	"go.uber.org/mock/gomock"
)

var (
	errorFromRepository = errors.New("error from repository")
)

func TestCreateDanger(t *testing.T) {
	tests := []struct {
		name        string
		danger      domains.Danger
		wantErr     bool
		waitErr     error
		prepareMock func(m *mocks.MockDangerRepository)
	}{
		{
			name:    "success test",
			danger:  domains.NewUnitializedDanger("MyDoom", "d41d8cd98f00b204e9800998ecf8427e", "Virus", "HIGH"),
			wantErr: false,
			waitErr: nil,
			prepareMock: func(m *mocks.MockDangerRepository) {
				m.EXPECT().CreateDanger(gomock.Any(), gomock.Any()).Return(domains.Danger{}, nil).Times(1)
			},
		},

		{
			name:    "fail test because name has char less 1",
			danger:  domains.NewUnitializedDanger("", "d41d8cd98f00b204e9800998ecf8427e", "Virus", "HIGH"),
			wantErr: true,
			waitErr: core_errors.ErrInvalidArgument,
			prepareMock: func(m *mocks.MockDangerRepository) {
				m.EXPECT().CreateDanger(gomock.Any(), gomock.Any()).Return(domains.Danger{}, nil).Times(0)
			},
		},

		{
			name:    "fail test because md5 has incorrect hash",
			danger:  domains.NewUnitializedDanger("MyDoom", "i41d8cd98f00b204e9800998ecf8427e", "Virus", "HIGH"),
			wantErr: true,
			waitErr: core_errors.ErrInvalidArgument,
			prepareMock: func(m *mocks.MockDangerRepository) {
				m.EXPECT().CreateDanger(gomock.Any(), gomock.Any()).Return(domains.Danger{}, nil).Times(0)
			},
		},

		{
			name:    "fail test because danger type has char less 1",
			danger:  domains.NewUnitializedDanger("MyDoom", "d41d8cd98f00b204e9800998ecf8427e", "", "HIGH"),
			wantErr: true,
			waitErr: core_errors.ErrInvalidArgument,
			prepareMock: func(m *mocks.MockDangerRepository) {
				m.EXPECT().CreateDanger(gomock.Any(), gomock.Any()).Return(domains.Danger{}, nil).Times(0)
			},
		},

		{
			name:    "fail test because danger level is incorrect",
			danger:  domains.NewUnitializedDanger("MyDoom", "d41d8cd98f00b204e9800998ecf8427e", "Virus", "haha"),
			wantErr: true,
			waitErr: core_errors.ErrInvalidArgument,
			prepareMock: func(m *mocks.MockDangerRepository) {
				m.EXPECT().CreateDanger(gomock.Any(), gomock.Any()).Return(domains.Danger{}, nil).Times(0)
			},
		},

		{
			name:    "fail test because Create Danger return err",
			danger:  domains.NewUnitializedDanger("MyDoom", "d41d8cd98f00b204e9800998ecf8427e", "Virus", "HIGH"),
			wantErr: true,
			waitErr: errorFromRepository,
			prepareMock: func(m *mocks.MockDangerRepository) {
				m.EXPECT().CreateDanger(gomock.Any(), gomock.Any()).Return(domains.Danger{}, errorFromRepository).Times(1)
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

			_, err := service.CreateDanger(ctx, tt.danger)
			if (err != nil) != tt.wantErr {
				t.Errorf("want err=%v, but returned=%v:", tt.wantErr, err)
			}

			if !errors.Is(err, tt.waitErr) {
				t.Errorf("wait err=%v, but returned=%v:", tt.waitErr, err)
			}
		})
	}
}
