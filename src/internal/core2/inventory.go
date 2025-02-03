package core2

type Inventory struct {
	RootModule *Module
}

func NewInventory() *Inventory {
	return &Inventory{}
}
