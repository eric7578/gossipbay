package cmd

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/eric7578/gossipbay/crawler"
	"github.com/spf13/cobra"
)

var (
	runType string
	runArgs []string
)

var runCmd = &cobra.Command{
	Use:   "run",
	Short: "Run crawler",
	Run: func(cmd *cobra.Command, args []string) {
		c := crawler.NewCrawler()

		argMap := make(map[string]string)
		for _, seg := range runArgs {
			pairs := strings.Split(seg, "=")
			argMap[pairs[0]] = pairs[1]
		}

		result, err := c.CreateJob(runType, argMap)
		exitOnError(err)

		buf, err := json.Marshal(result)
		fmt.Println(string(buf))
	},
}

func init() {
	rootCmd.AddCommand(runCmd)

	runCmd.Flags().StringVar(&runType, "type", "", "crawler type")
	runCmd.Flags().StringSliceVar(&runArgs, "args", []string{}, "arguments for crawler, ex: board=Gossiping,deviate=0.8")
	runCmd.MarkFlagRequired("type")
}
