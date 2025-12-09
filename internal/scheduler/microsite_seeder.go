package scheduler

import (
	"fmt"
	"log"
	"math/rand"
	"micro-site/api/model"
	"micro-site/database"
	"strconv"
	"time"

	"github.com/google/uuid"
	"github.com/gosimple/slug"
	"gorm.io/gorm"
)

/*
* Seed fake microsites for testing
* @return void
 */
func SeedFakeMicrosites() {
	fmt.Println("Starting microsite seeding...")
	count := 1000
	batchSize := 100

	// Fetch a list of existing user IDs to assign microsites to
	var userIDs []uint64
	if err := database.DB.Model(&model.User{}).Pluck("id", &userIDs).Error; err != nil {
		log.Printf("Error fetching user IDs: %v", err)
		return
	}

	if len(userIDs) == 0 {
		log.Println("No users found. Cannot seed microsites.")
		return
	}

	var microsites []model.MicroSite
	rand.Seed(time.Now().UnixNano())

	for i := 0; i < count; i++ {
		if i%100 == 0 {
			fmt.Println("Preparing microsite", i)
		}
		userID := userIDs[rand.Intn(len(userIDs))]
		title := "Microsite " + strconv.Itoa(i) + " " + strconv.Itoa(rand.Intn(100000))

		// Generate slug manually as BeforeCreate hook might not trigger on batch insert or simple struct creation without DB context
		// However, for bulk insert we might need a simple unique string.
		// Real helper.GenerateUniqueSlug needs DB. Let's make a simple unique slug here.
		slug := slug.Make(title) + "-" + uuid.New().String()

		ms := model.MicroSite{
			Uuid:        uuid.New(),
			UserId:      userID,
			FullName:    "Microsite Owner " + strconv.Itoa(i),
			Title:       title,
			Slug:        slug,
			Description: "This is a fake microsite.",
			Status:      "Active",
			IsDraft:     "0",
			AvatarIcon:  "default_avatar.png",
			BannerImage: "default_banner.png",
		}

		// Add Services
		numServices := rand.Intn(3) + 1
		for j := 0; j < numServices; j++ {
			ms.Services = append(ms.Services, model.Service{
				Name:   "Service " + strconv.Itoa(j),
				UserId: userID,
				Status: true,
			})
		}

		// Add Social Links
		numLinks := rand.Intn(2) + 1
		links := []string{"facebook", "twitter", "instagram", "linkedin"}
		for j := 0; j < numLinks; j++ {
			originalLink := links[rand.Intn(len(links))]
			ms.SocialLinks = append(ms.SocialLinks, model.SocialLink{
				Name:   originalLink,
				Url:    "https://" + originalLink + ".com/user" + strconv.Itoa(i),
				Type:   originalLink,
				UserId: userID,
			})
		}

		microsites = append(microsites, ms)

		// Insert in batches
		if len(microsites) >= batchSize {
			fmt.Println("Saving batch of", len(microsites))
			if err := saveMicrositesBatch(microsites); err != nil {
				log.Printf("Error seeding microsites batch: %v", err)
			} else {
				fmt.Printf("Seeded batch of %d microsites\n", batchSize)
			}
			microsites = []model.MicroSite{}
		}
	}

	// Insert remaining
	if len(microsites) > 0 {
		if err := saveMicrositesBatch(microsites); err != nil {
			log.Printf("Error seeding remaining microsites: %v", err)
		} else {
			fmt.Printf("Seeded remaining %d microsites\n", len(microsites))
		}
	}

	fmt.Println("Microsite seeding completed.")
}

/*
* Save batch of microsites
* @param microsites []model.MicroSite
* @return error
 */
func saveMicrositesBatch(microsites []model.MicroSite) error {
	// Note: CreateInBatches with associations can be tricky with large datasets.
	// If it fails, we might need to save parent first then children.
	// But GORM usually handles it.
	return database.DB.Session(&gorm.Session{CreateBatchSize: 100}).Create(&microsites).Error
}
