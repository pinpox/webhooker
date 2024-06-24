package main

import (
	"github.com/google/go-cmp/cmp"
	// "reflect"
	"testing"
)

func TestHookConfig_Run(t *testing.T) {
	type fields struct {
		Command string
		Script  string
		Token   string
	}
	type args struct {
		hookName string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := HookConfig{
				Command: tt.fields.Command,
				Script:  tt.fields.Script,
				Token:   tt.fields.Token,
			}
			h.Run(tt.args.hookName)
		})
	}
}

func TestHookConfig_Authorized(t *testing.T) {
	type fields struct {
		Command string
		Script  string
		Token   string
	}
	type args struct {
		token string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   bool
	}{
		// TODO: Add test cases.
		{name: "Authorized with specific token", want: true},
		{name: "Authorized global token token", want: true},
		{name: "Unauthorized", want: false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := HookConfig{
				Command: tt.fields.Command,
				Script:  tt.fields.Script,
				Token:   tt.fields.Token,
			}
			if got := h.Authorized(tt.args.token); got != tt.want {
				t.Errorf("HookConfig.Authorized() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestParseConfig(t *testing.T) {

	fullConfigExample := Config{
		Host:        "the-host",
		Port:        "the-port",
		GlobalToken: "the-defaulttoken",
		Hooks: map[string]HookConfig{
			"hook1": {
				Command: "command1",
				Token:   "token1",
			},
			"hook2": {
				Token:  "token2",
				Script: "script2",
			},
			"hook3": {
				Command: "command3",
				Script:  "script3",
				Token:   "token3",
			},
		},
	}

	tests := []struct {
		name string
		path string
		want Config
	}{
		{
			name: "Full config",
			path: "testdata/config.yaml",
			want: Config(fullConfigExample),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			got := ParseConfig(tt.path)

			if diff := cmp.Diff(tt.want, got); diff != "" {
				t.Errorf("ParseConfig() mismatch (-want +got):\n%s", diff)
			}

		})
	}
}

func Test_main(t *testing.T) {
	tests := []struct {
		name string
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			main()
		})
	}
}
