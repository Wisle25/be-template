package infrastructures

import (
	"context"
	"database/sql"
	"fmt"
	"reflect"

	_ "github.com/lib/pq"
	"github.com/redis/go-redis/v9"
	"github.com/wisle25/be-template/commons"
)

var (
	RedisClient *redis.Client
	ctx         context.Context
	DB          *sql.DB
)

// ConnectRedis initializes a connection to the Redis server using the provided configuration.
func ConnectRedis(config *commons.Config) {
	ctx = context.TODO()

	RedisClient = redis.NewClient(&redis.Options{
		Addr: config.RedisURL,
	})

	// Ping the Redis server to ensure the connection is established.
	if _, err := RedisClient.Ping(ctx).Result(); err != nil {
		panic(err)
	}

	fmt.Println("Successfully connected to Redis Client")
}

// ConnectDB initializes a connection to the PostgreSQL database using the provided configuration.
func ConnectDB(config *commons.Config) {
	var err error
	dbName := config.DBName

	// Use the test database if the application environment is set to 'dev'.
	if config.AppEnv == "dev" {
		dbName = config.DBNameTest
	}

	// Format the data source name (DSN) string for connecting to the database.
	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=Asia/Shanghai",
		config.DBHost,
		config.DBUserName,
		config.DBUserPassword,
		dbName,
		config.DBPort,
	)

	// Open a connection to the database.
	DB, err = sql.Open("postgres", dsn)
	if err != nil {
		panic(err)
	}

	// Ping the database to ensure the connection is established.
	if err = DB.Ping(); err != nil {
		panic(err)
	}

	fmt.Println("Successfully connected to Postgres!")
}

// GetTableDB parses SQL query results into a slice of the specified type T.
func GetTableDB[T any](rows *sql.Rows) []T {
	var table []T

	// Iterate through the rows returned by the query.
	for rows.Next() {
		var data T

		s := reflect.ValueOf(&data).Elem()
		numCols := s.NumField()
		columns := make([]interface{}, numCols)

		// Prepare the columns slice with pointers to the fields of the struct.
		for i := 0; i < numCols; i++ {
			field := s.Field(i)
			columns[i] = field.Addr().Interface()
		}

		// Scan the row into the struct fields.
		if err := rows.Scan(columns...); err != nil {
			fmt.Println("Case Read Error ", err)
		}

		// Append the struct to the table slice.
		table = append(table, data)
	}

	return table
}
