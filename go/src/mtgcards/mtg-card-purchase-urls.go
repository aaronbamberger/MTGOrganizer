package mtgcards

import "hash"
import "hash/fnv"

type MTGCardPurchaseURLs struct {
	Cardmarket string `json:"cardmarket"`
	TCGPlayer string `json:"tcgplayer"`
	MTGStocks string `json:"mtgstocks"`
}

func (purchaseURLs MTGCardPurchaseURLs) Hash() hash.Hash {
	hashRes := fnv.New128a()

	hashRes.Write([]byte(purchaseURLs.Cardmarket))
	hashRes.Write([]byte(purchaseURLs.TCGPlayer))
	hashRes.Write([]byte(purchaseURLs.MTGStocks))

	return hashRes
}
