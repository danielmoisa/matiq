package main

import (
	"flag"
	"log"

	"github.com/danielmoisa/matiq/internal/config"
	"github.com/danielmoisa/matiq/internal/driver/postgres"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

func main() {
	// Parse command line flags
	var seedType = flag.String("type", "workflows", "Type of data to seed (workflows, users, etc.)")
	flag.Parse()

	log.Println("Starting database seeder...")

	// Load configuration
	cfg := config.GetInstance()

	// Initialize logger
	logger, _ := zap.NewDevelopment()
	sugar := logger.Sugar()

	// Initialize database connection
	db, err := postgres.NewPostgresConnectionByGlobalConfig(cfg, sugar)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	log.Printf("Connected to database successfully")

	// Run appropriate seeder based on type
	switch *seedType {
	case "workflows":
		if err := seedWorkflows(db); err != nil {
			log.Fatalf("Failed to seed workflows: %v", err)
		}
		log.Println("Workflows seeded successfully!")
	default:
		log.Fatalf("Unknown seed type: %s", *seedType)
	}

	log.Println("Database seeding completed!")
}

// seedWorkflows inserts sample workflow data into the database
func seedWorkflows(db *gorm.DB) error {
	log.Println("Seeding sample workflows...")

	// Use GORM's raw SQL with proper JSON escaping
	sql := `
	INSERT INTO workflows (
		uid,
		team_id,
		workflow_id,
		version,
		resource_id,
		name,
		type,
		trigger_mode,
		transformer,
		template,
		config,
		created_at,
		created_by,
		updated_at,
		updated_by
	) VALUES (
		gen_random_uuid(),
		1,
		1001,
		0,
		1,
		'User Validation Workflow',
		6,
		'manually',
		?::jsonb,
		?::jsonb,
		?::jsonb,
		NOW(),
		'admin',
		NOW(),
		'admin'
	)`

	// Define the JSON objects as Go strings
	transformer := `{"enable": false, "rawData": ""}`

	template := `{
		"resourceID": 123,
		"runByAnonymous": true,
		"teamID": 1,
		"nodes": [
			{
				"id": "trigger-1",
				"type": "webhook",
				"triggerType": "webhook",
				"position": {"x": 100, "y": 100},
				"data": {
					"name": "Webhook Trigger",
					"description": "Receives incoming HTTP requests",
					"webhookUrl": "/api/webhooks/workflow-123",
					"method": "POST",
					"headers": {"Content-Type": "application/json"}
				}
			},
			{
				"id": "postgres-1",
				"type": "postgresql",
				"position": {"x": 300, "y": 100},
				"data": {
					"name": "Database Query",
					"description": "Query user data from PostgreSQL",
					"mode": "sql",
					"query": "SELECT * FROM users WHERE email = '{{trigger.body.email}}';",
					"resourceID": "postgres-resource-1"
				},
				"connections": ["trigger-1"]
			},
			{
				"id": "transformer-1",
				"type": "transformer",
				"position": {"x": 500, "y": 100},
				"data": {
					"name": "Data Transformer",
					"description": "Transform user data for response",
					"code": "const transformedData = { userId: data.postgres_1[0]?.id, fullName: data.postgres_1[0]?.first_name + ' ' + data.postgres_1[0]?.last_name, email: data.postgres_1[0]?.email, isActive: data.postgres_1[0]?.status === 'active' }; return transformedData;",
					"language": "javascript"
				},
				"connections": ["postgres-1"]
			}
		],
		"connections": [
			{"id": "conn-1", "sourceId": "trigger-1", "targetId": "postgres-1"},
			{"id": "conn-2", "sourceId": "postgres-1", "targetId": "transformer-1"}
		],
		"virtualResource": {"icon": "database", "category": "data-processing"},
		"metadata": {"version": "1.0.0", "description": "User validation workflow"}
	}`

	config := `{"public": false, "advancedConfig": {"runtime": "none", "pages": [], "delayWhenLoaded": "", "displayLoadingPage": false, "isPeriodically": false, "periodInterval": ""}}`

	// Execute the SQL with parameters
	result := db.Exec(sql, transformer, template, config)
	if result.Error != nil {
		return result.Error
	}

	log.Printf("Inserted %d workflow(s)", result.RowsAffected)
	return nil
}
