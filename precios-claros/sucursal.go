package preciosclaros

type Sucursal struct {
	Id string
	Comercio string `json:"comercioRazonSocial"`
	Nombre string `json:"sucursalNombre"`
	Direccion string
	Localidad string
	PreciosProducto *PreciosProducto

}

type PreciosProducto struct {
	PrecioLista float32
}
