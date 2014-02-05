package game

type Entity interface {
	Update(*State)
	Render(*State, float32)
}

type TextEntity interface {
	Entity
	Write(rune)
	Backspace()
	ShiftCursor(int)
}

func AddEntity(e Entity) {
	entities = append(entities, e)
	if t, ok := e.(TextEntity); ok {
		textEntities = append(textEntities, t)
	}
}
