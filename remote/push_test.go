package remote

import (
	"os"
	"testing"
)

func TestDo(t *testing.T) {
	type args struct {
		cmd string
		env []string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "normal",
			args: args{
				cmd: "ls -al",
				env: []string{"PATH=/bin"},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := Do(tt.args.cmd, tt.args.env, os.Stdout)
			if (err != nil) != tt.wantErr {
				t.Errorf("Do() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}
