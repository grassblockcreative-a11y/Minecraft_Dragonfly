package player

import (
	"encoding/json"
	"os"
	"strings"
	"sync"
	"time"
)

// BanEntry represents a ban on the server.
type BanEntry struct {
	Name      string    `json:"name"`
	XUID      string    `json:"xuid,omitempty"`
	UUID      string    `json:"uuid,omitempty"`
	IP        string    `json:"ip,omitempty"`
	BannedAt  time.Time `json:"banned_at"`
	ExpiresAt time.Time `json:"expires_at"` // Zero time means permanent
	Reason    string    `json:"reason"`
}

// BanManager handles player bans.
type BanManager struct {
	mu   sync.RWMutex
	bans map[string]BanEntry
	path string
}

// Bans is the global ban manager.
var Bans = &BanManager{
	bans: make(map[string]BanEntry),
	path: "bans.json",
}

// Load loads the banned players from bans.json.
func (bm *BanManager) Load() error {
	bm.mu.Lock()
	defer bm.mu.Unlock()
	data, err := os.ReadFile(bm.path)
	if err != nil {
		if os.IsNotExist(err) {
			return nil
		}
		return err
	}
	var list []BanEntry
	if err := json.Unmarshal(data, &list); err != nil {
		return err
	}
	bm.bans = make(map[string]BanEntry)
	for _, entry := range list {
		bm.bans[strings.ToLower(entry.Name)] = entry
	}
	return nil
}

// save writes bans to disk. Must be called WITHOUT holding bm.mu.
func (bm *BanManager) save() error {
	bm.mu.RLock()
	var list []BanEntry
	for _, entry := range bm.bans {
		list = append(list, entry)
	}
	bm.mu.RUnlock()

	data, err := json.MarshalIndent(list, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(bm.path, data, 0644)
}

// Save saves the banned players to bans.json.
func (bm *BanManager) Save() error {
	return bm.save()
}

// Ban bans a player by name, with optional XUID, UUID, and IP.
// IP bans are intentionally NOT stored to prevent blocking all
// players on the same local network (LAN/nethernet).
func (bm *BanManager) Ban(name string, duration time.Duration, reason string, xuid string, uuidStr string, ip string) {
	bm.mu.Lock()
	entry := BanEntry{
		Name:     name,
		XUID:     xuid,
		UUID:     uuidStr,
		IP:       "", // Do NOT store IP to prevent LAN-wide bans
		BannedAt: time.Now(),
		Reason:   reason,
	}
	if duration > 0 {
		entry.ExpiresAt = time.Now().Add(duration)
	}
	bm.bans[strings.ToLower(name)] = entry
	bm.mu.Unlock()
	_ = bm.save()
}

// Unban unbans a player by name.
func (bm *BanManager) Unban(name string) bool {
	bm.mu.Lock()
	nameLower := strings.ToLower(name)
	_, exists := bm.bans[nameLower]
	if exists {
		delete(bm.bans, nameLower)
	}
	bm.mu.Unlock()
	if exists {
		_ = bm.save()
	}
	return exists
}

// UnbanAll removes all bans. Useful for emergency reset.
func (bm *BanManager) UnbanAll() int {
	bm.mu.Lock()
	count := len(bm.bans)
	bm.bans = make(map[string]BanEntry)
	bm.mu.Unlock()
	if count > 0 {
		_ = bm.save()
	}
	return count
}

// IsBanned checks if a player is banned by name, XUID, or UUID.
// IP-based checking is intentionally excluded to prevent false
// positives on local networks where all players share the same IP.
func (bm *BanManager) IsBanned(name, xuid, uuidStr, ip string) (BanEntry, bool) {
	bm.mu.RLock()

	now := time.Now()
	var expiredKeys []string
	var foundEntry BanEntry
	var found bool

	for k, entry := range bm.bans {
		if !entry.ExpiresAt.IsZero() && now.After(entry.ExpiresAt) {
			expiredKeys = append(expiredKeys, k)
			continue
		}

		if !found {
			if strings.EqualFold(entry.Name, name) {
				foundEntry = entry
				found = true
			} else if xuid != "" && entry.XUID != "" && entry.XUID == xuid {
				foundEntry = entry
				found = true
			} else if uuidStr != "" && entry.UUID != "" && strings.EqualFold(entry.UUID, uuidStr) {
				foundEntry = entry
				found = true
			}
			// NOTE: IP matching is intentionally removed.
			// On local networks (LAN/nethernet), all players share the
			// same IP address. Matching by IP would block ALL players
			// from the same network when only one should be banned.
		}
	}
	bm.mu.RUnlock()

	// Clean up expired entries with a separate write lock.
	if len(expiredKeys) > 0 {
		bm.mu.Lock()
		for _, k := range expiredKeys {
			delete(bm.bans, k)
		}
		bm.mu.Unlock()
		go func() { _ = bm.save() }()
	}

	return foundEntry, found
}

// List returns all active bans.
func (bm *BanManager) List() []BanEntry {
	bm.mu.RLock()
	var active []BanEntry
	now := time.Now()
	var expiredKeys []string
	for nameLower, entry := range bm.bans {
		if !entry.ExpiresAt.IsZero() && now.After(entry.ExpiresAt) {
			expiredKeys = append(expiredKeys, nameLower)
			continue
		}
		active = append(active, entry)
	}
	bm.mu.RUnlock()

	if len(expiredKeys) > 0 {
		bm.mu.Lock()
		for _, k := range expiredKeys {
			delete(bm.bans, k)
		}
		bm.mu.Unlock()
		go func() { _ = bm.save() }()
	}
	return active
}
