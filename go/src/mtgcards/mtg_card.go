package mtgcards

import "encoding/binary"
import "fmt"
import "hash/fnv"
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

    commonHash string
    commonHashValid bool
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
	if !card.commonHashValid {
        hash := fnv.New128a()

        // Start with the common atomic card properties
        for _, colorIdentity := range card.ColorIdentity {
            hash.Write([]byte(colorIdentity))
        }
        for _, color := range card.Colors {
            hash.Write([]byte(color))
        }
        hash.Write([]byte(card.Layout))
        hash.Write([]byte(card.Power))
        hash.Write([]byte(card.ScryfallOracleId))
        for _, subtype := range card.Subtypes {
            hash.Write([]byte(subtype))
        }
        for _, supertype := range card.Supertypes {
            hash.Write([]byte(supertype))
        }
        hash.Write([]byte(card.Text))
        hash.Write([]byte(card.Toughness))
        hash.Write([]byte(card.Type))
        for _, cardType := range card.Types {
            hash.Write([]byte(cardType))
        }
        for _, colorIndicator := range card.ColorIndicator {
            hash.Write([]byte(colorIndicator))
        }
        hash.Write([]byte(card.Loyalty))
        hash.Write([]byte(card.Name))
        for _, name := range card.Names {
            hash.Write([]byte(name))
        }
        hash.Write([]byte(card.Side))

		// Next do the non-atomic common properties
		hash.Write([]byte(card.Artist))
		hash.Write([]byte(card.BorderColor))
		hash.Write([]byte(card.Number))
		hash.Write([]byte(card.ScryfallId))
		hash.Write([]byte(card.UUID))
		hash.Write([]byte(card.Watermark))
		binary.Write(hash, binary.BigEndian, card.IsOnlineOnly)
		hash.Write([]byte(card.ScryfallIllustrationId))

        card.commonHash = hashToHexString(hash)
		card.commonHashValid = true
	}

	return card.commonHash
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
	FlavorText string `json:"flavorText"`
	HasFoil bool `json:"hasFoil"`
	HasNonFoil bool `json:"hasNonFoil"`
	IsAlternative bool `json:"isAlternative"`
	IsArena bool `json:"isArena"`
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

	hash string
	hashValid bool
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
	if !card.hashValid {
        hash := fnv.New128a()

        // Start with the hash of the common properties
        hash.Write([]byte(card.MTGCardCommon.Hash()))

        // Next, do the rest of the atomic card properties
		binary.Write(hash, binary.BigEndian, card.ConvertedManaCost)
		binary.Write(hash, binary.BigEndian, card.FaceConvertedManaCost)
		for idx := range card.AlternateLanguageData {
            languageData := &card.AlternateLanguageData[idx]
            hash.Write([]byte(languageData.Hash()))
		}

		// Since go maps don't have a defined iteration order,
		// Ensure a repeatable hash by sorting the keyset, and using
		// that to define the iteration order
		legalityFormats := make([]string, 0, len(card.Legalities))
		for format, _ := range card.Legalities {
			legalityFormats = append(legalityFormats, format)
		}
		sort.Strings(legalityFormats)
		for _, format := range legalityFormats {
			hash.Write([]byte(format))
			hash.Write([]byte(card.Legalities[format]))
		}

		binary.Write(hash, binary.BigEndian, card.MTGStocksId)
		for _, printing := range card.Printings {
			hash.Write([]byte(printing))
		}

		// Since go maps don't have a defined iteration order,
		// Ensure a repeatable hash by sorting the keyset, and using
		// that to define the iteration order
		purchaseSites := make([]string, 0, len(card.PurchaseURLs))
		for site, _ := range card.PurchaseURLs {
			purchaseSites = append(purchaseSites, site)
		}
		sort.Strings(purchaseSites)
		for _, site := range purchaseSites {
			hash.Write([]byte(site))
			hash.Write([]byte(card.PurchaseURLs[site]))
		}

		for idx := range card.Rulings {
            ruling := &card.Rulings[idx]
            hash.Write([]byte(ruling.Hash()))
		}
		binary.Write(hash, binary.BigEndian, card.EDHRecRank)
		hash.Write([]byte(card.Hand))
		binary.Write(hash, binary.BigEndian, card.IsReserved)

		// Since go maps don't have a defined iteration order,
		// Ensure a repeatable hash by sorting the keyset, and using
		// that to define the iteration order
		leadershipFormats := make([]string, 0, len(card.LeadershipSkills))
		for format, _ := range card.LeadershipSkills {
			leadershipFormats = append(leadershipFormats, format)
		}
		sort.Strings(leadershipFormats)
		for _, format := range leadershipFormats {
			hash.Write([]byte(format))
			binary.Write(hash, binary.BigEndian, card.LeadershipSkills[format])
		}

		hash.Write([]byte(card.Life))
		hash.Write([]byte(card.Loyalty))
		hash.Write([]byte(card.ManaCost))

		// Last, do the rest of the non-atomic card properties
		for _, frameEffect := range card.FrameEffects {
			hash.Write([]byte(frameEffect))
		}
		hash.Write([]byte(card.FrameVersion))
		binary.Write(hash, binary.BigEndian, card.MCMId)
		binary.Write(hash, binary.BigEndian, card.MCMMetaId)
		binary.Write(hash, binary.BigEndian, card.MultiverseId)
		hash.Write([]byte(card.OriginalText))
		hash.Write([]byte(card.OriginalType))
		hash.Write([]byte(card.Rarity))
		binary.Write(hash, binary.BigEndian, card.TCGPlayerProductId)
		for _, variation := range card.Variations {
			hash.Write([]byte(variation))
		}
		hash.Write([]byte(card.DuelDeck))
		hash.Write([]byte(card.FlavorText))
		binary.Write(hash, binary.BigEndian, card.HasFoil)
		binary.Write(hash, binary.BigEndian, card.HasNonFoil)
		binary.Write(hash, binary.BigEndian, card.IsAlternative)
		binary.Write(hash, binary.BigEndian, card.IsArena)
		binary.Write(hash, binary.BigEndian, card.IsFullArt)
		binary.Write(hash, binary.BigEndian, card.IsMTGO)
		binary.Write(hash, binary.BigEndian, card.IsOnlineOnly)
		binary.Write(hash, binary.BigEndian, card.IsOversized)
		binary.Write(hash, binary.BigEndian, card.IsPaper)
		binary.Write(hash, binary.BigEndian, card.IsPromo)
		binary.Write(hash, binary.BigEndian, card.IsReprint)
		binary.Write(hash, binary.BigEndian, card.IsStarter)
		binary.Write(hash, binary.BigEndian, card.IsStorySpotlight)
		binary.Write(hash, binary.BigEndian, card.IsTextless)
		binary.Write(hash, binary.BigEndian, card.IsTimeshifted)
		binary.Write(hash, binary.BigEndian, card.MTGArenaId)
		binary.Write(hash, binary.BigEndian, card.MTGOFoilId)
		binary.Write(hash, binary.BigEndian, card.MTGOId)
		for _, otherFace := range card.OtherFaceIds {
			hash.Write([]byte(otherFace))
		}

        card.hash = hashToHexString(hash)
		card.hashValid = true
	}

	return card.hash
}

func (card *MTGCard) Canonicalize() {
    // First, canonicalize the common properties
    card.MTGCardCommon.Canonicalize()

	sort.Sort(ByLanguage(card.AlternateLanguageData))
	sort.Strings(card.Printings)
	sort.Sort(ByDate(card.Rulings))
	sort.Strings(card.FrameEffects)
	sort.Strings(card.Variations)
	sort.Strings(card.OtherFaceIds)
}

func (card *MTGCard) String() string {
	var builder strings.Builder
	fmt.Fprintf(&builder, "Card: %s\n", card.Name)
	fmt.Fprintf(&builder, "\tNames: %s\n", card.Names)
	fmt.Fprintf(&builder, "\tNumber: %s\n", card.Number)
	fmt.Fprintf(&builder, "\tPower: %s\n", card.Power)
	fmt.Fprintf(&builder, "\tColors: %s\n", card.Colors)
	fmt.Fprintf(&builder, "\tColorIdentity: %s\n", card.ColorIdentity)
	fmt.Fprintf(&builder, "\tColorIndicator: %s\n", card.ColorIndicator)
	fmt.Fprintf(&builder, "\tLayout: %s\n", card.Layout)
	fmt.Fprintf(&builder, "\tLoyalty: %s\n", card.Loyalty)
	fmt.Fprintf(&builder, "\tSide: %s\n", card.Side)
	fmt.Fprintf(&builder, "\tSubtypes: %s\n", card.Subtypes)
	fmt.Fprintf(&builder, "\tSupertypes: %s\n", card.Supertypes)
	fmt.Fprintf(&builder, "\tType: %s\n", card.Type)
	fmt.Fprintf(&builder, "\tTypes: %s\n", card.Types)
	fmt.Fprintf(&builder, "\tUUID: %s\n", card.UUID)
	fmt.Fprintf(&builder, "\tScryfallId: %s\n", card.ScryfallId)
	fmt.Fprintf(&builder, "\tScryfallOracleId: %s\n", card.ScryfallOracleId)
	fmt.Fprintf(&builder, "\tScryfallIllustrationId: %s\n", card.ScryfallIllustrationId)
	fmt.Fprintf(&builder, "\tOtherFaceIds: %s\n", card.OtherFaceIds)
	return builder.String()
}
