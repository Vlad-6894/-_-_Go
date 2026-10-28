package service

import (
	"errors"
	"security/internal/core/domains"
	core_errors "security/internal/core/errors"
	"security/internal/features/service/mocks"
	"testing"

	"go.uber.org/mock/gomock"
)

func TestUpdateDanger(t *testing.T) {
	tests := []struct {
		name        string
		put         domains.DangerPut
		md5         string
		wantErr     bool
		waitErr     error
		prepareMock func(m *mocks.MockDangerRepository)
	}{
		{
			name: "success test",
			put: domains.NewDangerPut(
				domains.Put[string]{
					Value: "MyDoom",
					Set:   true,
				},
				domains.Put[string]{
					Value: "d41d8cd98f00b204e9800998ecf8427e",
					Set:   true,
				},
				domains.Put[string]{
					Value: "Virus",
					Set:   true,
				},
				domains.Put[string]{
					Value: "HIGH",
					Set:   true,
				},
			),
			md5:     "e41d8cd98f00b204e9800998ecf8427e",
			wantErr: false,
			waitErr: nil,
			prepareMock: func(m *mocks.MockDangerRepository) {
				m.EXPECT().GetDanger(gomock.Any(), gomock.Any()).Return(domains.Danger{}, nil).Times(1)
				m.EXPECT().UpdateDanger(gomock.Any(), gomock.Any(), gomock.Any()).Return(domains.Danger{}, nil).Times(1)
			},
		},

		{
			name: "fail test because name has char less 1",
			put: domains.NewDangerPut(
				domains.Put[string]{
					Value: "",
					Set:   true,
				},
				domains.Put[string]{
					Value: "d41d8cd98f00b204e9800998ecf8427e",
					Set:   true,
				},
				domains.Put[string]{
					Value: "Virus",
					Set:   true,
				},
				domains.Put[string]{
					Value: "HIGH",
					Set:   true,
				},
			),
			md5:     "d41d8cd98f00b204e9800998ecf8427e",
			wantErr: true,
			waitErr: core_errors.ErrInvalidArgument,
			prepareMock: func(m *mocks.MockDangerRepository) {
				m.EXPECT().GetDanger(gomock.Any(), gomock.Any()).Return(domains.Danger{}, nil).Times(1)
				m.EXPECT().UpdateDanger(gomock.Any(), gomock.Any(), gomock.Any()).Return(domains.Danger{}, nil).Times(0)
			},
		},

		{
			name: "fail test because md5 has incorrect hash",
			put: domains.NewDangerPut(
				domains.Put[string]{
					Value: "MyDoom",
					Set:   true,
				},
				domains.Put[string]{
					Value: "i41d8cd98f00b204e9800998ecf8427e",
					Set:   true,
				},
				domains.Put[string]{
					Value: "Virus",
					Set:   true,
				},
				domains.Put[string]{
					Value: "HIGH",
					Set:   true,
				},
			),
			md5:     "d41d8cd98f00b204e9800998ecf8427e",
			wantErr: true,
			waitErr: core_errors.ErrInvalidArgument,
			prepareMock: func(m *mocks.MockDangerRepository) {
				m.EXPECT().GetDanger(gomock.Any(), gomock.Any()).Return(domains.Danger{}, nil).Times(1)
				m.EXPECT().UpdateDanger(gomock.Any(), gomock.Any(), gomock.Any()).Return(domains.Danger{}, nil).Times(0)
			},
		},

		{
			name: "fail test because danger type has char less 1",
			put: domains.NewDangerPut(
				domains.Put[string]{
					Value: "MyDoom",
					Set:   true,
				},
				domains.Put[string]{
					Value: "d41d8cd98f00b204e9800998ecf8427e",
					Set:   true,
				},
				domains.Put[string]{
					Value: "",
					Set:   true,
				},
				domains.Put[string]{
					Value: "HIGH",
					Set:   true,
				},
			),
			md5:     "e41d8cd98f00b204e9800998ecf8427e",
			wantErr: true,
			waitErr: core_errors.ErrInvalidArgument,
			prepareMock: func(m *mocks.MockDangerRepository) {
				m.EXPECT().GetDanger(gomock.Any(), gomock.Any()).Return(domains.Danger{}, nil).Times(1)
				m.EXPECT().UpdateDanger(gomock.Any(), gomock.Any(), gomock.Any()).Return(domains.Danger{}, nil).Times(0)
			},
		},

		{
			name: "fail test because danger level is incorrect",
			put: domains.NewDangerPut(
				domains.Put[string]{
					Value: "MyDoom",
					Set:   true,
				},
				domains.Put[string]{
					Value: "d41d8cd98f00b204e9800998ecf8427e",
					Set:   true,
				},
				domains.Put[string]{
					Value: "Virus",
					Set:   true,
				},
				domains.Put[string]{
					Value: "haha",
					Set:   true,
				},
			),
			md5:     "e41d8cd98f00b204e9800998ecf8427e",
			wantErr: true,
			waitErr: core_errors.ErrInvalidArgument,
			prepareMock: func(m *mocks.MockDangerRepository) {
				m.EXPECT().GetDanger(gomock.Any(), gomock.Any()).Return(domains.Danger{}, nil).Times(1)
				m.EXPECT().UpdateDanger(gomock.Any(), gomock.Any(), gomock.Any()).Return(domains.Danger{}, nil).Times(0)
			},
		},

		{
			name: "fail test because Create Danger return err",
			put: domains.NewDangerPut(
				domains.Put[string]{
					Value: "MyDoom",
					Set:   true,
				},
				domains.Put[string]{
					Value: "d41d8cd98f00b204e9800998ecf8427e",
					Set:   true,
				},
				domains.Put[string]{
					Value: "Virus",
					Set:   true,
				},
				domains.Put[string]{
					Value: "HIGH",
					Set:   true,
				},
			),
			md5:     "e41d8cd98f00b204e9800998ecf8427e",
			wantErr: true,
			waitErr: errorFromRepository,
			prepareMock: func(m *mocks.MockDangerRepository) {
				m.EXPECT().GetDanger(gomock.Any(), gomock.Any()).Return(domains.Danger{}, nil).Times(1)
				m.EXPECT().UpdateDanger(gomock.Any(), gomock.Any(), gomock.Any()).Return(domains.Danger{}, errorFromRepository).Times(1)
			},
		},

		{
			name: "fail test because Get Danger return err",
			put: domains.NewDangerPut(
				domains.Put[string]{
					Value: "MyDoom",
					Set:   true,
				},
				domains.Put[string]{
					Value: "d41d8cd98f00b204e9800998ecf8427e",
					Set:   true,
				},
				domains.Put[string]{
					Value: "Virus",
					Set:   true,
				},
				domains.Put[string]{
					Value: "HIGH",
					Set:   true,
				},
			),
			md5:     "e41d8cd98f00b204e9800998ecf8427e",
			wantErr: true,
			waitErr: errorFromRepository,
			prepareMock: func(m *mocks.MockDangerRepository) {
				m.EXPECT().GetDanger(gomock.Any(), gomock.Any()).Return(domains.Danger{}, errorFromRepository).Times(1)
				m.EXPECT().UpdateDanger(gomock.Any(), gomock.Any(), gomock.Any()).Return(domains.Danger{}, errorFromRepository).Times(0)
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

			_, err := service.PutDanger(ctx, tt.md5, tt.put)
			if (err != nil) != tt.wantErr {
				t.Errorf("want err=%v, but returned=%v:", tt.wantErr, err)
			}

			if !errors.Is(err, tt.waitErr) {
				t.Errorf("wait err=%v, but returned=%v:", tt.waitErr, err)
			}
		})
	}
}
