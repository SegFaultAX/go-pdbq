package puppetdb

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"time"
)

type (
	Host struct {
		Name string `json:"certname"`
	}

	Client struct {
		BaseURL *url.URL
	}
)

func NewClient(endpoint string) (*Client, error) {
	u, err := url.Parse(endpoint)
	if err != nil {
		return nil, err
	}
	return &Client{
		BaseURL: u,
	}, nil
}

func (c *Client) Hosts(q string) ([]Host, error) {
	cli := &http.Client{
		Timeout: 10 * time.Second,
	}

	expr := BinOp{Op: "~", LHS: Symbol{"certname"}, RHS: String{q}}

	query := map[string]string{
		"query": fmt.Sprintf("inventory[certname]{ %s }", expr.Eval()),
	}
	body, err := json.Marshal(query)
	if err != nil {
		return nil, err
	}

	endpoint, _ := c.BaseURL.Parse("/pdb/query/v4")

	req, err := http.NewRequest(http.MethodPost, endpoint.String(), bytes.NewReader(body))
	if err != nil {
		return nil, err
	}
	req.Header.Add("Content-Type", "application/json")

	resp, err := cli.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var hosts []Host
	err = json.NewDecoder(resp.Body).Decode(&hosts)
	if err != nil {
		return nil, err
	}

	return hosts, nil
}
