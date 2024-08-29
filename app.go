package main

import (
	"context"
	"time"

	"github.com/wailsapp/wails/v2/pkg/runtime"
)

type Types interface {
	Null | UserInfo | GameId | string | GameInfo | GameInfos | Stats | StatsList
}

type Message[T Types] struct {
	Code MsgCode `json:"code"`
	Msg  string  `json:"msg"`
	Data T       `json:"data"`
}

type Null = struct{}

//go:generate stringer -linecomment -output app_gen.go -type=MsgCode
type MsgCode int

func Msg[T Types](code MsgCode, data ...T) Message[T] {
	switch len(data) {
	case 0:
		return Message[T]{Code: code, Msg: code.String()}
	case 1:
		return Message[T]{Code: code, Msg: code.String(), Data: data[0]}
	default:
		panic(len(data))
	}
}

// wails enum bind require
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
	NotSelectGame                 // 请选择加速的游戏
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
	"NotSelectGame",
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
	NotSelectGame,
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
	LastActive      int64  `json:"last_active"` // utc 时间戳
	Duration        int64  `json:"duration"`
	Flow            int64  `json:"flow"` // 加速流量
}
type GameInfos = []GameInfo

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
	Name     string `json:"name"`
	Password string `json:"password"`
	Phone    string `json:"phone"`
	Expire   int64  `json:"expire"` // stamp
}

func (a *App) GetUser() (msg Message[UserInfo]) {
	return a.Mock.GetUser()
}

// GetUser 获取用户信息, 应用渲染完成即调用此函数, 如果msg.Code==NotLogin, 则弹出注册登录页面

// RegisterOrLogin 注册或登录,
func (a *App) RegisterOrLogin(user, pwd string) (msg Message[Null]) {
	return a.Mock.RegisterOrLogin(user, pwd)
}

// todo: 暂时不考虑
// Recharge 充值，返回一个字符二维码。参考 https://wails.io/zh-Hans/docs/reference/runtime/events
func (a *App) Recharge(months int, eventName string) (msg Message[string]) {
	return a.Mock.Recharge(months, func(m Message[Null]) {
		runtime.EventsEmit(a.ctx, eventName, m)
	})
}

// ListGames 获取已添加的游戏列表, selectedIdx 表示默认应该选中的游戏
func (a *App) ListGames() (msg Message[GameInfos]) {
	return a.Mock.ListGames()
}

// SelectGame 选中某个游戏
func (a *App) SelectGame(gameId GameId) Message[GameInfo] {
	return a.Mock.SelectGame(gameId)
}

// GetSelectedGame 获取当前选中的游戏
func (a *App) GetSelectedGame() Message[GameInfo] {
	return a.Mock.GetSelectedGame()
}

// SearchGame 根据关键字搜索游戏
func (a *App) SearchGame(keyword string) (msg Message[GameInfos]) {
	return a.Mock.SearchGame(keyword)
}

// AddGame 新增游戏
func (a *App) AddGame(gameId GameId) Message[Null] {
	return a.Mock.AddGame(gameId)
}

// SetGame 选择某个游戏
func (a *App) SetGame(gameId GameId) Message[Null] {
	return a.Mock.SetGame(gameId)
}

// SetGameServer 设置游戏区服
func (a *App) SetGameServer(id GameId, gameServer string) Message[Null] {
	return a.Mock.SetGameServer(id, gameServer)
}

// SetRouteMode 选择路由模式
func (a *App) SetRouteMode(id GameId, fixRoute bool) Message[Null] {
	return a.Mock.SetRouteMode(id, fixRoute)
}

// Accelerate 开始加速
func (a *App) Accelerate(id GameId) Message[Null] {
	return a.Mock.Accelerate(id)
}

// DisableAccelerate 停止加速
func (a *App) DisableAccelerate() Message[Null] {
	return a.Mock.DisableAccelerate()
}

type StatsList struct {
	List []Stats `json:"list"`
}

type Stats struct {
	MilliStamp int `json:"stamp"` // 本数据点对应的时间戳

	GatewayLoc string `json:"gateway_loc"` // gateway所在城市名
	ForwardLoc string `json:"forward_loc"` // forward所在城市名
	ServerLoc  string `json:"server_loc"`  // server所在城市名

	RttGateway time.Duration `json:"rtt_gateway"` // 到gateway的延时
	RttForward time.Duration `json:"rtt_forward"` // 到forward的延时

	LossClientUplink    float64 `json:"loss_client_uplink"`
	LossClientDownlink  float64 `json:"loss_client_downlink"`
	LossGatewayUplink   float64 `json:"loss_gateway_uplink"`
	LossGatewayDownlink float64 `json:"loss_gateway_downlink"`
}

// Stats 获取统计信息, 阻塞函数, 如果距上次调用时间短于3s, 会主动阻塞直到恰好相距3s
func (a *App) Stats() Message[StatsList] {
	return a.Mock.Stats()
}
