package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/fatih/color"
	"github.com/hokaccha/go-prettyjson"
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
		for _, arg := range args[1:] {
			if err := assignData(data, arg); err != nil {
				return errors.Trace(err)
			}
		}

		var reqBuf bytes.Buffer
		if err := json.NewEncoder(&reqBuf).Encode(data); err != nil {
			return errors.Trace(err)
		}

		parts := strings.Split(args[0], ".")
		if len(parts) < 3 {
			return errors.NotValidf("method should have package, service and method")
		}
		kingService := strings.Join(parts[:len(parts)-1], ".")
		kingMethod := parts[len(parts)-1]

		scheme := "https"
		isLocal, err := domain.IsLocal()
		if err != nil {
			return errors.Trace(err)
		}
		if isLocal {
			scheme = "http"
		}
		endpoint := fmt.Sprintf("%s://%s/_/%s/%s", scheme, domain.Hostname, kingService, kingMethod)
		req, err := http.NewRequest("POST", endpoint, &reqBuf)
		if err != nil {
			return errors.Trace(err)
		}
		req.Header.Set("Content-Type", "application/json; charset=utf-8")
		req.Header.Set("Accept", "application/json; charset=utf-8")

		if domain.Token != "" {
			req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", domain.Token))
		}

		content, err := prettyjson.Format(reqBuf.Bytes())
		if err != nil {
			return errors.Trace(err)
		}

		color.New(color.FgYellow, color.Bold).Printf("Method: ")
		color.New(color.FgMagenta).Printf("%s\n", args[0])
		color.New(color.FgYellow, color.Bold).Printf("Hostname: ")
		color.New(color.FgMagenta).Printf("%s\n", domain.Hostname)
		fmt.Println()
		fmt.Println(string(content))
		fmt.Println()
		fmt.Println()

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

		color.New(color.FgYellow, color.Bold).Printf("Status: ")
		if resp.StatusCode != http.StatusOK {
			color.New(color.FgRed).Println(resp.Status)
		} else {
			color.New(color.FgGreen).Println(resp.Status)
		}
		fmt.Println()

		content, err = prettyjson.Format(respBuf.Bytes())
		if err != nil {
			fmt.Println(respBuf.String())
			return nil
		}
		fmt.Println(string(content))
		fmt.Println()

		return nil
	},
}

func assignData(data map[string]interface{}, arg string) error {
	var plain bool
	parts := strings.Split(arg, ":=")
	if len(parts) == 1 {
		parts = strings.Split(arg, "=")
		plain = true
	}
	if len(parts) != 2 {
		return errors.NotValidf("incorrect key=value pair: %s", arg)
	}

	var value interface{} = parts[1]
	if !plain {
		n, err := strconv.ParseInt(parts[1], 10, 32)
		if err != nil {
			return errors.Trace(err)
		}

		value = n
	} else {
		switch parts[1] {
		case "true":
			value = true

		case "false":
			value = "false"
		}
	}

	assignDataRecursively(data, parts[0], value)

	return nil
}

func assignDataRecursively(data map[string]interface{}, key string, value interface{}) {
	parts := strings.Split(key, ".")

	if len(parts) == 1 {
		prev, ok := data[key]
		if ok {
			slice, ok := prev.([]interface{})
			if ok {
				slice = append(slice, value)
			} else {
				data[key] = []interface{}{prev, value}
			}
		} else {
			data[key] = value
		}

		return
	}

	key = parts[0]
	subkey := strings.Join(parts[1:], ".")

	sub, ok := data[key]
	if !ok {
		sub = map[string]interface{}{}
		data[key] = sub
	}

	assignDataRecursively(sub.(map[string]interface{}), subkey, value)
}
