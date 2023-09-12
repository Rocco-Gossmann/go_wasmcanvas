package go_wasmcanvas

type CanvasPanic struct {
	Msg     string
	Subject string
	Value   any
	Allowed any
}

func (p CanvasPanic) GetAllowed() any { return p.Allowed }

// Implement go_throwable.throwable
func (p CanvasPanic) GetValue() any      { return p.Value }
func (p CanvasPanic) GetMessage() string { return p.Msg }
func (p CanvasPanic) GetSubject() string { return p.Subject }
