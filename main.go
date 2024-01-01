package main

import (
	"fmt"
	"github.com/choyalpramod/fiberpostgress/models"
	"github.com/choyalpramod/fiberpostgress/storage"
	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
	"gorm.io/gorm"
	"log"
	"net/http"
	"os"
)

type Book struct {
	Author    string `json:"author"`
	Title     string `json:"title"`
	Publisher string `json:"publisher"`
}

type Repository struct {
	DB *gorm.DB
}

func (repo *Repository) CreateBook(c *fiber.Ctx) error {
	book := new(Book)
	fmt.Println("-----------------")
	fmt.Println("book: ", book)
	if err := c.BodyParser(&book); err != nil {
		fmt.Println(err)
		fmt.Println(book)
		return c.Status(400).SendString(err.Error())
	}
	fmt.Println("book after body parser update: ", book)
	fmt.Println("-----------------")

	if dbErr := repo.DB.Create(&book).Error; dbErr != nil {
		return c.Status(500).SendString("could not create book")
	}

	return c.Status(http.StatusOK).JSON(&fiber.Map{
		"message": "book created",
	})
}

func (repo *Repository) DeleteBook(c *fiber.Ctx) error {
	id := c.Params("id")
	book := new(models.Books)
	if id == "" {
		return c.Status(http.StatusBadRequest).JSON(&fiber.Map{
			"message": "id is required",
		})
	}

	response := repo.DB.Delete(book, id)
	if response.Error != nil {
		return c.Status(http.StatusBadRequest).JSON(&fiber.Map{
			"message": "could not delete book",
		})
	}
	return c.SendString("Book successfully deleted")
}

func (repo *Repository) GetBookByID(c *fiber.Ctx) error {
	id := c.Params("id")
	book := new(models.Books)

	if id == "" {
		return c.Status(http.StatusBadRequest).SendString("id is required")
	}

	err := repo.DB.Where("id = ?", id).First(book).Error
	if err != nil {
		return c.Status(http.StatusBadRequest).SendString("invalid id")
	}
	return c.JSON(book)
}

func (repo *Repository) GetBooks(c *fiber.Ctx) error {
	var books []models.Books
	err := repo.DB.Find(&books).Error
	if err != nil {
		return c.Status(500).SendString("could not get books")
	}
	return c.JSON(books)
}

func (repo *Repository) SetupRoutes(app *fiber.App) {
	api := app.Group("/api")
	api.Post("/create_books", repo.CreateBook)
	api.Delete("/delete_book/:id", repo.DeleteBook)
	api.Get("/get_books/:id", repo.GetBookByID)
	api.Get("/books", repo.GetBooks)
}

func main() {
	envErr := godotenv.Load(".env")
	if envErr != nil {
		log.Fatal(envErr)
	}
	config := &storage.Config{
		Host:     os.Getenv("DB_HOST"),
		Port:     os.Getenv("DB_PORT"),
		Password: os.Getenv("DB_PASS"),
		User:     os.Getenv("DB_USER"),
		DBName:   os.Getenv("DB_NAME"),
		SSLMode:  os.Getenv("DB_SSL_MODE"),
	}
	fmt.Println("config: ", config)

	db, dbErr := storage.NewConnection(config)
	if dbErr != nil {
		log.Fatal("could not load database")
	}

	migrateErr := models.MigrateBooks(db)
	if migrateErr != nil {
		log.Fatal("could not migrate books")
	}

	repo := Repository{
		DB: db,
	}
	app := fiber.New()
	repo.SetupRoutes(app)
	log.Fatal(app.Listen(":8080"))
}
