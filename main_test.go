package main

import (
	"github.com/google/go-cmp/cmp"
	// "reflect"
	"testing"
)

func TestHookConfig_Run(t *testing.T) {
	type fields struct {
		Command string
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
		{
			name:   "command",
			fields: fields{Command: "echo hello", Token: "token"},
			args:   args{hookName: "command"},
		},
		{
			name:   "script",
			fields: fields{Command: "./testdata/script.sh", Token: "token"},
			args:   args{hookName: "script"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := HookConfig{
				Command: tt.fields.Command,
				Token:   tt.fields.Token,
			}
			h.Run(tt.args.hookName)
		})
	}
}

func TestHookConfig_Authorized(t *testing.T) {

	config = Config{ GlobalToken: "the-globaltoken" }

	type fields struct {
		Command string
		Token   string
	}
	tests := []struct {
		name   string
		fields fields
		token  string
		want   bool
	}{
		{
			name:   "Authorized with specific token",
			want:   true,
			token:  "token1",
			fields: fields{Command: "echo hello", Token: "token1"},
		},
		{
			name:   "Authorized global token token",
			want:   true,
			token:  "the-globaltoken",
			fields: fields{Command: "echo hello", Token: "token1"},
		},
		{
			name:   "Unauthorized (no token)",
			want:   false,
			fields: fields{Command: "echo hello", Token: "token1"},
		},
		{
			name:   "Unauthorized (wrong token)",
			want:   false,
			token:  "invalid",
			fields: fields{Command: "echo hello", Token: "token1"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := HookConfig{
				Command: tt.fields.Command,
				Token:   tt.fields.Token,
			}
			if got := h.Authorized(tt.token); got != tt.want {
				t.Errorf("HookConfig.Authorized() = %v, want %v.\nToken request: %v\nToken hook: %v\nToken global: %v", got, tt.want, tt.token, h.Token, config.GlobalToken)
			}
		})
	}
}

func TestParseConfig(t *testing.T) {

	fullConfigExample := Config{
		Host:        "127.0.0.1",
		Port:        "9999",
		GlobalToken: "the-defaulttoken",
		Hooks: map[string]HookConfig{
			"hook1": {
				Command: "command1",
				Token:   "token1",
			},
			"hook2": {
				Token: "token2",
			},
			"hook3": {
				Command: "command3",
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
