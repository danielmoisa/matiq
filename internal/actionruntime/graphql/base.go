package graphql

import (
	"net/http"
	"net/url"

	"github.com/go-resty/resty/v2"
)

const (
	AUTH_NONE   = "none"
	AUTH_BASIC  = "basic"
	AUTH_BEARER = "bearer"
	AUTH_APIKEY = "apiKey"
)

func (g *Connector) doQuery(baseURL string, queryParams, headers, cookies map[string]string, authentication string,
	authContent map[string]string, query string, vars map[string]interface{}) (*resty.Response, error) {

	client := resty.New()

	// corner case
	if authentication == AUTH_APIKEY && authContent["addTo"] == "urlParams" {
		queryParams[authContent["key"]] = authContent["value"]
		authentication = AUTH_NONE
	}

	// build request url
	uri, err := url.Parse(baseURL)
	if err != nil {
		return nil, err
	}
	params := url.Values{}
	for k, v := range queryParams {
		if k != "" {
			params.Set(k, v)
		}
	}
	uri.RawQuery = params.Encode()
	reqURL := uri.String()

	// set authentication
	switch authentication {
	case AUTH_BASIC:
		client.SetBasicAuth(authContent["username"], authContent["password"])
		break
	case AUTH_BEARER:
		client.SetAuthToken(authContent["bearerToken"])
		break
	case AUTH_APIKEY:
		client.SetAuthScheme(authContent["headerPrefix"])
		client.SetAuthToken(authContent["value"])
		break
	case AUTH_NONE:
		break
	}

	queryClient := client.R()

	// set headers
	queryClient.SetHeaders(headers)
	queryClient.SetHeader("Content-Type", "application/json")

	// set cookies
	reqCookies := make([]*http.Cookie, 0, len(cookies))
	for k, v := range cookies {
		reqCookies = append(reqCookies, &http.Cookie{Name: k, Value: v})
	}
	queryClient.SetCookies(reqCookies)

	// set body
	queryClient.SetBody(map[string]interface{}{
		"query":     query,
		"variables": vars,
	})

	// do the query
	resp, err := queryClient.Post(reqURL)
	if err != nil {
		return nil, err
	}

	return resp, nil
}
