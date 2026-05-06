package admin

import (
	"errors"
	"fmt"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/shopspring/decimal"
	"github.com/spf13/cast"
	"github.com/v03413/bepusdt/app/handler/base"
	"github.com/v03413/bepusdt/app/model"
	"github.com/v03413/go-cache"
)

type Dashboard struct {
}

type homeReq struct {
	Fiat  string `json:"fiat" binding:"required"`
	Range string `json:"range"`
	TZ    string `json:"tz"`
	From  string `json:"from"`
	To    string `json:"to"`
	Force bool   `json:"force"`
}

type dashboardPoint struct {
	Date          string `json:"date"`
	OrdersTotal   int64  `json:"orders_total"`
	OrdersPaid    int64  `json:"orders_paid"`
	OrdersSuccess int64  `json:"orders_success"`
	OrdersFailed  int64  `json:"orders_failed"`
	GMVPaid       string `json:"gmv_paid"`
}

func (Dashboard) Home(ctx *gin.Context) {
	var req homeReq
	if err := ctx.ShouldBindJSON(&req); err != nil {
		base.BadRequest(ctx, err.Error())

		return
	}

	rangeKey, from, to, loc, err := resolveDashboardRange(req)
	if err != nil {
		base.BadRequest(ctx, err.Error())

		return
	}

	cacheKey := dashboardCacheKey(req.Fiat, rangeKey, loc.String(), from, to)
	if !req.Force {
		if data, ok := cache.Get(cacheKey); ok {
			base.Ok(ctx, data)

			return
		}
	}

	data := buildDashboardHome(req, rangeKey, from, to, loc)
	cache.Set(cacheKey, data, dashboardCacheTTL(rangeKey, to, loc))
	base.Ok(ctx, data)
}

func buildDashboardHome(req homeReq, rangeKey string, from, to time.Time, loc *time.Location) gin.H {
	var rows = make([]model.Order, 0)
	model.Db.Where("fiat = ? and created_at >= ? and created_at <= ?", req.Fiat, from, to).Find(&rows)

	var totalCount, pendingCount, confirmingCount, successCount, expiredCount, failedCount, notifyFailedCount int64
	gmvPaid := decimal.Zero
	var tokenMap = map[model.Crypto]float64{
		model.USDT: 0,
		model.USDC: 0,
		model.TRX:  0,
		model.BNB:  0,
		model.ETH:  0,
	}
	points := make([]dashboardPoint, 0)
	pointMap := make(map[string]int)
	for day := dayStart(from, loc); !day.After(dayStart(to, loc)); day = day.AddDate(0, 0, 1) {
		date := day.Format("2006-01-02")
		point := dashboardPoint{
			Date:    date,
			GMVPaid: "0.00",
		}
		points = append(points, point)
		pointMap[date] = len(points) - 1
	}

	for _, itm := range rows {
		totalCount++
		createdAt := time.Now().In(loc)
		if itm.CreatedAt != nil {
			createdAt = time.Time(*itm.CreatedAt).In(loc)
		}
		date := createdAt.Format("2006-01-02")
		pointIndex, ok := pointMap[date]
		if !ok {
			continue
		}
		point := &points[pointIndex]

		point.OrdersTotal++

		switch itm.Status {
		case model.OrderStatusWaiting:
			pendingCount++
		case model.OrderStatusConfirming:
			confirmingCount++
		case model.OrderStatusExpired:
			expiredCount++
		case model.OrderStatusSuccess:
			successCount++
			point.OrdersSuccess++
			point.OrdersPaid++

			money, _ := decimal.NewFromString(itm.Money)
			gmvPaid = gmvPaid.Add(money)
			pointMoney, _ := decimal.NewFromString(point.GMVPaid)
			point.GMVPaid = pointMoney.Add(money).Round(2).StringFixed(2)

			crypto := itm.Crypto
			if crypto == "" {
				if tradeCrypto, err := model.GetCrypto(itm.TradeType); err == nil {
					crypto = tradeCrypto
				}
			}
			if _, ok := tokenMap[crypto]; !ok {
				tokenMap[crypto] = 0
			}
			tokenAmount := cast.ToFloat64(itm.Amount)
			if tokenAmount == 0 {
				tokenAmount = cast.ToFloat64(itm.Money)
			}
			tokenMap[crypto] += tokenAmount

			if itm.NotifyState == model.OrderNotifyStateFail {
				notifyFailedCount++
			}
		case model.OrderStatusFailed:
			failedCount++
			point.OrdersFailed++
		}
	}

	successRate := "0.00"
	finishedCount := successCount + expiredCount + failedCount
	if finishedCount > 0 {
		rate := decimal.NewFromInt(successCount).Div(decimal.NewFromInt(finishedCount)).Mul(decimal.NewFromInt(100))
		successRate = rate.Round(2).StringFixed(2)
	}

	return gin.H{
		"range":    rangeKey,
		"from":     from.In(loc).Format(time.RFC3339),
		"to":       to.In(loc).Format(time.RFC3339),
		"timezone": loc.String(),
		"currency": req.Fiat,
		"kpi": gin.H{
			"orders_total":       totalCount,
			"orders_pending":     pendingCount,
			"orders_confirming":  confirmingCount,
			"gmv_paid":           gmvPaid.Round(2).StringFixed(2),
			"orders_success":     successCount,
			"orders_failed":      failedCount,
			"notify_failed":      notifyFailedCount,
			"order_success_rate": successRate,
		},
		"token_map": tokenMap,
		"points":    points,
	}
}

func dashboardCacheKey(fiat, rangeKey, timezone string, from, to time.Time) string {
	return fmt.Sprintf("dashboard:home:v1:%s:%s:%s:%d:%d", fiat, rangeKey, timezone, from.Unix(), to.Unix())
}

func dashboardCacheTTL(rangeKey string, to time.Time, loc *time.Location) time.Duration {
	switch rangeKey {
	case "today":
		return 10 * time.Second
	case "30d":
		return time.Minute
	case "custom":
		if to.Before(dayStart(time.Now().In(loc), loc)) {
			return 5 * time.Minute
		}

		return 30 * time.Second
	default:
		return 30 * time.Second
	}
}

func resolveDashboardRange(req homeReq) (string, time.Time, time.Time, *time.Location, error) {
	loc := time.Local
	if req.TZ != "" {
		if loadedLoc, err := time.LoadLocation(req.TZ); err == nil {
			loc = loadedLoc
		}
	}

	rangeKey := req.Range
	if rangeKey == "" {
		rangeKey = "7d"
	}

	now := time.Now().In(loc)
	switch rangeKey {
	case "today":
		start := dayStart(now, loc)
		return rangeKey, start, dayEnd(start), loc, nil
	case "30d":
		yesterday := dayStart(now, loc).AddDate(0, 0, -1)
		start := yesterday.AddDate(0, 0, -29)
		return rangeKey, start, dayEnd(yesterday), loc, nil
	case "custom":
		from, err := parseDashboardDate(req.From, loc, true)
		if err != nil {
			return "", time.Time{}, time.Time{}, nil, err
		}
		to, err := parseDashboardDate(req.To, loc, false)
		if err != nil {
			return "", time.Time{}, time.Time{}, nil, err
		}
		if to.Before(from) {
			return "", time.Time{}, time.Time{}, nil, errors.New("自定义统计周期结束时间不能早于开始时间")
		}

		return rangeKey, from, to, loc, nil
	case "7d":
		fallthrough
	default:
		yesterday := dayStart(now, loc).AddDate(0, 0, -1)
		start := yesterday.AddDate(0, 0, -6)
		return "7d", start, dayEnd(yesterday), loc, nil
	}
}

func parseDashboardDate(value string, loc *time.Location, start bool) (time.Time, error) {
	if value == "" {
		return time.Time{}, errors.New("请选择自定义统计周期")
	}

	layouts := []string{time.RFC3339, "2006-01-02 15:04:05", "2006-01-02"}
	for _, layout := range layouts {
		t, err := time.ParseInLocation(layout, value, loc)
		if err != nil {
			continue
		}
		if layout == "2006-01-02" {
			if start {
				return dayStart(t, loc), nil
			}
			return dayEnd(t), nil
		}

		return t.In(loc), nil
	}

	return time.Time{}, errors.New("自定义统计周期格式错误")
}

func dayStart(t time.Time, loc *time.Location) time.Time {
	local := t.In(loc)
	return time.Date(local.Year(), local.Month(), local.Day(), 0, 0, 0, 0, loc)
}

func dayEnd(t time.Time) time.Time {
	start := time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, t.Location())
	return start.Add(24*time.Hour - time.Nanosecond)
}
