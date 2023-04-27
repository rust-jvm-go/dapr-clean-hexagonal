package msgpack

import (
	"fmt"
	"github.com/vmihailenco/msgpack/v5"
	"reflect"
	"testing"
	"url-shortener/domain"
)

var mockRedirect1 *domain.Redirect
var mockBytes1, mockBytes2 []byte

func setup() {
	mockRedirect1 = &domain.Redirect{
		Code:      "ItsMe",
		URL:       "http://localhost:6000/dapr-clhex-url-shortener",
		CreatedAt: 1e9,
	}
	mockBytes1, _ = msgpack.Marshal(mockRedirect1)

	mockRedirect2 := &domain.Redirect{
		Code:      "NotThisOne",
		URL:       "http://localhost:6000/not-this-one",
		CreatedAt: 1000000,
	}
	mockBytes2, _ = msgpack.Marshal(mockRedirect2)
}

// TODO: Revisit in the future.
func TestMain(m *testing.M) {
	fmt.Println("##### Running TestMain()...")
	setup()
	m.Run()
	fmt.Println("##### Done TestMain()")
}

func TestRedirectMsgPackSerializer_Decode(t *testing.T) {
	r := &domain.Redirect{
		Code:      "ItsMe",
		URL:       "http://localhost:6000/dapr-clhex-url-shortener",
		CreatedAt: 1000000000,
	}
	arg, err := msgpack.Marshal(r)
	if err != nil {
		t.Errorf("Marshal() error = %v", err)
		return
	}

	type args struct {
		input []byte
	}
	tests := []struct {
		name    string
		args    args
		want    *domain.Redirect
		wantErr bool
	}{
		{
			name:    "Decode pass",
			args:    args{arg},
			want:    mockRedirect1,
			wantErr: false,
		},
		{
			name:    "Decode fail",
			args:    args{mockBytes2},
			want:    mockRedirect1,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &RedirectMsgPackSerializer{}

			got, err := m.Decode(tt.args.input)
			if err != nil {
				t.Errorf("Decode() error = %v, wantErr = %v", err, tt.wantErr)
				return
			}

			// fmt.Printf("##### Decode() got = %v, want = %v, wantErr = %v\n", got, tt.want, tt.wantErr)
			// fmt.Printf("##### RESULT 1 = %v, 2 = %v\n", !reflect.DeepEqual(got, tt.want), !tt.wantErr)

			if !reflect.DeepEqual(got, tt.want) && !tt.wantErr {
				t.Errorf("Decode() got = %v, want = %v", got, tt.want)
			}
		})
	}
}

func TestRedirectMsgPackSerializer_Encode(t *testing.T) {
	type args struct {
		input *domain.Redirect
	}
	tests := []struct {
		name    string
		args    args
		want    []byte
		wantErr bool
	}{
		{
			name:    "Encode pass",
			args:    args{mockRedirect1},
			want:    mockBytes1,
			wantErr: false,
		},
		{
			name:    "Encode fail",
			args:    args{mockRedirect1},
			want:    mockBytes2,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &RedirectMsgPackSerializer{}

			got, err := m.Encode(tt.args.input)
			if err != nil {
				t.Errorf("Encode() error = %v, wantErr = %v", err, tt.wantErr)
				return
			}

			// fmt.Printf("##### Encode() got = %v, want = %v, wantErr = %v\n", got, tt.want, tt.wantErr)
			// fmt.Printf("##### RESULT 1 = %v, 2 = %v\n", !reflect.DeepEqual(got, tt.want), !tt.wantErr)

			if !reflect.DeepEqual(got, tt.want) && !tt.wantErr {
				t.Errorf("Encode() got = %v, want = %v", got, tt.want)
			}
		})
	}
}
