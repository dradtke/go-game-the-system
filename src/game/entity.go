package game

type Entity interface {
	Update(*State)
	Render(*State, float32)
}

func AddEntity(e Entity) {
	entities = append(entities, e)
}
