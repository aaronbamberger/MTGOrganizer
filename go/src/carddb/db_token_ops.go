package carddb

import "database/sql"
import "mtgcards"
import "strings"

func GetTokenHashAndIdFromDB(
        tokenUUID string,
        queries *DBGetQueries) (bool, string, int64, error) {
	res := queries.GetTokenHashQuery.QueryRow(tokenUUID)

	var tokenHash string
	var tokenId int64
	err := res.Scan(&tokenHash, &tokenId)
	if err != nil {
		if err == sql.ErrNoRows {
			// This token isn't in the database
			return false, "", 0, nil
		} else {
			return false, "", 0, err
		}
	} else {
		return true, tokenHash, tokenId, nil
	}
}

func InsertTokenToDB(
        token *mtgcards.MTGToken,
        setId int64,
        queries *DBInsertQueries) (bool, error) {
	// Build all of the values that can be null
    var colorIdentity sql.NullString
	var colorIndicator sql.NullString
	var colors sql.NullString
	var loyalty sql.NullString
    var name sql.NullString
	var scryfallIllustrationId sql.NullString
	var side sql.NullString

	if len(token.ColorIdentity) > 0 {
		colorIdentity.String = strings.Join(token.ColorIdentity, ",")
		colorIdentity.Valid = true
	}

	if len(token.ColorIndicator) > 0 {
		colorIndicator.String = strings.Join(token.ColorIndicator, ",")
		colorIndicator.Valid = true
	}

	if len(token.Colors) > 0 {
		colors.String = strings.Join(token.Colors, ",")
		colors.Valid = true
	}

	if len(token.Loyalty) > 0 {
		loyalty.String = token.Loyalty
		loyalty.Valid = true
	}

	if len(token.Name) > 0 {
		name.String = token.Name
		name.Valid = true
	}

	if len(token.ScryfallIllustrationId) > 0 {
		scryfallIllustrationId.String = token.ScryfallIllustrationId
		scryfallIllustrationId.Valid = true
	}

	if len(token.Side) > 0 {
		side.String = token.Side
		side.Valid = true
	}

	res, err := queries.InsertTokenQuery.Exec(
        token.UUID,
        token.Hash(),
        token.Artist,
        token.BorderColor,
        token.Number,
        token.Power,
        token.Type,
		colorIdentity,
		colorIndicator,
		colors,
        token.IsOnlineOnly,
		token.Layout,
		loyalty,
		name,
        token.ScryfallId,
        scryfallIllustrationId,
		token.ScryfallOracleId,
        setId,
		side,
		token.Text,
		token.Toughness,
		token.Watermark)

	if err != nil {
		return false, err
	}

    // Since duplicate tokens are expected, it's possible that the insert query
    // doesn't actually insert anything (the query is written so a duplicate
    // insertion is a no-op).  For an existing token, we don't want to insert
    // the other token data, so bail out early if this is the case
    rowsAffected, err := res.RowsAffected()
    if err != nil {
        return false, err
    }
    if rowsAffected == 0 {
        return false, nil
    }

	tokenId, err := res.LastInsertId()
	if err != nil {
		return false, err
	}

    // Now, insert all of the token data that doesn't live in the all_tokens table
    err = InsertOtherTokenDataToDB(tokenId, token, queries)
    if err != nil {
        return false, err
    }

	return true, nil
}

func UpdateTokenInDB(
        tokenId int64,
        setId int64,
        token *mtgcards.MTGToken,
        updateQueries *DBUpdateQueries,
        deleteQueries *DBDeleteQueries,
        insertQueries *DBInsertQueries) error {
	// Build all of the values that can be null
    var colorIdentity sql.NullString
	var colorIndicator sql.NullString
	var colors sql.NullString
	var loyalty sql.NullString
    var name sql.NullString
	var scryfallIllustrationId sql.NullString
	var side sql.NullString

    if len(token.ColorIdentity) > 0 {
		colorIdentity.String = strings.Join(token.ColorIdentity, ",")
		colorIdentity.Valid = true
	}

	if len(token.ColorIndicator) > 0 {
		colorIndicator.String = strings.Join(token.ColorIndicator, ",")
		colorIndicator.Valid = true
	}

	if len(token.Colors) > 0 {
		colors.String = strings.Join(token.Colors, ",")
		colors.Valid = true
	}

	if len(token.Loyalty) > 0 {
		loyalty.String = token.Loyalty
		loyalty.Valid = true
	}

	if len(token.Name) > 0 {
		name.String = token.Name
		name.Valid = true
	}

	if len(token.ScryfallIllustrationId) > 0 {
		scryfallIllustrationId.String = token.ScryfallIllustrationId
		scryfallIllustrationId.Valid = true
	}

	if len(token.Side) > 0 {
		side.String = token.Side
		side.Valid = true
	}

    // First, update the main token record
	_, err := updateQueries.UpdateTokenQuery.Exec(
        token.Hash(),
        token.Artist,
        token.BorderColor,
        token.Number,
        token.Power,
        token.Type,
		colorIdentity,
		colorIndicator,
		colors,
        token.IsOnlineOnly,
		token.Layout,
		loyalty,
		name,
        token.ScryfallId,
        scryfallIllustrationId,
		token.ScryfallOracleId,
        setId,
		side,
		token.Text,
		token.Toughness,
		token.Watermark,
        token.UUID)
	if err != nil {
		return err
	}

    // Next, delete the rest of the old token data
    err = DeleteOtherTokenDataFromDB(tokenId, deleteQueries)
    if err != nil {
        return err
    }

    // Finally, insert the rest of the new token data
    err = InsertOtherTokenDataToDB(tokenId, token, insertQueries)
    if err != nil {
        return err
    }

    return nil
}

func InsertOtherTokenDataToDB(
        tokenId int64,
        token *mtgcards.MTGToken,
        queries *DBInsertQueries) error {
	// Calculate the set of "base" types, which I'm defining as the set
	// subtraction of card.Types - (card.Subtypes + card.Supertypes)
	tokenBaseTypes := make(map[string]bool)
	for _, tokenType := range token.Types {
		var inSubtype, inSupertype bool
		for _, subtype := range token.Subtypes {
			if subtype == tokenType {
				inSubtype = true
				break
			}
		}
		for _, supertype := range token.Supertypes {
			if supertype == tokenType {
				inSupertype = true
				break
			}
		}
		if !inSubtype && !inSupertype {
			tokenBaseTypes[tokenType] = true
		}
	}
	for baseType, _ := range tokenBaseTypes {
        err := InsertTokenBaseTypeToDB(tokenId, baseType, queries)
		if err != nil {
			return err
		}
	}

	// Subtypes
	for _, subtype := range token.Subtypes {
        err := InsertTokenSubtypeToDB(tokenId, subtype, queries)
		if err != nil {
			return err
		}
	}

	// Supertypes
	for _, supertype := range token.Supertypes {
        err := InsertTokenSupertypeToDB(tokenId, supertype, queries)
		if err != nil {
			return err
		}
	}

    // Reverse related
    for _, reverseRelated := range token.ReverseRelated {
        err := InsertTokenReverseRelatedToDB(tokenId, reverseRelated, queries)
        if err != nil {
            return err
        }
    }

    return nil
}

func DeleteOtherTokenDataFromDB(
        tokenId int64,
        queries *DBDeleteQueries) error {
    // Base types
    _, err := queries.DeleteTokenBaseTypesQuery.Exec(tokenId)
    if err != nil {
        return err
    }

	// Subtypes
    _, err = queries.DeleteTokenSubtypesQuery.Exec(tokenId)
    if err != nil {
        return err
    }

	// Supertypes
    _, err = queries.DeleteTokenSupertypesQuery.Exec(tokenId)
    if err != nil {
        return err
    }

    // Reverse related
    _, err = queries.DeleteTokenReverseRelatedQuery.Exec(tokenId)
    if err != nil {
        return err
    }

    return nil
}

func InsertTokenSubtypeToDB(
        tokenId int64,
        subtype string,
        queries *DBInsertQueries) error {
	subtypeId, err := getSubtypeOptionId(subtype)
	if err != nil {
		return err
	}

	_, err = queries.InsertTokenSubtypeQuery.Exec(tokenId, subtypeId)
	if err != nil {
		return err
	}

    return nil
}

func InsertTokenSupertypeToDB(
        tokenId int64,
        supertype string,
        queries *DBInsertQueries) error {
	supertypeId, err := getSupertypeOptionId(supertype)
	if err != nil {
		return err
	}

	_, err = queries.InsertTokenSupertypeQuery.Exec(tokenId, supertypeId)
	if err != nil {
		return err
	}

	return nil
}

func InsertTokenBaseTypeToDB(
        tokenId int64,
        baseTypeOption string,
        queries *DBInsertQueries) error {
	baseTypeOptionId, err := getBaseTypeOptionId(baseTypeOption)
	if err != nil {
		return err
	}

	_, err = queries.InsertTokenBaseTypeQuery.Exec(tokenId, baseTypeOptionId)
	if err != nil {
		return err
	}

	return nil
}

func InsertTokenReverseRelatedToDB(
        tokenId int64,
        reverseRelatedCard string,
        queries *DBInsertQueries) error {
    _, err := queries.InsertTokenReverseRelatedQuery.Exec(tokenId, reverseRelatedCard)
    if err != nil {
        return err
    }

    return nil
}
