package storage

import (
	"frederikhs/wolt/wolt"
	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"
)

func Connect() (*sqlx.DB, error) {
	db, err := sqlx.Connect("sqlite3", "wolt.db")
	if err != nil {
		return nil, err
	}

	return db, nil
}

func Setup(db *sqlx.DB) {
	db.MustExec("DROP VIEW IF EXISTS view_wolt_order")
	db.MustExec("DROP TABLE IF EXISTS wolt_order")
	db.MustExec("DROP TABLE IF EXISTS wolt_venue")
	db.MustExec(`
		CREATE TABLE wolt_venue (
		    venue_id TEXT PRIMARY KEY,
			venue_name TEXT,
			venue_product_line TEXT,
			venue_coordinate_x TEXT,
			venue_coordinate_y TEXT,
			venue_url TEXT
		)
	`)
	db.MustExec(`
		CREATE TABLE wolt_order (
		    order_id TEXT PRIMARY KEY,
		    client_pre_estimate TEXT NOT NULL,
			delivery_street TEXT,
			delivery_coordinate_x TEXT,
			delivery_coordinate_y TEXT,
			delivery_distance INT,
			delivery_eta TEXT,
			delivery_method TEXT,
			delivery_price INT,
			delivery_size_surcharge INT,
			delivery_time TEXT,
			driver_type TEXT,
			items_price INT,
			payment_amount INT,
			payment_time TEXT,
			status TEXT,
			service_fee INT,
			subscribed BOOLEAN,
			total_price INT,
			venue_id TEXT REFERENCES wolt_venue(venue_id),
			preorder_time TEXT,
			delivery_distance_surcharge INT
		)
	`)
	db.MustExec(`
		CREATE VIEW view_wolt_order AS
			SELECT * FROM wolt_order
			JOIN wolt_venue wv on wolt_order.venue_id = wv.venue_id
	`)
}

func SaveOrders(db *sqlx.DB, orders *[]wolt.FullOrder) error {
	var simpleOrders []wolt.SimpleOrder
	var simpleVenues []wolt.SimpleVenue
	for _, o := range *orders {
		simpleOrders = append(simpleOrders, o.ToSimpleOrder())
		simpleVenues = append(simpleVenues, o.ToSimpleVenue())
	}

	_, err := db.NamedExec(`
		INSERT OR IGNORE INTO wolt_venue (
			venue_id,
			venue_name,
			venue_product_line,
			venue_coordinate_x,
			venue_coordinate_y,
			venue_url
		) VALUES (
		    :venue_id,
			:venue_name,
			:venue_product_line,
			:venue_coordinate_x,
			:venue_coordinate_y,
			:venue_url
		)
	`, simpleVenues)
	_, err = db.NamedExec(`
		INSERT INTO wolt_order (
			order_id, 
			client_pre_estimate, 
			delivery_street, 
			delivery_coordinate_x, 
			delivery_coordinate_y, 
			delivery_distance, 
			delivery_eta, 
			delivery_method, 
			delivery_price, 
			delivery_size_surcharge, 
			delivery_time, 
			driver_type, 
			items_price, 
			payment_amount, 
			payment_time, 
			status, 
			service_fee, 
			subscribed, 
			total_price,
			venue_id,
			preorder_time, 
			delivery_distance_surcharge
		) VALUES (
		    :order_id,
			:client_pre_estimate,
			:delivery_street,
			:delivery_coordinate_x,
			:delivery_coordinate_y,
			:delivery_distance,
			:delivery_eta,
			:delivery_method,
			:delivery_price,
			:delivery_size_surcharge,
			:delivery_time,
			:driver_type,
			:items_price,
			:payment_amount,
			:payment_time,
			:status,
			:service_fee,
			:subscribed,
			:total_price,
		    :venue_id,
			:preorder_time,
		    :delivery_distance_surcharge
		)
	`, simpleOrders)

	if err != nil {
		return err
	}

	return nil
}

type VenueAgg struct {
	VenueName  string `db:"venue_name"`
	VenueValue int    `db:"venue_value"`
}

func GetTopVenuesByTotalSpend(db *sqlx.DB) (*[]VenueAgg, error) {
	sql := `
		SELECT DISTINCT vwo.venue_name, agg.total_spend / 100 as venue_value FROM view_wolt_order vwo
		JOIN (SELECT venue_id, SUM(payment_amount) as total_spend
			  FROM view_wolt_order
			  WHERE status = 'delivered'
			  GROUP BY venue_id) agg ON agg.venue_id = vwo.venue_id
		ORDER BY agg.total_spend
	`

	var rows []VenueAgg
	err := db.Select(&rows, sql)
	if err != nil {
		return nil, err
	}

	return &rows, nil
}

func GetTopVenuesByTotalNumberOfOrders(db *sqlx.DB) (*[]VenueAgg, error) {
	sql := `
		SELECT DISTINCT vwo.venue_name, agg.count as venue_value FROM view_wolt_order vwo
		JOIN (SELECT venue_id, COUNT(*) as count
			  FROM view_wolt_order
			  WHERE status = 'delivered'
			  GROUP BY venue_id) agg ON agg.venue_id = vwo.venue_id
		ORDER BY agg.count	
	`

	var rows []VenueAgg
	err := db.Select(&rows, sql)
	if err != nil {
		return nil, err
	}

	return &rows, nil
}

func GetTopVenuesByTotalSpendOnDelivery(db *sqlx.DB) (*[]VenueAgg, error) {
	sql := `
		SELECT DISTINCT vwo.venue_name, agg.count / 100 as venue_value FROM view_wolt_order vwo
		JOIN (SELECT venue_id, SUM(delivery_price) as count
			FROM view_wolt_order
			WHERE status = 'delivered'
			GROUP BY venue_id) agg ON agg.venue_id = vwo.venue_id
		ORDER BY agg.count
	`

	var rows []VenueAgg
	err := db.Select(&rows, sql)
	if err != nil {
		return nil, err
	}

	return &rows, nil
}

type TotalFoodAndDeliverySpend struct {
	TotalFood     int `db:"sum_food"`
	TotalDelivery int `db:"sum_delivery"`
}

func GetTotalFoodAndDeliverySpend(db *sqlx.DB) (*TotalFoodAndDeliverySpend, error) {
	sql := `
		SELECT (SUM(payment_amount) - SUM(delivery_price)) / 100 as sum_food,
			    SUM(delivery_price) / 100 as sum_delivery
		FROM view_wolt_order vwo
		WHERE status = 'delivered'
	`

	var rows []TotalFoodAndDeliverySpend
	err := db.Select(&rows, sql)
	if err != nil {
		return nil, err
	}

	return &rows[0], nil
}

type OrderDay struct {
	Date  string `db:"yearweek"`
	Count int    `db:"count"`
}

func GetNumberOfOrdersByDateRange(db *sqlx.DB) (*[]OrderDay, error) {
	sql := `
		SELECT a.yearweek, coalesce(b.count, 0) as count
		FROM (WITH RECURSIVE cnt(x) AS (SELECT 0
										UNION ALL
										SELECT x + 1
										FROM cnt
										LIMIT (SELECT ((julianday(date()) - julianday('2019-10-10'))) + 1))
			  SELECT DISTINCT strftime('%Y%W', julianday('2019-10-10'), '+' || x || ' days') as yearweek
			  FROM cnt) a
				 LEFT JOIN (SELECT strftime('%Y%W', payment_time) as yearweek, COUNT(*) as count
					   FROM view_wolt_order
					   WHERE status = 'delivered'
					   GROUP BY 1) b ON a.yearweek = b.yearweek
	`

	var rows []OrderDay
	err := db.Select(&rows, sql)
	if err != nil {
		return nil, err
	}

	return &rows, nil
}
