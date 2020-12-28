package cmd

import (
	"encoding/base64"
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	"github.com/AlecAivazis/survey/v2"
	"github.com/rs/zerolog/log"

	"github.com/davegallant/vpngate/pkg/vpn"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(connectCmd)
}

var connectCmd = &cobra.Command{
	Use:   "connect",
	Short: "Connect",
	Long:  `Connect to a vpn from a list of servers`,
	Run: func(cmd *cobra.Command, args []string) {

		vpnServers, err := vpn.GetList()

		if err != nil {
			log.Fatal().Msgf(err.Error())
			os.Exit(1)
		}

		serverSelection := []string{}
		serverSelected := vpn.Server{}

		for _, s := range *vpnServers {
			serverSelection = append(serverSelection, fmt.Sprintf("%s (%s)", s.HostName, s.CountryShort))
		}

		selection := ""
		prompt := &survey.Select{
			Message: "Choose a server:",
			Options: serverSelection,
		}
		survey.AskOne(prompt, &selection, survey.WithPageSize(10))

		// Server lookup from selection could be faster than this
		for _, s := range *vpnServers {
			if strings.Contains(selection, s.HostName) {
				serverSelected = s
			}
		}

		decodedConfig, err := base64.StdEncoding.DecodeString(serverSelected.OpenVpnConfigData)
		if err != nil {
			log.Fatal()
		}

		tmpfile, err := ioutil.TempFile("", "vpngate")
		if err != nil {
			log.Fatal()
		}

		defer os.Remove(tmpfile.Name())

		if _, err := tmpfile.Write(decodedConfig); err != nil {
			log.Fatal()
		}

		if err := tmpfile.Close(); err != nil {
			log.Fatal()
		}

		// log.Info().Msgf("Connecting to %s", serverSelected.CountryLong)

		vpn.Connect(tmpfile.Name())

	},
}
