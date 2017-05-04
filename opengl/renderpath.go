package opengl

import (
	. "github.com/wxdao/wengine"
)

type renderPath interface {
	init() error
	render(targetFBO uint32, lights []*LightComponent, meshes []*MeshComponent, scene *Scene, camera *CameraComponent) error
}
