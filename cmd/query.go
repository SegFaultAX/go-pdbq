package cmd

import (
	"fmt"
	"sort"
	"sync"

	"github.com/segfaultax/go-pdbq/puppetdb"
	"github.com/spf13/cobra"
)

var queryCmd = &cobra.Command{
	Use:   "query",
	Short: "query hosts out of PuppetDB",
	Long:  `TODO`,
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) < 1 {
			return fmt.Errorf("must supply a query string")
		}

		clus, _ := cmd.Flags().GetString("cluster")

		cfg, err := GetConfig()
		if err != nil {
			return err
		}

		var scope []Cluster

		cluster, ok := cfg.Clusters[clus]
		if clus == "" {
			for _, c := range cfg.Clusters {
				scope = append(scope, c)
			}
		} else if !ok {
			return fmt.Errorf("unknown cluster: %s", clus)
		} else {
			scope = append(scope, cluster)
		}

		var wg sync.WaitGroup

		ch := make(chan string)

		for _, cluster := range scope {
			wg.Add(1)
			go func(endpoint string) {
				defer wg.Done()
				c, err := puppetdb.NewClient(endpoint)
				if err != nil {
					fmt.Println(err)
				}

				hosts, err := c.Hosts(args[0])
				if err != nil {
					fmt.Println(err)
				}

				for _, h := range hosts {
					ch <- h.Name
				}
			}(cluster.URL)
		}

		go func() {
			wg.Wait()
			close(ch)
		}()

		var all []string
		for s := range ch {
			all = append(all, s)
		}

		sort.Strings(all)
		for _, h := range all {
			fmt.Println(h)
		}

		return nil
	},
}

func init() {
	rootCmd.AddCommand(queryCmd)
	queryCmd.Flags().StringP("cluster", "r", "", "which configured cluster to query")
}
