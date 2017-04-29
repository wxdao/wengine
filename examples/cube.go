package main

import (
	"github.com/go-gl/mathgl/mgl32"
	"github.com/wxdao/wengine"
	_ "github.com/wxdao/wengine/opengl"
	"os"
	"path"
)

func main() {
	config := &wengine.Config{
		Width:       800,
		Height:      600,
		WindowMode:  wengine.WINDOW_MODE_WINDOWED,
		WindowTitle: "wEngine Example: cube",
		FrameLimit:  60,
		Renderer:    "opengl",
	}
	app, _ := wengine.NewApp(config)
	context := app.Context()
	myScene := wengine.NewScene()
	context.RegisterScene("myScene", myScene)
	setupScene(context, myScene)
	context.ApplyScene("myScene")
	panic(app.Run())
}

func setupScene(context *wengine.Context, scene *wengine.Scene) {
	camera := &wengine.CameraComponent{}
	camera.Mode = wengine.CAMERA_MODE_PERSPECTIVE
	camera.FOV = mgl32.DegToRad(60)
	camera.Width = 15
	camera.FarPlane = 100
	camera.NearPlane = 0.3
	camera.Ambient = mgl32.Vec3{0.2, 0.2, 0.2}
	cameraObject := wengine.NewObject()
	cameraObject.Translate(mgl32.Vec3{0, 0, 10})
	cameraObject.AttachComponent(camera)
	cameraObject.SetEnabled(true)
	scene.RegisterObject("mainCamera", cameraObject)

	dirLight := &wengine.LightComponent{}
	dirLight.LightSource = wengine.LIGHT_SOURCE_DIRECTIONAL
	dirLight.ShadowType = wengine.LIGHT_SHADOW_TYPE_HARD
	dirLight.Diffuse = mgl32.Vec3{1, 1, 1}
	dirLight.Specular = mgl32.Vec3{0.5, 0.5, 0.5}
	dirLightObject := wengine.NewObject()
	dirLightObject.Translate(mgl32.Vec3{0, 10, 10})
	dirLightObject.Rotate(mgl32.DegToRad(-90), mgl32.Vec3{1, 0, 0})
	dirLightObject.AttachComponent(dirLight)
	dirLightObject.SetEnabled(false)
	scene.RegisterObject("dirLight", dirLightObject)

	pointLight := &wengine.LightComponent{}
	pointLight.LightSource = wengine.LIGHT_SOURCE_POINT
	pointLight.ShadowType = wengine.LIGHT_SHADOW_TYPE_HARD
	pointLight.Range = 20
	pointLight.Diffuse = mgl32.Vec3{1, 1, 1}
	pointLight.Specular = mgl32.Vec3{0.5, 0.5, 0.5}
	pointLightObject := wengine.NewObject()
	pointLightObject.Translate(mgl32.Vec3{0, 10, -5})
	pointLightObject.AttachComponent(pointLight)
	pointLightObject.SetEnabled(true)
	scene.RegisterObject("pointLight", pointLightObject)

	pointLight2 := &wengine.LightComponent{}
	pointLight2.LightSource = wengine.LIGHT_SOURCE_POINT
	pointLight2.ShadowType = wengine.LIGHT_SHADOW_TYPE_HARD
	pointLight2.Range = 20
	pointLight2.Diffuse = mgl32.Vec3{1, 1, 1}
	pointLight2.Specular = mgl32.Vec3{0.5, 0.5, 0.5}
	pointLightObject2 := wengine.NewObject()
	pointLightObject2.Translate(mgl32.Vec3{5, 10, -5})
	pointLightObject2.AttachComponent(pointLight2)
	pointLightObject2.SetEnabled(true)
	scene.RegisterObject("pointLight2", pointLightObject2)

	pointLight3 := &wengine.LightComponent{}
	pointLight3.LightSource = wengine.LIGHT_SOURCE_POINT
	pointLight3.ShadowType = wengine.LIGHT_SHADOW_TYPE_HARD
	pointLight3.Range = 20
	pointLight3.Diffuse = mgl32.Vec3{1, 1, 1}
	pointLight3.Specular = mgl32.Vec3{0.5, 0.5, 0.5}
	pointLightObject3 := wengine.NewObject()
	pointLightObject3.Translate(mgl32.Vec3{-5, 10, -5})
	pointLightObject3.AttachComponent(pointLight3)
	pointLightObject3.SetEnabled(true)
	scene.RegisterObject("pointLight3", pointLightObject3)

	cubeMesh := wengine.DefaultCubeAsset()
	context.RegisterAsset("cubeMesh", cubeMesh)

	cubeColorMaterial := &wengine.MaterialAsset{}
	cubeColorMaterial.DiffuseColor = mgl32.Vec4{0.7, 0.7, 0.7, 1}
	context.RegisterAsset("cubeColorMaterial", cubeColorMaterial)

	cubeMaterial := &wengine.MaterialAsset{}
	if len(os.Args) < 2 {
		exe, err := os.Executable()
		if err != nil {
			cubeMaterial.DiffuseMapPath = "cube.png"
		} else {
			cubeMaterial.DiffuseMapPath = path.Join(path.Dir(exe), "cube.png")
		}
	} else {
		cubeMaterial.DiffuseMapPath = os.Args[1]
	}
	context.RegisterAsset("cubeMaterial", cubeMaterial)

	cubeObject1 := wengine.NewObject()
	cubeObject1.Translate(mgl32.Vec3{0, 0, -5})
	cubeObject1.SetBehavior(&CubeBehavior{axis: mgl32.Vec3{1, 1, 1}})
	cubeObject1.AttachComponent(&wengine.MeshComponent{
		Mesh:          "cubeMesh",
		Material:      "cubeMaterial",
		CastShadow:    true,
		ReceiveShader: true,
	})
	cubeObject1.SetEnabled(true)
	scene.RegisterObject("simpleCube1", cubeObject1)

	cubeObject2 := wengine.NewObject()
	cubeObject2.SetParent(cubeObject1)
	cubeObject2.Translate(mgl32.Vec3{5, 0, 0})
	cubeObject2.SetBehavior(&CubeBehavior{axis: mgl32.Vec3{0, 1, 0}})
	cubeObject2.AttachComponent(&wengine.MeshComponent{
		Mesh:          "cubeMesh",
		Material:      "cubeMaterial",
		CastShadow:    true,
		ReceiveShader: true,
	})
	cubeObject2.SetEnabled(true)
	scene.RegisterObject("simpleCube2", cubeObject2)

	cubeObject3 := wengine.NewObject()
	cubeObject3.Translate(mgl32.Vec3{0, -5, 0})
	cubeObject3.Scale(mgl32.Vec3{100, 1, 100})
	cubeObject3.AttachComponent(&wengine.MeshComponent{
		Mesh:          "cubeMesh",
		Material:      "cubeColorMaterial",
		CastShadow:    false,
		ReceiveShader: true,
	})
	cubeObject3.SetEnabled(true)
	scene.RegisterObject("simpleCube3", cubeObject3)
}

type CubeBehavior struct {
	axis  mgl32.Vec3
	asset wengine.MaterialAsset
}

func (b *CubeBehavior) Start(bctx *wengine.BehaviorContext) {
}

func (b *CubeBehavior) Update(bctx *wengine.BehaviorContext) {
	bctx.Object.Rotate(0.3*float32(bctx.DeltaTime), b.axis)
}