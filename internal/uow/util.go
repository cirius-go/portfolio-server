package uow

import (
	"fmt"

	"gorm.io/gorm"
)

func lazyCache[Repo any](u *uow, key string, fn func(db *gorm.DB) *Repo) *Repo {
	u.mu.Lock()
	defer u.mu.Unlock()

	k, ok := u.caches[key]
	if !ok {
		u.caches[key] = fn(u.db)
		return u.caches[key].(*Repo)
	}

	switch v := k.(type) {
	case *Repo:
		return v
	default:
		// Should never happen.
		fmt.Printf("unexpected type %T of %s repo in lazyCache UnitOfWork.\n", v, key)
		return u.caches[key].(*Repo)
	}
}
