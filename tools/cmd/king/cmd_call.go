package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/fatih/color"
	"github.com/juju/errors"
	"github.com/spf13/cobra"

	"github.com/altipla-consulting/king/tools/pkg/auth"
)

func init() {
	CmdRoot.AddCommand(CmdCall)
}

var CmdCall = &cobra.Command{
	Use:     "call",
	Short:   "Llama a un servicio de King desde consola.",
	Example: "king call hotels.hotels.Search project=shs position.lat=123 position.lng=456",
	Args:    cobra.MinimumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		config, err := auth.ReadConfig()
		if err != nil {
			return errors.Trace(err)
		}

		// Defaults to localhost:8080 when there is no active domain
		if len(config.Domains) == 0 {
			config.ActiveDomain = "localhost:8080"
		}

		domain := config.Domain(config.ActiveDomain)

		data := map[string]interface{}{}

		var reqBuf bytes.Buffer
		if err := json.NewEncoder(&reqBuf).Encode(data); err != nil {
			return errors.Trace(err)
		}

		scheme := "https"
		if domain.IsLocal() {
			scheme = "http"
		}
		req, _ := http.NewRequest("POST", fmt.Sprintf("%s://%s/_/%s", scheme, domain.Hostname, args[0]), &reqBuf)
		req.Header.Set("Content-Type", "application/json; charset=utf-8")
		req.Header.Set("Accept", "application/json; charset=utf-8")

		client := &http.Client{
			Timeout: 30 * time.Second,
		}
		resp, err := client.Do(req)
		if err != nil {
			return errors.Trace(err)
		}
		defer resp.Body.Close()

		var respBuf bytes.Buffer
		if _, err := io.Copy(&respBuf, resp.Body); err != nil {
			return errors.Trace(err)
		}

		color.New(color.FgCyan, color.Bold).Printf("Status: ")
		color.New(color.FgYellow).Printf(resp.Status)

		if resp.StatusCode != http.StatusOK {
			return nil
		}

		return nil
	},
}
