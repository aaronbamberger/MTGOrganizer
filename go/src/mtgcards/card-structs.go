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

type MTGCard struct {
	Artist string "json:artist"
	BorderColor string "json:borderColor"
	ColorIdentity []string "json:colorIdentity"
	ColorIndicator []string "json:colorIndicator"
	Colors []string "json:colors"
	ConvertedManaCost float32 "json:convertedManaCost"
	DuelDeck string "json:duelDeck"
	HasFoil bool "json:hasFoil"
	IsReserved bool "json:IsReserved"
	Layout string "json:layout"
	MultiverseId int "json:multiverseId"
	Name string "json:name"
	Printings []string "json:printings"
	Rarity string "json:rarity"
	UUID string "json:uuid"
}
