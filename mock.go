package main

import (
	"crypto/md5"
	"encoding/gob"
	"encoding/hex"
	"fmt"
	"io"
	"log/slog"
	"os"
	"path/filepath"
	"slices"
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
}

func (i *Mock) init() *Mock {
	fh, err := os.Open(i.storePath())
	if err != nil && !os.IsNotExist(err) {
		panic(err)
	}

	if os.IsNotExist(err) {
		i.Games = map[GameId]GameInfo{
			1: {
				GameId:      1,
				Name:        "csgo",
				IconPath:    "./assets/images/images/csgo-icon.png",
				BgimgPath:   "./assets/images/images/csgo-bg.jpg",
				GameServers: []string{"美服", "欧服", "亚服"},
			},
			2: {
				GameId:      1,
				Name:        "地下城与勇士",
				IconPath:    "./assets/images/images/dnf-icon.png",
				BgimgPath:   "./assets/images/images/dnf-bg.jpg",
				GameServers: []string{"台服", "北美服", "日服"},
			},
			3: {
				GameId:    3,
				Name:      "绝地求生",
				IconPath:  "./assets/images/images/pubg-icon.png",
				BgimgPath: "./assets/images/images/pubg-bg.jpg",
			},
		}
		i.qrImgPath = "./assets/images/qr-alipay-img.png"
	} else {
		if err := gob.NewDecoder(fh).Decode(i); err != nil {
			panic(err)
		}
		i.AcceleratedGame = 0
	}

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
func (i *Mock) logined() bool { return i.User.Name == "" }

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

func (a *Mock) GetGame(gameId GameId) (GameInfo, Message) {
	if !a.logined() {
		return GameInfo{}, NotLogin.Message()
	} else if !a.gameIdValid(gameId) {
		return GameInfo{}, RequireGameId.Message()
	}

	return a.Games[gameId], OK.Message()
}

func (a *Mock) SearchGame(keyword string) (list []GameInfo, msg Message) {
	if !a.logined() {
		return nil, NotLogin.Message()
	}

	for _, e := range a.Games {
		list = append(list, e)
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

func (a *Mock) Stats() (s Stats) {
	return Stats{}
}
