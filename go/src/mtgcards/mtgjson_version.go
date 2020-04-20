package mtgcards

import "encoding/json"
import "fmt"
import "strconv"
import "strings"
import "time"

type MTGJSONVersion struct {
	BuildDate time.Time
	PricesDate time.Time
    VersionDate time.Time
	VersionMajor int
    VersionMinor int
    VersionPatch int
}

func (version MTGJSONVersion) String() string {
    var builder strings.Builder
    fmt.Fprintf(&builder, "Build Date: %v\n", version.BuildDate)
    fmt.Fprintf(&builder, "Prices Date: %v\n", version.PricesDate)
    fmt.Fprintf(&builder, "Version Date: %v\n", version.VersionDate)
    fmt.Fprintf(&builder, "Version: %d.%d.%d\n", version.VersionMajor,
        version.VersionMinor, version.VersionPatch)
    return builder.String()
}

type mtgjsonDummyVersion struct {
    BuildDate string `json:"date"`
    PricesDate string `json:"pricesDate"`
    Version string `json:"version"`
}

func (version *MTGJSONVersion) UnmarshalJSON(data []byte) error {
    // First, unmarshal into a dummy object
    var dummyVersion mtgjsonDummyVersion
    err := json.Unmarshal(data, &dummyVersion)
    if err != nil {
        return err
    }

    // Now, parse the various dates
    version.BuildDate, err = time.Parse("2006-01-02", dummyVersion.BuildDate)
    if err != nil {
        return err
    }
    version.PricesDate, err = time.Parse("2006-01-02", dummyVersion.PricesDate)
    if err != nil {
        return err
    }

    // Version can either be just a semver version number,
    // <major>.<minor>.<patch>
    // or a semver with a date offset,
    // <major>.<minor>.<patch>+<year><month><day>
    // Handle both of those potential formats, with a missing date offset
    // represented with a zero-valued time.Time object
    versionAndDate := strings.Split(dummyVersion.Version, "+")
    if len(versionAndDate) < 1 {
        return fmt.Errorf("Unexpected version string %s", dummyVersion.Version)
    }

    // Split and parse the parts of the semver version string
    semverParts := strings.Split(versionAndDate[0], ".")
    if len(semverParts) != 3 {
        return fmt.Errorf("Unexpected semver string %s", versionAndDate[0])
    }

    version.VersionMajor, err = strconv.Atoi(semverParts[0])
    if err != nil {
        return err
    }
    version.VersionMinor, err = strconv.Atoi(semverParts[1])
    if err != nil {
        return err
    }
    version.VersionPatch, err = strconv.Atoi(semverParts[2])
    if err != nil {
        return err
    }

    if len(versionAndDate) == 2 {
        version.VersionDate, err = time.Parse("20060102", versionAndDate[1])
        if err != nil {
            return err
        }
    } else {
        version.VersionDate = time.Time{}
    }

    return nil
}
