package data

import (
	"fmt"
	"github.com/rjansen/l"
	"github.com/rjansen/raizel"
	"strconv"
	"strings"
)

type Query struct {
	Restrictions []string
	Values       []interface{}
	SQL          string
	Order        string
}

type CardQuery struct {
	Query
	//Result Fields
	Result []Card
	//Filter Fields
	Hydrate      string
	RegexName    string
	RegexCost    string
	NotRegexCost string
	RegexType    string
	NotRegexType string
	RegexText    string
	NotRegexText string
	IDExpansion  string
	Number       string
	InventoryQtd string
}

func (q *CardQuery) Build() error {
	idxParam := 0
	if q.RegexName != "" {
		idxParam++
		q.Restrictions = append(q.Restrictions, fmt.Sprintf("c.name ~* $%d", idxParam))
		q.Values = append(q.Values, q.RegexName)
	}
	if q.RegexCost != "" {
		idxParam++
		q.Restrictions = append(q.Restrictions, fmt.Sprintf("c.manacost_label ~* $%d", idxParam))
		q.Values = append(q.Values, q.RegexCost)
	}
	if q.NotRegexCost != "" {
		idxParam++
		q.Restrictions = append(q.Restrictions, fmt.Sprintf("not c.manacost_label ~* $%d", idxParam))
		q.Values = append(q.Values, q.NotRegexCost)
	}
	if q.RegexType != "" {
		idxParam++
		q.Restrictions = append(q.Restrictions, fmt.Sprintf("c.type_label ~* $%d", idxParam))
		q.Values = append(q.Values, q.RegexType)
	}
	if q.NotRegexType != "" {
		idxParam++
		q.Restrictions = append(q.Restrictions, fmt.Sprintf("not c.type_label ~* $%d", idxParam))
		q.Values = append(q.Values, q.NotRegexType)
	}
	if q.RegexText != "" {
		idxParam++
		q.Restrictions = append(q.Restrictions, fmt.Sprintf("c.text ~* $%d", idxParam))
		q.Values = append(q.Values, q.RegexText)
	}
	if q.NotRegexText != "" {
		idxParam++
		q.Restrictions = append(q.Restrictions, fmt.Sprintf("not c.text ~* $%d", idxParam))
		q.Values = append(q.Values, q.NotRegexText)
	}
	if q.IDExpansion != "" {
		if idExpansion, convertErr := strconv.Atoi(q.IDExpansion); convertErr == nil {
			idxParam++
			q.Restrictions = append(q.Restrictions, fmt.Sprintf("c.id_expansion = $%d", idxParam))
			q.Values = append(q.Values, idExpansion)
		} else {
			l.Warn("AnonDeckHandler.ExpansionParamErr", l.String("Parameter", q.IDExpansion), l.Err(convertErr))
		}
	}
	if q.Number != "" {
		idxParam++
		q.Restrictions = append(q.Restrictions, fmt.Sprintf("c.multiverse_number = $%d", idxParam))
		q.Values = append(q.Values, q.Number)
	}
	if q.InventoryQtd != "" {
		idxParam++
		q.Restrictions = append(q.Restrictions, fmt.Sprintf("coalesce(i.quantity, 0) >= $%d", idxParam))
		q.Values = append(q.Values, q.InventoryQtd)
	}

	query :=
		`
            select c.id, c.multiverseid, c.multiverse_number, c.name, c.label, coalesce(c.text, ''),
                coalesce(c.manacost_label, ''), coalesce(c.combatpower_label, ''), c.type_label,
                c.id_rarity, coalesce(c.flavor, ''), c.artist, c.rate, c.rate_votes, c.id_asset,
                e.id, e.name, e.label, a.id_asset,
				coalesce(i.id_inventory, 0), coalesce(i.quantity, 0)
            from card c
                left join expansion e on c.id_expansion = e.id
                left join expansion_asset a on a.id_expansion = e.id and a.id_rarity = c.id_rarity
				left join inventory_card i on i.id_inventory = 0 and i.id_card = c.id
            where ` + strings.Join(q.Restrictions, " and ")

	if q.Order != "" {
		query += " order by " + q.Order
	} else {
		query += ` order by e.name, NULLIF(regexp_replace(multiverse_number, '\D', '', 'g'), '')::int, c.name`
	}

	q.SQL = query
	l.Debug("data.CardQuery.Built",
		l.String("Query", q.SQL),
		l.Struct("Params", q.Values),
	)
	return nil
}

func (q *CardQuery) Fetch(i raizel.Iterable) error {
	var resultCards []Card
	for i.Next() {
		var card Card
		if fetchErr := card.FetchFull(i); fetchErr != nil {
			return fetchErr
		}
		resultCards = append(resultCards, card)
	}
	q.Result = resultCards
	return nil
}

type TokenQuery struct {
	Query
	//Result Fields
	Result []Token
	//Filter Fields
	Hydrate      string
	RegexName    string
	RegexType    string
	NotRegexType string
	IDExpansion  string
}

func (q *TokenQuery) Build() error {
	idxParam := 0
	if q.RegexName != "" {
		idxParam++
		q.Restrictions = append(q.Restrictions, fmt.Sprintf("t.name ~* $%d", idxParam))
		q.Values = append(q.Values, q.RegexName)
	}
	if q.RegexType != "" {
		idxParam++
		q.Restrictions = append(q.Restrictions, fmt.Sprintf("t.type ~* $%d", idxParam))
		q.Values = append(q.Values, q.RegexType)
	}
	if q.NotRegexType != "" {
		idxParam++
		q.Restrictions = append(q.Restrictions, fmt.Sprintf("not t.type ~* $%d", idxParam))
		q.Values = append(q.Values, q.NotRegexType)
	}
	if q.IDExpansion != "" {
		if idExpansion, convertErr := strconv.Atoi(q.IDExpansion); convertErr == nil {
			idxParam++
			q.Restrictions = append(q.Restrictions, fmt.Sprintf("t.id_expansion = $%d", idxParam))
			q.Values = append(q.Values, idExpansion)
		} else {
			l.Warn("TokenQuery.Build.ExpansionParamErr", l.String("Parameter", q.IDExpansion), l.Err(convertErr))
		}
	}
	query :=
		`
		select t.id, t.name, t.label, coalesce(t.text, ''), coalesce(t.color, ''), 
			coalesce(t.combat_power, ''), coalesce(t.power, ''), coalesce(t.toughness, ''),
			t.type, t.artist, t.id_asset,
            e.id, e.name, e.label, a.id_asset
        from token t
            left join expansion e on t.id_expansion = e.id
            left join expansion_asset a on a.id_expansion = e.id and a.id_rarity = 0
		`

	if len(q.Restrictions) > 0 {
		query += "where " + strings.Join(q.Restrictions, " and ") + "\n"
	}

	if q.Order != "" {
		query += " order by " + q.Order
	} else {
		query += " order by e.name, t.type, t.name"
	}

	q.SQL = query
	l.Debug("data.TokenQuery.Built",
		l.String("Query", q.SQL),
		l.Struct("Params", q.Values),
	)
	return nil
}

func (q *TokenQuery) Fetch(i raizel.Iterable) error {
	var resultTokens []Token
	for i.Next() {
		var token Token
		if fetchErr := token.FetchFull(i); fetchErr != nil {
			return fetchErr
		}
		resultTokens = append(resultTokens, token)
	}
	q.Result = resultTokens
	return nil
}

type ExpansionQuery struct {
	Query
	//Result Fields
	Result []Expansion
	//Filter Fields
	Hydrate   string
	RegexName string
}

func (q *ExpansionQuery) Build() error {
	idxParam := 0
	if q.RegexName != "" {
		idxParam++
		q.Restrictions = append(q.Restrictions, fmt.Sprintf("e.name ~* $%d", idxParam))
		q.Values = append(q.Values, q.RegexName)
	}
	var selectFields string
	switch q.Hydrate {
	case "small":
		selectFields = "e.id, e.name, a.id_asset"
	default:
		selectFields = "e.id, e.name, e.label, a.id_asset"
	}
	query :=
		`
        select ` + selectFields + `
        from expansion e
        left join expansion_asset a on e.id = a.id_expansion and (a.id_rarity = 0 or a.id_rarity = 4)
        `
	if len(q.Restrictions) > 0 {
		query += "where " + strings.Join(q.Restrictions, " and ") + "\n"
	}
	if q.Order != "" {
		query += " order by " + q.Order
	} else {
		query += " order by e.name"
	}
	q.SQL = query
	l.Debug("data.ExpansionQuery.Built",
		l.String("Query", q.SQL),
		l.Struct("Params", q.Values),
	)
	return nil
}

func (q *ExpansionQuery) Fetch(i raizel.Iterable) error {
	var resultExpansions []Expansion
	for i.Next() {
		var expansion Expansion
		var fetchFunc func(raizel.Fetchable) error
		switch q.Hydrate {
		case "small":
			fetchFunc = expansion.FetchSmall
		default:
			fetchFunc = expansion.FetchFull
		}
		if fetchErr := fetchFunc(i); fetchErr != nil {
			return fetchErr
		}
		resultExpansions = append(resultExpansions, expansion)
	}
	q.Result = resultExpansions
	return nil
}

type DeckQuery struct {
	Query
	//Result Fields
	Result []Deck
	//Filter Fields
	Hydrate   string
	RegexName string
}

func (q *DeckQuery) Build() error {
	idxParam := 0
	if q.RegexName != "" {
		idxParam++
		q.Restrictions = append(q.Restrictions, fmt.Sprintf("d.name ~* $%d", idxParam))
		q.Values = append(q.Values, q.RegexName)
	}
	//TODO: Think better and create a mechanism that full hydration will fetch decks with cards
	// var selectFields string
	// switch q.Hydrate {
	// case "small":
	// 	selectFields = "e.id, e.name, a.id_asset"
	// default:
	// 	selectFields = "e.id, e.name, e.label, a.id_asset"
	// }
	query := `
    select d.id, d.name, d.id_player from deck d
    `
	if len(q.Restrictions) > 0 {
		query += "where " + strings.Join(q.Restrictions, " and ") + "\n"
	}
	if q.Order != "" {
		query += " order by " + q.Order
	} else {
		query += " order by d.name"
	}

	q.SQL = query
	l.Debug("data.DeckQuery.Built",
		l.String("Query", q.SQL),
		l.Struct("Params", q.Values),
	)
	return nil
}

func (q *DeckQuery) Fetch(i raizel.Iterable) error {
	var resultDecks []Deck
	for i.Next() {
		var deck Deck
		//TODO: Think better to fetch by hydrate parameter
		if fetchErr := deck.FetchSmall(i); fetchErr != nil {
			return fetchErr
		}
		resultDecks = append(resultDecks, deck)
	}
	q.Result = resultDecks
	return nil
}
