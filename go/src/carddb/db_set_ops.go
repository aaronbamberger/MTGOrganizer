package carddb

import "database/sql"
import "mtgcards"

func GetSetHashAndIdFromDB(
        setCode string,
        queries *DBGetQueries) (bool, string, int64, error) {
	res := queries.GetSetHashQuery.QueryRow(setCode)

	var setHash string
	var setId int64
	err := res.Scan(&setHash, &setId)
	if err != nil {
		if err == sql.ErrNoRows {
			// This set isn't in the database
			return false, "", 0, nil
		} else {
			return false, "", 0, err
		}
	} else {
		return true, setHash, setId, nil
	}
}

func InsertSetToDB(
        set *mtgcards.MTGSet,
        queries *DBInsertQueries) (int64, error) {
    setHash := HashToHexString(set.Hash())

    // Insert the set itself
	res, err := queries.InsertSetQuery.Exec(
        setHash,
        set.BaseSetSize,
        set.Block,
		set.Code,
        set.IsForeignOnly,
        set.IsFoilOnly,
        set.IsOnlineOnly,
		set.IsPartialPreview,
        set.KeyruneCode,
        set.MCMName,
        set.MCMId,
        set.MTGOCode,
		set.Name,
        set.ParentCode,
        set.ReleaseDate,
        set.TCGPlayerGroupId,
		set.TotalSetSize,
        set.Type)
	if err != nil {
		return 0, err
	}

    setId, err := res.LastInsertId()
    if err != nil {
        return 0, err
    }

	// Insert the set translations
	for lang, name := range set.Translations {
		err := InsertSetTranslationToDB(setId, lang, name, queries)
		if err != nil {
            return 0, err
		}
	}

	return setId, nil
}

func InsertSetTranslationToDB(
        setId int64,
        translationLang string,
		translatedName string,
        queries *DBInsertQueries) error {
	languageId, err := getSetTranslationLanguageId(translationLang)
	if err != nil {
		return err
	}

	_, err = queries.InsertSetTranslationQuery.Exec(setId, languageId, translatedName)
	if err != nil {
		return err
	}

	return nil
}

func UpdateSetInDB(
        setId int64,
        set *mtgcards.MTGSet,
        updateQueries *DBUpdateQueries,
        deleteQueries *DBDeleteQueries,
        insertQueries *DBInsertQueries) error {

    // First, update the main set record
    setHash := HashToHexString(set.Hash())
	_, err := updateQueries.UpdateSetQuery.Exec(setHash, set.BaseSetSize, set.Block,
		set.Code, set.IsForeignOnly, set.IsFoilOnly, set.IsOnlineOnly,
		set.IsPartialPreview, set.KeyruneCode, set.MCMName, set.MCMId, set.MTGOCode,
		set.Name, set.ParentCode, set.ReleaseDate, set.TCGPlayerGroupId,
		set.TotalSetSize, set.Type, setId)
	if err != nil {
		return err
	}

    // Next, delete the old set translations
    _, err = deleteQueries.DeleteSetTranslationsQuery.Exec(setId)
    if err != nil {
        return err
    }

    // Finally, insert the new set translations
	for lang, name := range set.Translations {
		err := InsertSetTranslationToDB(setId, lang, name, insertQueries)
		if err != nil {
            return err
		}
	}

    return nil
}
