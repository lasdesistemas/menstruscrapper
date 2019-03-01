package preciosclaros

import (
	"fmt"
	"net/http"
	"io/ioutil"
	"encoding/json"
)

type PreciosClarosClient struct {
	restClient RestClient
}

func NewClient(rc RestClient) *PreciosClarosClient {
	return &PreciosClarosClient{rc}
}

const (
	host = "https://d3e6htiiul5ek9.cloudfront.net/prueba"
	sucursales = "/sucursales"
)

func (pc *PreciosClarosClient) ObtenerSucursales() ([]string, error) {

	response, err := pc.restClient.Get(host + sucursales)
	defer response.Body.Close()

	// Valida el resultado
	if err != nil {
		return nil, err
	}

	if response.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("el pedido dio status: %v", response.StatusCode)
	}

	// Obtengo la respuesta
	bodyBytes, errRead := ioutil.ReadAll(response.Body)

	if errRead != nil {
		return nil, fmt.Errorf("error al leer la respuesta: %v", errRead)
	}

	respuesta := struct {
		TotalPagina int `json:"total_pagina"`
		Total int `json:"total"`
		Sucursales []*Sucursal
	}{0, 0, []*Sucursal{}}

	json.Unmarshal(bodyBytes, &respuesta)

	// Convierto a lista de string
	sucursales := []string{}
	for _, sucursal := range respuesta.Sucursales {
		sucursales = append(sucursales, sucursal.Id)
	}

	return sucursales, nil
}