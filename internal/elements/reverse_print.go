package elements

type ReversePrint struct {
	Value bool
}

func (p *ReversePrint) IsReversePrint() bool {
	return p.Value
}
