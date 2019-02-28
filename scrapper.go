package menstruscrapper

import (
	"github.com/lasdesistemas/menstruscrapper/precios-claros"
	"fmt"
)

type Client interface {
	ObtenerSucursales() ([]string, error)
	ObtenerListaDeTampones(sucursales []string) ([]int, error)
	ObtenerListaDeToallitas(sucursales []string) ([]int, error)
	ObtenerListaDePrecios(sucursales []string, productos []int) ([]*preciosclaros.Producto, error)
}

type Scrapper struct {
	client Client
}

func New(c Client) *Scrapper {
	return &Scrapper{c}
}

func (s *Scrapper) GenerarListaDePrecios() string {

	sucursales, errSucursales := s.client.ObtenerSucursales()
	if errSucursales != nil {
		fmt.Errorf("No se pudieron obtener las sucursales: %s", errSucursales.Error())
		return ""
	}

	tampones, errTampones := s.client.ObtenerListaDeTampones(sucursales)
	if errTampones != nil {
		fmt.Errorf("No se pudo obtener la lista de ids de tampones: %s", errTampones.Error())
		return ""
	}

	toallitas, errToallitas := s.client.ObtenerListaDeToallitas(sucursales)
	if errToallitas != nil {
		fmt.Errorf("No se pudo obtener la lista de ids de toallitas: %s", errToallitas.Error())
		return ""
	}

	preciosTampones, errPreciosTampones := s.client.ObtenerListaDePrecios(sucursales, tampones)
	if errPreciosTampones != nil {
		fmt.Errorf("No se pudo obtener la lista de precios de tampones: %s", errPreciosTampones.Error())
		return ""
	}

	preciosToallitas, errPreciosToallitas := s.client.ObtenerListaDePrecios(sucursales, toallitas)
	if errPreciosToallitas != nil {
		fmt.Errorf("No se pudo obtener la lista de precios de toallitas: %s", errPreciosToallitas.Error())
		return ""
	}

	rutaCsv, errCsv := s.generarCsv(preciosTampones, preciosToallitas)
	if errCsv != nil {
		fmt.Errorf("No se pudo generar el csv con los precios: %s", errCsv.Error())
		return ""
	}

	return rutaCsv
}

func (s *Scrapper) generarCsv(preciosTampones, preciosToallitas []*preciosclaros.Producto) (string, error) {

	return "", nil
}
