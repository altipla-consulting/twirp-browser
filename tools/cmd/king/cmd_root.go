package main

import (
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var debugApp bool

func init() {
	CmdRoot.PersistentFlags().BoolVarP(&debugApp, "debug", "d", false, "Activa el logging de depuración")
}

var CmdRoot = &cobra.Command{
	Use:   "king",
	Short: "Helper para llamar a servicios desde consola",
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		if debugApp {
			log.SetLevel(log.DebugLevel)
			log.Debug("DEBUG log level activated")
		}

		return nil
	},
}
