package domain

// Kind — вид объекта на балансе (ADR-0004).
type Kind struct {
	Code  string // префикс инвентарного номера: "SC"
	Title string // название для вывода: "электросамокат"
}

// Виды объектов на балансе «ВышКат».
var (
	KindElectricScooter = Kind{Code: "SC", Title: "электросамокат"}
	KindBicycle         = Kind{Code: "BK", Title: "велосипед"}
	KindElectricBike    = Kind{Code: "EB", Title: "электровелосипед"}
	KindHelmet          = Kind{Code: "HL", Title: "шлем"}
	KindDockStation     = Kind{Code: "DS", Title: "док-станция"}
	KindChargingCabinet = Kind{Code: "CC", Title: "зарядный шкаф"}
)
