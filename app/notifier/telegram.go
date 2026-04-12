package notifier

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
	"github.com/spf13/cast"
	"github.com/tidwall/gjson"
	"github.com/v03413/bepusdt/app"
	"github.com/v03413/bepusdt/app/conf"
	"github.com/v03413/bepusdt/app/log"
	"github.com/v03413/bepusdt/app/model"
	"github.com/v03413/bepusdt/app/utils"
	"github.com/v03413/tronprotocol/core"
)

type Telegram struct {
	api     *bot.Bot
	token   string
	chatID  int64
	topicID int
}

func (t *Telegram) Initialize(params string) error {
	info := gjson.Parse(params)

	t.token = info.Get("bot_token").String()
	t.chatID = info.Get("chat_id").Int()
	t.topicID = cast.ToInt(info.Get("topic_id").Int())

	b, err := bot.New(t.token)
	if err != nil {
		return err
	}

	t.api = b

	return nil
}

func (t *Telegram) Success(o model.Order) {
	if o.Status != model.OrderStatusSuccess {
		return
	}

	tradeType := string(o.TradeType)
	tokenType, err := model.GetCrypto(o.TradeType)
	if err != nil {
		t.sendMessage(&bot.SendMessageParams{Text: "❌交易类型不支持：" + tradeType})

		return
	}

	token := string(tokenType)

	text := `
\#收款成功 \#订单交易 \#` + token + `
\-\-\-
` + "```" + `
🚦商户订单：%v
💰请求金额：%v ` + string(o.Fiat) + `(%v)
💲支付数额：%v ` + tradeType + `
💎交易哈希：%s
✅收款地址：%s
⏱️创建时间：%s
️🎯️支付时间：%s
` + "```" + `
`
	text = fmt.Sprintf(text,
		o.OrderId,
		o.Money,
		o.Rate,
		o.Amount,
		utils.MaskHash(o.RefHash),
		utils.MaskAddress(o.Address),
		o.CreatedAt.Format(time.DateTime),
		o.UpdatedAt.Format(time.DateTime),
	)

	t.sendMessage(&bot.SendMessageParams{
		Text:      text,
		ParseMode: models.ParseModeMarkdown,
		ReplyMarkup: &models.InlineKeyboardMarkup{
			InlineKeyboard: [][]models.InlineKeyboardButton{
				{
					models.InlineKeyboardButton{Text: "📝查看交易明细", URL: o.GetTxUrl()},
				},
			},
		},
	})
}

func (t *Telegram) NotifyFail(o model.Order, reason string) {
	tradeType := string(o.TradeType)
	tokenT, err := model.GetCrypto(o.TradeType)
	if err != nil {
		t.sendMessage(&bot.SendMessageParams{Text: "❌交易类型不支持：" + tradeType})

		return
	}

	token := string(tokenT)

	text := fmt.Sprintf(`
\#回调失败 \#订单交易 \#`+token+`
\-\-\-
`+"```"+`
🚦商户订单：%v
💲支付数额：%v
💰请求金额：%v `+string(o.Fiat)+`(%v)
💍交易类别：%s
⚖️️确认时间：%s
⏰下次回调：%s
🗒️失败原因：%s
`+"```"+`
`,
		utils.Ec(o.OrderId),
		o.Amount,
		o.Money, o.Rate,
		strings.ToUpper(tradeType),
		o.ConfirmedAt.Format(time.DateTime),
		utils.CalcNextNotifyTime(*o.ConfirmedAt, o.NotifyNum+1).Format(time.DateTime),
		reason,
	)

	t.sendMessage(&bot.SendMessageParams{
		Text:      text,
		ParseMode: models.ParseModeMarkdown,
		ReplyMarkup: &models.InlineKeyboardMarkup{
			InlineKeyboard: [][]models.InlineKeyboardButton{
				{
					models.InlineKeyboardButton{Text: "📝查看收款详情", CallbackData: o.GetTxUrl()},
				},
			},
		},
	})
}

func (t *Telegram) NonOrderTransfer(trans model.TronTransfer, wa model.Wallet) {
	title := "收入"
	if trans.RecvAddress != wa.Address {
		title = "支出"
	}

	text := fmt.Sprintf(
		"\\#账户%s \\#非订单交易\n\\-\\-\\-\n```\n💲交易数额：%v \n💍交易类别："+strings.ToUpper(string(trans.TradeType))+"\n⏱️交易时间：%v\n✅接收地址：%v\n🅾️发送地址：%v```\n",
		title,
		trans.Amount.String(),
		trans.Timestamp.Format(time.DateTime),
		utils.MaskAddress(trans.RecvAddress),
		utils.MaskAddress(trans.FromAddress),
	)

	t.sendMessage(&bot.SendMessageParams{
		Text:      text,
		ParseMode: models.ParseModeMarkdown,
		ReplyMarkup: models.InlineKeyboardMarkup{
			InlineKeyboard: [][]models.InlineKeyboardButton{
				{
					models.InlineKeyboardButton{Text: "📝查看交易明细", URL: model.GetTxUrl(trans.TradeType, trans.TxHash)},
				},
			},
		},
	})
}

func (t *Telegram) TronResourceChange(res model.TronResource) {
	title := "代理"
	if res.Type == core.Transaction_Contract_UnDelegateResourceContract {
		title = "回收"
	}

	text := fmt.Sprintf(
		"\\#资源动态 \\#能量"+title+"\n\\-\\-\\-\n```\n🔋质押数量："+cast.ToString(res.Balance/1000000)+"\n⏱️交易时间：%v\n✅操作地址：%v\n🅾️资源来源：%v```\n",
		res.Timestamp.Format(time.DateTime),
		utils.MaskAddress(res.RecvAddress),
		utils.MaskAddress(res.FromAddress),
	)

	t.sendMessage(&bot.SendMessageParams{
		Text:      text,
		ParseMode: models.ParseModeMarkdown,
		ReplyMarkup: models.InlineKeyboardMarkup{
			InlineKeyboard: [][]models.InlineKeyboardButton{
				{
					models.InlineKeyboardButton{Text: "📝查看交易明细", URL: "https://tronscan.org/#/transaction/" + res.ID},
				},
			},
		},
	})
}

func (t *Telegram) Welcome() {
	text := `
👋 欢迎使用 BEpusdt，` + conf.Desc + `，如果您看到此消息，说明系统已启动成功！

📌当前版本：` + app.Version + `
🎉开源地址：` + conf.Github + `
---
`
	t.sendMessage(&bot.SendMessageParams{
		Text: text,
		ReplyMarkup: models.InlineKeyboardMarkup{
			InlineKeyboard: [][]models.InlineKeyboardButton{
				{
					{Text: "📢 关注频道", URL: "https://t.me/BEpusdtChannel"},
					{Text: "💬 社区交流", URL: "https://t.me/BEpusdtChat"},
				},
			},
		},
	})
}

func (t *Telegram) Test() error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()
	_, err := t.api.SendMessage(ctx, &bot.SendMessageParams{
		ChatID:          t.chatID,
		MessageThreadID: t.topicID,
		Text:            "✅ 这是一条测试消息，Telegram 通知配置成功！\n当前系统时间：" + time.Now().Format("2006-01-02 15:04:05"),
	})

	return err
}

func (t *Telegram) sendMessage(p *bot.SendMessageParams) {
	p.ChatID = t.chatID
	p.MessageThreadID = t.topicID

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	_, err := t.api.SendMessage(ctx, p)
	if err != nil {
		log.Warn("Bot Send Message Error:", err.Error())
	}
}
