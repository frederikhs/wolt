package wolt

import (
	"time"
)

type SimpleOrder struct {
	OrderId                   string     `json:"order_id" db:"order_id"`
	ClientPreEstimate         string     `json:"client_pre_estimate" db:"client_pre_estimate"`
	DeliveryStreet            string     `json:"delivery_street" db:"delivery_street"`
	DeliveryCoordinateX       float64    `json:"delivery_coordinate_x" db:"delivery_coordinate_x"`
	DeliveryCoordinateY       float64    `json:"delivery_coordinate_y" db:"delivery_coordinate_y"`
	DeliveryDistance          int        `json:"delivery_distance" db:"delivery_distance"`
	DeliveryEta               *time.Time `json:"delivery_eta" db:"delivery_eta"`
	DeliveryMethod            string     `json:"delivery_method" db:"delivery_method"`
	DeliveryPrice             int        `json:"delivery_price" db:"delivery_price"`
	DeliverySizeSurcharge     int        `json:"delivery_size_surcharge" db:"delivery_size_surcharge"`
	DeliveryDistanceSurcharge int        `json:"delivery_distance_surcharge" db:"delivery_distance_surcharge"`
	DeliveryTime              *time.Time `json:"delivery_time" db:"delivery_time"`
	DriverType                string     `json:"driver_type" db:"driver_type"`
	ItemsPrice                int        `json:"items_price" db:"items_price"`
	PaymentAmount             int        `json:"payment_amount" db:"payment_amount"`
	PaymentTime               *time.Time `json:"payment_time" db:"payment_time"`
	Status                    string     `json:"status" db:"status"`
	ServiceFee                int        `json:"service_fee" db:"service_fee"`
	Subscribed                bool       `json:"subscribed" db:"subscribed"`
	TotalPrice                int        `json:"total_price" db:"total_price"`
	VenueId                   string     `json:"venue_id" db:"venue_id"`
	PreorderTime              *time.Time `json:"preorder_time" db:"preorder_time"`
}

func UnixOrNil(i int64) *time.Time {
	if i == 0 {
		return nil
	}

	t := time.UnixMilli(i)

	return &t
}

func (fo *FullOrder) ToSimpleOrder() SimpleOrder {
	var dCoordX float64
	var dCoordY float64

	if fo.DeliveryMethod == "takeaway" {
		dCoordX = 0
		dCoordY = 0
	} else {
		dCoordX = fo.DeliveryLocation.Coordinates.Coordinates[0]
		dCoordY = fo.DeliveryLocation.Coordinates.Coordinates[1]
	}

	return SimpleOrder{
		ClientPreEstimate:         fo.ClientPreEstimate,
		DeliveryStreet:            fo.DeliveryLocation.Street,
		DeliveryCoordinateX:       dCoordX,
		DeliveryCoordinateY:       dCoordY,
		DeliveryDistance:          fo.DeliveryDistance,
		DeliveryEta:               UnixOrNil(fo.DeliveryEta.Date),
		DeliveryMethod:            fo.DeliveryMethod,
		DeliveryPrice:             fo.DeliveryPrice,
		DeliverySizeSurcharge:     fo.DeliverySizeSurcharge,
		DeliveryDistanceSurcharge: fo.DeliveryDistanceSurcharge,
		DeliveryTime:              UnixOrNil(fo.DeliveryTime.Date),
		DriverType:                fo.DriverType,
		ItemsPrice:                fo.ItemsPrice,
		OrderId:                   fo.OrderId,
		PaymentAmount:             fo.PaymentAmount,
		PaymentTime:               UnixOrNil(fo.PaymentTime.Date),
		Status:                    fo.Status,
		ServiceFee:                fo.ServiceFee,
		Subscribed:                fo.Subscribed,
		TotalPrice:                fo.TotalPrice,
		VenueId:                   fo.VenueId,
		PreorderTime:              UnixOrNil(fo.PreorderTime.Date),
	}
}

func (fo *FullOrder) ToSimpleVenue() SimpleVenue {
	return SimpleVenue{
		VenueId:          fo.VenueId,
		VenueName:        fo.VenueName,
		VenueProductLine: fo.VenueProductLine,
		VenueCoordinateX: fo.VenueCoordinates[0],
		VenueCoordinateY: fo.VenueCoordinates[1],
		VenueUrl:         fo.VenueUrl,
	}
}

type SimpleVenue struct {
	VenueId          string  `json:"venue_id" db:"venue_id"`
	VenueName        string  `json:"venue_name" db:"venue_name"`
	VenueProductLine string  `json:"venue_product_line" db:"venue_product_line"`
	VenueCoordinateX float64 `json:"venue_coordinate_x" db:"venue_coordinate_x"`
	VenueCoordinateY float64 `json:"venue_coordinate_y" db:"venue_coordinate_y"`
	VenueUrl         string  `json:"venue_url" db:"venue_url"`
}

type FullOrder struct {
	AutomaticRejectionTime struct {
		Date int64 `json:"$date"`
	} `json:"automatic_rejection_time"`
	ClientPreEstimate string `json:"client_pre_estimate"`
	Credits           int    `json:"credits"`
	Currency          string `json:"currency"`
	DeliveryBasePrice int    `json:"delivery_base_price"`
	DeliveryComment   string `json:"delivery_comment"`
	DeliveryDistance  int    `json:"delivery_distance,omitempty"`
	DeliveryEta       struct {
		Date int64 `json:"$date"`
	} `json:"delivery_eta"`
	DeliveryLocation struct {
		Address     string `json:"address"`
		Alias       string `json:"alias"`
		Apartment   string `json:"apartment,omitempty"`
		City        string `json:"city"`
		Coordinates struct {
			Coordinates []float64 `json:"coordinates"`
			Type        string    `json:"type"`
		} `json:"coordinates"`
		Street string `json:"street"`
	} `json:"delivery_location,omitempty"`
	DeliveryMethod        string `json:"delivery_method"`
	DeliveryPrice         int    `json:"delivery_price"`
	DeliveryPriceShare    int    `json:"delivery_price_share"`
	DeliverySizeSurcharge int    `json:"delivery_size_surcharge,omitempty"`
	DeliveryTime          struct {
		Date int64 `json:"$date"`
	} `json:"delivery_time"`
	DriverType      string        `json:"driver_type"`
	IsHostPaying    bool          `json:"is_host_paying"`
	IsMarketplaceV2 bool          `json:"is_marketplace_v2"`
	ItemChangeLog   []interface{} `json:"item_change_log"`
	Items           []struct {
		Count     int    `json:"count"`
		EndAmount int    `json:"end_amount"`
		Id        string `json:"id"`
		Name      string `json:"name"`
		Options   []struct {
			Id     string `json:"id"`
			Name   string `json:"name"`
			Type   string `json:"type"`
			Values []struct {
				Count int    `json:"count"`
				Id    string `json:"id"`
				Name  string `json:"name"`
				Price int    `json:"price"`
			} `json:"values"`
		} `json:"options"`
		Price                int  `json:"price"`
		RowNumber            int  `json:"row_number"`
		SkipOnRefill         bool `json:"skip_on_refill"`
		SubstitutionSettings struct {
			AllowedItems []interface{} `json:"allowed_items"`
			IsAllowed    bool          `json:"is_allowed"`
		} `json:"substitution_settings,omitempty"`
	} `json:"items"`
	ItemsPrice          int           `json:"items_price"`
	ListImage           string        `json:"list_image"`
	ListImageBlurhash   string        `json:"list_image_blurhash"`
	MainImage           string        `json:"main_image"`
	MainImageBlurhash   string        `json:"main_image_blurhash"`
	OrderAdjustmentRows []interface{} `json:"order_adjustment_rows"`
	OrderId             string        `json:"order_id"`
	OrderNumber         string        `json:"order_number"`
	PaymentAmount       int           `json:"payment_amount"`
	PaymentMethod       struct {
		Id       string `json:"id"`
		Provider string `json:"provider"`
		Type     string `json:"type"`
	} `json:"payment_method"`
	PaymentName string `json:"payment_name"`
	PaymentTime struct {
		Date int64 `json:"$date"`
	} `json:"payment_time"`
	ServiceFee                int       `json:"service_fee,omitempty"`
	Status                    string    `json:"status"`
	Subscribed                bool      `json:"subscribed"`
	Subtotal                  int       `json:"subtotal"`
	Tip                       int       `json:"tip"`
	TipShare                  int       `json:"tip_share"`
	Tokens                    int       `json:"tokens"`
	TotalPrice                int       `json:"total_price"`
	TotalPriceShare           int       `json:"total_price_share"`
	VenueAddress              string    `json:"venue_address"`
	VenueCoordinates          []float64 `json:"venue_coordinates"`
	VenueCountry              string    `json:"venue_country"`
	VenueFullAddress          string    `json:"venue_full_address"`
	VenueId                   string    `json:"venue_id"`
	VenueName                 string    `json:"venue_name"`
	VenueOpen                 bool      `json:"venue_open"`
	VenueOpenOnPurchase       bool      `json:"venue_open_on_purchase"`
	VenuePhone                string    `json:"venue_phone"`
	VenueProductLine          string    `json:"venue_product_line"`
	VenueTimezone             string    `json:"venue_timezone"`
	VenueUrl                  string    `json:"venue_url"`
	DeliveryDistanceSurcharge int       `json:"delivery_distance_surcharge,omitempty"`
	Comment                   string    `json:"comment,omitempty"`
	PreorderStatus            string    `json:"preorder_status,omitempty"`
	PreorderTime              struct {
		Date int64 `json:"$date"`
	} `json:"preorder_time,omitempty"`
	CancellableStatus struct {
		Reason    string `json:"reason"`
		ShowTimer bool   `json:"show_timer"`
		Start     struct {
			Date int64 `json:"$date"`
		} `json:"start"`
		Until struct {
			Date int64 `json:"$date"`
		} `json:"until"`
	} `json:"cancellable_status,omitempty"`
}
