package mtgcards

func (card *MTGCard) InsertRemainingAtomicPropertiesToDb(
        queries *dbQueries,
		atomicPropId int64) error {
    var err error

	// Alternate language data
	for _, altLangData := range card.AlternateLanguageData {
		err = altLangData.InsertAltLangDataToDb(queries, atomicPropId)
		if err != nil {
			return err
		}
	}

	// Leadership skills
	for leadershipFormat, leaderValid := range card.LeadershipSkills {
		err = InsertLeadershipSkillToDb(queries, atomicPropId, leadershipFormat, leaderValid)
		if err != nil {
			return err
		}
	}

	// Legalities
	for format, legality := range card.Legalities {
		err = InsertLegalityToDb(queries, atomicPropId, format, legality)
		if err != nil {
			return err
		}
	}

    // Printings
	for _, setCode := range card.Printings {
		err = InsertCardPrintingToDb(queries, atomicPropId, setCode)
		if err != nil {
			return err
		}
	}

	// Purchase URLs
	for site, url := range card.PurchaseURLs {
		err = InsertPurchaseURLToDb(queries, atomicPropId, site, url)
		if err != nil {
			return err
		}
	}

	// Rulings
	for _, ruling := range card.Rulings {
		err = ruling.InsertRulingToDb(queries, atomicPropId)
		if err != nil {
			return err
		}
	}

	// Subtypes
	for _, subtype := range card.Subtypes {
		err = InsertCardSubtypeToDb(queries, atomicPropId, subtype)
		if err != nil {
			return err
		}
	}

	// Supertypes
	for _, supertype := range card.Supertypes {
		err = InsertCardSupertypeToDb(queries, atomicPropId, supertype)
		if err != nil {
			return err
		}
	}

	// Calculate the set of "base" types, which I'm defining as the set
	// subtraction of card.Types - (card.Subtypes + card.Supertypes)
	cardBaseTypes := make(map[string]bool)
	for _, cardType := range card.Types {
		var inSubtype, inSupertype bool
		for _, subtype := range card.Subtypes {
			if subtype == cardType {
				inSubtype = true
				break
			}
		}
		for _, supertype := range card.Supertypes {
			if supertype == cardType {
				inSupertype = true
				break
			}
		}
		if !inSubtype && !inSupertype {
			cardBaseTypes[cardType] = true
		}
	}
	for baseType, _ := range cardBaseTypes {
		err = InsertBaseTypeToDb(queries, atomicPropId, baseType)
		if err != nil {
			return err
		}
	}

	return nil
}

func InsertCardSubtypeToDb(queries *dbQueries, atomicPropertiesId int64, subtype string) error {
	subtypeId, err := getSubtypeOptionId(subtype)
	if err != nil {
		return err
	}

	res, err := queries.InsertCardSubtypeQuery.Exec(atomicPropertiesId, subtypeId)
	if err != nil {
		return err
	}

	return checkRowsAffected(res, 1, "insert card subtype")
}

func InsertCardSupertypeToDb(queries *dbQueries, atomicPropertiesId int64, supertype string) error {
	supertypeId, err := getSupertypeOptionId(supertype)
	if err != nil {
		return err
	}

	res, err := queries.InsertCardSupertypeQuery.Exec(atomicPropertiesId, supertypeId)
	if err != nil {
		return err
	}

	return checkRowsAffected(res, 1, "insert card supertype")
}

func (altLangInfo *MTGCardAlternateLanguageInfo) InsertAltLangDataToDb(queries *dbQueries,
		atomicPropertiesId int64) error {
	res, err := queries.InsertAltLangDataQuery.Exec(atomicPropertiesId, altLangInfo.FlavorText,
		altLangInfo.Language, altLangInfo.MultiverseId, altLangInfo.Name,
		altLangInfo.Text, altLangInfo.Type)
	if err != nil {
		return err
	}

	return checkRowsAffected(res, 1, "insert alt lang info")
}

func (ruling *MTGCardRuling) InsertRulingToDb(queries *dbQueries, atomicPropertiesId int64) error {
	res, err := queries.InsertRulingQuery.Exec(atomicPropertiesId, ruling.Date, ruling.Text)
	if err != nil {
		return err
	}

	return checkRowsAffected(res, 1, "insert ruling")
}

func InsertBaseTypeToDb(queries *dbQueries, atomicPropertiesId int64,
		baseTypeOption string) error {
	baseTypeOptionId, err := getBaseTypeOptionId(baseTypeOption)
	if err != nil {
		return err
	}

	res, err := queries.InsertBaseTypeQuery.Exec(atomicPropertiesId, baseTypeOptionId)
	if err != nil {
		return err
	}

	return checkRowsAffected(res, 1, "insert base type")
}

func InsertLeadershipSkillToDb(queries *dbQueries, atomicPropertiesId int64,
		leadershipFormat string, leaderLegal bool) error {

	leadershipFormatId, err := getLeadershipFormatId(leadershipFormat)
	if err != nil {
		return err
	}

	res, err := queries.InsertLeadershipSkillQuery.Exec(atomicPropertiesId,
        leadershipFormatId, leaderLegal)
	if err != nil {
		return err
	}

	return checkRowsAffected(res, 1, "insert leadership skill")
}

func InsertLegalityToDb(queries *dbQueries, atomicPropertiesId int64, gameFormat string,
		legalityOption string) error {
	gameFormatId, err := getGameFormatId(gameFormat)
	if err != nil {
		return err
	}

	legalityOptionId, err := getLegalityOptionId(legalityOption)
	if err != nil {
		return err
	}

	res, err := queries.InsertLegalityQuery.Exec(atomicPropertiesId, gameFormatId, legalityOptionId)
	if err != nil {
		return err
	}

	return checkRowsAffected(res, 1, "insert legality")
}

func InsertCardPrintingToDb(queries *dbQueries, atomicPropertiesId int64, setCode string) error {
	res, err := queries.InsertCardPrintingQuery.Exec(atomicPropertiesId, setCode)
	if err != nil {
		return err
	}

	return checkRowsAffected(res, 1, "insert card printing")
}

func InsertPurchaseURLToDb(queries *dbQueries, atomicPropertiesId int64,
		purchaseSite string, purchaseURL string) error {
	purchaseSiteId, err := getPurchaseSiteId(purchaseSite)
	if err != nil {
		return err
	}

	res, err := queries.InsertPurchaseURLQuery.Exec(atomicPropertiesId, purchaseSiteId, purchaseURL)
	if err != nil {
		return err
	}

	return checkRowsAffected(res, 1, "insert purchase url")
}
