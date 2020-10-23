package cmd

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"sync"

	"github.com/eric7578/gossipbay/crawler"
	"github.com/spf13/cobra"
)

var runCmd = &cobra.Command{
	Use:   "run",
	Short: "Run crawlers",
	Run: func(cmd *cobra.Command, args []string) {
		var wg sync.WaitGroup
		c := crawler.NewCrawler()
		scanner := bufio.NewScanner(os.Stdin)
		for scanner.Scan() {
			wg.Add(1)
			go func(line string) {
				defer wg.Done()
				args := parseArgs(line)
				data, err := c.CreateJob(args)
				bytes, _ := json.Marshal(data)
				exitOnError(err)
				fmt.Println(string(bytes))
			}(scanner.Text())
		}
		wg.Wait()
	},
}

func init() {
	rootCmd.AddCommand(runCmd)
}

func parseArgs(s string) map[string]string {
	segs := strings.Split(s, " ")
	args := map[string]string{}
	for _, seg := range segs[1:] {
		pairs := strings.Split(seg, "=")
		args[pairs[0]] = pairs[1]
	}
	args["_type"] = segs[0]
	return args
}
