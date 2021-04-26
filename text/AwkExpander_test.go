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
			gotData, gotOpt, gotCmd, err := parseTemplateArg(tt.args.lines, tt.args.delimiter)
			if (err != nil) != tt.wantErr {
				t.Errorf("parseTemplateArg() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotData != tt.wantData {
				t.Errorf("parseTemplateArg() gotData = %v, want %v", gotData, tt.wantData)
			}
			if gotOpt != tt.wantOpt {
				t.Errorf("parseTemplateArg() gotOpt = %v, want %v", gotOpt, tt.wantOpt)
			}
			if gotCmd != tt.wantCmd {
				t.Errorf("parseTemplateArg() gotCmd = %v, want %v", gotCmd, tt.wantCmd)
			}
		})
	}
}

func TestExpend(t *testing.T) {
	type args struct {
		lines string
	}
	tests := []struct {
		name         string
		args         args
		wantExpended string
		wantErr      bool
	}{
		{
			name: "normal",
			args: args{
				lines: "a\nb\nc\n>>>{{$1}},",
			},
			wantExpended: "a,\nb,\nc,\n",
			wantErr:      false,
		},
		{
			name: "normal ,",
			args: args{
				lines: "a,1\nb,2\nc,3\n>>>-F ,\n>>>\n{{$1}}",
			},
			wantExpended: "a\nb\nc\n",
			wantErr:      false,
		},
		{
			name: "have new line 1",
			args: args{
				lines: "a\nb\nc\n>>>{{$1}}\n,",
			},
			wantExpended: "a\n,\nb\n,\nc\n,\n",
			wantErr:      false,
		},
		{
			name: "have new line 2",
			args: args{
				lines: "a\nb\nc\n>>>{{$1}}\n\n,",
			},
			wantExpended: "a\n\n,\nb\n\n,\nc\n\n,\n",
			wantErr:      false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotExpended, err := Expend(tt.args.lines)
			if (err != nil) != tt.wantErr {
				t.Errorf("Expend() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotExpended != tt.wantExpended {
				t.Errorf("Expend() gotExpended = %v, want %v", gotExpended, tt.wantExpended)
			}
		})
	}
}
