package main

import (
	"context"
	"time"

	"github.com/wailsapp/wails/v2/pkg/runtime"
)

type ApiResponse struct {
	Code MsgCode `json:"code"`
	Msg  string  `json:"msg"`
	Data any     `json:"data"`
}
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
	Expire   int64  `json:"expire"` // utc 时间戳
}

// GetUser 获取用户信息, 应用渲染完成即调用此函数, 如果msg.Code==NotLogin, 则弹出注册登录页面

func (a *App) GetUser() ApiResponse {
	user, message := a.Mock.GetUser()
	return ApiResponse{
		Code: message.Code,
		Msg:  message.Msg,
		Data: user,
	}
}

// RegisterOrLogin 注册或登录,
func (a *App) RegisterOrLogin(user, pwd string) (msg Message) {
	return a.Mock.RegisterOrLogin(user, pwd)
}

// todo: 暂时不考虑
// Recharge 充值，返回一个字符二维码、和一个全局事件。参考 https://wails.io/zh-Hans/docs/reference/runtime/events
// 回调返回结果是Message类型
func (a *App) Recharge(months int, eventName string) ApiResponse {
	path, message := a.Mock.Recharge(months, func(m Message) {
		runtime.EventsEmit(a.ctx, eventName, m)
	})
	return ApiResponse{
		Code: message.Code,
		Msg:  message.Msg,
		Data: path,
	}
}

// ListGames 获取已添加的游戏列表, selectedIdx 表示默认应该选中的游戏

func (a *App) ListGames() ApiResponse {
	list, idx, msg := a.Mock.ListGames()
	return ApiResponse{
		Code: msg.Code,
		Msg:  msg.Msg,
		Data: struct {
			List        []GameInfo `json:"list"`
			SelectedIdx int        `json:"selected_idx"`
		}{
			List:        list,
			SelectedIdx: idx,
		},
	}
}

// SelectGame 选中某个游戏
func (a *App) SelectGame(gameId GameId) ApiResponse {
	game, message := a.Mock.SelectGame(gameId)
	return ApiResponse{
		Code: message.Code,
		Msg:  message.Msg,
		Data: game,
	}
}

// GetSelectedGame 获取当前选中的游戏
func (a *App) GetSelectedGame() ApiResponse {
	game, message := a.Mock.GetSelectedGame()
	return ApiResponse{
		Code: message.Code,
		Msg:  message.Msg,
		Data: game,
	}
}

// SearchGame 根据关键字搜索游戏
func (a *App) SearchGame(keyword string) ApiResponse {
	game, message := a.Mock.SearchGame(keyword)
	return ApiResponse{
		Code: message.Code,
		Msg:  message.Msg,
		Data: game,
	}
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

type StatsList struct {
	List []Stats `json:"list"`
}

type Stats struct {
	Stamp   int    `json:"stamp"`   // 本数据点对应的时间戳
	Gateway string `json:"gateway"` // gateway所在城市名
	Forward string `json:"forward"` // forward所在城市名
	Server  string `json:"server"`  // server所在城市名

	PingGateway time.Duration `json:"ping_gateway"` // 到gateway的延时
	PingForward time.Duration `json:"ping_forward"` // 到forward的延时

	Uplink   Loss `json:"uplink"`   // 上行丢包
	Donwlink Loss `json:"donwlink"` // 下行丢包
	Total    Loss `json:"total"`    // 合计丢包
}

type Loss struct {
	Gateway float64 `json:"gateway"` // client-gateway 的丢包
	Forward float64 `json:"forward"` // gateway-forward的丢包
}

// Stats 获取统计信息, 阻塞函数, 如果距上次调用时间短于3s, 会主动阻塞直到恰好相距3s
func (a *App) Stats() StatsList {
	return a.Mock.Stats()
}
