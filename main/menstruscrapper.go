package main

import (
	"fmt"
	"time"

	"github.com/lasdesistemas/menstruscrapper"
	preciosclaros "github.com/lasdesistemas/menstruscrapper/precios-claros"
)

func main() {

	fmt.Printf("Comienzo a buscar precios de toallitas y tampones a las %v\n", time.Now().Format("2006-01-02T15:04:05"))
	cliente := preciosclaros.NewClient(&preciosclaros.Espera{}, &preciosclaros.PreciosClarosRestClient{})
	scrapper := menstruscrapper.New(cliente)
	scrapper.GenerarListaDePrecios()
	fmt.Printf("Terminé de buscar precios de toallitas y tampones a las %v\n", time.Now().Format("2006-01-02T15:04:05"))
}
