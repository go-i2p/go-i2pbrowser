package main

import (
	"errors"
	"log"
	"os"
	"path/filepath"

	goi2pbrowser "github.com/go-i2p/go-i2pbrowser"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func defaultDir() string {
	wd, err := os.Getwd()
	if err != nil {
		return ""
	}
	ret := filepath.Join(wd, "i2p-profiles")
	userHome, err := os.UserHomeDir()
	if err != nil {
		return ret
	}
	if wd == userHome {
		ret = filepath.Join(userHome, ".i2p/plugins/i2pbrowser")
	}
	return ret
}

var rootCmd = &cobra.Command{
	Use:   "i2pbrowser [url]",
	Short: "Launch a Firefox browser pre-configured for I2P",
	Args:  cobra.MaximumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		profileDir := viper.GetString("directory")
		if profileDir == "" {
			profileDir = defaultDir()
		}
		if profileDir == "" {
			return errors.New("could not determine profile directory: set --directory explicitly")
		}

		url := "http://127.0.0.1:7657"
		if len(args) > 0 {
			url = args[0]
		}

		// Build browser options from flags/env.
		opts := []goi2pbrowser.BrowserOption{
			goi2pbrowser.WithSAMAddr(viper.GetString("sam-addr")),
		}
		if viper.GetBool("metrics") {
			opts = append(opts, goi2pbrowser.WithMetricsAddr(viper.GetString("metrics-addr")))
		} else {
			opts = append(opts, goi2pbrowser.WithMetricsAddr(""))
		}

		browser, err := goi2pbrowser.NewI2PBrowser(profileDir, opts...)
		if err != nil {
			return err
		}
		defer func() {
			if stopErr := browser.Stop(); stopErr != nil {
				log.Printf("WARNING: failed to stop I2P tunnel: %v", stopErr)
			}
		}()

		switch {
		case viper.GetBool("application"):
			return browser.BrowseApp(url)
		case viper.GetBool("usability"):
			return browser.BrowseUsability(url)
		default:
			return browser.BrowseStrict(url)
		}
	},
}

func init() {
	flags := rootCmd.Flags()
	flags.Bool("usability", false, "Launch in usability mode")
	flags.Bool("application", false, "Launch in app mode")
	flags.StringP("directory", "d", "", "Directory to store profiles in")
	flags.String("sam-addr", "127.0.0.1:7656", "SAMv3 bridge address (host:port)")
	flags.Bool("metrics", false, "Enable the Prometheus metrics server")
	flags.String("metrics-addr", "127.0.0.1:9090", "Metrics server listen address (host:port)")

	if err := viper.BindPFlags(flags); err != nil {
		log.Fatalf("failed to bind flags to viper: %v", err)
	}
	viper.SetEnvPrefix("I2PBROWSER")
	viper.AutomaticEnv()
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}
