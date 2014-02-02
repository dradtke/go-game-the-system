package engine

type Entity struct {
	X, Y   float32 // coordinates of the entity
	Block  bool    // prevent other entities from moving through it
	prev_x, prev_y float32 // previous coordinates, unexported
}

func (e *Entity) Move(x, y float32) {
	e.X += x
	e.Y += y
}

// Delta() returns difference in x and y between the entity's
// current x-value and its previous one.
func (e *Entity) Delta() (x, y float32) {
	return e.X - e.prev_x, e.Y - e.prev_y
}
