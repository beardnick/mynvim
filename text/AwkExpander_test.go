package text

import (
	"regexp"
	"testing"
)

func Test_templateOf(t *testing.T) {
	type args struct {
		lines     string
		delimiter *regexp.Regexp
	}
	tests := []struct {
		name     string
		args     args
		wantData string
		wantOpt  string
		wantCmd  string
		wantErr  bool
	}{
		{
			name: "normal",
			args: args{
				lines:     "a\nb\nc\n>>>-F,\n>>>\n{{$1}}",
				delimiter: regexp.MustCompile(">>>+"),
			},
			wantData: "a\nb\nc",
			wantOpt:  "-F,",
			wantCmd:  "{{$1}}",
			wantErr:  false,
		},
		{
			name: "more than 2 >>>",
			args: args{
				lines:     "a\nb\nc\n>>>-F,\n>>>\n{{$1}}\n>>>",
				delimiter: regexp.MustCompile(">>>+"),
			},
			wantData: "a\nb\nc",
			wantOpt:  "-F,",
			wantCmd:  "{{$1}}\n>>>",
			wantErr:  false,
		},
		{
			name: "no enough >>>",
			args: args{
				lines:     "abc\ndef",
				delimiter: regexp.MustCompile(">>>+"),
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotData, gotOpt, gotCmd, err := templateOf(tt.args.lines, tt.args.delimiter)
			if (err != nil) != tt.wantErr {
				t.Errorf("templateOf() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotData != tt.wantData {
				t.Errorf("templateOf() gotData = %v, want %v", gotData, tt.wantData)
			}
			if gotOpt != tt.wantOpt {
				t.Errorf("templateOf() gotOpt = %v, want %v", gotOpt, tt.wantOpt)
			}
			if gotCmd != tt.wantCmd {
				t.Errorf("templateOf() gotCmd = %v, want %v", gotCmd, tt.wantCmd)
			}
		})
	}
}
