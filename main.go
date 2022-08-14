package main

import (
	"fmt"
	"frederikhs/wolt/storage"
	"frederikhs/wolt/wolt"
	"github.com/go-echarts/go-echarts/v2/charts"
	"github.com/go-echarts/go-echarts/v2/components"
	"github.com/go-echarts/go-echarts/v2/opts"
	"io"
	"log"
	"os"
	"time"
)

func GetOrders(client *wolt.Client) (*[]wolt.FullOrder, error) {
	if storage.JsonExists() {
		log.Println("json did exist, reusing")
		orders, err := storage.GetOrders()
		if err != nil {
			return nil, err
		}

		return orders, nil
	} else {
		log.Println("json did not exists, fetching orders")

		var orders []wolt.FullOrder

		limit := 50
		skip := 0
		emptyLast := false

		for !emptyLast {
			o, err := client.RequestOrders(limit, skip)
			if err != nil {
				return nil, err
			}

			orders = append(orders, *o...)
			if len(*o) == 0 {
				emptyLast = true
			}

			skip = skip + len(*o)
			log.Printf("requested orders, got %d back\n", len(*o))
			time.Sleep(time.Second)
		}

		err := storage.WriteOrders(&orders)
		if err != nil {
			return nil, err
		}

		return &orders, nil
	}
}

func main() {
	if len(os.Args) != 2 {
		fmt.Printf("usage: %s <TOKEN>\n", os.Args[0])
		os.Exit(1)
	}

	token := os.Args[1]
	client := wolt.NewClient(token)

	orders, err := GetOrders(client)
	if err != nil {
		panic(err)
	}

	var simpleOrders []wolt.SimpleOrder
	var simpleVenues []wolt.SimpleVenue
	for _, o := range *orders {
		simpleOrders = append(simpleOrders, o.ToSimpleOrder())
		simpleVenues = append(simpleVenues, o.ToSimpleVenue())
	}

	db, err := storage.Connect()
	if err != nil {
		panic(err)
	}

	storage.Setup(db)

	err = storage.SaveOrders(db, orders)
	if err != nil {
		panic(err)
	}

	venueSpends, err := storage.GetTopVenuesByTotalSpend(db)
	if err != nil {
		panic(err)
	}

	venueOrders, err := storage.GetTopVenuesByTotalNumberOfOrders(db)
	if err != nil {
		panic(err)
	}

	venueSpendsDelivery, err := storage.GetTopVenuesByTotalSpendOnDelivery(db)
	if err != nil {
		panic(err)
	}

	topVenueBySpendChartBar := CreateTopVenueChart("Venues by total spend", venueSpends)
	topVenueByOrdersChartBar := CreateTopVenueChart("Venues by total number of orders", venueOrders)
	topVenueBySpendDeliveryChartBar := CreateTopVenueChart("Venues by total spend on delivery", venueSpendsDelivery)

	totalSpends, err := storage.GetTotalFoodAndDeliverySpend(db)
	if err != nil {
		panic(err)
	}

	pie := charts.NewPie()
	pie.SetGlobalOptions(
		charts.WithTitleOpts(opts.Title{Title: "Total spends"}),
		charts.WithInitializationOpts(opts.Initialization{
			Width: "1500px",
		}),
	)

	pie.AddSeries("pie", []opts.PieData{
		{
			Name:  "Food sum",
			Value: totalSpends.TotalFood,
		}, {
			Name:  "Delivery sum",
			Value: totalSpends.TotalDelivery,
		},
	}).
		SetSeriesOptions(charts.WithLabelOpts(
			opts.Label{
				Show:      true,
				Formatter: "{b}: {c}",
			}),
		)

	orderDays, err := storage.GetNumberOfOrdersByDateRange(db)
	if err != nil {
		panic(err)
	}

	line := charts.NewLine()
	line.SetGlobalOptions(
		charts.WithTitleOpts(opts.Title{Title: "Orders per week"}),
		charts.WithInitializationOpts(opts.Initialization{
			Width: "1500px",
		}),
	)

	var dates []string
	for _, i := range *orderDays {
		dates = append(dates, i.Date)
	}

	orderCount := make([]opts.LineData, 0)
	for _, i := range *orderDays {
		orderCount = append(orderCount, opts.LineData{Value: i.Count})
	}

	line.SetXAxis(dates).
		AddSeries("Orders", orderCount)

	file, err := os.Create("wolt.html")
	if err != nil {
		panic(err)
	}

	page := components.NewPage()
	page.AddCharts(
		line,
		topVenueBySpendChartBar,
		topVenueByOrdersChartBar,
		topVenueBySpendDeliveryChartBar,
		pie,
	)

	err = page.Render(io.MultiWriter(file))
	if err != nil {
		panic(err)
	}

}

func CreateTopVenueChart(title string, data *[]storage.VenueAgg) *charts.Bar {
	// create a new bar instance
	bar := charts.NewBar()
	bar.SetGlobalOptions(charts.WithTitleOpts(opts.Title{
		Title: title,
	}), charts.WithInitializationOpts(opts.Initialization{
		Width:  "1500px",
		Height: "2000px",
	}))

	var names []string
	for _, i := range *data {
		names = append(names, i.VenueName)
	}

	spends := make([]opts.BarData, 0)
	for _, i := range *data {
		spends = append(spends, opts.BarData{Value: i.VenueValue})
	}

	// Put data into instance
	bar.SetXAxis(names).
		AddSeries("value", spends).
		SetSeriesOptions(
			charts.WithLabelOpts(opts.Label{
				Show:     true,
				Position: "right",
			}),
		)
	bar.XYReversal()

	return bar
}
