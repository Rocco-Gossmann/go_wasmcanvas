package go_wasmcanvas

type canvasPanic struct {
	msg     string
	subject string
	value   any
	allowed any
}

type CanvasPanic canvasPanic

func (p *canvasPanic) GetAllowed() any { return p.allowed }

// Implement go_throwable.throwable
func (p *canvasPanic) GetValue() any      { return p.value }
func (p *canvasPanic) GetMessage() string { return p.msg }
func (p *canvasPanic) GetSubject() string { return p.subject }
