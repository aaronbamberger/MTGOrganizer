package mtgcards

import "hash"
import "hash/fnv"

type MTGCardPurchaseURLs struct {
	Cardmarket string `json:"cardmarket"`
	TCGPlayer string `json:"tcgplayer"`
	MTGStocks string `json:"mtgstocks"`

	hash hash.Hash
	hashValid bool
}

func (purchaseURLs *MTGCardPurchaseURLs) Hash() hash.Hash {
	if !purchaseURLs.hashValid {
		purchaseURLs.hash = fnv.New128a()

		purchaseURLs.hash.Write([]byte(purchaseURLs.Cardmarket))
		purchaseURLs.hash.Write([]byte(purchaseURLs.TCGPlayer))
		purchaseURLs.hash.Write([]byte(purchaseURLs.MTGStocks))

		purchaseURLs.hashValid = true
	}

	return purchaseURLs.hash
}
