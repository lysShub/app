package main

import (
	"crypto/md5"
	"encoding/gob"
	"encoding/hex"
	"fmt"
	"io"
	"log/slog"
	"math"
	"math/rand"
	"os"
	"path/filepath"
	"slices"
	"strings"
	"time"
)

type Mock struct {
	SelectedGame    GameId // store GameId
	AcceleratedGame GameId // store GameId
	AddedGames      []GameId

	User  UserInfo
	Games map[GameId]GameInfo

	// temp
	qrImgPath string
	statsTime time.Time
	head      *Heap[Stats]
}

func (i *Mock) init() *Mock {
	fh, err := os.Open(i.storePath())
	if err != nil && !os.IsNotExist(err) {
		panic(err)
	}
	fmt.Println("store path", i.storePath())

	if os.IsNotExist(err) {
		i.Games = map[GameId]GameInfo{
			1: {
				GameId:      1,
				Name:        "csgo",
				IconPath:    "./assets/images/images/csgo-icon.png",
				BgimgPath:   "./assets/images/images/csgo-bg.jpg",
				GameServers: []string{"美服", "欧服", "亚服"},
				LastActive:  time.Now().Local().AddDate(0, 0, 0).Unix(),
				Duration:    1000,
				Flow:        1000,
			},
			2: {
				GameId:      1,
				Name:        "地下城与勇士",
				IconPath:    "./assets/images/images/dnf-icon.png",
				BgimgPath:   "./assets/images/images/dnf-bg.jpg",
				GameServers: []string{"台服", "北美服", "日服"},
				LastActive:  time.Now().Local().AddDate(-1, 0, 0).Unix(),
				Duration:    3000,
				Flow:        20000,
			},
			3: {
				GameId:      3,
				Name:        "绝地求生",
				IconPath:    "./assets/images/images/pubg-icon.png",
				BgimgPath:   "./assets/images/images/pubg-bg.jpg",
				GameServers: []string{"国际服", "国服", "日服"},
				LastActive:  time.Now().Local().AddDate(-1, -2, 0).Unix(),
				Duration:    2000,
				Flow:        33000,
			},
		}
	} else {
		if err := gob.NewDecoder(fh).Decode(i); err != nil {
			panic(err)
		}
		i.AcceleratedGame = 0
	}

	i.qrImgPath = "./assets/images/qr-alipay-img.png"
	i.statsTime = time.Time{}
	i.head = NewHeap[Stats]()
	return i
}

func (i *Mock) storePath() (path string) {
	if p, err := os.Executable(); err != nil {
		panic(err)
	} else {
		fh, err := os.Open(p)
		if err != nil {
			panic(err)
		}
		defer fh.Close()

		h := md5.New()
		io.Copy(h, fh)
		path = filepath.Join(os.TempDir(), fmt.Sprintf("app-%s.bin", hex.EncodeToString(h.Sum(nil))))
		return path
	}
}

func (i *Mock) sync() {
	path := i.storePath()
	w, err := os.Create(path)
	if err != nil {
		panic(err)
	}
	defer w.Close()

	if err := gob.NewEncoder(w).Encode(i); err != nil {
		panic(err)
	}
}
func (i *Mock) gameIdValid(id GameId) bool {
	_, has := i.Games[id]
	return has
}
func (i *Mock) logined() bool { return i.User.Name != "" }

func (a *Mock) GetUser() (info UserInfo, msg Message) {
	if !a.logined() {
		return UserInfo{}, NotLogin.Message()
	} else {
		return a.User, OK.Message()
	}
}

func (a *Mock) RegisterOrLogin(user, pwd string) (msg Message) {
	if a.logined() {
		return IsLogined.Message()
	} else {
		a.User.Name = user
		a.User.Password = pwd
		a.User.Phone = "1378494xxxxx"
		a.User.Expire = 0
		slog.Info("log", slog.String("user", user), slog.String("pwd", pwd))

		a.sync()
		return OK.Message()
	}
}

func (a *Mock) Recharge(months int, callback func(Message)) (qrImagePath string, msg Message) {
	if months <= 0 {
		return "", InvalidMonths.Message()
	} else if !a.logined() {
		return "", NotLogin.Message()
	}

	go func() {
		time.Sleep(time.Second * 3)
		if time.Now().UnixNano()%2 == 0 {
			t := time.Now().Local().AddDate(0, months, 0)
			a.User.Expire = t.Unix()
			a.sync()
			callback(Message{Code: OK, Msg: fmt.Sprintf("支付成功, 有效期至 %s", t.Format("2006-01-02T15:04:05"))})
		} else {
			callback(Message{Code: Unknown, Msg: "支付失败, xxxx"})
		}
	}()
	return a.qrImgPath, OK.Message()
}

func (a *Mock) ListGames() (list []GameInfo, selectedIdx int, msg Message) {
	if !a.logined() {
		return nil, -1, NotLogin.Message()
	}

	for i, id := range a.AddedGames {
		list = append(list, a.Games[id])
		if id == a.SelectedGame {
			selectedIdx = i
		}
	}
	return list, selectedIdx, OK.Message()
}

func (a *Mock) SelectGame(gameId GameId) (GameInfo, Message) {
	if !a.logined() {
		return GameInfo{}, NotLogin.Message()
	} else if !a.gameIdValid(gameId) {
		return GameInfo{}, RequireGameId.Message()
	}

	return a.Games[gameId], OK.Message()
}

func (a *Mock) GetSelectedGame() (GameInfo, Message) {
	if !a.logined() {
		return GameInfo{}, NotLogin.Message()
	}

	if a.SelectedGame == 0 {
		return GameInfo{}, NotSelectGame.Message()
	} else {
		return a.Games[a.SelectedGame], OK.Message()
	}

}

func (a *Mock) SearchGame(keyword string) (list []GameInfo, msg Message) {
	if !a.logined() {
		return nil, NotLogin.Message()
	}

	if keyword == "" {
		return nil, OK.Message()
	}

	for _, e := range a.Games {
		if strings.ContainsAny(e.Name, keyword) {
			list = append(list, e)
		}
	}
	return list, OK.Message()
}

func (a *Mock) AddGame(gameId GameId) Message {
	if !a.logined() {
		return NotLogin.Message()
	} else if !a.gameIdValid(gameId) {
		return RequireGameId.Message()
	}

	if slices.Contains(a.AddedGames, gameId) {
		return GameExist.Message()
	}

	a.AddedGames = append(a.AddedGames, gameId)
	a.sync()
	return OK.Message()
}

func (a *Mock) SetGame(gameId GameId) Message {
	if !a.logined() {
		return NotLogin.Message()
	} else if !a.gameIdValid(gameId) {
		return RequireGameId.Message()
	}

	a.SelectedGame = gameId
	a.sync()
	return OK.Message()
}

func (a *Mock) SetGameServer(gameId GameId, gameServer string) Message {
	if !a.logined() {
		return NotLogin.Message()
	} else if !a.gameIdValid(gameId) {
		return RequireGameId.Message()
	}

	g := a.Games[gameId]

	i := slices.Index(g.GameServers, gameServer)
	if i == -1 {
		panic(fmt.Sprintf("unknown game %s server %s", g.Name, gameServer))
	}
	g.CacheGameServer = gameServer

	a.Games[gameId] = g
	a.sync()
	return OK.Message()
}

func (a *Mock) SetRouteMode(gameId GameId, fixRoute bool) Message {
	if !a.logined() {
		return NotLogin.Message()
	} else if !a.gameIdValid(gameId) {
		return GameExist.Message()
	}

	g := a.Games[gameId]
	g.CacheFixRoute = fixRoute
	a.Games[gameId] = g

	a.sync()
	return OK.Message()
}

func (a *Mock) Accelerate(gameId GameId) Message {
	if !a.logined() {
		return NotLogin.Message()
	} else if a.User.Expire < time.Now().Unix() {
		return VIPExpired.Message()
	} else if !a.gameIdValid(gameId) {
		return RequireGameId.Message()
	} else if a.AcceleratedGame != 0 {
		return Message{Code: Accelerating, Msg: fmt.Sprintf("%s 正在加速", a.Games[a.AcceleratedGame].Name)}
	}

	a.AcceleratedGame = gameId

	info := a.Games[gameId]
	slog.Info("开始加速", slog.String("name", info.Name), slog.String("server", info.CacheGameServer), slog.Bool("fix-route", info.CacheFixRoute))

	a.sync()
	return OK.Message()
}

func (a *Mock) DisableAccelerate() Message {
	if !a.logined() {
		return NotLogin.Message()
	} else if a.AcceleratedGame == 0 {
		return NotAccelerated.Message()
	}

	a.AcceleratedGame = 0
	a.sync()
	return OK.Message()
}

var RandNew *rand.Rand = rand.New(rand.NewSource(time.Now().UnixNano()))

func (a *Mock) Stats() (s StatsList) {
	time.Sleep(time.Until(a.statsTime.Add(time.Second * 3)))
	a.statsTime = time.Now()

	a.head.Put(randStats())
	return StatsList{List: a.head.List()}
}

func randStats() Stats {
	var s = Stats{
		Stamp:   int(time.Now().Unix()),
		Gateway: "北京",
		Forward: "莫斯科",
		Server:  "圣彼得堡",
	}

	s.PingGateway = 60 + time.Duration(rand.Intn(60))
	s.PingForward = s.PingGateway + 90

	s.Uplink = Loss{
		Gateway: f(1.5 + rand.Float64()*3),
		Forward: f(rand.Float64()),
	}
	s.Donwlink = Loss{
		Gateway: f(0.5 + rand.Float64()),
		Forward: f(rand.Float64()),
	}
	s.Total = Loss{
		Gateway: f(s.Uplink.Gateway + s.Donwlink.Gateway),
		Forward: f(s.Uplink.Forward + s.Donwlink.Forward),
	}
	return s
}

func f(v float64) float64 { return math.Round(v*100) / 100 }

type Heap[T Stats | int] struct {
	cache []T
	i     int // head
}

func NewHeap[T Stats | int]() *Heap[T] {
	return &Heap[T]{
		cache: make([]T, 60),
	}
}

func (h *Heap[T]) Put(s T) {
	h.cache[h.i] = s

	n := len(h.cache)
	h.i = (h.i + 1) % n
}

func (h *Heap[T]) List() []T {
	var ss = make([]T, 0, len(h.cache))

	ss = append(ss, h.cache[h.i:]...)
	ss = append(ss, h.cache[:h.i]...)
	return ss
}
