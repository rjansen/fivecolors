package model

import "github.com/graphql-go/graphql"

var rarityIDType = graphql.NewEnum(
	graphql.EnumConfig{
		Name: "RarityID",
		Values: graphql.EnumValueConfigMap{
			"Common": &graphql.EnumValueConfig{
				Value: NewRarityID(RarityCommon),
			},
			"Uncommon": &graphql.EnumValueConfig{
				Value: NewRarityID(RarityUncommon),
			},
			"Rare": &graphql.EnumValueConfig{
				Value: NewRarityID(RarityRare),
			},
			"MythcRare": &graphql.EnumValueConfig{
				Value: NewRarityID(RarityMythicRare),
			},
		},
	},
)

var rarityType = graphql.NewObject(
	graphql.ObjectConfig{
		Name: "Rarity",
		Fields: graphql.Fields{
			"id": &graphql.Field{
				Type: graphql.String,
			},
			"name": &graphql.Field{
				Type: graphql.String,
			},
			"alias": &graphql.Field{
				Type: graphql.String,
			},
		},
	},
)

var setAssetType = graphql.NewObject(
	graphql.ObjectConfig{
		Name: "SetAsset",
		Fields: graphql.Fields{
			"Common": &graphql.Field{
				Type: graphql.String,
			},
			"Uncommon": &graphql.Field{
				Type: graphql.String,
			},
			"Rare": &graphql.Field{
				Type: graphql.String,
			},
			"MythcRare": &graphql.Field{
				Type: graphql.String,
			},
		},
	},
)

var setType = graphql.NewObject(
	graphql.ObjectConfig{
		Name: "Set",
		Fields: graphql.Fields{
			"id": &graphql.Field{
				Type: graphql.String,
			},
			"name": &graphql.Field{
				Type: graphql.String,
			},
			"alias": &graphql.Field{
				Type: graphql.String,
			},
			"asset": &graphql.Field{
				Type: setAssetType,
			},
		},
	},
)

var cardType = graphql.NewObject(
	graphql.ObjectConfig{
		Name: "Card",
		Fields: graphql.Fields{
			"id": &graphql.Field{
				Type: graphql.String,
			},
			"name": &graphql.Field{
				Type: graphql.String,
			},
			"idExternal": &graphql.Field{
				Type: graphql.String,
			},
			"numberCost": &graphql.Field{
				Type: graphql.Float,
			},
			"idAsset": &graphql.Field{
				Type: graphql.String,
			},
			"types": &graphql.Field{
				Type:    graphql.NewList(graphql.String),
				Resolve: cardResolveTypes,
			},
			"costs": &graphql.Field{
				Type:    graphql.NewList(graphql.String),
				Resolve: cardResolveCosts,
			},
			"rules": &graphql.Field{
				Type:    graphql.NewList(graphql.String),
				Resolve: cardResolveRules,
			},
			"rarity": &graphql.Field{
				Type:    rarityType,
				Resolve: cardResolveRarity,
			},
			"set": &graphql.Field{
				Type:    setType,
				Resolve: cardResolveSet,
			},
		},
	},
)

var cardFilterInputType = graphql.NewInputObject(
	graphql.InputObjectConfig{
		Name: "UserInput",
		Fields: graphql.InputObjectConfigFieldMap{
			"set": &graphql.InputObjectFieldConfig{
				Type: graphql.String,
			},
			"name": &graphql.InputObjectFieldConfig{
				Type: graphql.String,
			},
			"types": &graphql.InputObjectFieldConfig{
				Type: graphql.String,
			},
			"numberCost": &graphql.InputObjectFieldConfig{
				Type: graphql.Float,
			},
			"rarity": &graphql.InputObjectFieldConfig{
				Type: graphql.String,
			},
			"costs": &graphql.InputObjectFieldConfig{
				Type: graphql.String,
			},
			"rules": &graphql.InputObjectFieldConfig{
				Type: graphql.String,
			},
		},
	},
)
