package game

var shouldQuit bool

func ShouldQuit() bool {
	return shouldQuit
}

func TryQuit() {
	if scene.TryQuit() {
		shouldQuit = true
	}
}

func Quit() {
	shouldQuit = true
}
