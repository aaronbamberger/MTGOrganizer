package mtgcards

import "fmt"
import "reflect"
import "sort"
import "strings"

type MTGCardCommon struct {
	// Atomic properties
	// These don't change between different variations of the same card

	// Non-optional
	ColorIdentity []string `json:"colorIdentity"`
	Colors []string `json:"colors"`
	Layout string `json:"layout"`
	ManaCost string `json:"manaCost"`
	Power string `json:"power"`
	ScryfallOracleId string `json:"scryfallOracleId"`
	Subtypes []string `json:"subtypes"`
	Supertypes []string `json:"supertypes"`
	Text string `json:"text"`
	Toughness string `json:"toughness"`
	Type string `json:"type"`
	Types []string `json:"types"`

	// Optional
	ColorIndicator []string `json:"colorIndicator"`
	Loyalty string `json:"loyalty"`
	Name string `json:"name"`
	Names []string `json:"names"`
	Side string `json:"side"`

	// Non-atomic properties
	// These may differ between printings of the same card

	// Non-optional
	Artist string `json:"artist"`
	BorderColor string `json:"borderColor"`
	Number string `json:"number"`
	ScryfallId string `json:"scryfallId"`
	UUID string `json:"uuid"`
	Watermark string `json:"watermark"`

	// Optional
	IsOnlineOnly bool `json:"isOnlineOnly"`
	ScryfallIllustrationId string `json:"scryfallIllustrationId"`
}

func (card MTGCardCommon) String() string {
    var b strings.Builder

    fmt.Fprintf(&b, "Color Identity: %s\n", card.ColorIdentity)
    fmt.Fprintf(&b, "Colors: %s\n", card.Colors)
    fmt.Fprintf(&b, "Subtypes: %s\n", card.Subtypes)
    fmt.Fprintf(&b, "Supertypes: %s\n", card.Supertypes)
    fmt.Fprintf(&b, "Types: %s\n", card.Types)
    fmt.Fprintf(&b, "Color Indicator: %s\n", card.ColorIndicator)
    fmt.Fprintf(&b, "Names: %s\n", card.Names)
    fmt.Fprintf(&b, "Side: %s\n", card.Side)
    fmt.Fprintf(&b, "UUID: %s\n", card.UUID)
    fmt.Fprintf(&b, "Scryfall illustration ID: %s\n", card.ScryfallIllustrationId)

    return b.String()
}

func (card *MTGCardCommon) Diff(other *MTGCardCommon) string {
    var b strings.Builder

    cardType := reflect.TypeOf(*card)

    selfValue := reflect.ValueOf(*card)
    otherValue := reflect.ValueOf(*other)

    for i := 0; i < cardType.NumField(); i++ {
        cardField := cardType.Field(i)
        selfField := selfValue.FieldByName(cardField.Name)
        otherField := otherValue.FieldByName(cardField.Name)
        switch selfField.Kind() {
        case reflect.Bool:
            if selfField.Bool() != otherField.Bool() {
                fmt.Fprintf(&b, "%s (%t | %t)\n", cardField.Name,
                    selfField.Bool(), otherField.Bool())
            }
        case reflect.Int:
            if selfField.Int() != otherField.Int() {
                fmt.Fprintf(&b, "%s (%d | %d)\n", cardField.Name,
                    selfField.Int(), otherField.Int())
            }
        case reflect.String:
            if selfField.String() != otherField.String() {
                fmt.Fprintf(&b, "%s (%s | %s)\n", cardField.Name,
                    selfField.String(), otherField.String())
            }

        case reflect.Slice:
            if selfField.Len() != otherField.Len() {
                fmt.Fprintf(&b, "%s length (%d | %d)\n", cardField.Name,
                    selfField.Len(), otherField.Len())
            }
            // TODO: Diff values

        case reflect.Map:
            // TODO: Diff maps

        }
    }

    return b.String()
}

func (card *MTGCardCommon) Canonicalize() {
    sort.Strings(card.ColorIdentity)
	sort.Strings(card.Colors)
	sort.Strings(card.Subtypes)
	sort.Strings(card.Supertypes)
	sort.Strings(card.Types)
	sort.Strings(card.ColorIndicator)
	sort.Strings(card.Names)
}

func (card *MTGCardCommon) Hash() string {
    return objectHash(*card)
}

type MTGCard struct {
	// Card properties common between normal cards and tokens
	MTGCardCommon

	// Atomic properties
	// These don't change between different variations of the same card

	// Non-optional 
	ConvertedManaCost float32 `json:"convertedManaCost"`
	FaceConvertedManaCost float32 `json:"faceConvertedManaCost"`
	AlternateLanguageData []MTGCardAlternateLanguageInfo `json:"foreignData"`
	Legalities map[string]string `json:"legalities"`
	MTGStocksId int `"mtgstocksId"`
	Printings []string `json:"printings"`
	PurchaseURLs map[string]string `json:"purchaseUrls"`
	Rulings []MTGCardRuling `json:"rulings"`

	// Optional
    AsciiName string `json:"asciiName"`
	EDHRecRank int `json:"edhrecRank"`
	Hand string `json:"hand"`
	IsReserved bool `json:"isReserved"`
	LeadershipSkills map[string]bool `json:"leadershipSkills"`
	Life string `json:"life"`
	Loyalty string `json:"loyalty"`
	ManaCost string `json:"manaCost"`

	// Non-atomic properties
	// These may differ between printings of the same card

	// Non-optional
	FrameEffects []string `json:"frameEffects"`
	FrameVersion string `json:"frameVersion"`
	MCMId int `json:"mcmId"`
	MCMMetaId int `json:"mcmMetaId"`
	MultiverseId int `json:"multiverseId"`
	OriginalText string `json:"originalText"`
	OriginalType string `json:"originalType"`
	Rarity string `json:"rarity"`
	TCGPlayerProductId int `json:"tcgplayerProductId"`
	Variations []string `json:"variations"`

	// Optional
	DuelDeck string `json:"duelDeck"`
    FlavorName string `json:"flavorName"`
	FlavorText string `json:"flavorText"`
	HasFoil bool `json:"hasFoil"`
	HasNonFoil bool `json:"hasNonFoil"`
	IsAlternative bool `json:"isAlternative"`
	IsArena bool `json:"isArena"`
    IsBuyABox bool `json:"isBuyABox"`
    IsDateStamped bool `json:"isDateStamped"`
	IsFullArt bool `json:"isFullArt"`
	IsMTGO bool `json:"isMtgo"`
	IsOnlineOnly bool `json:"isOnlineOnly"`
	IsOversized bool `json:"isOversized"`
	IsPaper bool `json:"isPaper"`
	IsPromo bool `json:"isPromo"`
	IsReprint bool `json:"isReprint"`
	IsStarter bool `json:"isStarter"`
	IsStorySpotlight bool `json:"isStorySpotlight"`
	IsTextless bool `json:"isTextless"`
	IsTimeshifted bool `json:"isTimeshifted"`
	MTGArenaId int `json:"mtgArenaId"`
	MTGOFoilId int `json:"mtgoFoilId"`
	MTGOId int `json:"mtgoId"`
	OtherFaceIds []string `json:"otherFaceIds"`
}

func (card MTGCard) String() string {
    var b strings.Builder

    fmt.Fprintf(&b, "%s\n", card.Name)
    fmt.Fprintf(&b, "%s\n", card.MTGCardCommon.String())
    fmt.Fprintf(&b, "Printings: %s\n", card.Printings)
    fmt.Fprintf(&b, "Frame effects: %s\n", card.FrameEffects)
    fmt.Fprintf(&b, "Variations: %s\n", card.Variations)
    fmt.Fprintf(&b, "Other face IDs: %s\n", card.OtherFaceIds)
    for i, langInfo := range card.AlternateLanguageData {
        fmt.Fprintf(&b, "Alt lang info %d:\n%s", i, langInfo)
    }
    for i, ruling := range card.Rulings {
        fmt.Fprintf(&b, "Ruling %d:\n%s", i, ruling)
    }

    return b.String()
}

func (card *MTGCard) Diff(other *MTGCard) string {
    var b strings.Builder

    b.WriteString(card.MTGCardCommon.Diff(&other.MTGCardCommon))

    cardType := reflect.TypeOf(*card)

    selfValue := reflect.ValueOf(*card)
    otherValue := reflect.ValueOf(*other)

    for i := 0; i < cardType.NumField(); i++ {
        cardField := cardType.Field(i)
        selfField := selfValue.FieldByName(cardField.Name)
        otherField := otherValue.FieldByName(cardField.Name)
        switch selfField.Kind() {
        case reflect.Bool:
            if selfField.Bool() != otherField.Bool() {
                fmt.Fprintf(&b, "%s (%t | %t)\n", cardField.Name,
                    selfField.Bool(), otherField.Bool())
            }
        case reflect.Int:
            if selfField.Int() != otherField.Int() {
                fmt.Fprintf(&b, "%s (%d | %d)\n", cardField.Name,
                    selfField.Int(), otherField.Int())
            }
        case reflect.String:
            if selfField.String() != otherField.String() {
                fmt.Fprintf(&b, "%s (%s | %s)\n", cardField.Name,
                    selfField.String(), otherField.String())
            }

        case reflect.Slice:
            if selfField.Len() != otherField.Len() {
                fmt.Fprintf(&b, "%s length (%d | %d)\n", cardField.Name,
                    selfField.Len(), otherField.Len())
            }
            // TODO: Diff values

        case reflect.Map:
            // TODO: Diff maps

        }
    }

    return b.String()
}

func (card *MTGCard) Hash() string {
    return objectHash(*card)
}

func (card *MTGCard) Canonicalize() {
    // First, canonicalize the common properties
    card.MTGCardCommon.Canonicalize()

	sort.Sort(ByLanguageThenName(card.AlternateLanguageData))
	sort.Strings(card.Printings)
	sort.Sort(ByDateThenText(card.Rulings))
	sort.Strings(card.FrameEffects)
	sort.Strings(card.Variations)
	sort.Strings(card.OtherFaceIds)
}

