package domain

import (
	"reflect"
	"testing"
)

type mockRedirectRepository struct{}

func (m mockRedirectRepository) Find(code string) (*Redirect, error) {
	return &Redirect{Code: code, URL: "", CreatedAt: 0}, nil
}

func (m mockRedirectRepository) Store(*Redirect) error {
	return nil
}

func (m mockRedirectRepository) Info() string {
	return ""
}

var mock1 = &mockRedirectRepository{}

func TestNewRedirectService(t *testing.T) {
	type args struct {
		redirectRepository IRedirectRepository
	}

	tests := []struct {
		name string
		args args
		want IRedirectService
	}{
		{
			name: "New RedirectService",
			args: args{mock1},
			want: &redirectService{redirectRepository: mock1},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewRedirectService(tt.args.redirectRepository); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewRedirectService() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_redirectService_Find(t *testing.T) {
	type fields struct {
		redirectRepository IRedirectRepository
	}
	type args struct {
		code string
	}

	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *Redirect
		wantErr bool
	}{
		{
			name:    "Find",
			fields:  fields{mock1},
			args:    args{"ItsMe"},
			want:    &Redirect{Code: "ItsMe", URL: "", CreatedAt: 0},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &redirectService{
				redirectRepository: tt.fields.redirectRepository,
			}
			got, err := s.Find(tt.args.code)
			if (err != nil) != tt.wantErr {
				t.Errorf("Find() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Find() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_redirectService_Store(t *testing.T) {
	type fields struct {
		redirectRepository IRedirectRepository
	}
	type args struct {
		redirect *Redirect
	}

	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name:   "Store",
			fields: fields{mock1},
			args: args{&Redirect{
				Code:      "ItsMe",
				URL:       "http://localhost:6000/dapr-clhex-url-shortener",
				CreatedAt: 1e9,
			}},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &redirectService{
				redirectRepository: tt.fields.redirectRepository,
			}
			if err := s.Store(tt.args.redirect); (err != nil) != tt.wantErr {
				t.Errorf("Store() error = %v, wantErr = %v", err, tt.wantErr)
			}
		})
	}
}
