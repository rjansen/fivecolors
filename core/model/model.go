package model

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/rjansen/fivecolors/core/util"
)

const (
	// Rarity names
	RarityCommon     = "Common"
	RarityUncommon   = "Uncommon"
	RarityRare       = "Rare"
	RarityMythicRare = "Mythic Rare"
	RaritySpecial    = "Special"
	RarityLand       = "Land"
	RarityPromo      = "Promo"
	RarityBonus      = "Bonus"
	// Rarity aliases
	RarityAliasCommon     = "C"
	RarityAliasUncommon   = "U"
	RarityAliasRare       = "R"
	RarityAliasMythicRare = "M"
	RarityAliasSpecial    = "S"
	RarityAliasLand       = "L"
	RarityAliasPromo      = "P"
	RarityAliasBonus      = "B"
)

var (
	// TODO: Thinking better about this strategy, but for the MVP is enough
	// Rarities by Name
	raritiesByName = map[string]Rarity{
		RarityCommon:     NewRarity(RarityCommon, RarityAliasCommon),
		RarityUncommon:   NewRarity(RarityUncommon, RarityAliasUncommon),
		RarityRare:       NewRarity(RarityRare, RarityAliasRare),
		RarityMythicRare: NewRarity(RarityMythicRare, RarityAliasMythicRare),
		RaritySpecial:    NewRarity(RaritySpecial, RarityAliasSpecial),
		RarityLand:       NewRarity(RarityLand, RarityAliasLand),
		RarityPromo:      NewRarity(RarityPromo, RarityAliasPromo),
		RarityBonus:      NewRarity(RarityBonus, RarityAliasBonus),
	}
	// Rarities By Alias
	raritiesByAlias = map[string]Rarity{
		RarityAliasCommon:     raritiesByName[RarityCommon],
		RarityAliasUncommon:   raritiesByName[RarityUncommon],
		RarityAliasRare:       raritiesByName[RarityRare],
		RarityAliasMythicRare: raritiesByName[RarityMythicRare],
		RarityAliasSpecial:    raritiesByName[RaritySpecial],
		RarityAliasLand:       raritiesByName[RarityLand],
		RarityAliasPromo:      raritiesByName[RarityPromo],
		RarityAliasBonus:      raritiesByName[RarityBonus],
	}
)

type lifeCycle struct {
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt,omitempty"`
	DeletedAt time.Time `json:"deletedAt,omitempty"`
}

func (l lifeCycle) String() string {
	var updatedAt string
	if !l.UpdatedAt.IsZero() {
		updatedAt = l.UpdatedAt.Format(time.RFC3339)
	}
	var deletedAt string
	if !l.DeletedAt.IsZero() {
		deletedAt = l.DeletedAt.Format(time.RFC3339)
	}
	return fmt.Sprintf("CreatedAt:%s UpdatedAt:%s DeletedAt:%s", l.CreatedAt.Format(time.RFC3339), updatedAt, deletedAt)
}

type SetAsset map[string]string

func (a SetAsset) String() string {
	var b strings.Builder
	b.WriteString("model.SetAsset:{ ")
	for k, v := range a {
		fmt.Fprintf(&b, " %s:%s", k, v)
	}
	b.WriteString(" }")
	return b.String()
}

func (a SetAsset) Value() (driver.Value, error) {
	j, err := json.Marshal(a)
	return j, err
}

func (a *SetAsset) Scan(src interface{}) error {
	source, ok := src.([]byte)
	if !ok {
		return errors.New("InvalidSourceType: != []byte")
	}

	err := json.Unmarshal(source, a)
	if err != nil {
		return err
	}
	return nil
}

type Rarity struct {
	ID    string `json:"id"`
	Name  string `json:"name"`
	Alias string `json:"alias"`
	lifeCycle
}

func (r Rarity) String() string {
	return fmt.Sprintf("model.Rarity:{ ID:%s Name:%s Alias:%s %s }", r.ID, r.Name, r.Alias, r.lifeCycle)
}

type Set struct {
	ID    string   `json:"id"`
	Name  string   `json:"name"`
	Alias string   `json:"alias"`
	Asset SetAsset `json:"asset"`
	Cards []Card   `json:"cards"`
	lifeCycle
}

func (s Set) String() string {
	return fmt.Sprintf("model.Set:{ ID:%s Name:%s Alias:%s Asset:%s %s }", s.ID, s.Name, s.Alias, s.Asset, s.lifeCycle)
}

func (s Set) HasAsset(r Rarity) error {
	if _, exists := s.Asset[r.ID]; !exists {
		return ErrNotFound
	}
	return nil
}

func (s *Set) SetAsset(r Rarity, id string) {
	if s.Asset == nil {
		s.Asset = NewSetAsset()
	}
	s.Asset[r.ID] = NewAssetID(id)
}

func (s *Set) AddCard(c Card) {
	if s.Cards == nil {
		s.Cards = NewSetCards()
	}
	s.Cards = append(s.Cards, c)
}

type CardData map[string]string

func (d CardData) String() string {
	var b strings.Builder
	b.WriteString("model.CardData:{ ")
	for k, v := range d {
		fmt.Fprintf(&b, " %s:%s", k, v)
	}
	b.WriteString(" }")
	return b.String()
}

type Card struct {
	ID            string   `json:"id"`
	Name          string   `json:"name"`
	Types         []string `json:"types"`
	Costs         []string `json:"costs"`
	NumberCost    float32  `json:"numberCost"`
	IDExternal    string   `json:"idExternal"`
	OrderExternal string   `json:"orderExternal"`
	IDRarity      string   `json:"-"`
	Rarity        Rarity   `json:"rarity"`
	IDSet         string   `json:"-"`
	Set           Set      `json:"set"`
	IDAsset       string   `json:"idAsset"`
	Rules         []string `json:"rules"`
	Rate          float32  `json:"rate"`
	RateVotes     int      `json:"rateVotes"`
	Artist        string   `json:"artist"`
	Flavor        string   `json:"flavor"`
	Data          CardData `json:"data"`
	lifeCycle
}

func (c Card) String() string {
	return fmt.Sprintf("model.Card:{ ID:%s Name:%s Types:%q Costs:%q IDExternal:%s OrderExternal:%s Rarity:%s Set:%s IDAsset:%s Rules:%q Rate:%f RateVotes:%d Artist:%s Flavor:%s Data:%s %s }", c.ID, c.Name, c.Types, c.Costs, c.IDExternal, c.OrderExternal, c.Rarity, c.Set, c.IDAsset, c.Rules, c.Rate, c.RateVotes, c.Artist, c.Flavor, c.Data, c.lifeCycle)
}

func (c *Card) AddCost(o string) {
	if c.Costs == nil {
		c.Costs = newCardCosts()
	}
	c.Costs = append(c.Costs, o)
}

func (c *Card) AddRule(r string) {
	if c.Rules == nil {
		c.Rules = newCardRules()
	}
	c.Rules = append(c.Rules, r)
}

func NewRarity(name, alias string) Rarity {
	return Rarity{
		ID:    NewRarityID(name),
		Name:  name,
		Alias: alias,
		lifeCycle: lifeCycle{
			CreatedAt: time.Now(),
		},
	}
}

func NewRarityID(n string) string {
	return util.Sha1f("Rarity#%s", n)
}

func RarityByAlias(a string) (*Rarity, error) {
	if r, exists := raritiesByAlias[a]; exists {
		return &r, nil
	}
	return nil, ErrNotFound
}

func CreateSet(name string) *Set {
	return NewSet(name, strings.ToUpper(name[:3]))
}

func NewSet(name, alias string) *Set {
	return &Set{
		ID:    util.Sha1f("Set#%s", name),
		Name:  name,
		Alias: alias,
		Asset: NewSetAsset(),
		lifeCycle: lifeCycle{
			CreatedAt: time.Now(),
		},
	}
}

func NewSetAsset() SetAsset {
	return SetAsset{}
}

func NewAssetID(id string) string {
	return util.Sha1f("Asset#%s", id)
}

func NewCard(set Set, name, qsExternal, types, qsAsset, computedCost string, rarity Rarity) *Card {
	id := util.Sha1f("Card#%s#%s#%s", set.ID, name, qsExternal)
	numberCost, _ := strconv.ParseFloat(computedCost, 32)
	return &Card{
		ID:         id,
		Name:       name,
		Types:      []string{types},
		IDExternal: qsExternal,
		NumberCost: float32(numberCost),
		Rarity:     rarity,
		Set:        set,
		IDAsset:    NewAssetID(qsAsset),
		Costs:      newCardCosts(),
		Rules:      newCardRules(),
	}
}

func NewSetCards() []Card {
	return make([]Card, 0, 300)
}

func newCardCosts() []string {
	return make([]string, 0, 20)
}

func newCardRules() []string {
	return make([]string, 0, 20)
}
