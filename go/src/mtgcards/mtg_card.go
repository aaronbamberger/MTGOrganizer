package mtgcards

import "encoding/binary"
import "fmt"
import "hash"
import "hash/fnv"
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

	atomicPropertiesHash hash.Hash
	atomicPropertiesHashValid bool
	hash hash.Hash
	hashValid bool

}

func (card *MTGCard) AtomicPropertiesHash() hash.Hash {
	if !card.atomicPropertiesHashValid {
        card.atomicPropertiesHash = fnv.New128a()

        // Start with the common atomic card properties
        for _, colorIdentity := range card.ColorIdentity {
            card.atomicPropertiesHash.Write([]byte(colorIdentity))
        }
        for _, color := range card.Colors {
            card.atomicPropertiesHash.Write([]byte(color))
        }
        card.atomicPropertiesHash.Write([]byte(card.Layout))
        card.atomicPropertiesHash.Write([]byte(card.Power))
        card.atomicPropertiesHash.Write([]byte(card.ScryfallOracleId))
        for _, subtype := range card.Subtypes {
            card.atomicPropertiesHash.Write([]byte(subtype))
        }
        for _, supertype := range card.Supertypes {
            card.atomicPropertiesHash.Write([]byte(supertype))
        }
        card.atomicPropertiesHash.Write([]byte(card.Text))
        card.atomicPropertiesHash.Write([]byte(card.Toughness))
        card.atomicPropertiesHash.Write([]byte(card.Type))
        for _, cardType := range card.Types {
            card.atomicPropertiesHash.Write([]byte(cardType))
        }
        for _, colorIndicator := range card.ColorIndicator {
            card.atomicPropertiesHash.Write([]byte(colorIndicator))
        }
        card.atomicPropertiesHash.Write([]byte(card.Loyalty))
        card.atomicPropertiesHash.Write([]byte(card.Name))
        for _, name := range card.Names {
            card.atomicPropertiesHash.Write([]byte(name))
        }
        card.atomicPropertiesHash.Write([]byte(card.Side))

        // Next, do the rest of the atomic card properties
		binary.Write(card.atomicPropertiesHash, binary.BigEndian, card.ConvertedManaCost)
		binary.Write(card.atomicPropertiesHash, binary.BigEndian, card.FaceConvertedManaCost)
		for _, languageData := range card.AlternateLanguageData {
			languageDataHash := languageData.Hash()
			languageDataHashBytes := make([]byte, 0, languageDataHash.Size())
			card.atomicPropertiesHash.Write(languageDataHash.Sum(languageDataHashBytes))
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
			card.atomicPropertiesHash.Write([]byte(format))
			card.atomicPropertiesHash.Write([]byte(card.Legalities[format]))
		}

		binary.Write(card.atomicPropertiesHash, binary.BigEndian, card.MTGStocksId)
		for _, printing := range card.Printings {
			card.atomicPropertiesHash.Write([]byte(printing))
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
			card.atomicPropertiesHash.Write([]byte(site))
			card.atomicPropertiesHash.Write([]byte(card.PurchaseURLs[site]))
		}

		for _, ruling := range card.Rulings {
			rulingHash := ruling.Hash()
			rulingHashBytes := make([]byte, 0, rulingHash.Size())
			card.atomicPropertiesHash.Write(rulingHash.Sum(rulingHashBytes))
		}
		binary.Write(card.atomicPropertiesHash, binary.BigEndian, card.EDHRecRank)
		card.atomicPropertiesHash.Write([]byte(card.Hand))
		binary.Write(card.atomicPropertiesHash, binary.BigEndian, card.IsReserved)

		// Since go maps don't have a defined iteration order,
		// Ensure a repeatable hash by sorting the keyset, and using
		// that to define the iteration order
		leadershipFormats := make([]string, 0, len(card.LeadershipSkills))
		for format, _ := range card.LeadershipSkills {
			leadershipFormats = append(leadershipFormats, format)
		}
		sort.Strings(leadershipFormats)
		for _, format := range leadershipFormats {
			card.atomicPropertiesHash.Write([]byte(format))
			binary.Write(card.atomicPropertiesHash, binary.BigEndian, card.LeadershipSkills[format])
		}

		card.atomicPropertiesHash.Write([]byte(card.Life))
		card.atomicPropertiesHash.Write([]byte(card.Loyalty))
		card.atomicPropertiesHash.Write([]byte(card.ManaCost))

		card.atomicPropertiesHashValid = true
	}

	return card.atomicPropertiesHash
}

func (card *MTGCard) Hash() hash.Hash {
	if !card.hashValid {
		// Start with the hash of the atomic properties
        card.hash = fnv.New128a()
        atomicPropertiesHash := card.AtomicPropertiesHash()
        atomicPropertiesHashBytes := make([]byte, 0, atomicPropertiesHash.Size())
        atomicPropertiesHashBytes = atomicPropertiesHash.Sum(atomicPropertiesHashBytes)
        card.hash.Write(atomicPropertiesHashBytes)

		// Next do the card common properties
		card.hash.Write([]byte(card.Artist))
		card.hash.Write([]byte(card.BorderColor))
		card.hash.Write([]byte(card.Number))
		card.hash.Write([]byte(card.ScryfallId))
		card.hash.Write([]byte(card.UUID))
		card.hash.Write([]byte(card.Watermark))
		binary.Write(card.hash, binary.BigEndian, card.IsOnlineOnly)
		card.hash.Write([]byte(card.ScryfallIllustrationId))

		// Last, do the rest of the card properties
		for _, frameEffect := range card.FrameEffects {
			card.hash.Write([]byte(frameEffect))
		}
		card.hash.Write([]byte(card.FrameVersion))
		binary.Write(card.hash, binary.BigEndian, card.MCMId)
		binary.Write(card.hash, binary.BigEndian, card.MCMMetaId)
		binary.Write(card.hash, binary.BigEndian, card.MultiverseId)
		card.hash.Write([]byte(card.OriginalText))
		card.hash.Write([]byte(card.OriginalType))
		card.hash.Write([]byte(card.Rarity))
		binary.Write(card.hash, binary.BigEndian, card.TCGPlayerProductId)
		for _, variation := range card.Variations {
			card.hash.Write([]byte(variation))
		}
		card.hash.Write([]byte(card.DuelDeck))
		card.hash.Write([]byte(card.FlavorText))
		binary.Write(card.hash, binary.BigEndian, card.HasFoil)
		binary.Write(card.hash, binary.BigEndian, card.HasNonFoil)
		binary.Write(card.hash, binary.BigEndian, card.IsAlternative)
		binary.Write(card.hash, binary.BigEndian, card.IsArena)
		binary.Write(card.hash, binary.BigEndian, card.IsFullArt)
		binary.Write(card.hash, binary.BigEndian, card.IsMTGO)
		binary.Write(card.hash, binary.BigEndian, card.IsOnlineOnly)
		binary.Write(card.hash, binary.BigEndian, card.IsOversized)
		binary.Write(card.hash, binary.BigEndian, card.IsPaper)
		binary.Write(card.hash, binary.BigEndian, card.IsPromo)
		binary.Write(card.hash, binary.BigEndian, card.IsReprint)
		binary.Write(card.hash, binary.BigEndian, card.IsStarter)
		binary.Write(card.hash, binary.BigEndian, card.IsStorySpotlight)
		binary.Write(card.hash, binary.BigEndian, card.IsTextless)
		binary.Write(card.hash, binary.BigEndian, card.IsTimeshifted)
		binary.Write(card.hash, binary.BigEndian, card.MTGArenaId)
		binary.Write(card.hash, binary.BigEndian, card.MTGOFoilId)
		binary.Write(card.hash, binary.BigEndian, card.MTGOId)
		for _, otherFace := range card.OtherFaceIds {
			card.hash.Write([]byte(otherFace))
		}

		card.hashValid = true
	}

	return card.hash
}

func (card *MTGCard) Canonicalize() {
    sort.Strings(card.ColorIdentity)
	sort.Strings(card.Colors)
	sort.Strings(card.Subtypes)
	sort.Strings(card.Supertypes)
	sort.Strings(card.Types)
	sort.Strings(card.ColorIndicator)
	sort.Strings(card.Names)
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

type ByUUID []MTGCard

func (cards ByUUID) Len() int {
	return len(cards)
}

func (cards ByUUID) Less(i, j int) bool {
	return cards[i].UUID < cards[j].UUID
}

func (cards ByUUID) Swap(i, j int) {
	cards[i], cards[j] = cards[j], cards[i]
}
