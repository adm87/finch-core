package camera

func ScreenToWorld(cameraComp *CameraComponent, sx, sy float64) (float64, float64) {
	matrix := cameraComp.WorldMatrix()
	return matrix.Apply(float64(sx), float64(sy))
}

func WorldToScreen(cameraComp *CameraComponent, wx, wy float64) (float64, float64) {
	matrix := cameraComp.WorldMatrix()
	matrix.Invert()
	return matrix.Apply(float64(wx), float64(wy))
}
