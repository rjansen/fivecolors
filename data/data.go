package data

import (
	"database/sql"
	"encoding/json"
	"errors"
	"io"
	"log"
	"strings"
    "regexp"
)

var (
	selectLimit = 100
    primaryKeyViolation = regexp.MustCompile(`Duplicate.*PRIMARY`)
)

//Fetchable supply the sql.Scan interface for a struct
type Fetchable interface {
	Scan(dest ...interface{}) error
}

//Readable provides persistence actions for a struct
type Readable interface {
	Fetch(fetchable Fetchable) error
	Read() error
	Query(queryParameters map[string]interface{}, order string) ([]interface{}, error)
}

//Writable provides persistence actions for a struct
type Writable interface {
	Create() error
	Update() error
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
	tempDb, err := pool.GetConnection()
	if err != nil {
		return errors.New("data.Card.AttachError: Messages='" + err.Error() + "'")
	}
	return a.SetDB(tempDb)
}

func closeDB(a Attachable) error {
	if db, err := a.GetDB(); err == nil {
		log.Printf("data.ClosingDB: Attachable=%T", a)
		db.Close()
	} else {
		log.Printf("CloseDBError: Message='Error closing sql.DB' Cause='%s'", err.Error())
		return err
	}
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
	IDAsset          int           `json:"idAsset"`
	//db is a transient pointer to database connection
	db *sql.DB
}

//InventoryCard represents the INVENTORY_CARD entity
type InventoryCard struct {
	InventoryID int `json:"invetory"`
	Quantity    int `json:"quantity"`
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
		&c.InventoryCard.InventoryID, &c.InventoryCard.Quantity,
		&c.Expansion.ID, &c.Expansion.Name, &c.Expansion.Label, &c.Expansion.IDAsset)
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
        left join inventory_card i on i.id_inventory = 1 and i.id_card = c.id
        where c.id = ?`
	row := c.db.QueryRow(query, c.ID)
	return c.Fetch(row)
}

//Query querys CARDs by restrictions and create a list of Cards references
func (c *Card) Query(queryParameters map[string]interface{}, order string) ([]interface{}, error) {
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
		`select c.id, c.multiverse_number, c.name, c.label, coalesce(c.text, ''),
                coalesce(c.manacost_label, ''), coalesce(c.combatpower_label, ''), c.type_label,
                c.id_rarity, coalesce(c.flavor, ''), c.artist,
                c.rate, c.rate_votes, c.id_asset,
                coalesce(i.id_inventory, 1), coalesce(i.quantity, 0),
                e.id, e.name, e.label, a.id_asset
        from card c 
        left join expansion e on c.id_expansion = e.id
        left join expansion_asset a on a.id_expansion = e.id and a.id_rarity = c.id_rarity
        left join inventory_card i on i.id_inventory = 1 and i.id_card = c.id
        where ` + strings.Join(restrictions, " and ")

	if order != "" {
		query += " order by " + order
	}

	log.Printf("data.Card.Query: Query=%v Parameters=%v", query, values)
	if err := c.Attach(); err != nil {
		return nil, err
	}
	defer closeDB(c)
	rows, queryErr := c.db.Query(query, values...)
	if queryErr != nil {
		return nil, queryErr
	}
	//Limit the result
	//cards := make([]interface{}, selectLimit)
	var cards []interface{}
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
		from expansion e left join expansion_asset a on e.id = a.id_expansion and a.id_rarity = 0 
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
		`select e.id, e.name, e.label, a.id_asset 
		from expansion e left join expansion_asset a on e.id = a.id_expansion and a.id_rarity = 0`
	if len(restrictions) > 0 {
		query += "where " + strings.Join(restrictions, " and ")
	}

	if order != "" {
		query += " order by " + order
	}

	log.Printf("data.Card.Query: Query=%v Parameters=%v", query, values)
	if err := e.Attach(); err != nil {
		return nil, err
	}
	defer closeDB(e)
	rows, queryErr := e.db.Query(query, values...)
	if queryErr != nil {
		return nil, queryErr
	}
	//Limit the result
	//expansions := make([]interface{}, selectLimit)
	var expansions []interface{}
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

//Inventory represents the INVENTORY entity
type Inventory struct {
	ID    int    `json:"id"`
	Name  string `json:"name"`
	Label string `json:"label"`
	Cards []Card `json:"cards"`
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

//Create creates a new INVENTORY record with Inventory attributes values
func (i *Inventory) Create() error {
	i.Attach()
	defer closeDB(i)
	if i.ID > 0 {
		return errors.New("data.Inventory.CreateError: Message='Inventory.ID is not empty'")
	}
	if i.Name == "" {
		return errors.New("data.Inventory.CreateError: Message='Inventory.Name is empty'")
	}
	if i.Label == "" {
		return errors.New("data.Inventory.CreateError: Message='Inventory.Label is empty'")
	}
	createResult, createErr := i.db.Exec("insert into inventory (name, label) values (?, ?)", i.Name, i.Label)
	if createErr != nil {
		return createErr
	}
	createdID, lastIDErr := createResult.LastInsertId()
	if lastIDErr != nil {
		return lastIDErr
	}
	i.ID = int(createdID)
	for _, card := range i.Cards {
		insertQuery := `insert into inventory_card (id_inventory, id_card, multiverse_number, quantity) 
                        values (?, ?, ?, ?)`
		_, insertErr := i.db.Exec(insertQuery, i.ID, card.ID, card.Index, card.InventoryCard.Quantity)
		if insertErr != nil {
			return insertErr
		}
		card.InventoryCard.InventoryID = i.ID
	}
	return nil
}

//Update updates the INVENTORY record with Inventory attributes values
func (i *Inventory) Update() error {
	i.Attach()
	defer closeDB(i)
	if i.ID <= 0 {
		return errors.New("data.Inventory.UpdateError: Message='Inventory.ID is empty'")
	}
	// if i.Name == "" {
	// 	return errors.New("data.Inventory.UpdateError: Message='Inventory.Name is empty'")
	// }
	// if i.Label == "" {
	// 	return errors.New("data.Inventory.UpdateError: Message='Inventory.Label is empty'")
	// }
	// updateResult, updateErr := i.db.Exec("update inventory set name = ?, label = ? where id = ?", i.Name, i.Label, i.ID)
	// if updateErr != nil {
	// 	return updateErr
	// }
	// rowsUpdated, err := updateResult.RowsAffected()
	// if err == nil {
	// 	log.Printf("data.Inventory.UpdateGetRowsAffectedEx: Message='" + err.Error() + "'")
	// } else {
	// 	if rowsUpdated != 1 {
	// 		log.Printf("data.Inventory.UpdateMultipleEx: Message='%d Records was update for Inventory.ID=%d'", rowsUpdated, i.ID)
	// 	}
	// }
	insertQuery := `insert into inventory_card (id_inventory, id_card, multiverse_number, quantity) 
                        values (?, ?, ?, ?)`
	updateQuery := `update inventory_card 
                        set quantity = ?, multiverse_number = ? 
                        where id_inventory = ? and id_card = ?`

	for _, card := range i.Cards {
		_, insertErr := i.db.Exec(insertQuery, i.ID, card.ID, card.Index, card.InventoryCard.Quantity)
		if insertErr != nil {
			log.Printf("data.Inventory.InsertCardEx: Message='" + insertErr.Error() + "'")
            if !primaryKeyViolation.MatchString(insertErr.Error()) {
			    return insertErr
            }
			updateResult, updateErr := i.db.Exec(updateQuery, card.InventoryCard.Quantity, card.Index, i.ID, card.ID)
			if updateErr != nil {
				return updateErr
			}
			rowsUpdated, err := updateResult.RowsAffected()
			if err != nil {
				log.Printf("data.Inventory.UpdateGetRowsAffectedEx: Message='" + err.Error() + "'")
			} else {
				if rowsUpdated != 1 {
					log.Printf("data.Inventory.UpdateMultipleEx: Message='%d Records was update for Inventory.ID=%d'", rowsUpdated, i.ID)
				}
			}
		}
	}
	return nil
}

//Delete deletes the INVENTORY record references to Inventory
func (i *Inventory) Delete() error {
	i.Attach()
	defer closeDB(i)
	if i.ID <= 0 {
		return errors.New("data.Inventory.DeleteError: Message='Inventory.ID is empty'")
	}
	deleteCardsResult, deleteCardsErr := i.db.Exec("delete from inventory_card where id_inventory = ?", i.ID)
	if deleteCardsErr != nil {
		return deleteCardsErr
	}
	cardsRowsDeleted, cardsRowsDeletedErr := deleteCardsResult.RowsAffected()
	if cardsRowsDeletedErr != nil {
		log.Printf("data.Inventory.DeleteCardsGetRowsAffectedEx: Message='" + cardsRowsDeletedErr.Error() + "'")
	} else {
		log.Printf("data.Inventory.DeleteCards: DeletedCards=%v ID=%d", cardsRowsDeleted, i.ID)
	}
	deleteResult, deleteErr := i.db.Exec("delete from inventory where id = ?", i.ID)
	if deleteErr != nil {
		return deleteErr
	}
	rowsDeleted, err := deleteResult.RowsAffected()
	if err != nil {
		log.Printf("data.Inventory.DeleteGetRowsAffectedEx: Message='" + err.Error() + "'")
	} else {
		if rowsDeleted != 1 {
			log.Printf("data.Inventory.DeleteMultipleEx: Message='%d Records was delete for Inventory.ID=%d'", rowsDeleted, i.ID)
		}
	}
	return nil
}

//Fetch fetchs the Row and sets the values into Deck instance
func (i *Inventory) Fetch(fetchable Fetchable) error {
	return fetchable.Scan(&i.ID, &i.Name, &i.Label)
}

//Read gets the entity representation from the database.
//Inventory.ID must not empty to perform a Read operation
func (i *Inventory) Read() error {
	i.Attach()
	defer closeDB(i)
	if i.ID <= 0 {
		return errors.New("data.Inventory.ReadError: Message='Inventory.ID is empty'")
	}
	row := i.db.QueryRow("select i.id, i.name, i.label from inventory i where i.id = ?", i.ID)
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
		`select ic.quantity, ic.quantity, c.id, c.name, c.label, coalesce(c.text, ''),
                coalesce(c.manacost_label, ''), coalesce(c.combatpower_label, ''), c.type_label,
                c.id_rarity, coalesce(c.flavor, ''), c.artist,
                c.rate, c.rate_votes, c.id_asset,
                e.id, e.name, e.label, a.id_asset
        from inventory i
        left join inventory_card ic on i.id = ic.id_inventory
        left join card c on ic.id_card = c.id
        left join expansion e on c.id_expansion = e.id
        left join expansion_asset a on a.id_expansion = e.id and a.id_rarity = c.id_rarity
        where i.id = ?`
	cardRows, readCardsErr := i.db.Query(query, i.ID, selectLimit)
	if readCardsErr != nil {
		return readCardsErr
	}
	tempCards := make([]Card, selectLimit)
	cardIndex := 0
	for cardRows.Next() {
		nextCard := Card{Expansion: Expansion{}, InventoryCard: InventoryCard{}}
		if err := nextCard.Fetch(cardRows); err != nil {
			return err
		}
		tempCards[cardIndex] = nextCard
		cardIndex++
	}
	i.Cards = tempCards
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
	ID    int    `json:"id"`
	Name  string `json:"name"`
	Label string `json:"label"`
	Cards []Card `json:"cards"`
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

//Create creates a new DECK record with Deck attributes values
func (d *Deck) Create() error {
	d.Attach()
	defer closeDB(d)
	if d.ID > 0 {
		return errors.New("data.Deck.CreateError: Message='Deck.ID is not empty'")
	}
	if d.Name == "" {
		return errors.New("data.Deck.CreateError: Message='Deck.Name is empty'")
	}
	if d.Label == "" {
		return errors.New("data.Deck.CreateError: Message='Deck.Label is empty'")
	}
	createResult, createErr := d.db.Exec("insert into deck (name, label) values (?, ?)", d.Name, d.Label)
	if createErr != nil {
		return createErr
	}
	createdID, lastIDErr := createResult.LastInsertId()
	if lastIDErr != nil {
		return lastIDErr
	}
	d.ID = int(createdID)
	return nil
}

//Update updates the DECK record with Deck attributes values
func (d *Deck) Update() error {
	d.Attach()
	defer closeDB(d)
	if d.ID <= 0 {
		return errors.New("data.Deck.UpdateError: Message='Deck.ID is empty'")
	}
	if d.Name == "" {
		return errors.New("data.Deck.UpdateError: Message='Deck.Name is empty'")
	}
	if d.Label == "" {
		return errors.New("data.Deck.UpdateError: Message='Deck.Label is empty'")
	}
	updateResult, updateErr := d.db.Exec("update deck set name = ?, label = ? where id = ?", d.Name, d.Label, d.ID)
	if updateErr != nil {
		return updateErr
	}
	rowsUpdated, err := updateResult.RowsAffected()
	if err == nil {
		log.Printf("data.Deck.UpdateGetRowsAffectedEx: Message='" + err.Error() + "'")
	} else {
		if rowsUpdated != 1 {
			log.Printf("data.Deck.UpdateMultipleEx: Message='%d Records was update for Deck.ID=%d'", rowsUpdated, d.ID)
		}
	}
	return nil
}

//Delete deletes the DECK record references to Deck
func (d *Deck) Delete() error {
	d.Attach()
	defer closeDB(d)
	if d.ID <= 0 {
		return errors.New("data.Deck.DeleteError: Message='Deck.ID is empty'")
	}
	deleteResult, deleteErr := d.db.Exec("delete from deck where id = ?", d.ID)
	if deleteErr != nil {
		return deleteErr
	}
	rowsDeleted, err := deleteResult.RowsAffected()
	if err != nil {
		log.Printf("data.Deck.DeleteGetRowsAffectedEx: Message='" + err.Error() + "'")
	} else {
		if rowsDeleted != 1 {
			log.Printf("data.Deck.DeleteMultipleEx: Message='%d Records was delete for Deck.ID=%d'", rowsDeleted, d.ID)
		}
	}
	return nil
}

//Fetch fetchs the Row and sets the values into Deck instance
func (d *Deck) Fetch(fetchable Fetchable) error {
	return fetchable.Scan(&d.Name, &d.Label)
}

//Read gets the entity representation from the database.
//Deck.ID must not empty to perform a Read operation
func (d *Deck) Read() error {
	d.Attach()
	defer closeDB(d)
	if d.ID <= 0 {
		return errors.New("data.Asset.ReadError: Message='Deck.ID is empty'")
	}
	row := d.db.QueryRow("select d.name, d.label from deck d where d.id = ?", d.ID)
	return d.Fetch(row)
}

//ReadCards gets one related Cards page of this Deck
//Deck.ID must not empty to perform a ReadCards operation
func (d *Deck) ReadCards(page int) error {
	d.Attach()
	defer closeDB(d)
	if d.ID <= 0 {
		return errors.New("data.Deck.ReadCardsError: Message='Deck.ID is empty'")
	}
	cardRows, readCardsErr := d.db.Query("", d.ID, selectLimit)
	if readCardsErr != nil {
		return readCardsErr
	}
	tempCards := make([]Card, selectLimit)
	cardIndex := 0
	for cardRows.Next() {
		nextCard := Card{}
		if err := nextCard.Fetch(cardRows); err != nil {
			return err
		}
		//append(tempCards, nextCard)
		tempCards[cardIndex] = nextCard
		cardIndex++
	}
	d.Cards = tempCards
	return nil
}

//Marshal writes a json representation of Deck
func (d *Deck) Marshal() ([]byte, error) {
	return json.Marshal(d)
}

//Unmarshal reads a json representation of Deck
func (d *Deck) Unmarshal(reader io.Reader) error {
	return json.NewDecoder(reader).Decode(&d)
}
