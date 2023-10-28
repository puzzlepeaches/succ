package cmd

import (
	"net"
	"os"

	"bufio"
	"strings"

	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var (
	domain      string
	outputJson  bool
	socksProxy  string
	excludeSubs bool
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "succ [domain] [flags] || cat domains.txt | succ [flags]",
	Short: "succ up domains from MS",
	Args:  cobra.MaximumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		var domains []string
		if len(args) > 0 {
			domains = append(domains, args[0])
		} else {
			reader := bufio.NewReader(os.Stdin)
			for {
				input, err := reader.ReadString('\n')
				if err != nil {
					break
				}
				domains = append(domains, strings.TrimSpace(input))
			}
		}

		for _, domain := range domains {
			// Validate the domain
			if _, err := net.LookupIP(domain); err != nil {
				logrus.Fatalf("Invalid domain: %v", err)
			}

			outputFile, _ := cmd.Flags().GetString("output")
			socksProxy, _ := cmd.Flags().GetString("proxy")
			outputJson, _ := cmd.Flags().GetBool("json")
			excludeSubs, _ := cmd.Flags().GetBool("exclude-subs")

			if outputFile != "" {
				if _, err := os.Stat(outputFile); os.IsExist(err) {
					logrus.Fatalf("Output file already exists: %v", err)
				}

				file, err := os.Create(outputFile)
				if err != nil {
					logrus.Fatalf("Error creating output file: %v", err)
				}
				defer file.Close()

			}

			// Call the succer
			s := Succer{
				domain:      domain,
				output:      outputFile,
				outputJson:  outputJson,
				socksProxy:  socksProxy,
				excludeSubs: excludeSubs,
			}
			if err := s.Run(); err != nil {
				logrus.Fatalf("Error running Succer: %v", err)
			}
		}
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		logrus.Fatalf("Error executing root command: %v", err)
		os.Exit(1)
	}
}

func init() {

	// rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.generated code example.yaml)")
	rootCmd.PersistentFlags().StringP("output", "o", "", "Output file.")
	rootCmd.PersistentFlags().BoolP("json", "j", false, "Output to json.")
	rootCmd.PersistentFlags().StringP("proxy", "p", "", "SOCKS5 proxy to use.")
	rootCmd.PersistentFlags().BoolP("exclude-subs", "e", false, "Exclude subdomains from the results.")

}
