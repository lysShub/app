package main

import (
	"context"
	"time"

	"github.com/wailsapp/wails/v2/pkg/runtime"
)

type Message struct {
	Code MsgCode `json:"code"`
	Msg  string  `json:"msg"`
}

//go:generate stringer -linecomment -output app_gen.go -type=MsgCode
type MsgCode int

func (c MsgCode) Message() Message { return Message{Code: c, Msg: c.String()} }
func (c MsgCode) TSName() string {
	if c >= _end {
		panic(c)
	}
	return codeLits[c]
}

const (
	OK             MsgCode = iota // ok
	Notfound                      // not found
	NotLogin                      // 没有登录
	IsLogined                     // 已经登录
	VIPExpired                    // vip 已过期
	NotSetGame                    // 请先选择加速的游戏
	Accelerating                  // 游戏已在加速
	InvalidMonths                 // 无效月数
	GameExist                     // 游戏已存在
	NotAccelerated                // 没有加速
	RequireGameId                 // 请指定游戏
	Unknown                       // 未知
	_end
)

var codeLits = []string{
	"OK",
	"Notfound",
	"NotLogin",
	"IsLogined",
	"VIPExpired",
	"NotSetGame",
	"Accelerating",
	"InvalidMonths",
	"GameExist",
	"NotAccelerated",
	"RequireGameId",
	"Unknown",
}
var codes = []MsgCode{
	OK,
	Notfound,
	NotLogin,
	IsLogined,
	VIPExpired,
	NotSetGame,
	Accelerating,
	InvalidMonths,
	GameExist,
	NotAccelerated,
	RequireGameId,
	Unknown,
}

type GameId = int32 // 最小值为1

type GameInfo struct {
	GameId      GameId   `json:"game_id"`
	Name        string   `json:"name"`
	IconPath    string   `json:"icon_path"`
	BgimgPath   string   `json:"bgimg_path"`
	GameServers []string `json:"game_servers"`

	CacheGameServer string `json:"cache_game_server"`
	CacheFixRoute   bool   `json:"cache_fix_route"`
}

type App struct {
	ctx context.Context
	*Mock
}

func NewMockApp() *App {
	return &App{Mock: (&Mock{}).init()}
}

func (a *App) startup(ctx context.Context)                    { a.ctx = ctx }
func (a *App) domReady(ctx context.Context)                   {}
func (a *App) beforeClose(ctx context.Context) (prevent bool) { return false }
func (a *App) shutdown(ctx context.Context)                   {}

type UserInfo struct {
	Name     string
	Password string
	Phone    string
	Expire   int64 // utc 时间戳
}

// GetUser 获取用户信息, 应用渲染完成即调用此函数, 如果msg.Code==NotLogin, 则弹出注册登录页面
func (a *App) GetUser() (info UserInfo, msg Message) {
	return a.Mock.GetUser()
}

// RegisterOrLogin 注册或登录,
func (a *App) RegisterOrLogin(user, pwd string) (msg Message) {
	return a.Mock.RegisterOrLogin(user, pwd)
}

// todo: 暂时不考虑
// Recharge 充值，返回一个字符二维码、和一个全局事件。参考 https://wails.io/zh-Hans/docs/reference/runtime/events
// 回调返回结果是Message类型
func (a *App) Recharge(months int, eventName string) (qrImagePath string, msg Message) {
	return a.Mock.Recharge(months, func(m Message) {
		runtime.EventsEmit(a.ctx, eventName, m)
	})
}

// ListGames 获取已添加的游戏列表, selectedIdx 表示默认应该选中的游戏
func (a *App) ListGames() (list []GameInfo, selectedIdx int, msg Message) {
	return a.Mock.ListGames()
}

func (a *App) SelectGame(gameId GameId) (GameInfo, Message) {
	return a.Mock.SelectGame(gameId)
}

// SearchGame 根据关键字搜索游戏
func (a *App) SearchGame(keyword string) (list []GameInfo, msg Message) {
	return a.Mock.SearchGame(keyword)
}

// AddGame 新增游戏
func (a *App) AddGame(gameId GameId) Message {
	return a.Mock.AddGame(gameId)
}

// SetGame 选择某个游戏
func (a *App) SetGame(gameId GameId) Message {
	return a.Mock.SetGame(gameId)
}

// SetGameServer 设置游戏区服
func (a *App) SetGameServer(id GameId, gameServer string) Message {
	return a.Mock.SetGameServer(id, gameServer)
}

// SetRouteMode 选择路由模式
func (a *App) SetRouteMode(id GameId, fixRoute bool) Message {
	return a.Mock.SetRouteMode(id, fixRoute)
}

// Accelerate 开始加速
func (a *App) Accelerate(id GameId) Message {
	return a.Mock.Accelerate(id)
}

// DisableAccelerate 停止加速
func (a *App) DisableAccelerate() Message {
	return a.Mock.DisableAccelerate()
}

type Stats struct {
	GatewayLocation    string `json:"gateway_location"`
	DestinatioLocation string `json:"destinatio_location"`
	GameServerLocation string `json:"gameserver_location"`

	LossUplink1   float64       `json:"loss_uplink1"`   // 第一阶段上行丢包
	LossDownlink1 float64       `json:"loss_downlink1"` // 第一阶段下行丢包
	LossUplink2   float64       `json:"loss_uplink2"`   // 第二阶段上行丢包
	LossDownlink2 float64       `json:"loss_downlink2"` // 第二阶段下行丢包
	Ping1         time.Duration `json:"ping1"`          // 第一阶段延时
	Ping2         time.Duration `json:"ping2"`          // 整体延时
}

// Stats 获取统计信息, 阻塞函数, 如果距上次调用时间短于3s, 会主动阻塞直到恰好相距3s
func (a *App) Stats() (s Stats) {
	return a.Mock.Stats()
}
