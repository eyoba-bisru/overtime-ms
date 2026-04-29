package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/eyoba-bisru/overtime-backend/internal/config"
	"github.com/eyoba-bisru/overtime-backend/internal/models"
	"github.com/eyoba-bisru/overtime-backend/internal/repository"
	"github.com/eyoba-bisru/overtime-backend/internal/utils"
	"github.com/google/uuid"
)

func main() {
	fmt.Println("=== Overtime MS: Create Superuser (Admin) ===")

	// 1. Initialize DB
	config.DBConnect()
	defer config.CloseDB()

	reader := bufio.NewReader(os.Stdin)

	// 2. Prompt for Name
	fmt.Print("Name: ")
	name, _ := reader.ReadString('\n')
	name = strings.TrimSpace(name)

	// 3. Prompt for Email
	fmt.Print("Email address: ")
	email, _ := reader.ReadString('\n')
	email = strings.TrimSpace(email)

	// 4. Prompt for Password
	fmt.Print("Password: ")
	password, _ := reader.ReadString('\n')
	password = strings.TrimSpace(password)

	if name == "" || email == "" || password == "" {
		log.Fatal("All fields are required")
	}

	// 5. Hash Password
	hashedPassword, err := utils.HashPassword(password)
	if err != nil {
		log.Fatal("Failed to hash password:", err)
	}

	// 6. Create User Object
	user := &models.User{
		Base: models.Base{
			ID:        uuid.New(),
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
		Name:                name,
		Email:               email,
		PasswordHash:        hashedPassword,
		Role:                models.Admin,
		ForcePasswordChange: false,
	}

	// 7. Save to DB
	id, err := repository.CreateUserRepo(user)
	if err != nil {
		log.Fatalf("Error creating user: %v", err)
	}

	fmt.Printf("\nSuperuser created successfully with ID: %s\n", id)
}
