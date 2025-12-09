package scheduler

import (
	"fmt"
	"log"
	"math/rand"
	"micro-site/api/model"
	"micro-site/database"
	"strconv"
	"time"
)

/*
* Seed fake users for testing
* @return void
 */
func SeedFakeUsers() {
	fmt.Println("Starting user seeding...")
	count := 20000
	batchSize := 1000
	var users []model.User

	// Seed random generator
	rand.Seed(time.Now().UnixNano())

	for i := 0; i < count; i++ {
		user := model.User{
			UserName:     "User_" + strconv.Itoa(i) + "_" + strconv.Itoa(rand.Intn(100000)),
			Email:        "user_" + strconv.Itoa(i) + "_" + strconv.Itoa(rand.Intn(100000)) + "@example.com",
			MobileNumber: 1000000000 + rand.Intn(9000000000), // Random 10-digit number
			Status:       "Active",
			AboutMe:      "This is a fake user generated for testing purposes.",
		}
		users = append(users, user)

		// Insert in batches
		if len(users) >= batchSize {
			if err := database.DB.CreateInBatches(users, batchSize).Error; err != nil {
				log.Printf("Error seeding users batch: %v", err)
			} else {
				fmt.Printf("Seeded batch of %d users\n", batchSize)
			}
			users = []model.User{} // Reset slice
		}
	}

	// Insert remaining users
	if len(users) > 0 {
		if err := database.DB.CreateInBatches(users, len(users)).Error; err != nil {
			log.Printf("Error seeding remaining users: %v", err)
		} else {
			fmt.Printf("Seeded remaining %d users\n", len(users))
		}
	}

	fmt.Println("User seeding completed.")
}
