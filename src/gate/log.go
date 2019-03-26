package gate

import (
	"core/log"
	"core/share"
	"fmt"
)

var Logger = log.NewLogger(&log.Config{
	DevMode:    share.Env.DevMode,
	Filepath:   fmt.Sprintf("../log/gate-%d/", share.Env.BootID),
	CallerSkip: 3,
	MaxSize:    10,
})
