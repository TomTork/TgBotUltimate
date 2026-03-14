package Action

import (
	"TgBotUltimate/types/Database"
	"context"
	"github.com/mymmrac/telego"
)

type Action struct {
	ReqCtx   context.Context
	Ctx      context.Context
	Update   telego.Update
	Database *Database.DB
	Bot      *telego.Bot
}
