package cmd

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/spf13/cobra"
)

type (
	Host struct {
		Name string `json:"certname"`
	}

	Eval interface {
		Eval() string
	}

	BinOp struct {
		Op  string
		LHS Eval
		RHS Eval
	}

	Symbol struct {
		Name string
	}

	String struct {
		Val string
	}

	Or struct {
		LHS Eval
		RHS Eval
	}

	And struct {
		LHS Eval
		RHS Eval
	}
)

func (bo BinOp) Eval() string {
	return fmt.Sprintf("%s %s %s", bo.LHS.Eval(), bo.Op, bo.RHS.Eval())
}

func (s Symbol) Eval() string {
	return s.Name
}

func (s String) Eval() string {
	return strconv.Quote(s.Val)
}

func (o Or) Eval() string {
	return fmt.Sprintf("(%s) or (%s)", o.LHS.Eval(), o.RHS.Eval())
}

func (a And) Eval() string {
	return fmt.Sprintf("(%s) and (%s)", a.LHS.Eval(), a.RHS.Eval())
}

var queryCmd = &cobra.Command{
	Use:   "query",
	Short: "query hosts out of PuppetDB",
	Long:  `TODO`,
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) < 1 {
			return fmt.Errorf("must supply a query string")
		}

		cli := &http.Client{
			Timeout: 10 * time.Second,
		}

		expr := Or{
			LHS: BinOp{Op: "~", LHS: Symbol{"certname"}, RHS: String{args[0]}},
			RHS: BinOp{Op: "=", LHS: Symbol{"facts.osfamily"}, RHS: String{"Debian"}},
		}
		fmt.Println(expr.Eval())

		query := map[string]string{
			"query": fmt.Sprintf("inventory[certname]{ %s }", expr.Eval()),
		}
		body, err := json.Marshal(query)
		if err != nil {
			return err
		}
		fmt.Println(string(body))

		endpoint := "http://puppetdb-pp-rs.otenv.com:8080/pdb/query/v4"

		req, err := http.NewRequest(http.MethodPost, endpoint, bytes.NewReader(body))
		if err != nil {
			return err
		}
		req.Header.Add("Content-Type", "application/json")

		resp, err := cli.Do(req)
		if err != nil {
			return err
		}
		defer resp.Body.Close()

		var hosts []Host
		err = json.NewDecoder(resp.Body).Decode(&hosts)
		if err != nil {
			return err
		}

		for _, h := range hosts {
			fmt.Println(h.Name)
		}

		return nil
	},
}

func init() {
	rootCmd.AddCommand(queryCmd)
}
