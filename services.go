package opslevel

import (
	"context"
	"fmt"
)

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

func (c *Client) CreateService(ctx context.Context, service Service) (*Service, error) {
	args := map[string]string{
		"name":        service.Name,
		"product":     service.Product,
		"description": service.Description,
		"language":    service.Language,
		"framework":   service.Framework,
	}
	params := map[string]interface{}{
		"input": args,
	}
	var res createServiceResponse
	if err := c.Do(ctx, serviceCreateMutation, params, &res); err != nil {
		return nil, err
	}
	// Check for application level errors
	if err := handleGraphqlErrs(res.ServiceCreate.Errors); err != nil {
		return nil, err
	}
	return res.ServiceCreate.Service, nil
}

func (c *Client) GetServiceById(ctx context.Context, id string) (*Service, error) {
	params := map[string]interface{}{
		"serviceId": id,
	}
	return c.getService(ctx, serviceQueryById, params)
}

func (c *Client) GetServiceByAlias(ctx context.Context, alias string) (*Service, error) {
	params := map[string]interface{}{
		"serviceAlias": alias,
	}
	return c.getService(ctx, serviceQuery, params)
}

func (c *Client) getService(ctx context.Context, query string, params map[string]interface{}) (*Service, error) {
	var res serviceResponse
	if err := c.Do(ctx, query, params, &res); err != nil {
		return nil, fmt.Errorf("could not find service: %w", err)
	}
	if res.Account.Service == nil {
		return nil, fmt.Errorf("no service was found for params: %v", params)
	}
	return res.Account.Service, nil
}

func (c *Client) DeleteServiceByAlias(ctx context.Context, alias string) (*DeleteServiceResponse, error) {
	args := map[string]string{
		"alias": alias,
	}
	params := map[string]interface{}{
		"input": args,
	}
	res, err := c.deleteService(ctx, serviceDeleteMutation, params)
	return res, err
}

func (c *Client) DeleteServiceById(ctx context.Context, id string) (*DeleteServiceResponse, error) {
	args := map[string]string{
		"id": id,
	}
	params := map[string]interface{}{
		"input": args,
	}
	res, err := c.deleteService(ctx, serviceDeleteMutation, params)
	return res, err
}

func (c *Client) deleteService(ctx context.Context, query string, params map[string]interface{}) (*DeleteServiceResponse, error) {
	var res DeleteServiceResponse
	if err := c.Do(ctx, query, params, &res); err != nil {
		return nil, err
	}
	// Check for application level errors
	if err := handleGraphqlErrs(res.ServiceDelete.Errors); err != nil {
		return nil, err
	}
	return &res, nil
}

type serviceResponse struct {
	Account struct {
		Service *Service
	}
}

type createServiceResponse struct {
	ServiceCreate struct {
		Service *Service
		Errors  []graphqlError
	}
}

type DeleteServiceResponse struct {
	ServiceDelete struct {
		DeletedServiceAlias string
		DeletedServiceId    string
		Errors              []graphqlError
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

const serviceQueryById = `
query($serviceId: ID) {
  account {
    service(id: $serviceId){
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

const serviceCreateMutation = `
mutation create($input: ServiceCreateInput!){
  serviceCreate(input: $input){
    service{
      id
      name
      product
      description
      language
      framework
    }
    errors {
      path
      message
    }
  }
}
`

const serviceDeleteMutation = `
mutation delete($input: ServiceDeleteInput!){
  serviceDelete(input: $input){
    deletedServiceAlias
	deletedServiceId
    errors {
      path
      message
    }
  }
}
`
