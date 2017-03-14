package mtgo

import (
	"encoding/csv"
	"fmt"
	"github.com/rjansen/fivecolors/data"
	"github.com/rjansen/l"
	"github.com/rjansen/raizel"
	"io"
	"strconv"
	"strings"
)

//ImportCollection loads a MTGO collection definition into Fivecolors
//MTGO CSV: Card Name,Quantity,ID #,Rarity,Set,Collector #,Premium,
func ImportCollection(collectionReader io.Reader) error {
	_, err := ReadFile(collectionReader)
	return err
}

//ImportDeck loads a MTGO deck definition into Fivecolors
//Card Name,Quantity,ID #,Rarity,Set,Collector #,Premium,Sideboarded,
func ImportDeck(deckReader io.Reader) error {
	_, err := ReadFile(deckReader)
	return err
}

func ReadFile(fileReader io.Reader) ([]data.Card, error) {
	r := csv.NewReader(fileReader)
	r.FieldsPerRecord = 7
	// r.LazyQuotes = true
	//r.Comma = ';'
	// r.Comment = "#"

	// Buffered Read
	var (
		cards         []data.Card
		collectionQtd = 0
	)
	for {
		record, err := r.Read()
		if err == io.EOF {
			break
		}
		fmt.Printf("mtgo.MTGOCollection Record=%+v\n", record)
		if strings.TrimSpace(record[0]) == "" || strings.TrimSpace(record[0]) == "Card Name" {
			l.Info("mtgo.ImportDeck.IgnoringLines", l.String("Name", record[0]))
			continue
		}
		if err != nil {
			l.Error("mtgo.ImportDeck.ReadDeckErr", l.Err(err))
			return nil, err
		}
		var card data.Card
		card.Name = record[0]
		l.Info("mtgo.ImportDeck.FindingCard", l.String("Name", card.Name))
		if err := raizel.Execute(card.ReadByName); err != nil {
			l.Error("mtgo.ImportDeck.FindCardByNameErr", l.Err(err))
			return nil, err
		}
		cardQtd, _ := strconv.Atoi(record[1])
		card.DeckCard.Quantity = cardQtd
		collectionQtd += cardQtd
		fmt.Printf("mtgo.Found5ColorsCard Total=%d MTGO=%+v 5Colors=%+v\n", cardQtd, record, card)
		cards = append(cards, card)
	}
	fmt.Printf("mtgo.MTGOCollection Total=%d\n", collectionQtd)
	return cards, nil
}
