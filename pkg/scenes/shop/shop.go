package shop

import (
	"firefly-jam-2026/assets"
	"firefly-jam-2026/pkg/util"

	"github.com/firefly-zero/firefly-go/firefly"
)

type Shop struct {
	dpad4Old firefly.DPad4

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

	s.Items = append(s.Items, Item{
		Price:    0,
		Quantity: -1,
		Bg:       s.sellBG,
		Icon:     assets.ShopItem[7],
	})

	s.AddDrugItem(100, -1, assets.ShopItem[7])
	s.AddDrugItem(100, -1, assets.ShopItem[7])
	s.AddDrugItem(100, -1, assets.ShopItem[7])
}

func (s *Shop) Update() {
	if pad, ok := firefly.ReadPad(firefly.GetMe()); ok {
		dpad4 := pad.DPad4()
		justPressed := dpad4.JustPressed(s.dpad4Old)
		if justPressed != firefly.DPad4None {
			s.handleInput(justPressed)
		}
		s.dpad4Old = dpad4
	} else {
		s.dpad4Old = firefly.DPad4None
	}
	s.selectedAnim.Update()
}

func (s *Shop) handleInput(justPressed firefly.DPad4) {
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

func (s *Shop) Render() {
	startPos := firefly.P(130, 28)
	offset := firefly.S(27, 38)
	for i, item := range s.Items {
		pos := startPos.Add(firefly.P((i%4)*offset.W, (i/4)*offset.H))
		item.Bg.Draw(pos)
		item.Icon.Draw(pos)
		if s.Selected == i {
			s.selectedAnim.Draw(pos)
		}
	}
}

func (s *Shop) AddDrugItem(price, quantity int, icon firefly.SubImage) {
	s.Items = append(s.Items, Item{
		Price:    price,
		Quantity: quantity,
		Icon:     icon,
		Bg:       util.RandomSliceElem(s.drugBGs),
	})
}

func (s *Shop) AddItem(price, quantity int, icon firefly.SubImage) {
	s.Items = append(s.Items, Item{
		Price:    price,
		Quantity: quantity,
		Icon:     icon,
		Bg:       util.RandomSliceElem(s.itemBGs),
	})
}

type Item struct {
	Price    int
	Quantity int
	Icon     firefly.SubImage
	Bg       firefly.SubImage
}
