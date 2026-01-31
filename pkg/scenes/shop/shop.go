package shop

import (
	"slices"

	"github.com/applejag/epic-wizard-firefly-gladiators/assets"
	"github.com/applejag/epic-wizard-firefly-gladiators/pkg/scenes"
	"github.com/applejag/epic-wizard-firefly-gladiators/pkg/state"
	"github.com/applejag/epic-wizard-firefly-gladiators/pkg/util"

	"github.com/firefly-zero/firefly-go/firefly"
)

type Shop struct {
	Selected int
	Items    []Item

	selectedAnim util.AnimatedSheet
	drugBGs      []firefly.SubImage
	itemBGs      []firefly.SubImage
	sellBG       firefly.SubImage
}

func (s *Shop) Boot() {
	s.selectedAnim = assets.ShopItem[8:10].Animated(4)
	s.itemBGs = assets.ShopItem[0:4]
	s.drugBGs = assets.ShopItem[5:6]
	s.sellBG = assets.ShopItem[4]
}

func (s *Shop) Update() {
	if justPressed := state.Input.JustPressedDPad4(); justPressed != firefly.DPad4None {
		s.handleInputDPad4(justPressed)
	}
	if justPressed := state.Input.JustPressedButtons(); justPressed.Any() {
		s.handleInputButtons(justPressed)
	}

	s.selectedAnim.Update()
}

func (s *Shop) handleInputDPad4(justPressed firefly.DPad4) {
	switch justPressed {
	case firefly.DPad4Right:
		if s.Selected < len(s.Items)-1 {
			s.Selected++
		}
	case firefly.DPad4Left:
		if s.Selected > 0 {
			s.Selected--
		}
	case firefly.DPad4Down:
		s.Selected = min(s.Selected+4, len(s.Items)-1)
	case firefly.DPad4Up:
		s.Selected = max(s.Selected-4, 0)
	}
}

func (s *Shop) handleInputButtons(justPressed firefly.Buttons) {
	if justPressed.S && s.Selected >= 0 && s.Selected < len(s.Items) {
		item := s.Items[s.Selected]
		if item.Price > state.Game.Money {
			return
		}
		state.Game.Money -= item.Price

		var buf [len("buy 1x: ") + LongestItemKind]byte
		written := util.ConcatInto(buf[:], "buy 1x: ", item.Kind.String())
		util.LogDebugBytes(buf[:written])

		switch item.Kind {
		case ItemDrug:
		case ItemFirefly:
			state.Game.AddFirefly()
		case ItemHat:
		case ItemSell:
		default:
			panic("unexpected shop.ItemKind")
		}
		if s.Items[s.Selected].Quantity > 0 {
			s.Items[s.Selected].Quantity--
			if s.Items[s.Selected].Quantity <= 0 {
				s.Items = slices.Delete(s.Items, s.Selected, s.Selected+1)
			}
		}
		state.Game.Save()
	}
	if justPressed.E && len(state.Game.Fireflies) > 0 {
		scenes.SwitchScene(scenes.Field)
	}
}

func (s *Shop) Render() {
	startPos := firefly.P(130, 28)
	offset := firefly.S(27, 38)
	for i := range s.Items {
		item := &s.Items[i]
		pos := startPos.Add(firefly.P((i%4)*offset.W, (i/4)*offset.H))
		item.Bg.Draw(pos)
		item.Icon.Draw(pos)
		textWidth := len(item.priceStr) * assets.FontPico8_4x6.CharWidth()
		priceColor := firefly.ColorYellow
		if state.Game.Money >= item.Price {
			priceColor = firefly.ColorGreen
		}
		assets.FontPico8_4x6.Draw(item.priceStr, pos.Add(firefly.P(offset.W/2-textWidth/2, offset.H-6)), priceColor)
		if item.Quantity > 0 {
			var buf [5]byte
			buf[0] = 'x'
			written := util.FormatIntInto(buf[1:], item.Quantity)
			quantityStr := string(buf[:1+written])
			quantityWidth := len(quantityStr) * assets.FontPico8_4x6.CharWidth()
			assets.FontPico8_4x6.Draw(quantityStr, pos.Add(firefly.P(offset.W-quantityWidth-3, 7)), firefly.ColorDarkGray)
		}
		if s.Selected == i {
			s.selectedAnim.Draw(pos)
		}
	}
	if len(state.Game.Fireflies) > 0 {
		assets.Exit.Draw(firefly.P(firefly.Width-assets.Exit.Width()+4, 1))
	}
}

func (s *Shop) AddDrugItem(price, quantity int, icon firefly.SubImage) {
	s.Items = append(s.Items, Item{
		Kind:     ItemDrug,
		Price:    price,
		priceStr: formatPrice(price),
		Quantity: quantity,
		Icon:     icon,
		Bg:       util.RandomSliceElem(s.drugBGs),
	})
}

func (s *Shop) AddFireflyItem(price, quantity int, icon firefly.SubImage) {
	s.Items = append(s.Items, Item{
		Kind:     ItemFirefly,
		Price:    price,
		priceStr: formatPrice(price),
		Quantity: quantity,
		Icon:     icon,
		Bg:       util.RandomSliceElem(s.itemBGs),
	})
}

func (s *Shop) AddSellItem() {
	s.Items = append(s.Items, Item{
		Kind:     ItemSell,
		Price:    0,
		Quantity: -1,
		Bg:       s.sellBG,
		Icon:     assets.ShopItem[7],
	})
}

func formatPrice(price int) string {
	if price == 0 {
		return "FREE"
	}
	var buf [1 + 3]byte
	buf[0] = '$'
	written := util.FormatIntInto(buf[1:], price)
	return string(buf[:1+written])
}

type Item struct {
	Kind     ItemKind
	Price    int
	priceStr string
	Quantity int
	Icon     firefly.SubImage
	Bg       firefly.SubImage
}

type ItemKind byte

const (
	ItemNone ItemKind = iota
	ItemSell
	ItemFirefly
	ItemHat
	ItemDrug

	LongestItemKind = 7
)

var AllItemKinds = []ItemKind{
	ItemNone,
	ItemSell,
	ItemFirefly,
	ItemHat,
	ItemDrug,
}

func (k ItemKind) String() string {
	switch k {
	case ItemDrug:
		return "drug"
	case ItemFirefly:
		return "firefly"
	case ItemHat:
		return "hat"
	case ItemNone:
		return "none"
	case ItemSell:
		return "sell"
	default:
		panic("unexpected shop.ItemKind")
	}
}
