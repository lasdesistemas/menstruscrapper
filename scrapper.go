package menstruscrapper

type Client interface {
}

type Scrapper struct {
	client Client
}

func New(c Client) *Scrapper {
	return &Scrapper{}
}

func (s *Scrapper) GenerarListaDePrecios() string {
	/*
		sucursales := s.client.ObtenerSucursales()

		tampones := s.client.ObtenerListaDeTampones(sucursales)

		toallitas := s.client.ObtenerListaDeToallitas(sucursales)

		preciosTampones := s.client.ObtenerListaDePrecios(sucursales, tampones)
		preciosToallitas := s.client.ObtenerListaDePrecios(sucursales, toallitas)

		rutaCsv := s.generarCsv(preciosTampones, preciosToallitas)

		return rutaCsv*/

	return ""
}
