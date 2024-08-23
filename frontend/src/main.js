
function GetUser() {
    return window.go.main.App.GetUser()
}
function RegisterOrLogin(username, password) {
    return window.go.main.App.RegisterOrLogin(username, password)
}
function ListGames() {
    return window.go.main.App.ListGames()
}
function SearchGame(game) {
    return window.go.main.App.SearchGame(game)
}

function AddGame(game) {
    return window.go.main.App.AddGame(game)
}
function SelectGame(game) {
    return window.go.main.App.SelectGame(game)
}
function Accelerate(game) {
    return window.go.main.App.Accelerate(game)
}
function DisableAccelerate() {
    return window.go.main.App.DisableAccelerate()
}
function Recharge(arg1, arg2) {
    return window.go.main.App.Recharge(arg1, arg2)
}

function EventListener(eventName, callback) {
    return EventsOn(eventName, callback)
}
function GetStats() {
    return window.go.main.App.Stats()
}