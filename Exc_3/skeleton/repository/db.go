package repository

import (
	"errors"
	"fmt"
	"log/slog"
	"ordersystem/model"
	"os"
	"strings"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type DatabaseHandler struct {
	dbConn *gorm.DB
}

func NewDatabaseHandler() (*DatabaseHandler, error) {
	slog.Info("Connecting to database")
	// connect to db
	dsn, err := getDsn()
	if err != nil {
		return nil, err
	}

	dbConn, err := gorm.Open(postgres.New(postgres.Config{DSN: dsn}), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	// create tables and migrate
	if err := dbConn.AutoMigrate(&model.Drink{}, &model.Order{}); err != nil {
		return nil, err
	}

	// add test data to db
	if err := prepopulate(dbConn); err != nil {
		return nil, err
	}

	return &DatabaseHandler{dbConn: dbConn}, nil
}

func getDsn() (string, error) {
	dbUser, ok := os.LookupEnv("POSTGRES_USER")
	if !ok {
		return "", errors.New("environment variable 'POSTGRES_USER' is not set")
	}
	dbPw, ok := os.LookupEnv("POSTGRES_PASSWORD")
	if !ok {
		return "", errors.New("environment variable 'POSTGRES_PASSWORD' is not set")
	}
	dbName, ok := os.LookupEnv("POSTGRES_DB")
	if !ok {
		return "", errors.New("environment variable 'POSTGRES_DB' is not set")
	}
	dbPort, ok := os.LookupEnv("POSTGRES_TCP_PORT")
	if !ok {
		return "", errors.New("environment variable 'POSTGRES_TCP_PORT' is not set")
	}
	dbHost, ok := os.LookupEnv("DB_HOST")
	if !ok {
		return "", errors.New("environment variable 'DB_HOST' is not set")
	}

	// Sanitize: trim spaces and surrounding quotes that may come from .env files
	clean := func(s string) string {
		s = strings.TrimSpace(s)
		s = strings.Trim(s, `"'`)
		return s
	}
	dbUser = clean(dbUser)
	dbPw = clean(dbPw)
	dbName = clean(dbName)
	dbPort = clean(dbPort)
	dbHost = clean(dbHost)

	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=disable",
		dbHost, dbUser, dbPw, dbName, dbPort,
	)

	// ---- TEMP DEBUG LOGS (remove once it works) ----
	slog.Info("DB ENV CHECK",
		"user", dbUser,
		"password_length", len(dbPw),
		"db", dbName,
		"host", dbHost,
		"port", dbPort,
	)
	slog.Info("RAW PASSWORD (quoted)", "pw", fmt.Sprintf("%q", dbPw))
	slog.Info("FINAL DSN", "dsn", dsn)
	// ------------------------------------------------

	return dsn, nil
}

func prepopulate(dbConn *gorm.DB) error {
	// check if prepopulate has already run once
	var exists bool
	err := dbConn.Model(&model.Drink{}).
		Select("count(*) > 0").
		Find(&exists).
		Error
	if err != nil {
		return err
	}
	if exists {
		// don't prepopulate if has already run
		return nil
	}

	// create drink menu
	drinks := []model.Drink{
		{Name: "Espresso", Price: 2.20},
		{Name: "Latte", Price: 3.50},
		{Name: "Cappuccino", Price: 3.20},
		{Name: "Tea", Price: 1.70},
		{Name: "Cola", Price: 1.80},
	}
	if err := dbConn.Create(&drinks).Error; err != nil {
		return err
	}

	// create some demo orders
	orders := []model.Order{
		{DrinkID: drinks[0].ID, Amount: 2}, // 2x Espresso
		{DrinkID: drinks[1].ID, Amount: 1}, // 1x Latte
		{DrinkID: drinks[2].ID, Amount: 3}, // 3x Cappuccino
		{DrinkID: drinks[4].ID, Amount: 4}, // 4x Cola
	}
	if err := dbConn.Create(&orders).Error; err != nil {
		return err
	}

	return nil
}

func (db *DatabaseHandler) GetDrinks() (drinks []model.Drink, err error) {
	if err := db.dbConn.Find(&drinks).Error; err != nil {
		return nil, err
	}
	return drinks, nil
}

func (db *DatabaseHandler) GetOrders() (orders []model.Order, err error) {
	if err := db.dbConn.
		Preload("Drink"). // 👈 load the nested drink object
		Find(&orders).Error; err != nil {
		return nil, err
	}
	return orders, nil
}

const totalledStmt = `SELECT drink_id, SUM(amount) AS total_amount_ordered FROM orders WHERE deleted_at IS NULL GROUP BY drink_id ORDER BY drink_id;`

func (db *DatabaseHandler) GetTotalledOrders() (totals []model.DrinkOrderTotal, err error) {
	if err := db.dbConn.Raw(totalledStmt).Scan(&totals).Error; err != nil {
		return nil, err
	}
	return totals, nil
}

func (db *DatabaseHandler) AddOrder(order *model.Order) error {
	if err := db.dbConn.Create(order).Error; err != nil {
		return err
	}
	return nil
}
