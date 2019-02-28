package menstruscrapper

import (
	"github.com/lasdesistemas/menstruscrapper/precios-claros"
	"fmt"
	"strings"
	"io/ioutil"
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
		fmt.Printf("No se pudieron obtener las sucursales: %s", errSucursales.Error())
		return ""
	}

	tampones, errTampones := s.client.ObtenerListaDeTampones(sucursales)
	if errTampones != nil {
		fmt.Printf("No se pudo obtener la lista de ids de tampones: %s", errTampones.Error())
		return ""
	}

	toallitas, errToallitas := s.client.ObtenerListaDeToallitas(sucursales)
	if errToallitas != nil {
		fmt.Printf("No se pudo obtener la lista de ids de toallitas: %s", errToallitas.Error())
		return ""
	}

	preciosTampones, errPreciosTampones := s.client.ObtenerListaDePrecios(sucursales, tampones)
	if errPreciosTampones != nil {
		fmt.Printf("No se pudo obtener la lista de precios de tampones: %s", errPreciosTampones.Error())
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

	listaDePrecios := "Categoría,Marca,Nombre,Presentación,Comercio,Sucursal,Dirección,Localidad,Precio de lista\n"

	listaDePrecios = generarListaDePrecios(preciosTampones, listaDePrecios)
	listaDePrecios = generarListaDePrecios(preciosToallitas, listaDePrecios)

	errWrite := ioutil.WriteFile("precios-gestion-menstrual.csv", []byte(listaDePrecios), 0644)

	if errWrite != nil {
		return "", fmt.Errorf("No se pudo escribir en el csv: %s", errWrite.Error())
	}

	return "precios-gestion-menstrual.csv", nil
}

func generarListaDePrecios(productos []*preciosclaros.Producto, listaDePrecios string) string {

	for _, producto := range productos {

		linea := strings.Join([]string{producto.Categoria, producto.Marca, producto.Nombre, producto.Presentacion, producto.Comercio,
			producto.Sucursal, producto.Direccion, producto.Localidad, fmt.Sprintf("%.2f", producto.PrecioDeLista)}, ",")

		listaDePrecios = listaDePrecios + linea + "\n"
	}

	return listaDePrecios
}
