package carddb

import "database/sql"
import "mtgcards"
import "time"

func GetDbLastUpdate(db *sql.DB) (mtgcards.MTGJSONVersion, error) {
    var version mtgcards.MTGJSONVersion

    res := db.QueryRow(`SELECT
        last_card_update,
        last_prices_update,
        last_mtgjson_version_major,
        last_mtgjson_version_minor,
        last_mtgjson_version_patch
        FROM last_updates`)

    var lastCardUpdate time.Time
    var lastPricesUpdate time.Time
    var lastVersionMajor int
    var lastVersionMinor int
    var lastVersionPatch int
    if err := res.Scan(&lastCardUpdate, &lastPricesUpdate, &lastVersionMajor,
            &lastVersionMinor, &lastVersionPatch); err != nil {
        return version, err
    }

    version.BuildDate = lastCardUpdate
    version.PricesDate = lastPricesUpdate
    version.VersionMajor = lastVersionMajor
    version.VersionMinor = lastVersionMinor
    version.VersionPatch = lastVersionPatch

    return version, nil
}

func UpdateDbLastUpdate(db *sql.DB, newVersion mtgcards.MTGJSONVersion) error {
    _, err := db.Exec(`UPDATE last_updates
        SET
        last_card_update = ?,
        last_prices_update = ?,
        last_mtgjson_version_major = ?,
        last_mtgjson_version_minor = ?,
        last_mtgjson_version_patch = ?`,
        newVersion.BuildDate,
        newVersion.PricesDate,
        newVersion.VersionMajor,
        newVersion.VersionMinor,
        newVersion.VersionPatch)
    if err != nil {
        return err
    }

    return nil
}
