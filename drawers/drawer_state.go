package drawers

import "github.com/ingridhq/zebrash/elements"

type DrawerState struct {
	AutoPosX float64
	AutoPosY float64
}

func (state *DrawerState) UpdateAutomaticTextPosition(text *elements.TextField, w, scaleX float64) {
	if !text.Position.CalculateFromBottom {
		return
	}

	x := float64(text.Position.X)
	y := float64(text.Position.Y)

	if !text.Position.AutomaticPosition {
		state.AutoPosX = x
		state.AutoPosY = y
	}

	switch text.Font.Orientation {
	case elements.FieldOrientation90:
		state.AutoPosY += w * scaleX
	case elements.FieldOrientation180:
		state.AutoPosX -= w * scaleX
	case elements.FieldOrientation270:
		state.AutoPosY -= w * scaleX
	default:
		state.AutoPosX += w * scaleX
	}
}

func (state *DrawerState) GetTextPosition(text *elements.TextField) (float64, float64) {
	if text.Position.AutomaticPosition {
		return state.AutoPosX, state.AutoPosY
	}

	return float64(text.Position.X), float64(text.Position.Y)
}
