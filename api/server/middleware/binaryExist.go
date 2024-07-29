package middleware

import (
	"errors"
	"fmt"
	"gitlab.com/xiayesuifeng/gopanel/api/server/router"
	"os/exec"
	"time"
)

var cache = make(map[string]result)

type result struct {
	Exist     bool
	StoreTime time.Time
}

func getResult(name string) (bool, bool) {
	if result, ok := cache[name]; ok {
		now := time.Now()
		if now.Sub(result.StoreTime).Seconds() > 5 {
			return false, false
		} else {
			return result.Exist, true
		}
	} else {
		return false, false
	}
}

func BinaryExistMiddleware(name string) func(ctx *router.Context) error {
	return func(ctx *router.Context) error {
		exist, ok := false, false
		if exist, ok = getResult(name); !ok {
			_, err := exec.LookPath(name)
			if err != nil {
				if errors.Is(err, exec.ErrNotFound) {
					exist = false
				} else {
					return err
				}
			} else {
				exist = true
			}

			cache[name] = result{exist, time.Now()}
		}

		if !exist {
			ctx.Abort()
			return ctx.Error(503, fmt.Sprintf("binary %s not found", name))
		}

		ctx.Next()
		return nil
	}
}
