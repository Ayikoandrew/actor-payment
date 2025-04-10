package storage

import (
	"database/sql"
	"fmt"
	"log/slog"
	"os"
	"time"

	"github.com/Ayikoandrew/ap/msg/msg"
	_ "github.com/lib/pq"

	"github.com/anthdm/hollywood/actor"
)

type DatabaseProcessor struct {
	db *sql.DB
}

func NewDatabaseProcessor() actor.Producer {
	return func() actor.Receiver {
		return &DatabaseProcessor{}
	}
}

func (dp *DatabaseProcessor) Receive(ctx *actor.Context) {
	switch m := ctx.Message().(type) {
	case actor.Started:
		dp.OnStart(ctx)
	case actor.Initialized:
		dp.OnInit(ctx)
	case actor.Stopped:
		dp.OnStop(ctx)
	case *msg.PublicPayment:
		amount := m.Amount
		processed_amount := dp.handleCalculatePayment(float32(amount))
		processed_payment := ProcessedPayment{
			InitialPayment:   float64(amount),
			ProcessedPayment: float64(processed_amount),
			Status:           "public",
		}
		dp.savePaymentToDatabase(&processed_payment)
	case *msg.PrivatePayment:
		amount := m.Amount
		processed_amount := dp.handleCalculatePayment(float32(amount))
		processed_payment := ProcessedPayment{
			InitialPayment:   float64(amount),
			ProcessedPayment: float64(processed_amount),
			Status:           "private",
		}
		dp.savePaymentToDatabase(&processed_payment)
	default:
		fmt.Printf("[Database processor] Unknown stream %+v\n", m)
	}
}

func (dp *DatabaseProcessor) OnInit(ctx *actor.Context) {
	connStr := os.Getenv("DATABASE_URL")
	if connStr == "" {

		host := os.Getenv("DBHOST")
		if host == "" {
			host = "actor-payment"
		}
		port := os.Getenv("DBPORT")
		if port == "" {
			port = os.Getenv("DB_PORT")
		}
		name := os.Getenv("DBNAME")
		if name == "" {
			name = os.Getenv("POSTGRES_DB")
		}
		password := os.Getenv("DBPASSWORD")
		if password == "" {
			password = os.Getenv("POSTGRES_PASSWORD")
		}
		user := os.Getenv("DBUSER")
		if user == "" {
			user = os.Getenv("POSTGRES_USER")
		}
		sslmode := os.Getenv("SSLMODE")
		if sslmode == "" {
			sslmode = "disable"
		}

		// Log found/missing values for debugging
		fmt.Printf("Host %+v\n", host)
		fmt.Printf("Port %+v\n", port)

		if host == "" || port == "" || name == "" || password == "" || user == "" {
			fmt.Println("Missing configurations")
			return
		}

		connStr = fmt.Sprintf("host=%s port=%v dbname=%s password=%s user=%s sslmode=%s",
			host, port, name, password, user, sslmode)
	}

	var err error

	d, err := sql.Open("postgres", connStr)
	if err != nil {
		fmt.Printf("[Database Processor] error opening the database %+v\n", err)
	}
	dp.db = d

	dp.db.SetMaxOpenConns(25)
	dp.db.SetConnMaxIdleTime(10)
	dp.db.SetConnMaxLifetime(5 * time.Second)
	dp.db.SetMaxIdleConns(10)

	slog.Info("[Database Processor] testing database connection...")
	if err = dp.db.Ping(); err != nil {
		dp.db.Close()
		fmt.Printf("[Database Processor] error pinging the database %+v\n", err)
	}

	slog.Info("[Database Processor] successfully connected to the database")

	query := `CREATE TABLE IF NOT EXISTS payments(
		paymentID serial primary key,
		initial_payment DECIMAL(10,2),

		processed_payment DECIMAL(10,2),
		status varchar(255),
		timestamp TIMESTAMP DEFAULT CURRENT_TIMESTAMP
	)`
	if _, err = dp.db.Exec(query); err != nil {
		dp.db.Close()
		fmt.Printf("[Database Processor] error creating payments table: %+v\n", err)

	}
}

func (dp *DatabaseProcessor) OnStop(ctx *actor.Context) {
	fmt.Println("[Database Processor] closing database connection")
	if dp.db != nil {
		if err := dp.db.Close(); err != nil {
			fmt.Printf("[Database Processor] error closing database connection: %+v", err)
		}
	} else {
		fmt.Println("[Database Processor] no database connection to close")
	}
	fmt.Println("[Database Processor] initialization complete")
}

func (dp *DatabaseProcessor) OnStart(ctx *actor.Context) {
	fmt.Println("[Database processor] started")
}

func (dp *DatabaseProcessor) handleCalculatePayment(a float32) float32 {
	amount := a * 0.02
	return a - amount
}

func (dp *DatabaseProcessor) savePaymentToDatabase(pp *ProcessedPayment) {
	query := `INSERT INTO payments (initial_payment, processed_payment, status) VALUES ($1, $2, $3)`
	if _, err := dp.db.Exec(query,
		&pp.InitialPayment,
		&pp.ProcessedPayment,
		&pp.Status,
	); err != nil {
		fmt.Printf("[Database processor] failed to insert into the table %+v", err)
	}
	fmt.Println("[Database processor] successfully saved data in the database.")
}
