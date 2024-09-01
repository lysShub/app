package main

import (
	"context"
	"time"

	"github.com/wailsapp/wails/v2/pkg/runtime"
)

type Message struct {
	Code MsgCode `json:"code"`
	Msg  string  `json:"msg"`
	Data any     `json:"data"`
}

//go:generate stringer -linecomment -output app_gen.go -type=MsgCode
type MsgCode int

func (c MsgCode) Message(data any) Message {
	return Message{Code: c, Msg: c.String(), Data: data}
}
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
	GameNotExist                  // 游戏不存在
	NotAccelerated                // 没有加速
	RequireGameId                 // 未指定游戏
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
	"GameNotExist",
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

type UserInfo struct {
	Name     string `json:"name"`
	Password string `json:"password"`
	Phone    string `json:"phone"`
	Expire   int64  `json:"expire"` // utc 时间戳
	Icon     string `json:"icon"`
}

type GameId = int64 // 最小值为1

type GameInfo struct {
	GameId      GameId   `json:"game_id"`
	Name        string   `json:"name"`
	IconPath    string   `json:"icon_path"`
	BgimgPath   string   `json:"bgimg_path"`
	AdimgPath   string   `json:"adimg_path"`
	GameServers []string `json:"game_servers"` // 第一个为默认值
	FixRoute    bool     `json:"fix_route"`

	LastActive   int64 `json:"last_active"`   // utc 时间戳
	TotalSeconds int64 `json:"total_seconds"` // 累计加速时间(s)
	TotalBytes   int64 `json:"total_bytes"`   // 累计加速流量(B)
}
type GameInfos struct {
	SelectGame GameId     `json:"select_game"`
	AcceleGame GameId     `json:"accele_game"`
	Games      []GameInfo `json:"games"`
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

// GetUser 获取用户信息, 应用渲染完成即调用此函数, 如果msg.Code==NotLogin, 则弹出注册登录页面
func (a *App) GetUser() Message {
	return a.Mock.GetUser()
}

// RegisterOrLogin 注册或登录,
func (a *App) RegisterOrLogin(user, pwd string) (msg Message) {
	return a.Mock.RegisterOrLogin(user, pwd)
}

// todo: 暂时不考虑
// Recharge 充值，返回一个字符二维码、和一个全局事件。参考 https://wails.io/zh-Hans/docs/reference/runtime/events
// 回调返回结果是Message类型
func (a *App) Recharge(months int, eventName string) Message {
	return a.Mock.Recharge(months, func(m Message) {
		runtime.EventsEmit(a.ctx, eventName, m)
	})
}

// ListGames 获取已添加的游戏列表, selectedIdx 表示默认应该选中的游戏
func (a *App) ListGames() Message {
	return a.Mock.ListGames()
}

func (a *App) GetGame(id GameId) Message {
	return a.Mock.GetGame(id)
}

// SelectGame 选中某个游戏
func (a *App) SelectGame(gameId GameId) Message {
	return a.Mock.SelectGame(gameId)
}

// GetSelectedGame 获取当前选中的游戏
func (a *App) GetSelectedGame() Message {
	return a.Mock.GetSelectedGame()
}

// SearchGame 根据关键字搜索游戏
func (a *App) SearchGame(keyword string) Message {
	return a.Mock.SearchGame(keyword)
}

// AddGame 新增游戏
func (a *App) AddGame(gameId GameId) Message {
	return a.Mock.AddGame(gameId)
}

// DelGame 删除游戏
func (a *App) DelGame(gameId GameId) Message {
	return a.Mock.DelGame(gameId)
}

// SetGameServer 设置游戏区服
func (a *App) SetGameServer(gameServer string) Message {
	return a.Mock.SetGameServer(gameServer)
}

// SetRouteMode 选择路由模式
func (a *App) SetRouteMode(fixRoute bool) Message {
	return a.Mock.SetRouteMode(fixRoute)
}

// StartAccele 开始加速
func (a *App) StartAccele(id GameId) Message {
	return a.Mock.StartAccele(id)
}

// StopAccele 停止加速
func (a *App) StopAccele() Message {
	return a.Mock.StopAccele()
}

type Stats struct {
	Stamp   int64 `json:"stamp"`   // 本数据点对应的毫秒时间戳
	Seconds int64 `json:"seconds"` // 本次加速时长
	Bytes   int64 `json:"bytes"`   // 本次加速流量

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
func (a *App) Stats() Stats {
	return a.Mock.Stats()
}
