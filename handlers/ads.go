package handlers

import (
	"gomatri/storage"
)

type AdsHandler struct {
	store storage.AdSqlliteStore
}

func CreateAdsHandler(store storage.AdSqlliteStore) *AdsHandler {
	return &AdsHandler{store: store}
}
