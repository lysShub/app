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
	SelectedGame GameId // store GameId
	AcceleedGame GameId // store GameId
	AddedGames   []GameId

	User  UserInfo
	Games map[GameId]GameInfo

	// temp
	qrImgPath string
	statsTime time.Time
}

func (i *Mock) init() *Mock {
	fh, err := os.Open(i.storePath())
	if err != nil && !os.IsNotExist(err) {
		panic(err)
	}
	fmt.Println("store path", i.storePath())

	if os.IsNotExist(err) {
		i.User = UserInfo{
			Name:     "test",
			Password: "1234567",
			Phone:    "13448612544",
			Icon:     "./assets/images/logo-universal.png",
			Expire:   time.Now().AddDate(1, 0, 0).Unix(),
		}

		i.Games = map[GameId]GameInfo{
			1: {
				GameId:       1,
				Name:         "反恐精英",
				IconPath:     "./assets/images/images/csgo-icon.png",
				BgimgPath:    "./assets/images/images/csgo-bg.png",
				AdimgPath:    "./assets/images/images/csgo-ad.jpg",
				GameServers:  []string{"美服", "欧服", "亚服"},
				LastActive:   time.Now().Local().AddDate(0, 0, 0).Unix(),
				TotalSeconds: 5645,
				TotalBytes:   145234,
			},
			2: {
				GameId:       2,
				Name:         "地下城与勇士阿巴阿巴阿巴阿巴",
				IconPath:     "./assets/images/images/dnf-icon.png",
				BgimgPath:    "./assets/images/images/dnf-bg.jpg",
				AdimgPath:    "./assets/images/images/dnf-ad.jpg",
				GameServers:  []string{"台服", "北美服", "日服"},
				LastActive:   time.Now().Local().AddDate(-1, 0, 0).Unix(),
				TotalSeconds: 13435,
				TotalBytes:   6733412,
			},
			3: {
				GameId:       3,
				Name:         "绝地求生",
				IconPath:     "./assets/images/images/pubg-icon.png",
				BgimgPath:    "./assets/images/images/pubg-bg.jpg",
				AdimgPath:    "./assets/images/images/pubg-ad.png",
				GameServers:  []string{"国际服", "国服", "日服"},
				LastActive:   time.Now().Local().AddDate(-1, -2, 0).Unix(),
				TotalSeconds: 3452,
				TotalBytes:   457234,
			},
		}
		i.AddedGames = []GameId{
			2, 3,
		}
		i.SelectedGame = i.AddedGames[0]
	} else {
		if err := gob.NewDecoder(fh).Decode(i); err != nil {
			panic(err)
		}
		i.AcceleedGame = 0
	}

	i.qrImgPath = "./assets/images/qr-alipay-img.png"
	i.statsTime = time.Time{}
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

func (a *Mock) GetUser() (msg Message) {
	if !a.logined() {
		return NotLogin.Message(nil)
	} else {
		return OK.Message(a.User)
	}
}

func (a *Mock) RegisterOrLogin(user, pwd string) (msg Message) {
	if a.logined() {
		return IsLogined.Message(nil)
	} else {
		a.User.Name = user
		a.User.Password = pwd
		a.User.Phone = "1378494xxxxx"
		a.User.Expire = 0
		slog.Info("log", slog.String("user", user), slog.String("pwd", pwd))

		a.sync()
		return OK.Message(nil)
	}
}

func (a *Mock) Recharge(months int, callback func(Message)) (msg Message) {
	if months <= 0 {
		return InvalidMonths.Message("")
	} else if !a.logined() {
		return NotLogin.Message("")
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
	return OK.Message(a.qrImgPath)
}

func (a *Mock) ListGames() (msg Message) {
	if !a.logined() {
		return NotLogin.Message(nil)
	}

	var gs = GameInfos{SelectGame: a.SelectedGame, AcceleGame: a.AcceleedGame}
	for _, id := range a.AddedGames {
		gs.Games = append(gs.Games, a.Games[id])
	}
	return OK.Message(gs)
}

func (a *Mock) GetGame(id GameId) Message {
	if !a.logined() {
		return NotLogin.Message(nil)
	} else if !a.gameIdValid(id) {
		return RequireGameId.Message(nil)
	}
	return OK.Message(a.Games[id])
}

func (a *Mock) SelectGame(id GameId) Message {
	if !a.logined() {
		return NotLogin.Message(nil)
	} else if !a.gameIdValid(id) {
		return RequireGameId.Message(nil)
	} else if !slices.Contains(a.AddedGames, id) {
		return Notfound.Message(nil)
	}

	a.SelectedGame = id
	a.sync()
	return a.ListGames()
}

func (a *Mock) GetSelectedGame() Message {
	if !a.logined() {
		return NotLogin.Message(nil)
	}

	if a.SelectedGame == 0 {
		return NotSelectGame.Message(nil)
	} else {
		return OK.Message(a.Games[a.SelectedGame])
	}
}

func (a *Mock) SearchGame(keyword string) (msg Message) {
	if !a.logined() {
		return NotLogin.Message(nil)
	}

	var list []GameInfo
	for _, e := range a.Games {
		if !slices.Contains(a.AddedGames, e.GameId) &&
			(keyword == "" || strings.ContainsAny(e.Name, keyword)) {

			list = append(list, e)
		}
	}

	return OK.Message(list)
}

func (a *Mock) AddGame(gameId GameId) Message {
	if !a.logined() {
		return NotLogin.Message(nil)
	} else if !a.gameIdValid(gameId) {
		return RequireGameId.Message(nil)
	} else if slices.Contains(a.AddedGames, gameId) {
		return GameExist.Message(nil)
	}

	a.AddedGames = append(a.AddedGames, gameId)
	a.SelectedGame = gameId
	a.sync()
	return a.ListGames()
}

func (a *Mock) DelGame(gameId GameId) Message {
	if !a.logined() {
		return NotLogin.Message(nil)
	} else if !a.gameIdValid(gameId) {
		return RequireGameId.Message(nil)
	} else if !slices.Contains(a.AddedGames, gameId) {
		return GameNotExist.Message(nil)
	}

	a.AddedGames = slices.DeleteFunc(a.AddedGames, func(id GameId) bool { return id == gameId })
	if a.SelectedGame == gameId {
		if len(a.AddedGames) > 0 {
			a.SelectedGame = a.AddedGames[0]
		} else {
			a.SelectedGame = 0
		}
	}

	a.sync()
	return a.ListGames()
}

func (a *Mock) SetGameServer(gameServer string) Message {
	if !a.logined() {
		return NotLogin.Message(nil)
	} else if !a.gameIdValid(a.SelectedGame) {
		return RequireGameId.Message(nil)
	}

	g := a.Games[a.SelectedGame]

	i := slices.Index(g.GameServers, gameServer)
	if i == -1 {
		panic(fmt.Sprintf("unknown game %s server %s", g.Name, gameServer))
	} else if i != 0 {
		g.GameServers[i], g.GameServers[0] = g.GameServers[0], g.GameServers[i]
		a.Games[a.SelectedGame] = g
		a.sync()
	}

	return OK.Message(nil)
}

func (a *Mock) SetRouteMode(fixRoute bool) Message {
	if !a.logined() {
		return NotLogin.Message(nil)
	} else if !a.gameIdValid(a.SelectedGame) {
		return GameExist.Message(nil)
	}

	g := a.Games[a.SelectedGame]
	if g.FixRoute != fixRoute {
		g.FixRoute = fixRoute
		a.Games[a.SelectedGame] = g
		a.sync()
	}
	return OK.Message(g)
}

func (a *Mock) StartAccele(gameId GameId) Message {
	if !a.logined() {
		return NotLogin.Message(nil)
	} else if a.User.Expire < time.Now().Unix() {
		return VIPExpired.Message(nil)
	} else if !a.gameIdValid(gameId) {
		return RequireGameId.Message(nil)
	} else if a.AcceleedGame != 0 {
		return Message{Code: Accelerating, Msg: fmt.Sprintf("%s 正在加速", a.Games[a.AcceleedGame].Name)}
	}

	a.AcceleedGame = gameId

	info := a.Games[gameId]
	slog.Info("开始加速", slog.String("name", info.Name), slog.String("server", info.GameServers[0]), slog.Bool("fix-route", info.FixRoute))

	a.sync()
	return a.ListGames()
}

func (a *Mock) StopAccele() Message {
	if !a.logined() {
		return NotLogin.Message(nil)
	} else if a.AcceleedGame == 0 {
		return NotAccelerated.Message(nil)
	}

	a.AcceleedGame = 0
	a.sync()
	return a.ListGames()
}

var RandNew *rand.Rand = rand.New(rand.NewSource(time.Now().UnixNano()))

func (a *Mock) Stats() (s Stats) {
	time.Sleep(time.Until(a.statsTime.Add(time.Second * 3)))
	a.statsTime = time.Now()

	return randStats()
}

func randStats() Stats {
	var s = Stats{
		Stamp:   time.Now().UnixMilli(),
		Seconds: 8234,
		Bytes:   345254,

		GatewayLoc: "北京",
		ForwardLoc: "莫斯科",
		ServerLoc:  "圣彼得堡",
	}

	s.RttGateway = 30 + time.Duration(rand.Intn(30))
	s.RttForward = s.RttGateway + 90

	s.LossClientUplink = f(1.5 + rand.Float64()*3)
	s.LossClientDownlink = f(rand.Float64())
	s.LossGatewayUplink = f(rand.Float64())
	s.LossGatewayDownlink = f(rand.Float64())
	return s
}

func f(v float64) float64 { return math.Round(v*100) / 100 }
