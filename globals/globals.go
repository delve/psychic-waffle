package globals

const (
	EdgeScrollSpeed     = KeyboardScrollSpeed
	EdgeWidth           = 20
	ZoomSpeed           = -0.125
	WorldWidth          = 800
	WorldHeight         = 800
	KeyboardScrollSpeed = (WorldWidth * WorldHeight) / 600
	// todo needs more intelligence
	HUDHeight = WorldHeight / 4
	HUDWidth  = WorldWidth / 4
	HUDZ      = 1000
)
