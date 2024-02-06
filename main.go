package main

import (
	"fmt"
	"log"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type App struct {
	db *gorm.DB
}

var app *App

type Person struct {
	ID   int `gorm:"primaryKey"`
	Name string
	Age  int
}

func init() {

	dsn := "root:root@tcp(192.168.1.200:3306)/test?charset=utf8mb4&parseTime=True&loc=Local"
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	// Auto-migrate the schema
	err = db.AutoMigrate(&Person{})
	if err != nil {
		log.Fatal("Failed to auto-migrate schema:", err)
	}

	app = &App{db: db}
}

func main() {

	// Create a new person
	person := &Person{Name: "Bob", Age: 35}
	app.db.Create(person)
	fmt.Println("Inserted person:", person)

	// Read the person by ID
	var retrievedPerson Person
	app.db.First(&retrievedPerson, person.ID)
	fmt.Println("Retrieved person:", retrievedPerson)

	// Update the person's age
	app.db.Model(&retrievedPerson).Update("Age", 38)
	fmt.Println("Updated person:", retrievedPerson)

	// Delete the person
	// app.db.Delete(&retrievedPerson)
	// fmt.Println("Deleted person:", retrievedPerson)
}
