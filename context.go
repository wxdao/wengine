package wengine

import (
	"errors"
)

type AssetMap map[string]Asset
type SceneMap map[string]*Scene

type Context struct {
	scrWidth, scrHeight int

	assets       AssetMap
	scenes       SceneMap
	currentScene *Scene

	renderer        Renderer
	rendererSetting RendererSetting

	assetsToFinalize AssetMap
}

func NewContext() *Context {
	return &Context{assets: make(AssetMap), scenes: make(SceneMap)}
}

func (ctx *Context) AccessRenderSetting() *RendererSetting {
	return &ctx.rendererSetting
}

func (ctx *Context) ScreenSize() (width, height int) {
	width = ctx.scrWidth
	height = ctx.scrHeight
	return
}

func (ctx *Context) Assets() AssetMap {
	return ctx.assets
}

func (ctx *Context) RegisterAsset(name string, asset Asset) {
	ctx.assets[name] = asset
}

func (ctx *Context) RegisterScene(name string, scene *Scene) {
	ctx.scenes[name] = scene
}

func (ctx *Context) ApplyScene(name string) {
	scene, exists := ctx.scenes[name]
	if !exists {
		return
	}
	ctx.currentScene = scene
}

func (ctx *Context) SetScreenSize(width, height int) {
	ctx.scrWidth = width
	ctx.scrHeight = height
}

func (ctx *Context) LoadAssets(assets []string) error {
	for _, name := range assets {
		asset, exists := ctx.assets[name]
		if !exists {
			return errors.New("no such asset")
		}
		if asset.Loaded() {
			continue
		}
		err := asset.load()
		if err != nil {
			return err
		}
		println("loaded asset: " + name)
	}
	ctx.renderer.NotifyInstall(assets)
	return nil
}

func (ctx *Context) asyncLoadScene(scene *Scene) (chan error, error) {
	// figure out all assets
	assetsToLoad := []string{}
	for _, obj := range scene.objects {
		for _, compo := range obj.components {
			switch compo.Type() {
			case COMPO_MESH:
				meshCompo, ok := compo.(*MeshComponent)
				if !ok {
					return nil, errors.New("found invalid component")
				}

				meshAsset := ctx.assets[meshCompo.Mesh]
				if meshAsset != nil {
					assetsToLoad = append(assetsToLoad, meshCompo.Mesh)
				} else {
					return nil, errors.New("found invalid component")
				}

				materialAsset := ctx.assets[meshCompo.Material]
				if materialAsset != nil {
					assetsToLoad = append(assetsToLoad, meshCompo.Material)
				}

				shaderAsset := ctx.assets[meshCompo.Shader]
				if shaderAsset != nil {
					assetsToLoad = append(assetsToLoad, meshCompo.Shader)
				}
			}
		}
	}
	result := make(chan error, 1)
	go func() {
		err := ctx.LoadAssets(assetsToLoad)
		result <- err
	}()
	return result, nil
}
