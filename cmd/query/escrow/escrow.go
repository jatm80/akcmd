package escrow

import (
	"github.com/gookit/gcli/v3"
)

func Cmd() *gcli.Command {
	cmd := &gcli.Command{
		Name: "escrow",
		Desc: "Escrow query commands",
		Func: func(cmd *gcli.Command, args []string) error {
			cmd.ShowHelp()
			return nil
		},
	}

	cmd.Add(blocksRemainingCMD())

	return cmd
}

func blocksRemainingCMD() *gcli.Command {
	cmd := &gcli.Command{
		Name: "blocks-remaining",
		Desc: "Compute the number of blocks remaining for an escrow account",
		Func: func(cmd *gcli.Command, args []string) error {
			cmd.ShowHelp()
			return nil
		},
	}

	return cmd
}
