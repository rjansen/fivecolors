package data

import (
	"database/sql"
	"errors"
	"regexp"
	"strings"

	"github.com/rjansen/l"
	"github.com/rjansen/raizel"
	// "github.com/rjansen/avalon/identity"
)

const (
	//MainBoard identifies the deck main board
	MainBoard = 1 << iota
	//SideBoard identifies the deck side board
	SideBoard    = 1 << iota
	HydrateSmall = "small"
	HydrateFull  = "full"
)

var (
	selectLimit             = 100
	primaryKeyViolation     = regexp.MustCompile(`Duplicate.*PRIMARY`)
	primaryKeyViolationByID = regexp.MustCompile(`duplicate key value`)
	NotFoundErr             = sql.ErrNoRows
)

type Card struct {
	ID               int           `json:"id"`
	MultiverseID     string        `json:"multiverseid"`
	Index            string        `json:"index"`
	Name             string        `json:"name"`
	Label            string        `json:"label"`
	Rate             float32       `json:"rate"`
	RateVotes        int           `json:"rateVotes"`
	Text             string        `json:"text"`
	ManacostLabel    string        `json:"manacostLabel"`
	CombatpowerLabel string        `json:"combatpowerLabel"`
	TypeLabel        string        `json:"typeLabel"`
	IDRarity         int           `json:"idRarity"`
	Flavor           string        `json:"flavor"`
	Artist           string        `json:"artist"`
	Expansion        Expansion     `json:"expansion"`
	InventoryCard    InventoryCard `json:"inventoryCard"`
	DeckCard         DeckCard      `json:"deckCard"`
	IDAsset          int           `json:"idAsset"`
}

type InventoryCard struct {
	IDInventory int `json:"idInvetory"`
	Quantity    int `json:"quantity"`
}

type DeckCard struct {
	IDDeck   int `json:"idDeck"`
	IDBoard  int `json:"idBoard"`
	Quantity int `json:"quantity"`
}

//Fetch fetchs the Row and sets the values into Card instance
// func (c *Card) Fetch(fetchable Fetchable) error {
// 	return fetchable.Scan(&c.ID, &c.Index, &c.Name, &c.Label, &c.Text,
// 		&c.ManacostLabel, &c.CombatpowerLabel, &c.TypeLabel,
// 		&c.IDRarity, &c.Flavor, &c.Artist,
// 		&c.Rate, &c.RateVotes, &c.IDAsset,
// 		&c.InventoryCard.IDInventory, &c.InventoryCard.Quantity,
// 		&c.Expansion.ID, &c.Expansion.Name, &c.Expansion.Label, &c.Expansion.IDAsset)
// }

//FetchWithDeckCard fetchs the Row and sets the values into Card instance with a DeckCard instance attached
// func (c *Card) FetchWithDeckCard(fetchable Fetchable) error {
// 	return fetchable.Scan(&c.ID, &c.Index, &c.Name, &c.Label, &c.Text,
// 		&c.ManacostLabel, &c.CombatpowerLabel, &c.TypeLabel,
// 		&c.IDRarity, &c.Flavor, &c.Artist,
// 		&c.Rate, &c.RateVotes, &c.IDAsset,
// 		&c.InventoryCard.IDInventory, &c.InventoryCard.Quantity,
// 		&c.DeckCard.IDDeck, &c.DeckCard.IDBoard, &c.DeckCard.Quantity,
// 		&c.Expansion.ID, &c.Expansion.Name, &c.Expansion.Label, &c.Expansion.IDAsset)
// }

func (c *Card) FetchFullWithDeckCard(fetchable raizel.Fetchable) error {
	return fetchable.Scan(&c.ID, &c.Index, &c.Name, &c.Label, &c.Text,
		&c.ManacostLabel, &c.CombatpowerLabel, &c.TypeLabel,
		&c.IDRarity, &c.Flavor, &c.Artist,
		&c.Rate, &c.RateVotes, &c.IDAsset,
		&c.DeckCard.IDDeck, &c.DeckCard.IDBoard, &c.DeckCard.Quantity,
		&c.Expansion.ID, &c.Expansion.Name, &c.Expansion.IDAsset)
}

func (c *Card) FetchFull(fetchable raizel.Fetchable) error {
	return fetchable.Scan(&c.ID, &c.MultiverseID, &c.Index, &c.Name, &c.Label, &c.Text,
		&c.ManacostLabel, &c.CombatpowerLabel, &c.TypeLabel,
		&c.IDRarity, &c.Flavor, &c.Artist, &c.Rate, &c.RateVotes, &c.IDAsset,
		&c.Expansion.ID, &c.Expansion.Name, &c.Expansion.Label, &c.Expansion.IDAsset,
		&c.InventoryCard.IDInventory, &c.InventoryCard.Quantity)
}

func (c *Card) ReadByID(client raizel.Client) error {
	if c.ID <= 0 {
		return errors.New("data.Card.ReadErr: Message='Card.ID is empty'")
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
        where c.id = $1`
	return client.QueryOne(query, c.FetchFull, c.ID)
}

func (c *Card) ReadByName(client raizel.Client) error {
	if strings.TrimSpace(c.Name) == "" {
		return errors.New("data.Card.ReadErr: Message='Card.ID is empty'")
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
        where c.name = $1`
	return client.QueryOne(query, c.FetchFull, c.Name)
}

func (c Card) Query(client raizel.Client, args ...interface{}) error {
	builder := args[0].(*CardQuery)
	if err := builder.Build(); err != nil {
		return err
	}
	queryErr := client.Query(builder.SQL, builder.Fetch, builder.Values...)
	if queryErr != nil {
		return queryErr
	}
	return nil
}

type Expansion struct {
	ID      int    `json:"id"`
	Name    string `json:"name"`
	Label   string `json:"label"`
	IDAsset int    `json:"idAsset"`
}

func (e *Expansion) FetchSmall(fetchable raizel.Fetchable) error {
	return fetchable.Scan(&e.ID, &e.Name, &e.IDAsset)
}

func (e *Expansion) FetchFull(fetchable raizel.Fetchable) error {
	return fetchable.Scan(&e.ID, &e.Name, &e.Label, &e.IDAsset)
}

func (e *Expansion) ReadByID(client raizel.Client) error {
	if e.ID <= 0 {
		return errors.New("data.Expansion.ReadErr: Message='Expansion.ID is empty'")
	}
	query :=
		`select e.id, e.name, e.label, a.id_asset
		from expansion e
            left join expansion_asset a on e.id = a.id_expansion and (a.id_rarity = 0 or a.id_rarity = 4)
		where e.id = $1`
	return client.QueryOne(query, e.FetchFull, e.ID)
}

func (e *Expansion) ReadByName(client raizel.Client) error {
	if strings.TrimSpace(e.Name) == "" {
		return errors.New("data.Expansion.ReadErr: Message='Expansion.Name is empty'")
	}
	query :=
		`select e.id, e.name, e.label, a.id_asset
		from expansion e
            left join expansion_asset a on e.id = a.id_expansion and (a.id_rarity = 0 or a.id_rarity = 4)
		where e.name = $1`
	return client.QueryOne(query, e.FetchFull, e.Name)
}

func (e Expansion) Query(client raizel.Client, args ...interface{}) error {
	builder := args[0].(*ExpansionQuery)
	if err := builder.Build(); err != nil {
		return err
	}
	queryErr := client.Query(builder.SQL, builder.Fetch, builder.Values...)
	if queryErr != nil {
		return queryErr
	}
	return nil
}

func GetPlayer(username string) (*Player, error) {
	player := &Player{Username: username}
	if readErr := raizel.Execute(player.ReadByUsername); readErr == NotFoundErr {
		//Creates a new player and associates the login with him
		if persistsErr := raizel.Execute(player.Persist); persistsErr != nil {
			return nil, persistsErr
		}
		if readErr := raizel.Execute(player.ReadByUsername); readErr != nil {
			return nil, readErr
		}
	} else if readErr != nil {
		return nil, readErr
	}
	return player, nil
}

type Player struct {
	ID          int    `json:"id"`
	Username    string `json:"username"`
	IDInventory int    `json:"idInventory"`
	IDDecks     []int  `json:"idDecks"`
}

//FillFromSession loads the player attributes from identity.Session object
// func (p *Player) FillFromSession(session *identity.Session) error {
// 	l.Infof("FillingPlayerFromSession: Session=%+v", session)
// 	if session.Data == nil ||
// 		session.Data["id"] == nil || session.Data["id"].(float64) <= 0 ||
// 		session.Data["username"] == nil || strings.TrimSpace(session.Data["username"].(string)) == "" ||
// 		session.Data["idInventory"] == nil || session.Data["idInventory"].(float64) <= 0 ||
// 		session.Data["idDecks"] == nil || len(session.Data["idDecks"].([]interface{})) < 0 {
// 		return errors.New("Some required attributes to fill a player is missing")
// 	}
// 	p.ID = int(session.Data["id"].(float64))
// 	p.Username = session.Data["username"].(string)
// 	p.IDInventory = int(session.Data["idInventory"].(float64))
// 	idDecks := session.Data["idDecks"].([]interface{})
// 	p.IDDecks = make([]int, len(idDecks))
// 	for i := range idDecks {
// 		p.IDDecks[i] = int(idDecks[i].(float64))
// 	}
// 	return nil
// }

func (p *Player) SessionData() map[string]interface{} {
	return map[string]interface{}{
		"id":          p.ID,
		"username":    p.Username,
		"idInventory": p.IDInventory,
		"idDecks":     p.IDDecks,
	}
}

func (p *Player) FetchFull(fetchable raizel.Fetchable) error {
	return fetchable.Scan(&p.ID, &p.Username, &p.IDInventory)
}

func (p *Player) ReadByUsername(client raizel.Client) error {
	if strings.TrimSpace(p.Username) == "" {
		return errors.New("data.Player.ReadError: Message='Player.Username is empty'")
	}
	query := `
        select p.id, p.username, (select max(i.id) from inventory i where i.id_player = p.id) as id_inventory
        from player p
        where p.username = ?;
    `
	err := client.QueryOne(query, p.FetchFull, p.Username)
	if err != nil {
		return err
	}
	//Read Fully
	return p.ReadDecks(client, -1)
}

func (p *Player) ReadDecks(client raizel.Client, page int) error {
	if p.ID <= 0 {
		return errors.New("data.Player.ReadDecksError: Message='Player.ID is empty'")
	}

	query := `
        select d.id
        from deck d
        where d.id_player = ? order by d.id
    `
	iterFunc := func(i raizel.Iterable) error {
		tempIDDecks := []int{}
		//var tempIDDecks []int
		for i.Next() {
			var idDeck int
			if err := i.Scan(&idDeck); err != nil {
				return err
			}
			tempIDDecks = append(tempIDDecks, idDeck)
		}
		p.IDDecks = tempIDDecks
		return nil
	}
	err := client.Query(query, iterFunc, p.ID)
	if err != nil {
		return err
	}
	return nil
}

func (p *Player) Persist(client raizel.Client) error {
	if strings.TrimSpace(p.Username) == "" {
		return errors.New("data.Player.PersistError: Message='Player.Username is empty'")
	}

	//Checks if the player already exists in the database
	countQuery := `select count(1) from player p where p.username = ?`
	var countPlayer int
	fetchFunc := func(f raizel.Fetchable) error {
		return f.Scan(&countPlayer)
	}
	err := client.QueryOne(countQuery, fetchFunc, p.Username)
	if err != nil {
		return err
	}

	if countPlayer <= 0 {
		insert := `insert into player (id, username) values (nextval('sq_player'), ?) returning id`

		var insertedID int
		idFetchFunc := func(f raizel.Fetchable) error {
			return f.Scan(&insertedID)
		}

		_, insertErr := client.Exec(insert, idFetchFunc, p.Username)
		if insertErr != nil {
			return insertErr
		}
		p.ID = int(insertedID)

		insertInventory := `insert into inventory (name, id_player) values (?, ?)`
		insertInvetoryResult, insertInvetoryErr := client.Exec(insertInventory, p.Username, p.ID)
		if insertInvetoryErr != nil {
			return insertInvetoryErr
		}
		inventoryID, lastIDInventoryErr := insertInvetoryResult.LastInsertId()
		if lastIDInventoryErr != nil {
			return lastIDInventoryErr
		}
		p.IDInventory = int(inventoryID)

		l.Infof("Player.PersistedNewPlayer: ID=%v Username='%v' IDInventory=%v", p.ID, p.Username, p.IDInventory)
	} else {
		update := `update player set dt_lastlogin = current_timestamp() where username = ?`

		_, updateErr := client.Exec(update, p.Username)
		if updateErr != nil {
			return updateErr
		}
		l.Infof("data.Player.UpdatedPlayer: ID=%v Username='%v' IDInventory=%v", p.ID, p.Username, p.IDInventory)
	}
	return nil
}

type Inventory struct {
	ID       int    `json:"id"`
	Name     string `json:"name"`
	Label    string `json:"label"`
	IDPlayer int    `json:"idPlayer"`
	Cards    []Card `json:"cards"`
}

func (i *Inventory) Persist(client raizel.Client) error {
	if i.IDPlayer > 0 {
		return errors.New("data.Inventory.UpdateError: Message='Inventory.IDPlayer is not zero'")
	}
	if i.ID > 0 {
		return errors.New("data.Inventory.UpdateError: Message='Inventory.ID is not zero'")
	}

	insertCardQuery :=
		`insert into inventory_card (id_inventory, id_card, quantity)
        values ($1, $2, $3)`
	updateCardQuery :=
		`update inventory_card
            set quantity = $1
        where id_inventory = $2 and id_card = $3`

	for _, card := range i.Cards {
		_, insertErr := client.Exec(insertCardQuery, i.ID, card.ID, card.InventoryCard.Quantity)
		if insertErr != nil {
			l.Error("Inventory.InsertCardEx", l.Err(insertErr))
			if !primaryKeyViolationByID.MatchString(insertErr.Error()) {
				return insertErr
			}
			updateResult, updateErr := client.Exec(updateCardQuery, card.InventoryCard.Quantity, i.ID, card.ID)
			if updateErr != nil {
				return updateErr
			}
			rowsUpdated, err := updateResult.RowsAffected()
			if err != nil {
				l.Errorf("Inventory.UpdateGetRowsAffectedEx: Message='%v'", err.Error())
			} else {
				if rowsUpdated != 1 {
					l.Errorf("Inventory.UpdateMultipleEx: Message='%d Card Records was update for Inventory.ID=%d'", rowsUpdated, i.ID)
				}
			}
		}
	}
	l.Infof("Inventory.Persisted: ID=%v IDPlayer=%v", i.ID, i.IDPlayer)
	return nil
}

//Delete deletes the INVENTORY record references to Inventory
// func (i *Inventory) Delete() error {
// 	i.Attach()
// 	defer closeDB(i)
// 	if i.ID <= 0 {
// 		return errors.New("data.Inventory.DeleteError: Message='Inventory.ID is empty'")
// 	}
// 	if i.IDPlayer <= 0 {
// 		return errors.New("data.Inventory.DeleteError: Message='Inventory.IDPlayer is empty'")
// 	}
// 	deleteCardsResult, deleteCardsErr := i.db.Exec("delete from inventory_card where id_inventory = ?", i.ID)
// 	if deleteCardsErr != nil {
// 		return deleteCardsErr
// 	}
// 	cardsRowsDeleted, cardsRowsDeletedErr := deleteCardsResult.RowsAffected()
// 	if cardsRowsDeletedErr != nil {
// 		l.Errorf("Inventory.DeleteCardsGetRowsAffectedEx: Message='%v'", cardsRowsDeletedErr.Error())
// 	} else {
// 		l.Debugf("Inventory.DeleteCards: DeletedCards=%v ID=%d IDPlayer=%v", cardsRowsDeleted, i.ID, i.IDPlayer)
// 	}
// 	deleteResult, deleteErr := i.db.Exec("delete from inventory where id = ? and id_player = ?", i.ID, i.IDPlayer)
// 	if deleteErr != nil {
// 		return deleteErr
// 	}
// 	rowsDeleted, err := deleteResult.RowsAffected()
// 	if err != nil {
// 		l.Errorf("data.Inventory.DeleteGetRowsAffectedEx: Message='%v'", err.Error())
// 	} else {
// 		if rowsDeleted != 1 {
// 			l.Errorf("data.Inventory.DeleteMultipleEx: Message='%d Records was delete for Inventory.ID=%d and Inventory.IDPlayer=%d'", rowsDeleted, i.ID, i.IDPlayer)
// 		}
// 	}
// 	l.Infof("Inventory.Deleted: ID=%d IDPlayer=%v", i.ID, i.IDPlayer)
// 	return nil
// }

// func (i *Inventory) FetchFull(fetchable raizel.Fetchable) error {
// 	return fetchable.Scan(&i.ID, &i.Name, &i.IDPlayer)
// }

//Read gets the entity representation from the database.
//Inventory.ID must not empty to perform a Read operation
// func (i *Inventory) Read() error {
// 	i.Attach()
// 	defer closeDB(i)
// 	if i.ID <= 0 {
// 		return errors.New("data.Inventory.ReadError: Message='Inventory.ID is empty'")
// 	}
// 	if i.IDPlayer <= 0 {
// 		return errors.New("data.Inventory.ReadError: Message='Inventory.IDPlayer is empty'")
// 	}
// 	row := i.db.QueryRow("select i.id, i.name, i.id_player from inventory i where i.id = ? and i.id_player = ?", i.ID, i.IDPlayer)
// 	return i.Fetch(row)
// }

//ReadCards gets one related Cards page of this Inventory
//Inventory.ID must not empty to perform a ReadCards operation
// func (i *Inventory) ReadCards(page int) error {
// 	i.Attach()
// 	defer closeDB(i)
// 	if i.ID <= 0 {
// 		return errors.New("data.Inventory.ReadCardsError: Message='Inventory.ID is empty'")
// 	}
// 	query :=
// 		`select c.id, c.multiverse_number, c.name, c.label, coalesce(c.text, ''),
//             coalesce(c.manacost_label, ''), coalesce(c.combatpower_label, ''), c.type_label,
//             c.id_rarity, coalesce(c.flavor, ''), c.artist,
//             c.rate, c.rate_votes, c.id_asset,
//             coalesce(i.id_inventory, 1), coalesce(i.quantity, 0),
//             e.id, e.name, e.label, a.id_asset
//         from card c
//             left join expansion e on c.id_expansion = e.id
//             left join expansion_asset a on a.id_expansion = e.id and a.id_rarity = c.id_rarity
//             left join inventory_card i on i.id_card = c.id
//         where i.id_inventory = ?
//         order by e.name, abs(c.multiverse_number)`
// 	rows, readCardsErr := i.db.Query(query, i.ID)
// 	if readCardsErr != nil {
// 		return readCardsErr
// 	}
// 	//cards := make([]Card, 0)
// 	cards := []Card{}
// 	//var cards []Card
// 	for rows.Next() {
// 		nextCard := Card{Expansion: Expansion{}, InventoryCard: InventoryCard{}}
// 		if err := nextCard.Fetch(rows); err != nil {
// 			return err
// 		}
// 		cards = append(cards, nextCard)
// 	}
// 	i.Cards = cards
// 	return nil
// }

type Deck struct {
	ID          int    `json:"id"`
	Name        string `json:"name"`
	IDPlayer    int    `json:"idPlayer,omitempty"`
	IDInventory int    `json:"idInventory,omitempty"`
	Cards       []Card `json:"cards"`
}

//Delete deletes the DECK record references to Deck
// func (d *Deck) Delete() error {
// 	d.Attach()
// 	defer closeDB(d)
// 	if d.ID <= 0 {
// 		return errors.New("data.Deck.DeleteError: Message='Deck.ID is empty'")
// 	}
// 	if d.IDPlayer <= 0 {
// 		return errors.New("data.Deck.DeleteError: Message='Deck.IDPlayer is empty'")
// 	}
// 	deleteCardsResult, deleteCardsErr := d.db.Exec("delete from deck_card where id_deck = ?", d.ID)
// 	if deleteCardsErr != nil {
// 		return deleteCardsErr
// 	}
// 	cardsRowsDeleted, err := deleteCardsResult.RowsAffected()
// 	if err != nil {
// 		l.Errorf("Deck.DeleteCardsGetRowsAffectedEx: Message='%v'", err.Error())
// 	}
// 	l.Debugf("Deck.DeletedDeckCards: RowsDeleted=%d Deck.ID=%d and Deck.IDPlayer=%d", cardsRowsDeleted, d.ID, d.IDPlayer)
// 	deleteResult, deleteErr := d.db.Exec("delete from deck where id = ? and id_player = ?", d.ID, d.IDPlayer)
// 	if deleteErr != nil {
// 		return deleteErr
// 	}
// 	rowsDeleted, err := deleteResult.RowsAffected()
// 	if err != nil {
// 		l.Errorf("Deck.DeleteGetRowsAffectedEx: Message='%v'", err.Error())
// 	} else {
// 		if rowsDeleted != 1 {
// 			l.Errorf("Deck.DeleteMultipleEx: Message='%d Records was delete for Deck.ID=%d and Deck.IDPlayer=%d'", rowsDeleted, d.ID, d.IDPlayer)
// 		}
// 	}
// 	l.Infof("Deck.Deleted: Deck.ID=%d and Deck.IDPlayer=%d", d.ID, d.IDPlayer)
// 	return nil
// }

func (d *Deck) FetchSmall(fetchable raizel.Fetchable) error {
	return fetchable.Scan(&d.ID, &d.Name, &d.IDPlayer)
}

func (d *Deck) Persist(client raizel.Client) error {
	if d.Name == "" {
		return errors.New("data.Deck.PersistError: Message='Deck.Name is empty'")
	}

	if d.ID == 0 {
		fetchID := func(f raizel.Fetchable) error {
			return f.Scan(&d.ID)
		}
		createErr := client.QueryOne("insert into deck (id, name, id_player) values (nextval('sq_deck'), $1, $2) returning id", fetchID, d.Name, d.IDPlayer)
		if createErr != nil {
			return createErr
		}
		l.Debug("data.Deck.InsertNewDeck",
			l.Int("ID", d.ID),
			l.Int("IDPlayer", d.IDPlayer),
			l.String("Name", d.Name),
		)
	} else {
		updateResult, updateErr := client.Exec("update deck set name = $1, id_player = $2 where id = $3", d.Name, d.IDPlayer, d.ID)
		if updateErr != nil {
			return updateErr
		}
		rowsUpdated, err := updateResult.RowsAffected()
		if err != nil {
			l.Warn("data.Deck.UpdateGetRowsAffectedErr", l.Err(err))
		} else {
			if rowsUpdated != 1 {
				l.Warn("data.Deck.UpdateMultipleErr",
					l.Int64("RowsUpdated", rowsUpdated),
					l.Int("ID", d.ID),
				)
			}
		}
		l.Debug("data.Deck.UpdateOldDeck",
			l.Int("ID", d.ID),
			l.Int("IDPlayer", d.IDPlayer),
			l.String("Name", d.Name),
		)
	}

	if _, deleteErr := client.Exec("delete from deck_card where id_deck = $1", d.ID); deleteErr != nil {
		return deleteErr
	}

	insertCardQuery :=
		`insert into deck_card (id_deck, id_card, id_board, quantity)
        values ($1, $2, $3, $4)`
	updateCardQuery :=
		`update deck_card
            set quantity = $1
        where id_deck = $2 and id_card = $3 and id_board = $4`

	for _, card := range d.Cards {
		_, insertErr := client.Exec(insertCardQuery, d.ID, card.ID, card.DeckCard.IDBoard, card.DeckCard.Quantity)
		if insertErr != nil {
			l.Error("data.Deck.InsertCardErr", l.Err(insertErr))
			if !primaryKeyViolationByID.MatchString(insertErr.Error()) {
				return insertErr
			}
			updateResult, updateErr := client.Exec(updateCardQuery, card.DeckCard.Quantity, d.ID, card.ID, card.DeckCard.IDBoard)
			if updateErr != nil {
				return updateErr
			}
			rowsUpdated, err := updateResult.RowsAffected()
			if err != nil {
				l.Warn("data.Deck.UpdateGetRowsAffectedWarn", l.Err(err))
			} else {
				if rowsUpdated != 1 {
					l.Warn("data.Deck.UpdateMultipleErr",
						l.Int64("RowsUpdated", rowsUpdated),
						l.Int("ID", d.ID),
					)
				}
			}
		}
	}
	l.Info("data.Deck.Persisted",
		l.Int("ID", d.ID),
		l.Int("IDPlayer", d.IDPlayer),
		l.String("Name", d.Name),
	)
	return nil
}

func (d *Deck) ReadByID(client raizel.Client) error {
	if d.ID <= 0 {
		return errors.New("data.Asset.ReadByIDError: Message='Deck.ID is empty'")
	}
	err := client.QueryOne("select d.id, d.name, d.id_player from deck d where d.id = $1", d.FetchSmall, d.ID)
	if err != nil {
		return err
	}
	//Read Fully
	return d.ReadCards(client, -1)
}

func (d *Deck) ReadByName(client raizel.Client) error {
	if strings.TrimSpace(d.Name) == "" {
		return errors.New("data.Deck.ReadByNameErr: Message='Deck.Name is empty'")
	}
	err := client.QueryOne("select d.id, d.name, d.id_player from deck d where d.name = $1", d.FetchSmall, d.Name)
	if err != nil {
		return err
	}
	//Read Fully
	return d.ReadCards(client, -1)
}

func (d *Deck) ReadCards(client raizel.Client, page int) error {
	if d.ID <= 0 {
		return errors.New("data.Deck.ReadCardsRaizelErr: Message='Deck.ID is empty'")
	}
	query :=
		`select c.id, c.multiverse_number, c.name, c.label, coalesce(c.text, ''),
            coalesce(c.manacost_label, ''), coalesce(c.combatpower_label, ''), c.type_label,
            c.id_rarity, coalesce(c.flavor, ''), c.artist,
            c.rate, c.rate_votes, c.id_asset,
            d.id_deck, d.id_board, coalesce(d.quantity, 0) as deck_quantity,
            e.id, e.name, a.id_asset
        from card c
            left join expansion e on c.id_expansion = e.id
            left join expansion_asset a on a.id_expansion = e.id and a.id_rarity = c.id_rarity
            left join deck_card d on d.id_card = c.id
        where d.id_deck = $1 
        order by d.id_board, c.type_label, e.name`

	cardsFetchFunc := func(i raizel.Iterable) error {
		//tempCards := make([]Card, selectLimit)
		var tempCards []Card
		for i.Next() {
			nextCard := Card{}
			if err := nextCard.FetchFullWithDeckCard(i); err != nil {
				return err
			}
			tempCards = append(tempCards, nextCard)
		}
		// if tempCards == nil {
		// 	tempCards = make([]Card, 0)
		// }
		d.Cards = tempCards
		return nil
	}
	readCardsErr := client.Query(query, cardsFetchFunc, d.ID)
	if readCardsErr != nil {
		return readCardsErr
	}
	return nil
}

func (d Deck) Query(client raizel.Client, args ...interface{}) error {
	builder := args[0].(*DeckQuery)
	if err := builder.Build(); err != nil {
		return err
	}
	queryErr := client.Query(builder.SQL, builder.Fetch, builder.Values...)
	if queryErr != nil {
		return queryErr
	}
	return nil
}
