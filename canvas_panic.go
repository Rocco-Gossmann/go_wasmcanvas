package go_wasmcanvas

type canvasPanic struct {
	msg string
}

func (p *canvasPanic) GetValue() any      { return nil }
func (p *canvasPanic) GetMessage() string { return p.msg }
func (p *canvasPanic) GetSubject() string { return "" }
