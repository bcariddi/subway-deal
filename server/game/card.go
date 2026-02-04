package game

// Card is the interface all card types implement
type Card interface {
	GetID() string
	GetType() CardType
	GetName() string
	GetValue() int // Money value for banking
}

type CardType string

const (
	CardTypeProperty CardType = "property"
	CardTypeWildcard CardType = "wildcard"
	CardTypeAction   CardType = "action"
	CardTypeRent     CardType = "rent"
	CardTypeMoney    CardType = "money"
)

// BaseCard contains fields common to all cards
type BaseCard struct {
	ID    string   `json:"id"`
	Type  CardType `json:"type"`
	Name  string   `json:"name"`
	Value int      `json:"value"`
}

func (c BaseCard) GetID() string     { return c.ID }
func (c BaseCard) GetType() CardType { return c.Type }
func (c BaseCard) GetName() string   { return c.Name }
func (c BaseCard) GetValue() int     { return c.Value }

// PropertyCard represents a subway line property
type PropertyCard struct {
	BaseCard
	Color    string      `json:"color"`
	ColorHex string      `json:"colorHex"`
	SetSize  int         `json:"setSize"`
	Position int         `json:"position"`
	Rent     map[int]int `json:"rent"` // count -> rent amount
}

func (c PropertyCard) GetRent(propertyCount int) int {
	if rent, ok := c.Rent[propertyCount]; ok {
		return rent
	}
	return 0
}

// WildcardCard can be used as multiple colors
type WildcardCard struct {
	BaseCard
	Colors       []string `json:"colors"`
	CurrentColor string   `json:"currentColor"`
	Description  string   `json:"description"`
}

func (c *WildcardCard) FlipColor() {
	for i, color := range c.Colors {
		if color == c.CurrentColor {
			c.CurrentColor = c.Colors[(i+1)%len(c.Colors)]
			return
		}
	}
}

func (c WildcardCard) CanBeColor(color string) bool {
	for _, clr := range c.Colors {
		if clr == color {
			return true
		}
	}
	return false
}

// ActionCard performs special effects
type ActionCard struct {
	BaseCard
	Effect   string `json:"effect"`
	MTATheme string `json:"mtaTheme"`
}

// RentCard charges rent to other players
type RentCard struct {
	BaseCard
	Colors     []string `json:"colors"`
	ColorNames []string `json:"colorNames"`
	Target     string   `json:"target"` // "all" or "one"
}

func (c RentCard) IsWildRent() bool {
	return len(c.Colors) == 0 || c.Name == "Wild Rent"
}

// MoneyCard is pure currency
type MoneyCard struct {
	BaseCard
	Denomination int    `json:"denomination"`
	DisplayValue string `json:"displayValue"`
	Theme        string `json:"theme"`
}
