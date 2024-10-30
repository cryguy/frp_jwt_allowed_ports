package main

import (
	"fmt"
	"github.com/cryguy/frp_jwt_allowed_ports/pkg/server"
	"log"
	"os"

	"github.com/spf13/cobra"
)

const version = "0.0.1"

var (
	showVersion bool

	bindAddr  string
	tokenFile string
)

func init() {
	rootCmd.PersistentFlags().BoolVarP(&showVersion, "version", "v", false, "version")
	rootCmd.PersistentFlags().StringVarP(&bindAddr, "bind_addr", "l", "127.0.0.1:7200", "bind address")
	rootCmd.PersistentFlags().StringVarP(&tokenFile, "secret_file", "k", "./secret", "secret file")
}

var rootCmd = &cobra.Command{
	Use:   "frp_jwt_allowed_ports",
	Short: "frp_jwt_allowed_ports is the server plugin of frp to support multiple users using jwt and by extension grant of port use.",
	RunE: func(cmd *cobra.Command, args []string) error {
		if showVersion {
			fmt.Println(version)
			return nil
		}
		secret, err := getSecretFromFile(tokenFile)
		if err != nil {
			log.Printf("parse tokens from file %s error: %v", tokenFile, err)
			return err
		}
		s, err := server.New(server.Config{
			BindAddress: bindAddr,
			Secret:      secret,
		})
		if err != nil {
			return err
		}
		s.Run()
		return nil
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}

func getSecretFromFile(filePath string) ([]byte, error) {
	secret, err := os.ReadFile(filePath)
	if err != nil {
		return nil, err
	}
	return secret, nil
}
