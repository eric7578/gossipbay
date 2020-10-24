package cmd

import (
	"errors"
	"fmt"
	"os"
	"path"
	"path/filepath"

	"github.com/eric7578/gossipbay/daemon"
	"github.com/spf13/cobra"
)

var (
	port    string
	pubKey  string
	privKey string
	dirKeys string
	cert    bool
)

var rootCmd = &cobra.Command{
	Use:   "gossipbay",
	Short: "run gossipbay daemon server",
	Run: func(cmd *cobra.Command, args []string) {
		var (
			pub  string
			priv string
		)
		if pubKey != "" {
			pub, _ = filepath.Abs(pubKey)
		}
		if privKey != "" {
			priv, _ = filepath.Abs(privKey)
		}
		if dirKeys != "" {
			dirKeys, _ = filepath.Abs(dirKeys)
			if pub == "" {
				pub = path.Join(dirKeys, "id_rsa.pub")
			}
			if priv == "" {
				priv = path.Join(dirKeys, "id_rsa")
			}
		}

		if pub == "" {
			exitOnError(errors.New("public key is required"))
		} else if cert && priv == "" {
			exitOnError(errors.New("private key is required"))
		}

		opt := daemon.DaemonOption{
			PublicKey:  pub,
			PrivateKey: priv,
		}
		d := daemon.NewDaemon(opt)
		switch {
		case cert:
			token, err := d.Sign()
			exitOnError(err)
			fmt.Println(token)

		default:
			d.Run(port)
		}
	},
}

func init() {
	rootCmd.Flags().StringVar(&port, "port", ":8080", "daemon port")
	rootCmd.Flags().StringVar(&dirKeys, "keys", "", "read pub/private key from path")
	rootCmd.Flags().StringVar(&pubKey, "pub", "", "public key path")
	rootCmd.Flags().StringVar(&privKey, "priv", "", "private key path")
	rootCmd.Flags().BoolVar(&cert, "cert", false, "create access token")
}

func Execute() {
	exitOnError(rootCmd.Execute())
}

func exitOnError(err error) {
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
