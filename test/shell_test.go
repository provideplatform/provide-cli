package test

import (
	"testing"

	"github.com/provideplatform/provide-cli/prvd/shell"
	"github.com/spf13/cobra"
)

func TestShell(t *testing.T) {
	type args struct {
		cmd  *cobra.Command
		args []string
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "prvd wallets",
			args: args{
				cmd:  shell.ShellCmd,
				args: []string{},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.args.cmd.SetArgs(tt.args.args)
			tt.args.cmd.Execute()
			// tt.args.cmd.RunE(tt.args.cmd, tt.args.args)
		})
	}
}
