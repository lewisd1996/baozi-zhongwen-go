package testutils

import (
	"context"
	"database/sql"
	"fmt"
	"testing"

	embeds "github.com/lewisd1996/baozi-zhongwen"
	_ "github.com/lib/pq" // Import the PostgreSQL driver
	"github.com/pressly/goose/v3"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
)

const (
	TestUserId    = "0759efb4-9d6c-48b6-bead-0c55f2920e04"
	TestUserEmail = "testuser@example.com"
	TestDeckId    = "1432792e-be69-4ad2-a9f5-7350d103a80b"
)

// InitializeTestDB starts a PostgreSQL container and returns a *sql.DB connection to it, along with a cleanup function.
func InitializeTestDB(ctx context.Context, t *testing.T) (*sql.DB, func()) {
	req := testcontainers.ContainerRequest{
		Image:        "postgres:latest",
		ExposedPorts: []string{"5432/tcp"},
		Env: map[string]string{
			"POSTGRES_USER":     "postgres",
			"POSTGRES_PASSWORD": "postgres",
			"POSTGRES_DB":       "postgres",
		},
		WaitingFor: wait.ForListeningPort("5432/tcp"),
	}

	postgresContainer, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: req,
		Started:          true,
	})
	if err != nil {
		t.Fatalf("Failed to start container: %v", err)
	}

	// Construct DSN for the test database
	containerPort, err := postgresContainer.MappedPort(ctx, "5432")
	if err != nil {
		t.Fatalf("Failed to get container port: %v", err)
	}
	host, err := postgresContainer.Host(ctx)
	if err != nil {
		t.Fatalf("Failed to get container host: %v", err)
	}
	dsn := fmt.Sprintf("host=%s port=%s user=postgres password=postgres dbname=postgres sslmode=disable", host, containerPort.Port())

	// Connect to the database
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		t.Fatalf("Failed to connect to the database: %v", err)
	}

	// Migration
	if err := goose.SetDialect("postgres"); err != nil {
		panic(err)
	}
	if err != nil {
		t.Fatal(err)
	}

	goose.SetBaseFS(embeds.Migrations)
	if err := goose.Up(db, "sql/migrations"); err != nil {
		t.Fatal(err)
	}

	t.Log("migration done")

	// Seed the database
	if err := seedDatabase(db); err != nil {
		t.Fatal(err)
	}

	// Return the DB connection and a cleanup function
	return db, func() {
		db.Close()
		postgresContainer.Terminate(ctx)
	}
}

// seedDatabase populates the test database with sample data
func seedDatabase(db *sql.DB) error {
	// Insert a sample user
	_, err := db.Exec(fmt.Sprintf(`
        INSERT INTO "user" (id, email) 
        VALUES ('%s', '%s');
    `, TestUserId, TestUserEmail))
	if err != nil {
		return fmt.Errorf("failed to insert into user: %w", err)
	}

	// Insert a sample deck
	_, err = db.Exec(fmt.Sprintf(`
        INSERT INTO deck (id, name, description, owner_id)
        VALUES ('%s', 'Sample Deck', 'This is a sample deck description.', (SELECT id FROM "user" WHERE email = '%s'));
    `, TestDeckId, TestUserEmail))
	if err != nil {
		return fmt.Errorf("failed to insert into deck: %w", err)
	}

	// Insert sample cards
	_, err = db.Exec(`
        INSERT INTO card (deck_id, content, translation) 
        VALUES ((SELECT id FROM deck WHERE name = 'Sample Deck'), 'Content of card 1', 'Translation of card 1'),
               ((SELECT id FROM deck WHERE name = 'Sample Deck'), 'Content of card 2', 'Translation of card 2');
    `)
	if err != nil {
		return fmt.Errorf("failed to insert into card: %w", err)
	}

	// Insert a sample learning session
	_, err = db.Exec(`
        INSERT INTO learning_session (user_id, deck_id) 
        VALUES ((SELECT id FROM "user" WHERE email = 'testuser@example.com'), (SELECT id FROM deck WHERE name = 'Sample Deck'));
    `)
	if err != nil {
		return fmt.Errorf("failed to insert into learning_session: %w", err)
	}

	// Insert sample card learning progress
	// This assumes you have at least one card in the `card` table
	_, err = db.Exec(`
        INSERT INTO card_learning_progress (session_id, card_id, user_id, review_count, success_count, last_reviewed_at) 
        SELECT (SELECT id FROM learning_session WHERE user_id = (SELECT id FROM "user" WHERE email = 'testuser@example.com')),
               id,
               (SELECT id FROM "user" WHERE email = 'testuser@example.com'),
               0, 0, CURRENT_TIMESTAMP
        FROM card
        WHERE deck_id = (SELECT id FROM deck WHERE name = 'Sample Deck');
    `)
	if err != nil {
		return fmt.Errorf("failed to insert into card_learning_progress: %w", err)
	}

	return nil
}
