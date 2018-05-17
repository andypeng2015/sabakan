package client

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"path"

	"github.com/cybozu-go/cmd"
	"github.com/cybozu-go/sabakan"
)

// Client is a sabakan client
type Client struct {
	endpoint string
	http     *cmd.HTTPClient
}

// NewClient creates a new sabakan client
func NewClient(endpoint string, http *cmd.HTTPClient) *Client {
	return &Client{
		endpoint: endpoint,
		http:     http,
	}
}

// IPAMConfigGet retrieves IPAM configurations
func (c *Client) IPAMConfigGet(ctx context.Context) (*sabakan.IPAMConfig, *Status) {
	var conf sabakan.IPAMConfig
	err := c.getJSON(ctx, "/config/ipam", nil, &conf)
	if err != nil {
		return nil, err
	}
	return &conf, nil
}

// IPAMConfigSet sets a remote config
func (c *Client) IPAMConfigSet(ctx context.Context, conf *sabakan.IPAMConfig) *Status {
	return c.sendRequestWithJSON(ctx, "PUT", "/config/ipam", conf)
}

// MachinesGet get machine information from sabakan server
func (c *Client) MachinesGet(ctx context.Context, params map[string]string) ([]sabakan.Machine, *Status) {
	var machines []sabakan.Machine
	err := c.getJSON(ctx, "/machines", params, &machines)
	if err != nil {
		return nil, err
	}
	return machines, nil
}

// MachinesCreate create machines information to sabakan server
func (c *Client) MachinesCreate(ctx context.Context, machines []sabakan.Machine) *Status {
	return c.sendRequestWithJSON(ctx, "POST", "/machines", machines)
}

// MachinesRemove removes machine information from sabakan server
func (c *Client) MachinesRemove(ctx context.Context, serial string) *Status {
	return c.sendRequest(ctx, "DELETE", path.Join("/machines", serial))
}

func (c *Client) getJSON(ctx context.Context, path string, params map[string]string, data interface{}) *Status {
	req, err := http.NewRequest("GET", c.endpoint+"/api/v1"+path, nil)
	if err != nil {
		return ErrorStatus(err)
	}
	q := req.URL.Query()
	for k, v := range params {
		q.Add(k, v)
	}
	req.URL.RawQuery = q.Encode()
	req = req.WithContext(ctx)
	res, err := c.http.Do(req)
	if err != nil {
		return ErrorStatus(err)
	}
	defer res.Body.Close()

	errorStatus := ErrorHTTPStatus(res)
	if errorStatus != nil {
		return errorStatus
	}

	err = json.NewDecoder(res.Body).Decode(data)
	if err != nil {
		return ErrorStatus(err)
	}

	return nil
}

func (c *Client) sendRequestWithJSON(ctx context.Context, method string, path string, data interface{}) *Status {
	b := new(bytes.Buffer)
	err := json.NewEncoder(b).Encode(data)
	if err != nil {
		return ErrorStatus(err)
	}

	req, err := http.NewRequest(method, c.endpoint+"/api/v1"+path, b)
	if err != nil {
		return ErrorStatus(err)
	}
	req = req.WithContext(ctx)
	res, err := c.http.Do(req)
	if err != nil {
		return ErrorStatus(err)
	}
	defer res.Body.Close()

	return ErrorHTTPStatus(res)
}

func (c *Client) sendRequest(ctx context.Context, method, path string) *Status {
	req, err := http.NewRequest(method, c.endpoint+"/api/v1"+path, nil)
	if err != nil {
		return ErrorStatus(err)
	}
	req = req.WithContext(ctx)
	res, err := c.http.Do(req)
	if err != nil {
		return ErrorStatus(err)
	}
	defer res.Body.Close()

	return ErrorHTTPStatus(res)
}