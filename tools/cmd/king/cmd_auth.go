package main

import (
	"github.com/juju/errors"
	"github.com/spf13/cobra"

	"github.com/altipla-consulting/king/tools/pkg/auth"
)

func init() {
	CmdRoot.AddCommand(CmdAuth)
}

var CmdAuth = &cobra.Command{
	Use:     "auth",
	Short:   "Guarda el servidor y la autenticación que tenemos que usar en las próximas llamadas.",
	Args:    cobra.RangeArgs(1, 2),
	Example: "king auth api.altipla.consulting FOO_TOKEN_1234\nking auth api.altipla.consulting",
	RunE: func(cmd *cobra.Command, args []string) error {
		config, err := auth.ReadConfig()
		if err != nil {
			return errors.Trace(err)
		}

		if len(args) < 2 {
			args = append(args, "")
		}

		config.SetDomain(args[0], args[1])

		if err := config.Write(); err != nil {
			return errors.Trace(err)
		}

		return nil
	},
}
