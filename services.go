package opslevel

import (
	"context"
	"fmt"
)

func (c *Client) GetService(ctx context.Context, alias string) (*Service, error) {
	params := map[string]interface{}{
		"serviceAlias": alias,
	}
	var res serviceResponse
	if err := c.Do(ctx, serviceQuery, params, &res); err != nil {
		return nil, fmt.Errorf("could not find service: %w", err)
	}
	if res.Account.Service == nil {
		return nil, fmt.Errorf("no service was found by alias: %s", alias)
	}
	return res.Account.Service, nil
}

type Service struct {
	Id          string
	Name        string
	Aliases     []string
	Description string
	Framework   string
	Language    string
	Owner       ServiceOwner
	Product     string
	Tier        ServiceTier
}

type ServiceOwner struct {
	Id string
}

type ServiceTier struct {
	Id          string
	Alias       string
	Description string
	Index       int
	Name        string
}

type serviceResponse struct {
	Account struct {
		Service *Service
	}
}

const serviceQuery = `
query($serviceAlias: String) {
  account {
    service(alias: $serviceAlias){
      id
      name
      aliases
      description
      framework
      language
      owner {
        id
      }
      product
      tier {
        id
      }
    }
  }
}
`
