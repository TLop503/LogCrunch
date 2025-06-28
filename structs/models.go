package structs

type Mode int

const (
	NormalMode Mode = iota
	EditingMode
)

type Hud_model struct {
	Ips       []string
	Aliases   map[string]string
	Mode      Mode
	InputText string
	Cursor    int
}
