package data

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"io"
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
	SideBoard = 1 << iota
)

var (
	selectLimit         = 100
	primaryKeyViolation = regexp.MustCompile(`Duplicate.*PRIMARY`)
	NotFoundErr         = sql.ErrNoRows
)

//Fetchable supply the sql.Scan interface for a struct
type Fetchable interface {
	Scan(dest ...interface{}) error
}

//Readable provides read actions for a struct
type Readable interface {
	Fetch(fetchable Fetchable) error
	Read() error
}

//Queryable provides query actions for a struct
type Queryable interface {
	Readable
	Query(queryParameters map[string]interface{}, order string) ([]interface{}, error)
}

//Writable provides persistence actions for a struct
type Writable interface {
	Persist() error
	Delete() error
}

//Attachable creates a interface for structs do database actions
type Attachable interface {
	SetDB(db *sql.DB) error
	GetDB() (*sql.DB, error)
	Attach() error
}

//JSONSerializable provides to a struct json external representation
type JSONSerializable interface {
	Marshal() ([]byte, error)
	Unmarshal(reader io.Reader) error
}

//AttachToDB binds a new database connection to Attachable reference
func attachToDB(a Attachable) error {
	if _, err := a.GetDB(); err != nil {
		tempDB, err := pool.GetConnection()
		if err != nil {
			l.Error("data.GetConnectionErr", l.Err(err))
			return fmt.Errorf("data.Card.AttachError: Messages='%v'", err.Error())
		}
		return a.SetDB(tempDB)
	}
	return nil
}

func closeDB(a Attachable) error {
	return a.SetDB(nil)
}

//Card represents the CARD entity
type Card struct {
	ID               int           `json:"id"`
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
	//db is a transient pointer to database connection
	db *sql.DB
}

//InventoryCard represents the INVENTORY_CARD virtual entity
type InventoryCard struct {
	IDInventory int `json:"idInvetory"`
	Quantity    int `json:"quantity"`
}

//DeckCard represents the DECK_CARD virtual entity
type DeckCard struct {
	IDDeck   int `json:"idDeck"`
	IDBoard  int `json:"idBoard"`
	Quantity int `json:"quantity"`
}

//SetDB attachs a database connection to Card
func (c *Card) SetDB(db *sql.DB) error {
	if db == nil {
		return errors.New("NullDBReferenceError: Message='The db parameter is required'")
	}
	c.db = db
	return nil
}

//GetDB returns the Card attached connection
func (c *Card) GetDB() (*sql.DB, error) {
	if c.db == nil {
		return nil, errors.New("NotAttachedError: Message='The db context is null'")
	}
	return c.db, nil
}

//Attach binds a new database connection to Card reference
func (c *Card) Attach() error {
	return attachToDB(c)
}

//Fetch fetchs the Row and sets the values into Card instance
func (c *Card) Fetch(fetchable Fetchable) error {
	return fetchable.Scan(&c.ID, &c.Index, &c.Name, &c.Label, &c.Text,
		&c.ManacostLabel, &c.CombatpowerLabel, &c.TypeLabel,
		&c.IDRarity, &c.Flavor, &c.Artist,
		&c.Rate, &c.RateVotes, &c.IDAsset,
		&c.InventoryCard.IDInventory, &c.InventoryCard.Quantity,
		&c.Expansion.ID, &c.Expansion.Name, &c.Expansion.Label, &c.Expansion.IDAsset)
}

//FetchWithDeckCard fetchs the Row and sets the values into Card instance with a DeckCard instance attached
func (c *Card) FetchWithDeckCard(fetchable Fetchable) error {
	return fetchable.Scan(&c.ID, &c.Index, &c.Name, &c.Label, &c.Text,
		&c.ManacostLabel, &c.CombatpowerLabel, &c.TypeLabel,
		&c.IDRarity, &c.Flavor, &c.Artist,
		&c.Rate, &c.RateVotes, &c.IDAsset,
		&c.InventoryCard.IDInventory, &c.InventoryCard.Quantity,
		&c.DeckCard.IDDeck, &c.DeckCard.IDBoard, &c.DeckCard.Quantity,
		&c.Expansion.ID, &c.Expansion.Name, &c.Expansion.Label, &c.Expansion.IDAsset)
}

func (c *Card) FetchByIDWithDeckCard(fetchable raizel.Fetchable) error {
	return fetchable.Scan(&c.ID, &c.Index, &c.Name, &c.Label, &c.Text,
		&c.ManacostLabel, &c.CombatpowerLabel, &c.TypeLabel,
		&c.IDRarity, &c.Flavor, &c.Artist,
		&c.Rate, &c.RateVotes, &c.IDAsset,
		&c.DeckCard.IDDeck, &c.DeckCard.IDBoard, &c.DeckCard.Quantity,
		&c.Expansion.ID, &c.Expansion.Name, &c.Expansion.IDAsset)
}

//Read gets the entity representation from the database.
//Card.ID must not empty to perform a Read operation
func (c *Card) Read() error {
	if err := c.Attach(); err != nil {
		return err
	}
	defer closeDB(c)
	if c.ID <= 0 {
		return errors.New("data.Card.ReadError: Message='Card.ID is empty'")
	}
	if c.InventoryCard.IDInventory <= 0 {
		return errors.New("data.Card.ReadError: Message='Card.InventoryCard.IDInventory is empty'")
	}
	query :=
		`
		select c.id, c.multiverse_number, c.name, c.label, coalesce(c.text, ''),
            coalesce(c.manacost_label, ''), coalesce(c.combatpower_label, ''), c.type_label,
            c.id_rarity, coalesce(c.flavor, ''), c.artist,
            c.rate, c.rate_votes, c.id_asset,
            coalesce(i.id_inventory, ?), coalesce(i.quantity, 0),
            e.id, e.name, e.label, a.id_asset
        from card c
            left join expansion e on c.id_expansion = e.id
            left join expansion_asset a on a.id_expansion = e.id and a.id_rarity = c.id_rarity
            left join inventory_card i on i.id_inventory = ? and i.id_card = c.id
        where c.id = ?`
	row := c.db.QueryRow(query, c.InventoryCard.IDInventory, c.ID, c.InventoryCard.IDInventory)
	return c.Fetch(row)
}

//Query querys CARDs by restrictions and create a list of Cards references
func (c *Card) Query(queryParameters map[string]interface{}, order string) ([]interface{}, error) {
	if c.InventoryCard.IDInventory <= 0 {
		return nil, errors.New("data.Card.QueryError: Message='Card.InventoryCard.IDInventory is empty'")
	}
	if queryParameters == nil {
		return nil, errors.New("data.Card.QueryError: Message='QueryParameter is empty'")
	}
	parameterSize := len(queryParameters)
	if parameterSize <= 0 {
		return nil, errors.New("data.Card.QueryError: Message='QueryParameter is empty'")
	}
	restrictions := make([]string, parameterSize)
	values := make([]interface{}, parameterSize)
	paramIndex := 0
	for k, v := range queryParameters {
		restrictions[paramIndex] = k
		values[paramIndex] = v
		paramIndex++
	}
	query :=
		`
            select c.id, c.multiverse_number, c.name, c.label, coalesce(c.text, ''),
                coalesce(c.manacost_label, ''), coalesce(c.combatpower_label, ''), c.type_label,
                c.id_rarity, coalesce(c.flavor, ''), c.artist,
                c.rate, c.rate_votes, c.id_asset,
                coalesce(i.id_inventory, ?), coalesce(i.quantity, 0),
                e.id, e.name, e.label, a.id_asset
            from card c
                left join expansion e on c.id_expansion = e.id
                left join expansion_asset a on a.id_expansion = e.id and a.id_rarity = c.id_rarity
                left join inventory_card i on i.id_inventory = ? and i.id_card = c.id
            where ` + strings.Join(restrictions, " and ")

	if order != "" {
		query += " order by " + order
	} else {
		query += " order by e.name, abs(c.multiverse_number)"
	}

	//Prepend the IDInventory parameter into parameters value
	values = append([]interface{}{c.InventoryCard.IDInventory, c.InventoryCard.IDInventory}, values...)

	l.Debugf("Card.Query: Query=%v Parameters=%v", query, values)
	if err := c.Attach(); err != nil {
		return nil, err
	}
	defer closeDB(c)
	rows, queryErr := c.db.Query(query, values...)
	if queryErr != nil {
		return nil, queryErr
	}
	//Limit the result
	//cards := make([]interface{}, 0)
	cards := []interface{}{}
	//var cards []interface{}
	for rows.Next() {
		nextCard := Card{Expansion: Expansion{}, InventoryCard: InventoryCard{}}
		if err := nextCard.Fetch(rows); err != nil {
			return nil, err
		}
		cards = append(cards, nextCard)
	}
	return cards, nil
}

//Marshal writes a json representation of Card
func (c *Card) Marshal() ([]byte, error) {
	return json.Marshal(c)
}

//Unmarshal reads a json representation of Card
func (c *Card) Unmarshal(reader io.Reader) error {
	return json.NewDecoder(reader).Decode(&c)
}

//Expansion represents the EXPANSION entity
type Expansion struct {
	ID      int    `json:"id"`
	Name    string `json:"name"`
	Label   string `json:"label"`
	IDAsset int    `json:"idAsset"`
	//db is a transient pointer to database connection
	db *sql.DB
}

//SetDB attachs a database connection to Expansion
func (e *Expansion) SetDB(db *sql.DB) error {
	if db == nil {
		return errors.New("NullDBReferenceError: Message='The db parameter is required'")
	}
	e.db = db
	return nil
}

//GetDB returns the Expansion attached connection
func (e *Expansion) GetDB() (*sql.DB, error) {
	if e.db == nil {
		return nil, errors.New("NotAttachedError: Message='The db context is null'")
	}
	return e.db, nil
}

//Attach binds a new database connection to Expansion reference
func (e *Expansion) Attach() error {
	return attachToDB(e)
}

//Fetch fetchs the Row and sets the values into Expansion instance
func (e *Expansion) Fetch(fetchable Fetchable) error {
	return fetchable.Scan(&e.ID, &e.Name, &e.Label, &e.IDAsset)
}

//Read gets the entity representation from the database.
//Expansion.ID must not empty to perform a Read operation
func (e *Expansion) Read() error {
	e.Attach()
	defer closeDB(e)
	if e.ID <= 0 {
		return errors.New("data.Expansion.ReadError: Message='Expansion.ID is empty'")
	}
	query :=
		`select e.id, e.name, e.label, a.id_asset
		from expansion e
            left join expansion_asset a on e.id = a.id_expansion and (a.id_rarity = 0 or a.id_rarity = 4)
		where e.id = ?`
	row := e.db.QueryRow(query, e.ID)
	return e.Fetch(row)
}

//Query querys CARDs by restrictions and create a list of Cards references
func (e *Expansion) Query(queryParameters map[string]interface{}, order string) ([]interface{}, error) {
	var parameterSize int
	if queryParameters == nil {
		parameterSize = 0
	} else {
		parameterSize = len(queryParameters)
	}
	restrictions := make([]string, parameterSize)
	values := make([]interface{}, parameterSize)
	if parameterSize > 0 {
		paramIndex := 0
		for k, v := range queryParameters {
			restrictions[paramIndex] = k
			values[paramIndex] = v
			paramIndex++
		}
	}
	query :=
		`
        select e.id, e.name, e.label, a.id_asset
		from expansion e
        left join expansion_asset a on e.id = a.id_expansion and (a.id_rarity = 0 or a.id_rarity = 4)
        `
	if len(restrictions) > 0 {
		query += "where " + strings.Join(restrictions, " and ") + "\n"
	}

	if order != "" {
		query += " order by " + order
	} else {
		query += " order by e.name"
	}

	l.Debugf("Expansion.Query: Query=%v Parameters=%v", query, values)
	if err := e.Attach(); err != nil {
		return nil, err
	}
	defer closeDB(e)
	rows, queryErr := e.db.Query(query, values...)
	if queryErr != nil {
		return nil, queryErr
	}
	//Limit the result
	//expansions := make([]interface{}, 0)
	expansions := []interface{}{}
	//var expansions []interface{}
	for rows.Next() {
		nextExpansion := Expansion{}
		if err := nextExpansion.Fetch(rows); err != nil {
			return nil, err
		}
		expansions = append(expansions, nextExpansion)
	}
	return expansions, nil
}

//Marshal writes a json representation of Expansion
func (e *Expansion) Marshal() ([]byte, error) {
	return json.Marshal(e)
}

//Unmarshal reads a json representation of Expansion
func (e *Expansion) Unmarshal(reader io.Reader) error {
	return json.NewDecoder(reader).Decode(&e)
}

//Asset represents the ASSET entity
type Asset struct {
	ID         int    `json:"id"`
	Label      string `json:"label"`
	BinaryData []byte `json:"binaryData"`
	//db is a transient pointer to database connection
	db *sql.DB
}

//SetDB attachs a database connection to Asset
func (a *Asset) SetDB(db *sql.DB) error {
	if db == nil {
		return errors.New("NullDBReferenceError: Message='The db parameter is required'")
	}
	a.db = db
	return nil
}

//GetDB returns the Asset attached connection
func (a *Asset) GetDB() (*sql.DB, error) {
	if a.db == nil {
		return nil, errors.New("NotAttachedError: Message='The db context is null'")
	}
	return a.db, nil
}

//Attach binds a new database connection to Asset reference
func (a *Asset) Attach() error {
	return attachToDB(a)
}

//Fetch fetchs the Row and sets the values into Asset instance
func (a *Asset) Fetch(fetchable Fetchable) error {
	return fetchable.Scan(&a.ID, &a.Label, &a.BinaryData)
}

//Read gets the entity representation from the database.
//Asset.ID must not empty to perform a Read operation
func (a *Asset) Read() error {
	a.Attach()
	defer closeDB(a)
	if a.ID <= 0 {
		return errors.New("data.Asset.ReadError: Message='Asset.ID is empty'")
	}
	row := a.db.QueryRow("select a.id, a.label, a.binarydata from asset a where a.id = ?", a.ID)
	return a.Fetch(row)
}

//Marshal writes a json representation of Asset
func (a *Asset) Marshal() ([]byte, error) {
	return json.Marshal(a)
}

//Unmarshal reads a json representation of Asset
func (a *Asset) Unmarshal(reader io.Reader) error {
	return json.NewDecoder(reader).Decode(&a)
}

//GetPlayer loads the player from database or creates a new one and associate the login eith him
func GetPlayer(username string) (*Player, error) {
	player := &Player{Username: username}
	if readErr := player.Read(); readErr == NotFoundErr {
		//Creates a new player and associates the login with him
		if persistsErr := player.Persist(); persistsErr != nil {
			return nil, persistsErr
		}
		if readErr := player.Read(); readErr != nil {
			return nil, readErr
		}
	} else if readErr != nil {
		return nil, readErr
	}
	return player, nil
}

//Player represents the PLAYER_DECK and PLAYER_INVENTORY entity
type Player struct {
	ID          int    `json:"id"`
	Username    string `json:"username"`
	IDInventory int    `json:"idInventory"`
	IDDecks     []int  `json:"idDecks"`
	//db is a transient pointer to database connection
	db *sql.DB
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

//Marshal writes a json representation of Inventory
func (p *Player) Marshal() ([]byte, error) {
	return json.Marshal(p)
}

//Unmarshal reads a json representation of Inventory
func (p *Player) Unmarshal(data []byte) error {
	return json.Unmarshal(data, &p)
}

//SetDB attachs a database connection to Expansion
func (p *Player) SetDB(db *sql.DB) error {
	if db == nil {
		return errors.New("NullDBReferenceError: Message='The db parameter is required'")
	}
	p.db = db
	return nil
}

//GetDB returns the Expansion attached connection
func (p *Player) GetDB() (*sql.DB, error) {
	if p.db == nil {
		return nil, errors.New("NotAttachedError: Message='The db context is null'")
	}
	return p.db, nil
}

//Attach binds a new database connection to Expansion reference
func (p *Player) Attach() error {
	return attachToDB(p)
}

//Fetch fetchs the Row and sets the values into Expansion instance
func (p *Player) Fetch(fetchable Fetchable) error {
	return fetchable.Scan(&p.ID, &p.Username, &p.IDInventory)
}

//Read gets the entity representation from the database.
//Player.Username must not empty to perform a Read operation
func (p *Player) Read() error {
	if err := p.Attach(); err != nil {
		return err
	}
	defer closeDB(p)
	if strings.TrimSpace(p.Username) == "" {
		return errors.New("data.Player.ReadError: Message='Player.Username is empty'")
	}
	query := `
        select p.id, p.username, (select max(i.id) from inventory i where i.id_player = p.id) as id_inventory
        from player p
        where p.username = ?;
    `
	row := p.db.QueryRow(query, p.Username)
	if err := p.Fetch(row); err != nil {
		return err
	}
	//Read Fully
	return p.ReadDecks(-1)
}

//ReadDecks gets one related Cards page of this Deck
//Player.ID must not empty to perform a ReadDecks operation
func (p *Player) ReadDecks(page int) error {
	p.Attach()
	defer closeDB(p)
	if p.ID <= 0 {
		return errors.New("data.Player.ReadDecksError: Message='Player.ID is empty'")
	}

	query := `
        select d.id
        from deck d
        where d.id_player = ? order by d.id
    `

	deckRows, readDecksErr := p.db.Query(query, p.ID)
	if readDecksErr != nil {
		return readDecksErr
	}
	//tempIDDecks := make([]int, 0)
	tempIDDecks := []int{}
	//var tempIDDecks []int
	for deckRows.Next() {
		var nextIDDeck int
		if err := deckRows.Scan(&nextIDDeck); err != nil {
			return err
		}
		tempIDDecks = append(tempIDDecks, nextIDDeck)
	}
	p.IDDecks = tempIDDecks
	return nil
}

//Persist persists the PLAYER record with Player attributes values
func (p *Player) Persist() error {
	p.Attach()
	defer closeDB(p)
	if strings.TrimSpace(p.Username) == "" {
		return errors.New("data.Player.PersistError: Message='Player.Username is empty'")
	}

	//Checks if the player already exists in the database
	countQuery := `select count(1) from player p where p.username = ?`
	countRow := p.db.QueryRow(countQuery, p.Username)
	var countPlayer int
	if countErr := countRow.Scan(&countPlayer); countErr != nil {
		return countErr
	}

	if countPlayer <= 0 {
		insert := `insert into player (username) values (?)`

		insertResult, insertErr := p.db.Exec(insert, p.Username)
		if insertErr != nil {
			return insertErr
		}
		insertedID, lastIDErr := insertResult.LastInsertId()
		if lastIDErr != nil {
			return lastIDErr
		}
		p.ID = int(insertedID)

		insertInventory := `insert into inventory (name, id_player) values (?, ?)`
		insertInvetoryResult, insertInvetoryErr := p.db.Exec(insertInventory, p.Username, p.ID)
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

		_, updateErr := p.db.Exec(update, p.Username)
		if updateErr != nil {
			return updateErr
		}
		l.Infof("data.Player.UpdatedPlayer: ID=%v Username='%v' IDInventory=%v", p.ID, p.Username, p.IDInventory)
	}
	return nil
}

//Inventory represents the INVENTORY entity
type Inventory struct {
	ID       int    `json:"id"`
	Name     string `json:"name"`
	Label    string `json:"label"`
	IDPlayer int    `json:"idPlayer"`
	Cards    []Card `json:"cards"`
	//db is a transient pointer to database connection
	db *sql.DB
}

//SetDB attachs a database connection to Inventory
func (i *Inventory) SetDB(db *sql.DB) error {
	if db == nil {
		return errors.New("NullDBReferenceError: Message='The db parameter is required'")
	}
	i.db = db
	return nil
}

//GetDB returns the Inventory attached connection
func (i *Inventory) GetDB() (*sql.DB, error) {
	if i.db == nil {
		return nil, errors.New("NotAttachedError: Message='The db context is null'")
	}
	return i.db, nil
}

//Attach binds a new database connection to Inventory reference
func (i *Inventory) Attach() error {
	return attachToDB(i)
}

//Persist persists the INVENTORY record with Inventory attributes values
func (i *Inventory) Persist() error {
	i.Attach()
	defer closeDB(i)
	if i.IDPlayer <= 0 {
		return errors.New("data.Inventory.UpdateError: Message='Inventory.IDPlayer is empty'")
	}
	if i.ID <= 0 {
		insertInventoryQuery := `insert into inventory (name, id_player) values (?, ?)`
		insertResult, insertErr := i.db.Exec(insertInventoryQuery, i.Name, i.IDPlayer)
		if insertErr != nil {
			return insertErr
		}
		insertedID, lastIDErr := insertResult.LastInsertId()
		if lastIDErr != nil {
			return lastIDErr
		}
		l.Debugf("Inventory.PersistedNewInventory: ID=%v IDPlayer=%v Name='%v'", insertedID, i.IDPlayer, i.Name)
		i.ID = int(insertedID)
	}

	insertCardQuery :=
		`insert into inventory_card (id_inventory, id_card, quantity)
        values (?, ?, ?)`
	updateCardQuery :=
		`update inventory_card
            set quantity = ?
        where id_inventory = ? and id_card = ?`

	for _, card := range i.Cards {
		_, insertErr := i.db.Exec(insertCardQuery, i.ID, card.ID, card.InventoryCard.Quantity)
		if insertErr != nil {
			l.Errorf("Inventory.InsertCardEx: Message='%v'", insertErr.Error())
			if !primaryKeyViolation.MatchString(insertErr.Error()) {
				return insertErr
			}
			updateResult, updateErr := i.db.Exec(updateCardQuery, card.InventoryCard.Quantity, i.ID, card.ID)
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
func (i *Inventory) Delete() error {
	i.Attach()
	defer closeDB(i)
	if i.ID <= 0 {
		return errors.New("data.Inventory.DeleteError: Message='Inventory.ID is empty'")
	}
	if i.IDPlayer <= 0 {
		return errors.New("data.Inventory.DeleteError: Message='Inventory.IDPlayer is empty'")
	}
	deleteCardsResult, deleteCardsErr := i.db.Exec("delete from inventory_card where id_inventory = ?", i.ID)
	if deleteCardsErr != nil {
		return deleteCardsErr
	}
	cardsRowsDeleted, cardsRowsDeletedErr := deleteCardsResult.RowsAffected()
	if cardsRowsDeletedErr != nil {
		l.Errorf("Inventory.DeleteCardsGetRowsAffectedEx: Message='%v'", cardsRowsDeletedErr.Error())
	} else {
		l.Debugf("Inventory.DeleteCards: DeletedCards=%v ID=%d IDPlayer=%v", cardsRowsDeleted, i.ID, i.IDPlayer)
	}
	deleteResult, deleteErr := i.db.Exec("delete from inventory where id = ? and id_player = ?", i.ID, i.IDPlayer)
	if deleteErr != nil {
		return deleteErr
	}
	rowsDeleted, err := deleteResult.RowsAffected()
	if err != nil {
		l.Errorf("data.Inventory.DeleteGetRowsAffectedEx: Message='%v'", err.Error())
	} else {
		if rowsDeleted != 1 {
			l.Errorf("data.Inventory.DeleteMultipleEx: Message='%d Records was delete for Inventory.ID=%d and Inventory.IDPlayer=%d'", rowsDeleted, i.ID, i.IDPlayer)
		}
	}
	l.Infof("Inventory.Deleted: ID=%d IDPlayer=%v", i.ID, i.IDPlayer)
	return nil
}

//Fetch fetchs the Row and sets the values into Deck instance
func (i *Inventory) Fetch(fetchable Fetchable) error {
	return fetchable.Scan(&i.ID, &i.Name, &i.IDPlayer)
}

//Read gets the entity representation from the database.
//Inventory.ID must not empty to perform a Read operation
func (i *Inventory) Read() error {
	i.Attach()
	defer closeDB(i)
	if i.ID <= 0 {
		return errors.New("data.Inventory.ReadError: Message='Inventory.ID is empty'")
	}
	if i.IDPlayer <= 0 {
		return errors.New("data.Inventory.ReadError: Message='Inventory.IDPlayer is empty'")
	}
	row := i.db.QueryRow("select i.id, i.name, i.id_player from inventory i where i.id = ? and i.id_player = ?", i.ID, i.IDPlayer)
	return i.Fetch(row)
}

//ReadCards gets one related Cards page of this Inventory
//Inventory.ID must not empty to perform a ReadCards operation
func (i *Inventory) ReadCards(page int) error {
	i.Attach()
	defer closeDB(i)
	if i.ID <= 0 {
		return errors.New("data.Inventory.ReadCardsError: Message='Inventory.ID is empty'")
	}
	query :=
		`select c.id, c.multiverse_number, c.name, c.label, coalesce(c.text, ''),
            coalesce(c.manacost_label, ''), coalesce(c.combatpower_label, ''), c.type_label,
            c.id_rarity, coalesce(c.flavor, ''), c.artist,
            c.rate, c.rate_votes, c.id_asset,
            coalesce(i.id_inventory, 1), coalesce(i.quantity, 0),
            e.id, e.name, e.label, a.id_asset
        from card c
            left join expansion e on c.id_expansion = e.id
            left join expansion_asset a on a.id_expansion = e.id and a.id_rarity = c.id_rarity
            left join inventory_card i on i.id_card = c.id
        where i.id_inventory = ?
        order by e.name, abs(c.multiverse_number)`
	rows, readCardsErr := i.db.Query(query, i.ID)
	if readCardsErr != nil {
		return readCardsErr
	}
	//cards := make([]Card, 0)
	cards := []Card{}
	//var cards []Card
	for rows.Next() {
		nextCard := Card{Expansion: Expansion{}, InventoryCard: InventoryCard{}}
		if err := nextCard.Fetch(rows); err != nil {
			return err
		}
		cards = append(cards, nextCard)
	}
	i.Cards = cards
	return nil
}

//Marshal writes a json representation of Inventory
func (i *Inventory) Marshal() ([]byte, error) {
	return json.Marshal(i)
}

//Unmarshal reads a json representation of Inventory
func (i *Inventory) Unmarshal(reader io.Reader) error {
	return json.NewDecoder(reader).Decode(&i)
}

//Deck represents the DECK entity
type Deck struct {
	ID          int    `json:"id"`
	Name        string `json:"name"`
	IDPlayer    int    `json:"idPlayer,omitempty"`
	IDInventory int    `json:"idInventory,omitempty"`
	Cards       []Card `json:"cards,omitempty"`

	//db is a transient pointer to database connection
	db *sql.DB
}

//SetDB attachs a database connection to Deck
func (d *Deck) SetDB(db *sql.DB) error {
	if db == nil {
		return errors.New("NullDBReferenceError: Message='The db parameter is required'")
	}
	d.db = db
	return nil
}

//GetDB returns the Deck attached connection
func (d *Deck) GetDB() (*sql.DB, error) {
	if d.db == nil {
		return nil, errors.New("NotAttachedError: Message='The db context is null'")
	}
	return d.db, nil
}

//Attach binds a new database connection to Deck reference
func (d *Deck) Attach() error {
	return attachToDB(d)
}

//Persist persists the DECK record with Deck attributes values
func (d *Deck) Persist() error {
	d.Attach()
	defer closeDB(d)
	if d.Name == "" {
		return errors.New("data.Deck.PersistError: Message='Deck.Name is empty'")
	}
	if d.IDPlayer <= 0 {
		return errors.New("data.Deck.PersistError: Message='Deck.IDPlayer is empty'")
	}

	if d.ID == 0 {
		createResult, createErr := d.db.Exec("insert into deck (name, id_player) values (?, ?)", d.Name, d.IDPlayer)
		if createErr != nil {
			return createErr
		}
		createdID, lastIDErr := createResult.LastInsertId()
		if lastIDErr != nil {
			return lastIDErr
		}
		l.Debugf("Deck.PersistNewDeck: ID=%v IDPlayer=%v Name='%v'", createdID, d.IDPlayer, d.Name)
		d.ID = int(createdID)
	} else {
		updateResult, updateErr := d.db.Exec("update deck set name = ? where id = ? and id_player = ?", d.Name, d.ID, d.IDPlayer)
		if updateErr != nil {
			return updateErr
		}
		rowsUpdated, err := updateResult.RowsAffected()
		if err != nil {
			l.Errorf("Deck.UpdateGetRowsAffectedEx: Message='%v'", err.Error())
		} else {
			if rowsUpdated != 1 {
				l.Errorf("data.Deck.UpdateMultipleEx: Message='%d Records was update for Deck.ID=%d Deck.IDPlayer=%d'", rowsUpdated, d.ID, d.IDPlayer)
			}
		}
		l.Debugf("data.Deck.PersistOldDeck: ID=%v IDPlayer=%v Name='%v'", d.ID, d.IDPlayer, d.Name)
	}

	if _, deleteErr := d.db.Exec("delete from deck_card where id_deck = ?", d.ID); deleteErr != nil {
		return deleteErr
	}

	insertCardQuery :=
		`insert into deck_card (id_deck, id_card, id_board, quantity)
        values (?, ?, ?, ?)`
	updateCardQuery :=
		`update deck_card
            set quantity = ?
        where id_deck = ? and id_card = ? and id_board = ?`

	for _, card := range d.Cards {
		_, insertErr := d.db.Exec(insertCardQuery, d.ID, card.ID, card.DeckCard.IDBoard, card.DeckCard.Quantity)
		if insertErr != nil {
			l.Errorf("Deck.InsertCardEx: Message='%v'", insertErr.Error())
			if !primaryKeyViolation.MatchString(insertErr.Error()) {
				return insertErr
			}
			updateResult, updateErr := d.db.Exec(updateCardQuery, card.DeckCard.Quantity, d.ID, card.ID, card.DeckCard.IDBoard)
			if updateErr != nil {
				return updateErr
			}
			rowsUpdated, err := updateResult.RowsAffected()
			if err != nil {
				l.Errorf("Deck.UpdateGetRowsAffectedEx: Message='%v'", err.Error())
			} else {
				if rowsUpdated != 1 {
					l.Errorf("Deck.UpdateMultipleEx: Message='%d Records was update for Inventory.ID=%d'", rowsUpdated, d.ID)
				}
			}
		}
	}
	l.Infof("Deck.Persisted: ID=%v IDPlayer=%v", d.ID, d.IDPlayer)
	return nil
}

//Delete deletes the DECK record references to Deck
func (d *Deck) Delete() error {
	d.Attach()
	defer closeDB(d)
	if d.ID <= 0 {
		return errors.New("data.Deck.DeleteError: Message='Deck.ID is empty'")
	}
	if d.IDPlayer <= 0 {
		return errors.New("data.Deck.DeleteError: Message='Deck.IDPlayer is empty'")
	}
	deleteCardsResult, deleteCardsErr := d.db.Exec("delete from deck_card where id_deck = ?", d.ID)
	if deleteCardsErr != nil {
		return deleteCardsErr
	}
	cardsRowsDeleted, err := deleteCardsResult.RowsAffected()
	if err != nil {
		l.Errorf("Deck.DeleteCardsGetRowsAffectedEx: Message='%v'", err.Error())
	}
	l.Debugf("Deck.DeletedDeckCards: RowsDeleted=%d Deck.ID=%d and Deck.IDPlayer=%d", cardsRowsDeleted, d.ID, d.IDPlayer)
	deleteResult, deleteErr := d.db.Exec("delete from deck where id = ? and id_player = ?", d.ID, d.IDPlayer)
	if deleteErr != nil {
		return deleteErr
	}
	rowsDeleted, err := deleteResult.RowsAffected()
	if err != nil {
		l.Errorf("Deck.DeleteGetRowsAffectedEx: Message='%v'", err.Error())
	} else {
		if rowsDeleted != 1 {
			l.Errorf("Deck.DeleteMultipleEx: Message='%d Records was delete for Deck.ID=%d and Deck.IDPlayer=%d'", rowsDeleted, d.ID, d.IDPlayer)
		}
	}
	l.Infof("Deck.Deleted: Deck.ID=%d and Deck.IDPlayer=%d", d.ID, d.IDPlayer)
	return nil
}

//Fetch fetchs the Row and sets the values into Deck instance
func (d *Deck) Fetch(fetchable Fetchable) error {
	return fetchable.Scan(&d.ID, &d.Name, &d.IDPlayer)
}

//Read gets the entity representation from the database.
//Deck.ID must not empty to perform a Read operation
func (d *Deck) Read() error {
	d.Attach()
	defer closeDB(d)
	if d.ID <= 0 {
		return errors.New("data.Asset.ReadError: Message='Deck.ID is empty'")
	}
	if d.IDPlayer <= 0 {
		return errors.New("data.Asset.ReadError: Message='Deck.IDPlayer is empty'")
	}
	row := d.db.QueryRow("select d.id, d.name, d.id_player from deck d where d.id = ? and d.id_player = ?", d.ID, d.IDPlayer)
	if err := d.Fetch(row); err != nil {
		return err
	}
	//Read Fully
	return d.ReadCards(-1)
}

//ReadCards gets one related Cards page of this Deck
//Deck.ID must not empty to perform a ReadCards operation
func (d *Deck) ReadCards(page int) error {
	d.Attach()
	defer closeDB(d)
	if d.ID <= 0 {
		return errors.New("data.Deck.ReadCardsError: Message='Deck.ID is empty'")
	}
	if d.IDInventory <= 0 {
		return errors.New("data.Deck.ReadCardsError: Message='Deck.IDInventory is empty'")
	}

	query :=
		`select c.id, c.multiverse_number, c.name, c.label, coalesce(c.text, ''),
            coalesce(c.manacost_label, ''), coalesce(c.combatpower_label, ''), c.type_label,
            c.id_rarity, coalesce(c.flavor, ''), c.artist,
            c.rate, c.rate_votes, c.id_asset,
            coalesce(i.id_inventory, ?) as id_inventory, coalesce(i.quantity, 0) as inventory_quantity,
            d.id_deck, d.id_board, coalesce(d.quantity, 0) as deck_quantity,
            e.id, e.name, e.label, a.id_asset
        from card c
            left join expansion e on c.id_expansion = e.id
            left join expansion_asset a on a.id_expansion = e.id and a.id_rarity = c.id_rarity
            left join inventory_card i on i.id_inventory = ? and i.id_card = c.id
            left join deck_card d on d.id_card = c.id
        where d.id_deck = ? 
        order by c.type_label, e.name, abs(c.multiverse_number)`

	cardRows, readCardsErr := d.db.Query(query, d.IDInventory, d.IDInventory, d.ID)
	if readCardsErr != nil {
		return readCardsErr
	}
	//tempCards := make([]Card, selectLimit)
	var tempCards []Card
	for cardRows.Next() {
		nextCard := Card{Expansion: Expansion{}, InventoryCard: InventoryCard{}, DeckCard: DeckCard{}}
		if err := nextCard.FetchWithDeckCard(cardRows); err != nil {
			return err
		}
		tempCards = append(tempCards, nextCard)
	}
	if tempCards == nil {
		tempCards = make([]Card, 0)
	}
	d.Cards = tempCards
	return nil
}

//FetchByID fetchs the Row and sets the values into Deck instance
func (d *Deck) FetchByID(fetchable raizel.Fetchable) error {
	return fetchable.Scan(&d.ID, &d.Name, &d.IDPlayer)
}

func (d *Deck) PersistV2(client raizel.Client) error {
	return errors.New("NotImplementedErr")
}

//Persist persists the DECK record with Deck attributes values
// func (d *Deck) PersistV2(client raizel.Client) error {
// 	if d.Name == "" {
// 		return errors.New("data.Deck.PersistError: Message='Deck.Name is empty'")
// 	}
// 	if d.IDPlayer <= 0 {
// 		return errors.New("data.Deck.PersistError: Message='Deck.IDPlayer is empty'")
// 	}

// 	if d.ID == 0 {
// 		createResult, createErr := client.Exec("insert into deck (name, id_player) values (?, ?)", d.Name, d.IDPlayer)
// 		if createErr != nil {
// 			return createErr
// 		}
// 		createdID, lastIDErr := createResult.LastInsertId()
// 		if lastIDErr != nil {
// 			return lastIDErr
// 		}
// 		l.Debugf("Deck.PersistNewDeck: ID=%v IDPlayer=%v Name='%v'", createdID, d.IDPlayer, d.Name)
// 		d.ID = int(createdID)
// 	} else {
// 		updateResult, updateErr := client.Exec("update deck set name = ? where id = ? and id_player = ?", d.Name, d.ID, d.IDPlayer)
// 		if updateErr != nil {
// 			return updateErr
// 		}
// 		rowsUpdated, err := updateResult.RowsAffected()
// 		if err != nil {
// 			l.Errorf("Deck.UpdateGetRowsAffectedEx: Message='%v'", err.Error())
// 		} else {
// 			if rowsUpdated != 1 {
// 				l.Errorf("data.Deck.UpdateMultipleEx: Message='%d Records was update for Deck.ID=%d Deck.IDPlayer=%d'", rowsUpdated, d.ID, d.IDPlayer)
// 			}
// 		}
// 		l.Debugf("data.Deck.PersistOldDeck: ID=%v IDPlayer=%v Name='%v'", d.ID, d.IDPlayer, d.Name)
// 	}

// 	if _, deleteErr := d.db.Exec("delete from deck_card where id_deck = ?", d.ID); deleteErr != nil {
// 		return deleteErr
// 	}

// 	insertCardQuery :=
// 		`insert into deck_card (id_deck, id_card, id_board, quantity)
//         values (?, ?, ?, ?)`
// 	updateCardQuery :=
// 		`update deck_card
//             set quantity = ?
//         where id_deck = ? and id_card = ? and id_board = ?`

// 	for _, card := range d.Cards {
// 		_, insertErr := d.db.Exec(insertCardQuery, d.ID, card.ID, card.DeckCard.IDBoard, card.DeckCard.Quantity)
// 		if insertErr != nil {
// 			l.Errorf("Deck.InsertCardEx: Message='%v'", insertErr.Error())
// 			if !primaryKeyViolation.MatchString(insertErr.Error()) {
// 				return insertErr
// 			}
// 			updateResult, updateErr := d.db.Exec(updateCardQuery, card.DeckCard.Quantity, d.ID, card.ID, card.DeckCard.IDBoard)
// 			if updateErr != nil {
// 				return updateErr
// 			}
// 			rowsUpdated, err := updateResult.RowsAffected()
// 			if err != nil {
// 				l.Errorf("Deck.UpdateGetRowsAffectedEx: Message='%v'", err.Error())
// 			} else {
// 				if rowsUpdated != 1 {
// 					l.Errorf("Deck.UpdateMultipleEx: Message='%d Records was update for Inventory.ID=%d'", rowsUpdated, d.ID)
// 				}
// 			}
// 		}
// 	}
// 	l.Infof("Deck.Persisted: ID=%v IDPlayer=%v", d.ID, d.IDPlayer)
// 	return nil
// }

func (d *Deck) ReadByID(client raizel.Client) error {
	if d.ID <= 0 {
		return errors.New("data.Asset.ReadByIDError: Message='Deck.ID is empty'")
	}
	err := client.QueryOne("select d.id, d.name, d.id_player from deck d where d.id = $1", d.FetchByID, d.ID)
	if err != nil {
		return err
	}
	//Read Fully
	return d.ReadCardsByID(client, -1)
}

func (d *Deck) ReadByName(client raizel.Client) error {
	if strings.TrimSpace(d.Name) == "" {
		return errors.New("data.Deck.ReadByNameErr: Message='Deck.Name is empty'")
	}
	err := client.QueryOne("select d.id, d.name, d.id_player from deck d where d.name = $1", d.FetchByID, d.Name)
	if err != nil {
		return err
	}
	//Read Fully
	return d.ReadCardsByID(client, -1)
}

func (d *Deck) ReadCardsByID(client raizel.Client, page int) error {
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

	//tempCards := make([]Card, selectLimit)
	var tempCards []Card
	cardsFetchFunc := func(i raizel.Iterable) error {
		for i.Next() {
			nextCard := Card{}
			if err := nextCard.FetchByIDWithDeckCard(i); err != nil {
				return err
			}
			tempCards = append(tempCards, nextCard)
		}
		return nil
	}

	readCardsErr := client.Query(query, cardsFetchFunc, d.ID)
	if readCardsErr != nil {
		return readCardsErr
	}
	if tempCards == nil {
		tempCards = make([]Card, 0)
	}
	d.Cards = tempCards
	return nil
}

func (d *Deck) QueryByName(client raizel.Client) ([]Deck, error) {
	var resultDecks []Deck
	iterFunc := func(i raizel.Iterable) error {
		for i.Next() {
			var deck Deck
			if fetchErr := deck.FetchByID(i); fetchErr != nil {
				return fetchErr
			}
			resultDecks = append(resultDecks, deck)
		}
		return nil
	}
	var err error
	if strings.TrimSpace(d.Name) != "" {
		err = client.Query("select d.id, d.name, d.id_player from deck d where d.name ~* $1", iterFunc, d.Name)
	} else {
		err = client.Query("select d.id, d.name, d.id_player from deck d", iterFunc)
	}

	if err != nil {
		return nil, err
	}
	l.Infof("ResultDecks=%+v", resultDecks)
	return resultDecks, nil
}

//Query querys DECKs by restrictions and create a list of Decks references
func (d *Deck) Query(queryParameters map[string]interface{}, order string) ([]interface{}, error) {
	if d.IDPlayer <= 0 {
		return nil, errors.New("data.Deck.ReadCardsError: Message='Deck.IDPlayer is empty'")
	}
	var parameterSize int
	if queryParameters == nil {
		parameterSize = 0
	} else {
		parameterSize = len(queryParameters)
	}
	restrictions := make([]string, parameterSize)
	values := make([]interface{}, parameterSize)
	if parameterSize > 0 {
		paramIndex := 0
		for k, v := range queryParameters {
			restrictions[paramIndex] = k
			values[paramIndex] = v
			paramIndex++
		}
	}
	query :=
		`
    select d.id, d.name, d.id_player from deck d where d.id_player = ?
    `
	if len(restrictions) > 0 {
		query += " and " + strings.Join(restrictions, " and ") + "\n"
	}

	if order != "" {
		query += "order by " + order
	} else {
		query += "order by d.id"
	}

	values = append([]interface{}{d.IDPlayer}, values...)
	l.Debugf("Deck.Query: Query=%v Parameters=%v", query, values)
	if err := d.Attach(); err != nil {
		return nil, err
	}
	defer closeDB(d)
	rows, queryErr := d.db.Query(query, values...)
	if queryErr != nil {
		return nil, queryErr
	}
	//Limit the result
	//decks := make([]interface{}, 0)
	decks := []interface{}{}
	//var decks []interface{}
	for rows.Next() {
		nextDeck := Deck{}
		if err := nextDeck.Fetch(rows); err != nil {
			return nil, err
		}
		decks = append(decks, nextDeck)
	}
	return decks, nil
}

//Marshal writes a json representation of Deck
func (d *Deck) Marshal() ([]byte, error) {
	return json.Marshal(d)
}

//Unmarshal reads a json representation of Deck
func (d *Deck) Unmarshal(reader io.Reader) error {
	return json.NewDecoder(reader).Decode(&d)
}
