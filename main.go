package main

import (
	"database/sql"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	_ "github.com/lib/pq"
)

type Activity struct {
	ID        int    `json:"id"`
	Title     string `json:"title" validate:"required,min=3,max=100"`
	Category  string `json:"category" validate:"required,min=3,max=100,oneof=TASK EVENT"`
	Description string `json:"description" validate:"required,min=3,max=100"`
	ActivityDate time.Time `json:"activity_date" validate:"required"`
	Status string `json:"status"`
	CreatedAt time.Time `json:"created_at"`
}

func initDB() (*sql.DB, error) {
	dns := "user=postgres.phrxssottrrxjphlltfs password=ZBPzCwb!u6!T!pU host=aws-0-ap-southeast-1.pooler.supabase.com port=6543 dbname=postgres"
	db, err := sql.Open("postgres", dns)
	if err != nil {
		return nil, err
	}

	err = db.Ping()
	if err != nil {
		db.Close()
		return nil, err
	}
	return db, nil
}
func main() { 
	db, err := initDB()
	if err != nil {
		panic(err)
	}
	defer db.Close()

	app := fiber.New()
	validate := validator.New()

	app.Get("/activites", func(c *fiber.Ctx) error {	
		rows, err := db.Query("SELECT * FROM activities")
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"status": "error", "message": err.Error(), "data": nil})//c.Status(500).SendString(err.Error())
		}
		defer rows.Close()
		var activities []Activity
		for rows.Next() {
			var activity Activity
			if err := rows.Scan(&activity.ID, &activity.Title, &activity.Category, &activity.Description, &activity.ActivityDate, &activity.Status, &activity.CreatedAt); err != nil {
				return c.Status(500).SendString(err.Error())
			}
			activities = append(activities, activity)
		}
		return c.JSON(activities)
	})
	app.Post("/activites", func(c *fiber.Ctx) error {
		var activity Activity
		if err := c.BodyParser(&activity); err != nil {
			return c.Status(500).SendString(err.Error())
		}
		
		if err = validate.Struct(&activity); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": "error", "message": err.Error(),})//c.Status(500).SendString(err.Error())
		}
		sqlStatement := `INSERT INTO activities (title, category, description, activity_date, status, created_at) VALUES ($1, $2, $3, $4, $5, $6) RETURNING id`
		err = db.QueryRow(sqlStatement, activity.Title, activity.Category, activity.Description, activity.ActivityDate, activity.Status, activity.CreatedAt).Scan(&activity.ID)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"status": "error", "message": err.Error(), "data": nil})//c.Status(500).SendString(err.Error())
		}
		return c.Status(fiber.StatusCreated).JSON(fiber.Map{"status": "success", "message": "Activity created successfully", "data": activity})
	})
	app.Put("/activites/:id", func(c *fiber.Ctx) error {
		id := c.Params("id")
		var activity Activity
		if err := c.BodyParser(&activity); err != nil {
			return c.Status(500).SendString(err.Error())
		}
		sqlStatement := `UPDATE activities SET title = $1, category = $2, description = $3, activity_date = $4, status = $5, created_at = $6 WHERE id = $7 RETURNING id`
		err = db.QueryRow(sqlStatement, activity.Title, activity.Category, activity.Description, activity.ActivityDate, activity.Status, activity.CreatedAt, id).Scan(&activity.ID)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"status": "error", "message": err.Error(), "data": nil})//c.Status(500).SendString(err.Error())
		}
		return c.Status(fiber.StatusOK).JSON(fiber.Map{"status": "success", "message": "Activity updated successfully", "data": activity})
	})

	app.Delete("/activites/:id", func(c *fiber.Ctx) error {
		id := c.Params("id")
		sqlStatement := `DELETE FROM activities WHERE id = $1`
		_, err := db.Exec(sqlStatement, id)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"status": "error", "message": err.Error(), "data": nil})//c.Status(500).SendString(err.Error())
		}
		return c.Status(fiber.StatusOK).JSON(fiber.Map{"status": "success", "message": "Activity deleted successfully", "data": nil})
	})





	app.Listen(":8084")
}