package main

import (
	"context"
	"fmt"
	"log/slog"
	"slices"
	"sync/atomic"
	"time"

	"github.com/wailsapp/wails/v2/pkg/runtime"
)

type ApiResponse struct {
	Code MsgCode   `json:"code"`
	Msg  string    `json:"msg"`
	data *UserInfo `json:"data"`
}
type Message struct {
	Code MsgCode `json:"code"`
	Msg  string  `json:"msg"`
}

//go:generate stringer -linecomment -output acc_gen.go -type=MsgCode
type MsgCode int

func (c MsgCode) Message() Message { return Message{Code: c, Msg: c.String()} }
func (c MsgCode) TSName() string {
	if c >= _end {
		panic(c)
	}
	return codeLits[c]
}

const (
	OK           MsgCode = iota // ok
	Notfound                    // not found
	NotLogin                    // 没有登录
	IsLogined                   // 已经登录
	VIPExpired                  // vip 已过期
	NotSetGame                  // 请先选择加速的游戏
	Accelerating                // 游戏已在加速
	_end
)

var codeLits = []string{"OK", "Notfound", "NotLogin", "IsLogined", "VIPExpired", "NotSetGame", "Accelerating"}
var codes = []MsgCode{OK, Notfound, NotLogin, IsLogined, VIPExpired, NotSetGame, Accelerating}

type App struct {
	ctx context.Context

	login                     atomic.Bool
	username, password, phone string
	expire                    time.Time

	game             atomic.Int32 // store GameId
	games            map[GameId]GameInfo
	addedGames       []GameInfo
	cacheSelectedIdx int
}

func NewMockApp() *App { return &App{} }

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
	if !a.login.Load() {
		return ApiResponse{
			Code: NotLogin,
			Msg:  NotLogin.String(),
			data: &UserInfo{
				Name: "123",
			},
		}
	} else {
		return ApiResponse{
			Code: OK,
			Msg:  OK.String(),
			data: &UserInfo{
				Name:     a.username,
				Password: a.password,
				Phone:    a.phone,
				Expire:   a.expire.Unix(),
			},
		}
	}
}

// RegisterOrLogin 注册或登录,
func (a *App) RegisterOrLogin(user, pwd string) (msg Message) {
	if a.login.Swap(true) {
		return IsLogined.Message()
	} else {
		a.username = user
		a.password = pwd
		a.phone = "1378494xxxxx"
		a.expire = time.Now().Add(time.Hour * 24)
		slog.Info("log", slog.String("user", user), slog.String("pwd", pwd))
		return OK.Message()
	}
}

// todo: 暂时不考虑
// Recharge 充值，返回一个字符二维码、和一个全局事件。参考 https://wails.io/zh-Hans/docs/reference/runtime/events
// 回调返回结果是Message类型
func (a *App) Recharge(months int, eventName string) (qrImagePath string) {
	go func() {
		// 等待支付结果
		runtime.EventsEmit(context.Background(), eventName, Message{Code: OK, Msg: "支付成功"})
	}()

	return "./xx/xx/qr.png"
}

type GameId = int32 // 最小值为1

type GameInfo struct {
	GameId      GameId   `json:"game_id"`
	Name        string   `json:"name"`
	IconPath    string   `json:"icon_path"`
	BgimgPath   string   `json:"bgimg_path"`
	GameServers []string `json:"game_servers"`

	// 表示上次代理配置
	CacheGameServer string `json:"cache_game_server"`
	CacheFixRoute   bool   `json:"cache_fix_route"`
}

// ListGames 获取已添加的游戏列表, selectedIdx 表示默认应该选中的游戏
func (a *App) ListGames() (list []GameInfo, selectedIdx int, msg Message) {
	if !a.login.Load() {
		return nil, -1, NotLogin.Message()
	}
	return a.addedGames, a.cacheSelectedIdx, OK.Message()
}

// SearchGame 根据关键字搜索游戏
func (a *App) SearchGame(keyword string) (list []GameInfo, msg Message) {
	for _, e := range a.games {
		list = append(list, e)
	}
	return list, OK.Message()
}

// AddGame 新增游戏
func (a *App) AddGame(gameId GameId) Message {
	info, has := a.games[gameId]
	if !has {
		invalidGameID(gameId)
	}

	a.addedGames = append(a.addedGames, info)
	return OK.Message()
}

// SetGame 选择某个游戏
func (a *App) SetGame(gameId GameId) Message {
	if _, has := a.games[gameId]; has {
		invalidGameID(gameId)
	}
	a.game.Store(gameId)
	return OK.Message()
}

func invalidGameID(id GameId) {
	panic(fmt.Sprintf("not found game id %d", id))
}

// SetGameServer 设置游戏区服
func (a *App) SetGameServer(id GameId, gameServer string) Message {
	var idx int = -1
	var info GameInfo
	for i, e := range a.addedGames {
		if e.GameId == id {
			idx, info = i, e
			break
		}
	}
	if idx == -1 {
		invalidGameID(id)
	}
	if !slices.Contains(info.GameServers, gameServer) {
		return Message{Code: Notfound, Msg: fmt.Sprintf("不支持 %s 区服：%s", info.Name, gameServer)}
	}
	a.addedGames[idx].CacheGameServer = gameServer

	return OK.Message()
}

// SetRouteMode 选择路由模式
func (a *App) SetRouteMode(id GameId, fixRoute bool) Message {
	var idx int = -1
	for i, e := range a.addedGames {
		if e.GameId == id {
			idx = i
			break
		}
	}
	if idx == -1 {
		invalidGameID(id)
	}
	a.addedGames[idx].CacheFixRoute = fixRoute

	return OK.Message()
}

// Accelerate 开始加速
func (a *App) Accelerate(id GameId) Message {
	if !a.game.CompareAndSwap(0, id) {
		info, has := a.games[a.game.Load()]
		if !has {
			invalidGameID(a.game.Load())
		}
		return Message{Code: Accelerating, Msg: fmt.Sprintf("%s 正在加速", info.Name)}
	}

	info, has := a.games[id]
	if !has {
		invalidGameID(id)
	}
	slog.Info("开始加速", slog.String("name", info.Name), slog.String("server", info.CacheGameServer), slog.Bool("fix-route", info.CacheFixRoute))

	return OK.Message()
}

// DisableAccelerate 停止加速
func (a *App) DisableAccelerate() Message {
	a.game.Swap(0)
	return OK.Message()
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

// Stats 获取统计信息
func (a *App) Stats() (s Stats) {
	return Stats{}
}
