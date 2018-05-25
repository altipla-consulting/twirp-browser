package main

import (
	"github.com/juju/errors"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"

	"github.com/altipla-consulting/king/tools/pkg/auth"
)

func init() {
	CmdRoot.AddCommand(CmdActive)
}

var CmdActive = &cobra.Command{
	Use:   "active",
	Short: "Imprime el dominio que está configurado actualmente.",
	RunE: func(cmd *cobra.Command, args []string) error {
		config, err := auth.ReadConfig()
		if err != nil {
			return errors.Trace(err)
		}

		domain := config.Domain(config.ActiveDomain)

		if domain.Hostname == "" {
			log.Warning("No hay ningún dominio active. Usa `king auth domain` para rellenar el primero")
			return nil
		}

		log.Info("Active domain:")
		log.Info("\tHostname: ", domain.Hostname)
		log.Info("\tToken: ", domain.Token)
		log.Info("\tLastUpdated: ", domain.LastUpdated)

		return nil
	},
}
