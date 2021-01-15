package middleware

import (
	"errors"
	"fmt"
	"runtime"
	"strings"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/maybgit/glog"
	"github.com/spf13/cast"
)

//校验渠道id和来源，并写入context
func Auth() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			c.Set("channel_id", cast.ToInt32(c.Request().Header.Get("channel_id")))
			c.Set("user_agent", cast.ToInt32(c.Request().Header.Get("user_agent")))
			return next(c)
		}
	}
}

func MyRecover(config middleware.RecoverConfig) echo.MiddlewareFunc {
	// Defaults
	if config.Skipper == nil {
		config.Skipper = middleware.DefaultRecoverConfig.Skipper
	}
	if config.StackSize == 0 {
		config.StackSize = middleware.DefaultRecoverConfig.StackSize
	}

	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			if config.Skipper(c) {
				return next(c)
			}

			defer func() {
				if rErr := recover(); rErr != nil {
					err, ok := rErr.(error)
					if !ok {
						err = fmt.Errorf("%v", rErr)
					}
					stack := make([]byte, config.StackSize)
					length := runtime.Stack(stack, !config.DisableStackAll)
					if !config.DisablePrintStack {
						glog.Errorf("[PANIC RECOVER] %v %s\n", err, stack[:length])
					}
					c.Error(err)
				}
			}()
			return next(c)
		}
	}
}

//返回错误信息处理，屏蔽内部服务调用错误
func MyErrorHandle() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			err := next(c)

			if err != nil && strings.Contains(err.Error(), "rpc error") {
				glog.Errorf("[内部错误]，%v，%v", c.Path(), err)
				c.Error(errors.New("内部错误"))
			}

			return err
		}
	}
}
