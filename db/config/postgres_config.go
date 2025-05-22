package db

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"os"

	"github.com/Adarsh-Kmt/EncryptedChat/db/sqlc"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

var (
	Client     *sqlc.Queries
	PostgresDB *sql.DB
	connPool   *pgxpool.Pool
	logger     = log.New(os.Stdout, "GLINT DATABASE CONFIG >> ", 0)
)

type postgresConfig struct {
	host        string
	port        string
	username    string
	password    string
	database    string
	sslMode     string
	sslRootCert string
}

func postgresConfiguration() (*postgresConfig, error) {

	config := &postgresConfig{}

	//err := godotenv.Load("/prod/.env") // Specify the filename if it's named differently

	// EnvFilePath := os.Getenv("ENV_FILE_PATH")

	// if EnvFilePath == "" {
	// 	return nil, fmt.Errorf("ENV_FILE_PATH not configured")
	// }
	err := godotenv.Load(".env")
	if err != nil {
		logger.Fatalf("Error loading .env file: %v", err)
	}
	config.password = os.Getenv("DB_PASSWORD")
	config.username = os.Getenv("DB_USERNAME")

	config.port = os.Getenv("DB_PORT")

	config.host = os.Getenv("DB_HOST")

	config.database = os.Getenv("DB_DATABASE")

	config.sslMode = os.Getenv("DB_SSLMODE")

	logger.Printf("password : %s", config.password)
	logger.Printf("username : %s", config.username)
	logger.Printf("port : %s", config.port)
	logger.Printf("host : %s", config.host)
	logger.Printf("database : %s", config.database)
	logger.Printf("sslmode : %s", config.sslMode)
	return config, nil

}

func PostgresDBClientInit() error {

	config, err := postgresConfiguration()

	if err != nil {
		return err
	}

	postgresConnStringFormat := "postgresql://%s:%s@%s:%s/%s?sslmode=%s"
	//testPostgresConnString := "postgresql://postgres:password@localhost:8087/glint?sslmode=disable"
	connString := fmt.Sprintf(postgresConnStringFormat, config.username, config.password, config.host, config.port, config.database, config.sslMode)
	ctx := context.Background()
	connPool, err = pgxpool.New(ctx, connString)

	if err != nil {
		return err
	}

	Client = sqlc.New(connPool)
	if PostgresDB, err = sql.Open("postgres", connString); err != nil {
		return err
	}
	return nil
}
