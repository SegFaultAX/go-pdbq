package cmd

import (
	"fmt"

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

		for _, cluster := range scope {
			c, err := puppetdb.NewClient(cluster.URL)
			if err != nil {
				return err
			}

			hosts, err := c.Hosts(args[0])
			if err != nil {
				return err
			}

			for _, h := range hosts {
				fmt.Println(h.Name)
			}
		}

		return nil
	},
}

func init() {
	rootCmd.AddCommand(queryCmd)
	queryCmd.Flags().StringP("cluster", "r", "", "which configured cluster to query")
}
