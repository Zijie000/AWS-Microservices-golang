package main

import (
	"webapp/usercrud"

	"gorm.io/gorm"

	"log"

	"time"
)

func clearExpiredRecords(db *gorm.DB, interval, maxAge time.Duration) {
	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			expirationTime := time.Now().Add(-maxAge)
			result := db.Where("account_created < ?", expirationTime).Delete(&usercrud.UserCache{})
			log.Printf("user cache is running")
			if result.Error != nil {
				log.Printf("failed to clear expired records: %v", result.Error)
			} else {
				log.Printf("cleared %d expired records", result.RowsAffected)
			}
		}
	}
}
