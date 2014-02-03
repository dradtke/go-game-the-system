package game

func Update() {
	defer func() {
		sync(&state.current, &state.prev)
	}()

	if state.sceneLoaded {
		for _, e := range entities {
			e.Update(state)
		}
	}
	scene.Update(state)

	select {
	case <-loading:
		state.sceneLoaded = true
		return
	default:
		// not yet loaded
	}

	if state.current.MouseLeftDown && !state.prev.MouseLeftDown {
		scene.OnLeftPress(state)
	} else if !state.current.MouseLeftDown && state.prev.MouseLeftDown {
		scene.OnLeftRelease(state)
	}

	if state.current.MouseRightDown && !state.prev.MouseRightDown {
		scene.OnRightPress(state)
	} else if !state.current.MouseRightDown && state.prev.MouseRightDown {
		scene.OnRightRelease(state)
	}
}
