package nametransform

import (
	"sync"
)

// Unbound DirIV cache. Stores entries of directory IV and the encrypted
// path.
type dirIVCache struct {
	cache map[string]dirIVCacheEntry
	sync.Mutex
}

// DirIV cache entry. Stores the directory IV and the encrypted
// path.
type dirIVCacheEntry struct {
	// The DirIV
	iv []byte
	// Encrypted version of "dir"
	cDir string
}

// NewDirIVCache returns a new dirIVCache.
func NewDirIVCache() *dirIVCache {
	return &dirIVCache{
		cache: map[string]dirIVCacheEntry{},
	}
}

// lookup - fetch entry for "dir" from the cache
func (c *dirIVCache) lookup(dir string) ([]byte, string) {
	c.Lock()
	defer c.Unlock()
	entry, hit := c.cache[dir]
	if hit {
		return entry.iv, entry.cDir
	}
	return nil, ""
}

// store - write entry for "dir" into the cache
func (c *dirIVCache) store(dir string, iv []byte, cDir string) {
	c.Lock()
	defer c.Unlock()
	c.cache[dir] = dirIVCacheEntry{iv, cDir}
}

// Clear ... clear the cache.
// Exported because it is called from fusefrontend when directories are
// renamed or deleted.
func (c *dirIVCache) Clear() {
	c.Lock()
	defer c.Unlock()
	c.cache = map[string]dirIVCacheEntry{}
}
