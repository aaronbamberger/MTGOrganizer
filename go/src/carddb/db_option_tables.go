package carddb

import "context"
import "database/sql"
import "sync"

var gameFormatsCache = map[string]int64{}
var legalityOptionsCache = map[string]int64{}
var purchaseSitesCache = map[string]int64{}
var leadershipFormatsCache = map[string]int64{}
var setTranslationLanguagesCache = map[string]int64{}
var baseTypeOptionsCache = map[string]int64{}
var frameEffectOptionsCache = map[string]int64{}
var subtypeOptionsCache = map[string]int64{}
var supertypeOptionsCache = map[string]int64{}

var gameFormatsCacheMutex sync.Mutex
var legalityOptionsCacheMutex sync.Mutex
var purchaseSitesCacheMutex sync.Mutex
var leadershipFormatsCacheMutex sync.Mutex
var setTranslationLanguagesCacheMutex sync.Mutex
var baseTypeOptionsCacheMutex sync.Mutex
var frameEffectOptionsCacheMutex sync.Mutex
var subtypeOptionsCacheMutex sync.Mutex
var supertypeOptionsCacheMutex sync.Mutex

var insertGameFormatQuery *sql.Stmt
var insertLegalityOptionQuery *sql.Stmt
var insertPurchaseSiteQuery *sql.Stmt
var insertLeadershipFormatQuery *sql.Stmt
var insertSetTranslationLanguageQuery *sql.Stmt
var insertBaseTypeOptionQuery *sql.Stmt
var insertFrameEffectOptionQuery *sql.Stmt
var insertSubtypeOptionQuery *sql.Stmt
var insertSupertypeOptionQuery *sql.Stmt

var dbConn *sql.Conn

func prepareOptionTables(db *sql.DB) error {
	var err error
	ctx := context.Background()

	dbConn, err = db.Conn(ctx)
	if err != nil {
		return err
	}

	// The following queries are used to potentially update the db
	// if options not in the caches are seen, so are left open upon
	// the exit of this function.  They're closed in a separate cleanup
	// function that must be called
	insertGameFormatQuery, err = dbConn.PrepareContext(ctx,
		`INSERT INTO game_formats (game_format_name) VALUES (?)`)
	if err != nil {
		return err
	}

	insertLegalityOptionQuery, err = dbConn.PrepareContext(ctx,
		`INSERT INTO legality_options (legality_option_name) VALUES (?)`)
	if err != nil {
		return err
	}

	insertPurchaseSiteQuery, err = dbConn.PrepareContext(ctx,
		`INSERT INTO purchase_sites (purchase_site_name) VALUES (?)`)
	if err != nil {
		return err
	}

	insertLeadershipFormatQuery, err = dbConn.PrepareContext(ctx,
		`INSERT INTO leadership_formats (leadership_format_name) VALUES (?)`)
	if err != nil {
		return err
	}

	insertSetTranslationLanguageQuery, err = dbConn.PrepareContext(ctx,
		`INSERT INTO set_translation_languages (set_translation_language) VALUES (?)`)
	if err != nil {
		return err
	}

	insertBaseTypeOptionQuery, err = dbConn.PrepareContext(ctx,
		`INSERT INTO base_type_options (base_type_option) VALUES (?)`)
	if err != nil {
		return err
	}

	insertFrameEffectOptionQuery, err = dbConn.PrepareContext(ctx,
		`INSERT INTO frame_effect_options (frame_effect_option) VALUES (?)`)
	if err != nil {
		return err
	}

	insertSubtypeOptionQuery, err = dbConn.PrepareContext(ctx,
		`INSERT INTO card_subtype_options (subtype_option) VALUES (?)`)
	if err != nil {
		return err
	}

	insertSupertypeOptionQuery, err = dbConn.PrepareContext(ctx,
		`INSERT INTO card_supertype_options (supertype_option) VALUES (?)`)
	if err != nil {
		return err
	}

	// The following queries are only used to initially populate the option
	// caches, so we're fine to close them when this function exits
	retrieveGameFormatsQuery, err := dbConn.PrepareContext(ctx,
		`SELECT game_format_id, game_format_name FROM game_formats`)
	if err != nil {
		return err
	}
	defer retrieveGameFormatsQuery.Close()

	retrieveLegalityOptionsQuery, err := dbConn.PrepareContext(ctx,
		`SELECT legality_option_id, legality_option_name FROM legality_options`)
	if err != nil {
		return err
	}
	defer retrieveLegalityOptionsQuery.Close()

	retrievePurchaseSitesQuery, err := dbConn.PrepareContext(ctx,
		`SELECT purchase_site_id, purchase_site_name FROM purchase_sites`)
	if err != nil {
		return err
	}
	defer retrievePurchaseSitesQuery.Close()

	retrieveSetTranslationLanguagesQuery, err := dbConn.PrepareContext(ctx,
		`SELECT set_translation_language_id, set_translation_language FROM set_translation_languages`)
	if err != nil {
		return err
	}
	defer retrieveSetTranslationLanguagesQuery.Close()

	retrieveBaseTypeOptionsQuery, err := dbConn.PrepareContext(ctx,
		`SELECT base_type_option_id, base_type_option FROM base_type_options`)
	if err != nil {
		return err
	}
	defer retrieveBaseTypeOptionsQuery.Close()

	retrieveLeadershipFormatsQuery, err := dbConn.PrepareContext(ctx,
		`SELECT leadership_format_id, leadership_format_name FROM leadership_formats`)
	if err != nil {
		return err
	}
	defer retrieveLeadershipFormatsQuery.Close()

	retrieveFrameEffectOptionsQuery, err := dbConn.PrepareContext(ctx,
		`SELECT frame_effect_option_id, frame_effect_option FROM frame_effect_options`)
	if err != nil {
		return err
	}
	defer retrieveFrameEffectOptionsQuery.Close()

	retrieveSubtypeOptionsQuery, err := dbConn.PrepareContext(ctx,
		`SELECT subtype_option_id, subtype_option FROM card_subtype_options`)
	if err != nil {
		return err
	}
	defer retrieveSubtypeOptionsQuery.Close()

	retrieveSupertypeOptionsQuery, err := dbConn.PrepareContext(ctx,
		`SELECT supertype_option_id, supertype_option FROM card_supertype_options`)
	if err != nil {
		return err
	}
	defer retrieveSupertypeOptionsQuery.Close()

	err = populateOptionsCache(retrieveGameFormatsQuery, gameFormatsCache)
	if err != nil {
		return err
	}

	err = populateOptionsCache(retrieveLegalityOptionsQuery, legalityOptionsCache)
	if err != nil {
		return err
	}

	err = populateOptionsCache(retrievePurchaseSitesQuery, purchaseSitesCache)
	if err != nil {
		return err
	}

	err = populateOptionsCache(retrieveLeadershipFormatsQuery, leadershipFormatsCache)
	if err != nil {
		return err
	}

	err = populateOptionsCache(retrieveSetTranslationLanguagesQuery, setTranslationLanguagesCache)
	if err != nil {
		return err
	}

	err = populateOptionsCache(retrieveBaseTypeOptionsQuery, baseTypeOptionsCache)
	if err != nil {
		return err
	}

	err = populateOptionsCache(retrieveFrameEffectOptionsQuery, frameEffectOptionsCache)
	if err != nil {
		return err
	}

	err = populateOptionsCache(retrieveSubtypeOptionsQuery, subtypeOptionsCache)
	if err != nil {
		return err
	}

	err = populateOptionsCache(retrieveSupertypeOptionsQuery, supertypeOptionsCache)
	if err != nil {
		return err
	}

	return nil
}

func cleanupOptionTables() {
	if insertGameFormatQuery != nil {
		insertGameFormatQuery.Close()
	}

	if insertLegalityOptionQuery != nil {
		insertLegalityOptionQuery.Close()
	}

	if insertPurchaseSiteQuery != nil {
		insertPurchaseSiteQuery.Close()
	}

	if insertLeadershipFormatQuery != nil {
		insertLeadershipFormatQuery.Close()
	}

	if insertSetTranslationLanguageQuery != nil {
		insertSetTranslationLanguageQuery.Close()
	}

	if insertBaseTypeOptionQuery != nil {
		insertBaseTypeOptionQuery.Close()
	}

	if insertFrameEffectOptionQuery != nil {
		insertFrameEffectOptionQuery.Close()
	}

	if insertSubtypeOptionQuery != nil {
		insertSubtypeOptionQuery.Close()
	}

	if insertSupertypeOptionQuery != nil {
		insertSupertypeOptionQuery.Close()
	}

	if dbConn != nil {
		dbConn.Close()
	}
}

func populateOptionsCache(getOptionsQuery *sql.Stmt, optionsCache map[string]int64) error {
	optionsRows, err := getOptionsQuery.Query()
	if err != nil {
		return  err
	}
	defer optionsRows.Close()

	for optionsRows.Next() {
		if err := optionsRows.Err(); err != nil {
			return err
		}

		var optionId int64
		var option string
		err := optionsRows.Scan(&optionId, &option)
		if err != nil {
			return err
		}
		optionsCache[option] = optionId
	}

	return nil
}

func getGameFormatId(gameFormat string) (int64, error) {
	return getOptionIdFromOption(gameFormat, gameFormatsCache,
		&gameFormatsCacheMutex, insertGameFormatQuery)
}

func getLegalityOptionId(legalityOption string) (int64, error) {
	return getOptionIdFromOption(legalityOption, legalityOptionsCache,
		&legalityOptionsCacheMutex, insertLegalityOptionQuery)
}

func getPurchaseSiteId(purchaseSite string) (int64, error) {
	return getOptionIdFromOption(purchaseSite, purchaseSitesCache,
		&purchaseSitesCacheMutex, insertPurchaseSiteQuery)
}

func getLeadershipFormatId(leadershipFormat string) (int64, error) {
	return getOptionIdFromOption(leadershipFormat, leadershipFormatsCache,
		&leadershipFormatsCacheMutex, insertLeadershipFormatQuery)
}

func getSetTranslationLanguageId(setTranslationLanguage string) (int64, error) {
	return getOptionIdFromOption(setTranslationLanguage, setTranslationLanguagesCache,
		&setTranslationLanguagesCacheMutex, insertSetTranslationLanguageQuery)
}

func getBaseTypeOptionId(baseTypeOption string) (int64, error) {
	return getOptionIdFromOption(baseTypeOption, baseTypeOptionsCache,
		&baseTypeOptionsCacheMutex, insertBaseTypeOptionQuery)
}

func getFrameEffectId(frameEffect string) (int64, error) {
	return getOptionIdFromOption(frameEffect, frameEffectOptionsCache,
		&frameEffectOptionsCacheMutex, insertFrameEffectOptionQuery)
}

func getSubtypeOptionId(subtypeOption string) (int64, error) {
	return getOptionIdFromOption(subtypeOption, subtypeOptionsCache,
		&subtypeOptionsCacheMutex, insertSubtypeOptionQuery)
}

func getSupertypeOptionId(supertypeOption string) (int64, error) {
	return getOptionIdFromOption(supertypeOption, supertypeOptionsCache,
		&supertypeOptionsCacheMutex, insertSupertypeOptionQuery)
}

func getOptionIdFromOption(optionValue string, optionsCache map[string]int64,
		optionsCacheMutex *sync.Mutex, insertOptionQuery *sql.Stmt) (int64, error) {
	// Normally, a sync.Map would be a better fit here for the options cache,
	// since it's specifically optimized for threadsafe access to a map where
	// the entries get read often and the map is only appended to, but we need
	// to use an external mutex here anyway, so no point in having two locking
	// mechanisms at the same time.
	// The reason we need to use an external mutex, is because if we see a value
	// that isn't already in the cache, we have to add it to the database
	// (to get the proper ID), and then add it to the cache for later use.
	// We want to avoid race conditions where multiple threads try to concurrently
	// add the same option to the database before updating the shared cache
	optionsCacheMutex.Lock()
	defer optionsCacheMutex.Unlock()
	optionId, exists := optionsCache[optionValue]
	if !exists {
		// Unlikely case that the value isn't in the cache, so we need to update
		// the db, and update the cache
		res, err := insertOptionQuery.Exec(optionValue)
		if err != nil {
			return 0, err
		}

		optionId, err = res.LastInsertId()
		if err != nil {
			return 0, err
		}
		optionsCache[optionValue] = optionId
	}

	return optionId, nil
}
