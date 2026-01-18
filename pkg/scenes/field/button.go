package field

import (
	"github.com/applejag/firefly-jam-2026/assets"
	"github.com/firefly-zero/firefly-go/firefly"
	"github.com/orsinium-labs/tinymath"
)

const ButtonShakeDuration = 45

type Button struct {
	text     string
	Disabled bool
	shake    int
}

func NewButton(text string) Button {
	return Button{
		text: text,
	}
}

func (b *Button) Update() {
	if b.shake > 0 {
		b.shake--
	}
}

func (b *Button) Render(point firefly.Point, isFocused bool) {
	prefix := "- "
	color := firefly.ColorGray
	if isFocused {
		prefix = "> "
		color = firefly.ColorBlack
	}
	if b.Disabled {
		color = firefly.ColorLightGray
	}
	assets.FontPico8_4x6.Draw(prefix, point, color)
	if b.text != "" {
		if b.shake > 0 {
			t := float32(b.shake) / ButtonShakeDuration
			point = point.Add(firefly.P(int(tinymath.Sin(t*45)*t*4), 0))
		}
		textPoint := point.Add(firefly.P(assets.FontPico8_4x6.LineWidth(prefix), 0))
		assets.FontPico8_4x6.Draw(b.text, textPoint, color)
		if b.Disabled {
			// Draw strikethrough
			firefly.DrawLine(textPoint, textPoint.Add(firefly.P(
				assets.FontPico8_4x6.LineWidth(b.text),
				-assets.FontEG_6x9.CharHeight()/2,
			)), firefly.L(firefly.ColorGray, 1))
		}
	}
}

func (b *Button) Shake() {
	b.shake = ButtonShakeDuration
}
