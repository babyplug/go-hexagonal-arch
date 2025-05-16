package background

import (
	"context"
	"log"
	"time"

	"clean-arch/internal/core/port"
)

func StartUserCountLogger(repo port.UserRepository, stopCh <-chan struct{}) {
	ticker := time.NewTicker(10 * time.Second)
	ctx := context.Background()
	go func() {
		for {
			select {
			case <-ticker.C:
				count, err := repo.Count(ctx)
				if err == nil {
					log.Printf("User count: %d", count)
				} else {
					log.Printf("Failed to count users: %v", err)
				}
			case <-stopCh:
				ticker.Stop()
				return
			}
		}
	}()
}
