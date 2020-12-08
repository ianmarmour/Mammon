package api

// Auction represents a single auction returned from the Blizzard auctions API.
type Auction struct {
	ID        int64  `json:"id"`
	Item      Item   `json:"item"`
	Quantity  int64  `json:"quantity"`
	UnitPrice int64  `json:"unit_price"`
	TimeLeft  string `json:"time_left"`
}
