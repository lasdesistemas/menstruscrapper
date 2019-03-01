package preciosclaros

import "net/http"

type RestClient interface {
	Get(url string) (*http.Response, error)
}

type PreciosClarosRestClient struct {

}

func (rc *PreciosClarosRestClient) Get(url string) (*http.Response, error) {
	return http.Get(url)
}



