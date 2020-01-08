package mtgcards

type MTGSet struct {
	BaseSetSize int "json:baseSetSize"
	Block string "json:block"
	Cards []MTGCard "json:cards"
	Code string "json:code"
	IsForeignOnly bool "json:isForeignOnly"
	IsFoilOnly bool "json:isFoilOnly"
	IsOnlineOnly bool "json:isOnlineOnly"
	IsPartialPreview bool "json:isPartialPreview"
	KeyruneCode string "json:keyruneCode"
	MCMName string "json:mcmName"
	MCMId int "json:mcmId"
	MTGOCode string "json:mtgoCode"
	Name string "json:name"
	ParentCode string "json:parentCode"
	ReleaseDate string "json:releaseDate"
	TCGPlayerGroupId int "json:tcgplayerGroupId"
	TotalSetSize int "json:totalSetSize"
	Translations MTGSetNameTranslations "json:translations"
	Type string "json:type"
}

type MTGSetNameTranslations struct {
	ChineseSimplified string "json:Chinese Simplified"
	ChineseTraditional string "json:Chinese Traditional"
	French string "json:French"
	German string "json:German"
	Italian string "json:Italian"
	Japanese string "json:Japanese"
	Korean string "json:Korean"
	PortugeseBrazil string "json:Portugese (Brazil)"
	Russian string "json:Russian"
	Spanish string "json:Spanish"
}

type MTGJSONVersion struct {
	BuildDate string "json:date"
	PricesDate string "json:pricesDate"
	Version string "json:version"
}

type MTGCardAlternateLanguageInfo struct {
	FlavorText string "json:flavorText"
	Language string "json:language"
	MultiverseId int "json:multiverseId"
	Name string "json:name"
	Text string "json:text"
	Type string "json:text"
}

type MTGCardPurchaseURLs struct {
	Cardmarket string "json:cardmarket"
	TCGPlayer string "json:tcgplayer"
	MTGStocks string "json:mtgstocks"
}

type MTGCardRuling struct {
	Date string "json:date"
	Text string "json:text"
}

type MTGCardCommon struct {
	Artist string "json:artist"
	BorderColor string "json:borderColor"
	ColorIdentity []string "json:colorIdentity"
	ColorIndicator []string "json:colorIndicator"
	Colors []string "json:colors"
	IsOnlineOnly bool "json:isOnlineOnly"
	Layout string "json:layout"
	Loyalty string "json:loyalty"
	Name string "json:name"
	Names []string "json:names"
	Number string "json:number"
	Power string "json:power"
	ScryfallId string "json:scryfallId"
	ScryfallOracleId string "json:scryfallOracleId"
	ScryfallIllustrationId string "json:scryfallIllustrationId"
	Side string "json:side"
	Subtypes []string "json:subtypes"
	Supertypes []string "json:supertypes"
	Text string "json:text"
	Toughness string "json:toughness"
	Type string "json:type"
	Types []string "json:types"
	UUID string "json:uuid"
	Watermark string "json:watermark"
}

type MTGToken struct {
	MTGCardCommon
	ReverseRelated []string "json:reverseRelated"
}

type MTGCard struct {
	MTGCardCommon
	ConvertedManaCost float32 "json:convertedManaCost"
	Count int "json:count"
	DuelDeck string "json:duelDeck"
	EDHRecRank int "json:edhrecRank"
	FaceConvertedManaCost float32 "json:faceConvertedManaCost"
	FlavorText string "json:flavorText"
	AlternateLanguageData []MTGCardAlternateLanguageInfo "json:foreignData"
	FrameEffects []string "json:frameEffects"
	FrameVersion string "json:frameVersion"
	Hand string "json:hand"
	HasFoil bool "json:hasFoil"
	HasNonFoil bool "json:hasNonFoil"
	IsAlternative bool "json:isAlternative"
	IsArena bool "json:isArena"
	IsFullArt bool "json:isFullArt"
	IsMTGO bool "json:isMtgo"
	IsOnlineOnly bool "json:isOnlineOnly"
	IsOversized bool "json:isOversized"
	IsPaper bool "json:isPaper"
	IsPromo bool "json:isPromo"
	IsReprint bool "json:isReprint"
	IsReserved bool "json:IsReserved"
	IsStarter bool "json:isStarter"
	IsStorySpotlight bool "json:isStorySpotlight"
	IsTextless bool "json:isTextless"
	IsTimeshifted bool "json:isTimeshifted"
	Layout string "json:layout"
	LeadershipSkills map[string]bool "json:leadershipSkills"
	Legalities map[string]string "json:legalities"
	Life string "json:life"
	Loyalty string "json:loyalty"
	ManaCost string "json:manaCost"
	MCMId int "json:mcmId"
	MCMMetaId int "json:mcmMetaId"
	MTGArenaId int "json:mtgArenaId"
	MTGOFoilId int "json:mtgoFoilId"
	MTGOId int "json:mtgoId"
	MTGStocksId int "mtgstocksId"
	MultiverseId int "json:multiverseId"
	OriginalText string "json:originalText"
	OriginalType string "json:originalType"
	OtherFaceIds []string "json:otherFaceIds"
	Printings []string "json:printings"
	PurchaseURLS MTGCardPurchaseURLs "json:purchaseUrls"
	Rarity string "json:rarity"
	Rulings []MTGCardRuling "json:rulings"
	TCGPlayerProductId int "json:tcgplayerProductId"
	Variations []string "json:variations"
	UUID string "json:uuid"
}
