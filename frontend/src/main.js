

var app

window.onload = function () {
  app = new App()

  // var s = new Stats()
  // s.start()
  // document.getElementById('accelerate').style.display = 'flex'
  // var s = new Stats()
  // s.start()




  //   window.go.main.App.SearchGame("").then(res => {
  // 
  //     b.renderSearchGame("test", res.data)
  //   })

  //   window.go.main.App.GetGame(3).then(res => {
  // 
  //     app.renderGameDetail(res.data)
  //   })

}


class Base {
  constructor() {
    this.gamesItemTemplate = 'games-item-template'

    this.userInfoIcon = 'user-info-icon'
    this.userInfoNickname = 'user-info-nickname'
    this.userInfoSubscribe = 'user-info-subscribe'

    this.searchGame = 'search-game'
    this.searchGameInput = 'search-game-input'
    this.searchGameGridItemTemplate = 'search-game-grid-item-template'

    this.gameDetail = 'game-detail'
    this.gameDetailBg = 'game-detail-bg'
    this.gameDetailBodyTitleName = 'game-detail-body-title-name'
    this.gameDetailBodyTitleInfoRecent = 'game-detail-body-title-info-recent'
    this.gameDetailBodyTitleInfoTotalHours = 'game-detail-body-title-info-total-hours'
    this.gameDetailBodyTitleInfoTotalBytes = 'game-detail-body-title-info-total-bytes'
    this.gameDetailBodyAccSettingZone = 'game-detail-body-acc-setting-zone'
    this.gameDetailBodyAccRouteAuto = 'game-detail-body-acc-route-auto'
    this.gameDetailBodyAccRouteFix = 'game-detail-body-acc-route-fix'

    this.accelerate = 'accelerate'
    this.accelerateDetailImg = 'accelerate-detail-img'
    this.accelerateDetailBaseInfoTitleName = 'accelerate-detail-base-info-title-name'
    this.accelerateChartDelay = 'accelerate-chart-delay'
    this.accelerateChartLoss = 'accelerate-chart-loss'
    this.routeLinkGateway = 'route-link-gateway'
    this.routeLinkForward = 'route-link-forward'
    this.routeLinkServer = 'route-link-server'

    this.mainBodyIds = [this.searchGame, this.gameDetail, this.accelerate]
  }

  clearItem(templateItem) {
    var parent = templateItem.parentElement
    for (let i = parent.children.length - 1; i > 0; i--) {
      parent.removeChild(parent.lastChild)
    }
  }
  appendItem(templateItem, callback) {
    var e = templateItem.cloneNode(true)
    e.removeAttribute('id')
    if (e.style.length == 1) {
      e.removeAttribute('style')
    } else if (e.style.length > 1) {
      e.style.removeProperty('display')
    }

    callback(e)
    templateItem.parentElement.appendChild(e)
  }

  renderUserInfo(user) {
    document.getElementById(this.userInfoIcon).src = user.icon
    document.getElementById(this.userInfoNickname).innerText = user.name
    var t = new Date(user.expire * 1000)
    document.getElementById(this.userInfoSubscribe).innerText = t.getFullYear() + "年" + t.getMonth() + "月" + t.getDay() + "日"
  }
  renderGameList(infos) {
    var template = document.getElementById(this.gamesItemTemplate)
    this.clearItem(template)
    for (let i = 0; infos != null && i < infos.games.length; i++) {
      this.appendItem(
        template,
        function (e) {
          e.dataset.game = infos.games[i].game_id
          e.children[0].children[0].src = infos.games[i].icon_path
          e.children[1].innerText = infos.games[i].name

          if (infos.games[i].game_id == infos.select_game) {
            e.dataset.select = 'true'
            e.style.backgroundColor = 'rgb(64, 64, 64)'
            e.scrollIntoView({ behavior: 'smooth', block: 'nearest' })
          }
          if (infos.games[i].game_id == infos.accele_game) {
            e.dataset.accele = 'true'
            e.children[1].style.width = '5rem'
            e.children[2].style.display = 'flex'
            e.children[2].children[0].style.display = 'none'
            e.children[2].children[1].style.display = 'block'
          }
        }
      )
    }
  }

  renderSearchGame(key, list) {
    if (key != null && key.length > 0) {
      document.getElementById(this.searchGameInput).value = key
    }

    var template = document.getElementById(this.searchGameGridItemTemplate)
    this.clearItem(template)
    for (let i = 0; list != null && i < list.length; i++) {
      this.appendItem(
        template,
        function (e) {
          e.dataset.game = list[i].game_id
          e.children[0].src = list[i].adimg_path
          e.children[1].innerText = list[i].name
        }
      )
    }

    this.mainBodyIds.forEach(id => {
      id == this.searchGame ?
        document.getElementById(id).style.display = 'flex' :
        document.getElementById(id).style.display = 'none'
    })
  }

  renderGameDetail(info) {
    if (info == null) { return }
    document.getElementById(this.gameDetail).dataset.game = info.game_id
    document.getElementById(this.gameDetailBg).src = info.bgimg_path
    document.getElementById(this.gameDetailBodyTitleName).innerText = info.name

    document.getElementById(this.gameDetailBodyTitleInfoRecent).innerText = (function (date) {
      var year = date.getFullYear();
      var month = date.getMonth() + 1;
      var day = date.getDate();
      var hours = date.getHours();
      var minutes = date.getMinutes();
      var seconds = date.getSeconds();
      function pad(number) { return number < 10 ? '0' + number : number }

      return year + '年' + pad(month) + '月' + pad(day) + '日 ' +
        pad(hours) + ':' + pad(minutes) + ':' + pad(seconds);
    })(new Date(info.last_active * 1000))


    document.getElementById(this.gameDetailBodyTitleInfoTotalHours).innerText = (info.total_seconds / 3600).toFixed(1).toString()
    document.getElementById(this.gameDetailBodyTitleInfoTotalBytes).innerText = (function (bytes) {
      if (bytes === 0) return '0B';
      const k = 1024;
      const sizes = ['B', 'KB', 'MB', 'GB', 'TB', 'PB'];

      const i = Math.floor(Math.log(bytes) / Math.log(k));
      return parseFloat((bytes / Math.pow(k, i)).toFixed(1)) + sizes[i];
    })(info.total_bytes)

    var opt = document.getElementById(this.gameDetailBodyAccSettingZone)
    opt.innerHTML = ''
    info.game_servers.forEach(e => {
      opt.appendChild(new Option(e, e))
    })

    if (info.fix_route) {
      document.getElementById(this.gameDetailBodyAccRouteAuto).style.listStyleType = 'circle';
      document.getElementById(this.gameDetailBodyAccRouteFix).style.listStyleType = 'disc';
    } else {
      document.getElementById(this.gameDetailBodyAccRouteFix).style.listStyleType = 'circle';
      document.getElementById(this.gameDetailBodyAccRouteAuto).style.listStyleType = 'disc';
    }

    this.mainBodyIds.forEach(id => {
      id == this.gameDetail ?
        document.getElementById(id).style.display = 'flex' :
        document.getElementById(id).style.display = 'none'
    })
  }

  renderGameAccele(info, stats) {
    document.getElementById(this.accelerate).dataset.game = info.game_id
    document.getElementById(this.accelerateDetailImg).src = info.adimg_path
    document.getElementById(this.accelerateDetailBaseInfoTitleName).innerText = info.name

    this.mainBodyIds.forEach(id => {
      id == this.accelerate ?
        document.getElementById(id).style.display = 'flex' :
        document.getElementById(id).style.display = 'none'
    })
    setTimeout(500, stats.start())
  }
}



class App extends Base {
  constructor() {
    super()
    this.stats = null

    window.go.main.App.GetUser().then(res => {
      if (res.code != 0) {
        alert(res.msg) // todo: pop login/register
      } else {
        this.renderUserInfo(res.data)
      }
    })

    window.go.main.App.ListGames().then(res => {
      if (res.code != 0) {
        alert(res.msg) // todo: exit app
      } else if (res.data.games != null) {
        this.renderGameList(res.data)
        this.renderGameDetail(
          res.data.games.find(e => e.game_id == res.data.select_game)
        )
      }
    })
  }

  SearchGame(key) {
    window.go.main.App.SearchGame(key).then(res => {
      if (res.code != 0) {
        alert(res.msg)
      } else {
        this.renderSearchGame(key, res.data)
      }
    })
  }
  AddGame(gameId) {
    window.go.main.App.AddGame(gameId).then(res => {
      if (res.code != 0) {
        alert(res.msg)
      } else {
        this.renderGameList(res.data)
        this.renderGameDetail(
          res.data.games.find(e => e.game_id == gameId)
        )
      }
    })
  }
  DelGame(gameId, event) {
    event.stopPropagation();
    window.go.main.App.DelGame(gameId).then(res => {
      if (res.code != 0) {
        alert(res.msg)
      } else {
        this.renderGameList(res.data)
        this.renderGameDetail(
          res.data.games.find(e => e.game_id == res.data.select_game)
        )
      }
    })
  }
  SelectGame(gameId) {
    window.go.main.App.SelectGame(gameId).then(res => {
      if (res.code != 0) {
        alert(res.msg)
      } else {
        this.renderGameList(res.data)

        if (res.data.accele_game == gameId) {
          this.renderGameAccele(
            res.data.games.find(e => e.game_id == gameId),
            this.stats,
          )
        } else {
          this.renderGameDetail(
            res.data.games.find(e => e.game_id == gameId)
          )
        }
      }
    })
  }
  SetGameServer(server) {
    window.go.main.App.SetGameServer(server).then(res => {
      if (res.code != 0) { alert(res.msg) }
    })
  }
  SetRouteMode(fixRoute) {
    window.go.main.App.SetRouteMode(fixRoute).then(res => {
      if (res.code != 0) {
        alert(res.msg)
      } else {
        this.renderGameDetail(res.data)
      }
    })
  }
  StartAccele(gameId) {
    window.go.main.App.StartAccele(gameId).then(res => {
      if (res.code != 0) {
        alert(res.msg)
      } else {
        this.stats = new Stats()
        this.renderGameList(res.data)
        this.renderGameAccele(
          res.data.games.find(info => info.game_id == gameId), this.stats
        )
      }
    })
  }
  StopAccele(gameId) {
    window.go.main.App.StopAccele().then(res => {
      if (res.code != 0) {
        alert(res.msg)
      } else {
        this.stats.stop()
        this.renderGameList(res.data)
        this.renderGameDetail(
          res.data.games.find(info => info.game_id == gameId)
        )
      }
    })
  }
}


class Stats extends Base {
  constructor() {
    super()
    this.period = 3000 // 数据点间隔 3s
    this.count = 20
    this.color = ['rgb(5,115,205)', 'rgb(255,137,2)']
    this.delay = null
    this.loss = null
    this.started = false

    this.worker = new Worker('worker.js')
    this.worker.onmessage = function (ctx) {
      var ctx = ctx
      return function () {
        window.go.main.App.Stats().then(stats => {
          { // title message
            var f = function (v) { return v == null ? " ? ? " : v }
            stats.gateway_loc = f(stats.gateway_loc)
            stats.forward_loc = f(stats.forward_loc)
            stats.server_loc = f(stats.server_loc)

            document.getElementById(ctx.routeLinkGateway).innerText = stats.gateway_loc
            document.getElementById(ctx.routeLinkForward).innerText = stats.forward_loc
            document.getElementById(ctx.routeLinkServer).innerText = stats.server_loc
          }
          { // delay chart
            var opt = ctx.delay.getOption()
            opt.xAxis = {
              min: new Date(stats.stamp - ctx.count * ctx.period),
            }

            var f = function (v) {
              return v == null || v < 0.000001 || v > 300 ? null : v.toFixed(1)
            }
            opt.series[0].data.push([new Date(stats.stamp), f(stats.rtt_gateway), stats.gateway_loc])
            opt.series[1].data.push([new Date(stats.stamp), f(stats.rtt_forward), stats.forward_loc])
            ctx.delay.setOption(opt)
          }
          { // packet loss chart
            var opt = ctx.loss.getOption()
            opt.xAxis = {
              min: new Date(stats.stamp - ctx.count * ctx.period),
            }

            var f = function (v) {
              return v == null || v < 0.000001 || v > 100 ? null : v.toFixed(1)
            }
            opt.series[0].data.push([new Date(stats.stamp), f(stats.loss_client_uplink), stats.gateway_loc])
            opt.series[1].data.push([new Date(stats.stamp), f(stats.loss_client_downlink), stats.forward_loc])
            opt.series[2].data.push([new Date(stats.stamp), f(stats.loss_gateway_uplink), stats.gateway_loc])
            opt.series[3].data.push([new Date(stats.stamp), f(stats.loss_gateway_downlink), stats.forward_loc])
            ctx.loss.setOption(opt)
          }
        })
      }
    }(this)

  }

  start() {
    if (this.started) { return }
    this.started = true

    this.delay = (function (ctx) {
      var chart = echarts.init(document.getElementById(ctx.accelerateChartDelay))

      var option = {
        title: {
          text: '延迟 (ms)',
          textStyle: { color: "#fafafa", fontSize: '12', fontWeight: 'normal' },
          bottom: '0px', left: 'center',
        },
        color: ctx.color,
        tooltip: {
          trigger: 'axis',
          confine: true,
          formatter: function (params) {
            var time = params[0].data[0]
            var h = time.getHours().toString().padStart(2, '0')
            var m = time.getMinutes().toString().padStart(2, '0')
            var s = time.getSeconds().toString().padStart(2, '0')

            var gateway = params[0].data[1] === null ? 'null' : (params[0].data[1]).toString()
            var forward = params[1].data[1] === null ? 'null' : (params[1].data[1]).toString()
            var gateway_loc = params[0].data[2]
            var forward_loc = params[1].data[2]

            return `
              <div style="font-size:xx-small;">
                <div>${h}:${m}:${s}</div>
                <div style="color:${ctx.color[1]}">本机&nbsp;-&nbsp;${forward_loc}: ${forward}&nbsp;ms</div>
                <div style="color:${ctx.color[0]}">本机&nbsp;-&nbsp;${gateway_loc}: ${gateway}&nbsp;ms</div>
              </div>
            `
          },
        },
        xAxis: { type: 'time', show: false, },
        yAxis: {
          type: 'value',
          splitLine: { show: true, lineStyle: { color: '#525359' } },
          axisLabel: { formatter: '{value}ms' }
        },
        series: [],
        grid: { left: '50px', right: '6px', top: '15px', bottom: '30px', }
      };


      var series = ['gateway', 'forward']
      for (let i = 0; i < series.length; i++) {
        option.series.push(
          {
            name: series[i], type: 'line', animationDuration: 400,
            smooth: 0.4, data: [],
            lineStyle: { width: 1 },
            itemStyle: { opacity: 0.2 },
          }
        )
      }
      chart.setOption(option)
      return chart
    })(this)

    this.loss = (function (ctx) {
      var chart = echarts.init(document.getElementById(ctx.accelerateChartLoss));
      var names = ['上行丢包', '下行丢包']

      var option = {
        color: ctx.color,
        tooltip: {
          trigger: 'axis',
          confine: true,
          formatter: function (params) {
            var time = params[0].data[0]
            var h = time.getHours().toString().padStart(2, '0')
            var m = time.getMinutes().toString().padStart(2, '0')
            var s = time.getSeconds().toString().padStart(2, '0')

            var loss1 = params[0].data[1] === null ? 'null' : (params[0].data[1]).toString()
            var loss2 = params[1].data[1] === null ? 'null' : (params[1].data[1]).toString()
            var gateway_loc = params[0].data[2]
            var forward_loc = params[1].data[2]

            if (params[0].seriesName == names[0]) {
              return `
                <div style="font-size:xx-small;">
                  <div>${h}:${m}:${s}</div>
                  <div style="color:${ctx.color[0]}">本机&nbsp;-&nbsp;${gateway_loc}: ${loss1}%</div>
                  <div style="color:${ctx.color[1]}">${gateway_loc}&nbsp;-&nbsp;${forward_loc}: ${loss2}%</div>
                </div>
              `
            } else {
              return `
                <div style="font-size:xx-small;">
                  <div>${h}:${m}:${s}</div>
                  <div style="color:${ctx.color[1]}">${forward_loc}&nbsp;-&nbsp;${gateway_loc}: ${loss2}&nbsp;%</div>
                  <div style="color:${ctx.color[0]}">${gateway_loc}&nbsp;-&nbsp;本机: ${loss1}&nbsp;%</div>
                </div>
              `
            }
          },
        },
        legend: {
          type: 'plain', data: names,
          bottom: '0px', icon: 'none',
          textStyle: { color: '#fafafa', fontSize: '12' },
          selected: { "下行丢包": false },
          selectedMode: 'single',
        },

        xAxis: { type: 'time', show: false, },
        yAxis: {
          type: 'value', min: 0, //max: 100,
          splitLine: { show: true, lineStyle: { color: '#525359' } },
          axisLabel: { formatter: '{value}%' }
        },
        series: [],
        grid: { left: '50px', right: '6px', top: '15px', bottom: '30px', }
      };

      for (let i = 0; i < names.length; i++) {
        for (let j = 0; j < 2; j++) {
          option.series.push(
            {
              name: names[i], type: 'line', animationDuration: 400,
              smooth: 0.4, data: [],
              lineStyle: { width: 1 },
              itemStyle: { opacity: 0.2, color: ctx.color[j] },
            }
          )
        }
      }

      chart.setOption(option)
      return chart
    })(this)

    this.worker.postMessage(this.period)
  }
  stop() {
    if (this.started) {
      this.worker.postMessage(0)
      this.started = false
    }
  }
}



// 调用方法
// var stats = new Stats()
// worker = new Worker('worker.js');
// worker.onmessage = function (e) {
//   var rtt_gateay = 400
//   var loss_g_d = 0
//   if ((new Date().getSeconds() % 10 != 9)) {
//     rtt_gateay = Math.random() * 100
//     loss_g_d = Math.random()
//   }
//
//   stats.push({
//     time: (new Date()).getTime(),
//     gateway_loc: '北京',
//     forward_loc: '莫斯科',
//     server_loc: '圣彼得堡',
//
//     rtt_gateway: rtt_gateay,
//     rtt_forward: (100 + Math.random() * 50),
//
//     loss_client_uplink: Math.random() * 3,
//     loss_client_downlink: Math.random() * 1,
//     loss_gateway_uplink: Math.random(),
//     loss_gateway_downlink: loss_g_d,
//   })
// }
// worker.postMessage(stats.period) // start