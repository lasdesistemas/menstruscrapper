package menstruscrapper

import (
	"fmt"
	"io/ioutil"
	"strings"

	preciosclaros "github.com/lasdesistemas/menstruscrapper/precios-claros"
)

type Client interface {
	ObtenerSucursales() ([]string, error)
	ObtenerListaDeTampones(sucursales []string) ([]int, error)
	ObtenerListaDeToallitas(sucursales []string) ([]int, error)
	ObtenerListaDePrecios(sucursales []string, productos []int, categoria string) ([]*preciosclaros.Producto, error)
}

type Scrapper struct {
	client Client
}

func New(c Client) *Scrapper {
	return &Scrapper{c}
}

func (s *Scrapper) GenerarListaDePrecios() string {

	fmt.Println("Obteniendo sucursales de todo el país..")
	sucursales, errSucursales := s.client.ObtenerSucursales()
	if errSucursales != nil {
		fmt.Printf("No se pudieron obtener las sucursales: %s\n", errSucursales.Error())
		return ""
	}
	fmt.Printf("Se obtuvieron %v sucursales\n", len(sucursales))

	fmt.Println("Obteniendo los ids de todos los productos de la categoría tampones..")
	tampones, errTampones := s.client.ObtenerListaDeTampones(sucursales)
	if errTampones != nil {
		fmt.Printf("No se pudo obtener la lista de ids de tampones: %s\n", errTampones.Error())
		return ""
	}
	fmt.Printf("Se obtuvieron %v ids de productos de la categoría tampones\n", len(tampones))

	fmt.Println("Obteniendo los ids de todos los productos de la categoría toallitas..")
	toallitas, errToallitas := s.client.ObtenerListaDeToallitas(sucursales)
	if errToallitas != nil {
		fmt.Printf("No se pudo obtener la lista de ids de toallitas: %s\n", errToallitas.Error())
		return ""
	}
	fmt.Printf("Se obtuvieron %v ids de productos de la categoría toallitas\n", len(toallitas))

	fmt.Printf("Obteniendo la lista de precios para %v sucursales de %v productos de la categoría tampones\n", len(sucursales), len(tampones))
	preciosTampones, errPreciosTampones := s.client.ObtenerListaDePrecios(sucursales, tampones, "tampones")
	if errPreciosTampones != nil {
		fmt.Printf("No se pudo obtener la lista de precios de tampones: %s\n", errPreciosTampones.Error())
		return ""
	}
	fmt.Printf("Se obtuvieron %v precios de productos de la categoría tampones\n", len(preciosTampones))

	fmt.Printf("Obteniendo la lista de precios para %v sucursales de %v productos de la categoría toallitas\n", len(sucursales), len(toallitas))
	preciosToallitas, errPreciosToallitas := s.client.ObtenerListaDePrecios(sucursales, toallitas, "toallitas")
	if errPreciosToallitas != nil {
		fmt.Printf("No se pudo obtener la lista de precios de toallitas: %s\n", errPreciosToallitas.Error())
		return ""
	}
	fmt.Printf("Se obtuvieron %v precios de productos de la categoría toallitas\n", len(preciosToallitas))

	fmt.Println("Generando archivo con resultados...")
	rutaCsv, errCsv := s.generarCsv(preciosTampones, preciosToallitas)
	if errCsv != nil {
		fmt.Printf("No se pudo generar el csv con los precios: %s\n", errCsv.Error())
		return ""
	}
	fmt.Printf("Archivo generado: %s\n", rutaCsv)

	return rutaCsv
}

func (s *Scrapper) generarCsv(preciosTampones, preciosToallitas []*preciosclaros.Producto) (string, error) {

	listaDePrecios := "Categoría,Marca,Nombre,Presentación,Comercio,Sucursal,Dirección,Localidad,Provincia,Precio de lista\n"

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
			producto.Sucursal, producto.Direccion, producto.Localidad, producto.Provincia, fmt.Sprintf("%.2f", producto.PrecioDeLista)}, ",")

		listaDePrecios = listaDePrecios + linea + "\n"
	}

	return listaDePrecios
}
