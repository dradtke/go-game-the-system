This repository is just the result of me playing around with game development patterns and architecture in Go. It's based off of the concept of "scenes", which is essentially just an implementation of the state pattern. The idea is that each state the game can be in is a `Scene`, and that all of the game's logic can be encapsulated within its respective scene.

To that end, the Allegro backend and game state are abstracted away, and not intended to be modified at all by the fictional developer who's actually using this framework to build a game. Game state is double-buffered into "previous" and "current" states, and the current state is made available to scenes via a struct which grants read-only access via instance methods while the fields themselves are kept hidden.

Scenes should be based off of `BaseScene` by including it as an embedded field; this allows new `Scene` implementations to fulfill the interface requirements without needing to implement every single available method. Any of them may be overridden, giving full flexibility for how scenes pay attention to game events.

Maybe someday the code here will actually be useful to somebody. =)

Building
--------

This code isn't structured like a normal Go program because it's not intended to be used as a package. If you'd like to try building it, clone the repo into a new directory, cd into it, and then run

```
$ GOPATH=${GOPATH}:`pwd` go build main.go
```
