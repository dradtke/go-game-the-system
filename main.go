package main

// TODO:
//  1. Implement double-pass updates
//  2. Come up with a consistent scheme for rounding

import (
	"flag"
	"game"
	scenes "scenes/menu"
	"github.com/dradtke/go-allegro/allegro"
	"github.com/dradtke/go-allegro/allegro/font"
	"github.com/dradtke/go-allegro/allegro/image"
	"os"
	"runtime"
	"runtime/pprof"
	"time"
)

const FPS = 60

var (
	cpuprofile = flag.String("cpuprofile", "", "output a cpuprofile to...")
	memprofile = flag.String("memprofile", "", "output a memprofile to...")
)

func main() {
	runtime.LockOSThread()
	defer runtime.UnlockOSThread()

	flag.Parse()

	var (
		display        *allegro.Display
		eventQueue     *allegro.EventQueue
		fpsTimer       *allegro.Timer
		fpsTimerSource *allegro.EventSource
		err            error
	)

	// Event Queue
	if eventQueue, err = allegro.CreateEventQueue(); err == nil {
		defer eventQueue.Destroy()
	} else {
		panic(err)
	}

	// Display
	allegro.SetNewDisplayFlags(allegro.WINDOWED)
	if display, err = allegro.CreateDisplay(640, 480); err == nil {
		defer display.Destroy()
		display.SetWindowTitle("My Game")
		eventQueue.RegisterEventSource(display.EventSource())
	} else {
		panic(err)
	}

	// Mouse
	if err = allegro.InstallMouse(); err == nil {
		var mouseEventSource *allegro.EventSource
		if mouseEventSource, err = allegro.MouseEventSource(); err == nil {
			eventQueue.RegisterEventSource(mouseEventSource)
		} else {
			panic(err)
		}
	} else {
		panic(err)
	}

	// Keyboard
	if err = allegro.InstallKeyboard(); err == nil {
		var keyboardEventSource *allegro.EventSource
		if keyboardEventSource, err = allegro.KeyboardEventSource(); err == nil {
			eventQueue.RegisterEventSource(keyboardEventSource)
		} else {
			panic(err)
		}
	} else {
		panic(err)
	}

	// FPS Timer
	secondsPerFrame := 1 / float64(FPS)
	if fpsTimer, err = allegro.CreateTimer(secondsPerFrame); err == nil {
		defer fpsTimer.Destroy()
		fpsTimerSource = fpsTimer.EventSource()
		eventQueue.RegisterEventSource(fpsTimerSource)
	} else {
		panic(err)
	}

	// Submodules
	font.Init()
	image.Init()

	game.Init(display)
	game.GoTo(eventQueue, &scenes.MenuScene{Name: "main"})

	if *cpuprofile != "" {
		f, err := os.Create(*cpuprofile)
		if err != nil {
			panic(err)
		}
		pprof.StartCPUProfile(f)
		defer pprof.StopCPUProfile()
	}

	var (
		event       allegro.Event
		running     bool          = true
		needsUpdate bool          = false
		lastUpdate  time.Time     = time.Now()
		step        time.Duration = time.Duration(secondsPerFrame * float64(time.Second))
		lag         time.Duration
		now         time.Time
		elapsed     time.Duration
	)

	fpsTimer.Start()

	for {
		eventQueue.WaitForEvent(&event)
		switch event.Type {
		case allegro.EVENT_TIMER:
			if event.Source == fpsTimerSource {
				needsUpdate = true
			}

		case allegro.EVENT_DISPLAY_CLOSE:
			running = false

		default:
			game.HandleEvent(&event)
		}

		if !running {
			break
		}

		if needsUpdate && eventQueue.IsEmpty() {
			now = time.Now()
			elapsed = now.Sub(lastUpdate)
			lastUpdate = now
			lag += elapsed
			for lag >= step {
				game.Update()
				lag -= step
			}
			game.Render(float32(lag / step))
			needsUpdate = false
		}
	}
}
