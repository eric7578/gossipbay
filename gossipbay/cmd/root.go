package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/eric7578/gossipbay/daemon"
	"github.com/spf13/cobra"
)

var (
	port    string
	pubKey  string
	privKey string
	cert    bool
)

var rootCmd = &cobra.Command{
	Use:   "gossipbay",
	Short: "run gossipbay daemon server",
	Run: func(cmd *cobra.Command, args []string) {
		pub, _ := filepath.Abs(pubKey)
		priv, _ := filepath.Abs(privKey)
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
	rootCmd.Flags().StringVar(&pubKey, "pub", "", "public key path")
	rootCmd.Flags().StringVar(&privKey, "priv", "", "private key path")
	rootCmd.Flags().BoolVar(&cert, "cert", false, "create access token")
	rootCmd.MarkFlagRequired("pub")
	rootCmd.MarkFlagRequired("priv")
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
