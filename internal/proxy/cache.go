package proxy

import "time"

// cacheCleaner запускается в горутине и управляет размером кэша
func (ps *ProxyServer) cacheCleaner() {
	ticker := time.NewTicker(1 * time.Minute) // Проверяем кэш каждые 1 минуту
	defer ticker.Stop()

	for range ticker.C {
		ps.cacheMutex.Lock()

		// Удаляем устаревшие записи
		for domain, entry := range ps.cache {
			if time.Now().After(entry.expiresAt) {
				delete(ps.cache, domain)
			}
		}

		// Если кэш превышает лимит, удаляем самые старые записи
		if len(ps.cache) > ps.maxCache {
			var oldestDomain string
			var oldestTime time.Time = time.Now()

			for domain, entry := range ps.cache {
				if entry.expiresAt.Before(oldestTime) {
					oldestDomain = domain
					oldestTime = entry.expiresAt
				}
			}
			if oldestDomain != "" {
				delete(ps.cache, oldestDomain)
			}
		}

		ps.cacheMutex.Unlock()
	}
}
