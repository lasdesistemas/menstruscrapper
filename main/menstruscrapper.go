package main

import (
	"github.com/lasdesistemas/menstruscrapper"
	preciosclaros "github.com/lasdesistemas/menstruscrapper/precios-claros"
)

func main() {

	cliente := preciosclaros.NewClient(&preciosclaros.Espera{}, &preciosclaros.PreciosClarosRestClient{})
	scrapper := menstruscrapper.New(cliente)
	scrapper.GenerarListaDePrecios()
}
