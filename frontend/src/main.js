function getUser() {
    return window.go.main.App.GetUser()
}
function registerOrLogin(username, password) {
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