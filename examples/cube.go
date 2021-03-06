package main

import (
	"flag"
	"github.com/go-gl/mathgl/mgl32"
	"github.com/wxdao/wengine"
	_ "github.com/wxdao/wengine/opengl"
	"log"
	"math"
	"os"
	"path"
	"runtime/pprof"
)

var cpuprofile = flag.String("cpuprofile", "", "write cpu profile `file`")
var texfile = flag.String("tex", "", "texture file for cube `file`")

var app *wengine.App

func main() {
	flag.Parse()
	if *cpuprofile != "" {
		f, err := os.Create(*cpuprofile)
		if err != nil {
			log.Fatal("could not create CPU profile: ", err)
		}
		if err := pprof.StartCPUProfile(f); err != nil {
			log.Fatal("could not start CPU profile: ", err)
		}
		defer pprof.StopCPUProfile()
	}

	config := &wengine.Config{
		Width:       800,
		Height:      600,
		WindowMode:  wengine.WINDOW_MODE_WINDOWED,
		WindowTitle: "wEngine Example: cube",
		FrameLimit:  30,
		VSync:       true,
	}
	app, _ = wengine.NewApp(config)
	context := app.Context()
	myScene := wengine.NewScene()
	context.RegisterScene("myScene", myScene)

	setupInput(context)
	setupScene(context, myScene)
	context.ApplyScene("myScene")

	panic(app.Run())
}

func setupInput(context *wengine.Context) {
	context.Input().BindAxis("angle", wengine.AxisMeta{
		Source:      wengine.AXIS_SOURCE_KEY,
		PositiveKey: "w",
		NegativeKey: "s",
		Gravity:     0.5,
		Sensitivity: 0.5,
	})
	context.Input().BindAxis("intensity", wengine.AxisMeta{
		Source:      wengine.AXIS_SOURCE_KEY,
		PositiveKey: "e",
		NegativeKey: "d",
		Gravity:     0.5,
		Sensitivity: 0.5,
	})

	context.Input().BindAxis("mouse x", wengine.AxisMeta{
		Source:      wengine.AXIS_SOURCE_MOUSE,
		From:        wengine.AXIS_FROM_X,
		Sensitivity: 0.01,
	})
	context.Input().BindAxis("mouse y", wengine.AxisMeta{
		Source:      wengine.AXIS_SOURCE_MOUSE,
		From:        wengine.AXIS_FROM_Y,
		Sensitivity: 0.01,
	})
}

func setupScene(context *wengine.Context, scene *wengine.Scene) {
	camera := &wengine.CameraComponent{}
	camera.ViewportX, camera.ViewportY, camera.ViewportW, camera.ViewportH = 0, 0, 1, 1
	camera.Depth = 0
	camera.ClearColor, camera.ClearDepth = true, true
	camera.Mode = wengine.CAMERA_MODE_PERSPECTIVE
	camera.FOV = mgl32.DegToRad(60)
	camera.Width = 15
	camera.FarPlane = 100
	camera.NearPlane = 0.3
	camera.Ambient = mgl32.Vec3{0.2, 0.2, 0.2}
	cameraObject := wengine.NewObject()
	cameraObject.Translate(mgl32.Vec3{0, 0, 10})
	cameraObject.SetBehavior(&CameraBehavior{})
	cameraObject.AttachComponent(camera)
	cameraObject.SetEnabled(true)
	scene.RegisterObject("mainCamera", cameraObject)

	camera2 := &wengine.CameraComponent{}
	camera2.ViewportX, camera2.ViewportY, camera2.ViewportW, camera2.ViewportH = 0, 0.6, 0.4, 0.4
	camera2.Depth = -1
	camera2.ClearColor, camera2.ClearDepth = true, true
	camera2.Mode = wengine.CAMERA_MODE_PERSPECTIVE
	camera2.FOV = mgl32.DegToRad(60)
	camera2.Width = 15
	camera2.FarPlane = 100
	camera2.NearPlane = 0.3
	camera2.Ambient = mgl32.Vec3{0.2, 0.2, 0.2}
	cameraObject2 := wengine.NewObject()
	cameraObject2.Translate(mgl32.Vec3{0, 20, -5})
	cameraObject2.Rotate(mgl32.DegToRad(-90), mgl32.Vec3{1, 0, 0})
	cameraObject2.AttachComponent(camera2)
	cameraObject2.SetEnabled(true)
	scene.RegisterObject("secondaryCamera", cameraObject2)

	dirLight := &wengine.LightComponent{}
	dirLight.LightSource = wengine.LIGHT_SOURCE_DIRECTIONAL
	dirLight.ShadowType = wengine.LIGHT_SHADOW_TYPE_HARD
	dirLight.Diffuse = mgl32.Vec3{0.8, 0.8, 0.8}
	dirLight.Specular = mgl32.Vec3{0.2, 0.2, 0.2}
	dirLightObject := wengine.NewObject()
	dirLightObject.Translate(mgl32.Vec3{0, 10, 10})
	dirLightObject.Rotate(mgl32.DegToRad(-90), mgl32.Vec3{1, 0, 0})
	dirLightObject.AttachComponent(dirLight)
	dirLightObject.SetEnabled(false)
	scene.RegisterObject("dirLight", dirLightObject)

	pointLight := &wengine.LightComponent{}
	pointLight.LightSource = wengine.LIGHT_SOURCE_POINT
	pointLight.ShadowType = wengine.LIGHT_SHADOW_TYPE_NONE
	pointLight.Range = 20
	pointLight.Diffuse = mgl32.Vec3{1, 1, 1}
	pointLight.Specular = mgl32.Vec3{0.5, 0.5, 0.5}
	pointLightObject := wengine.NewObject()
	pointLightObject.Translate(mgl32.Vec3{5, 8, -5})
	pointLightObject.AttachComponent(pointLight)
	pointLightObject.SetEnabled(false)
	scene.RegisterObject("pointLight", pointLightObject)

	spotLight := &wengine.LightComponent{}
	spotLight.LightSource = wengine.LIGHT_SOURCE_SPOT
	spotLight.ShadowType = wengine.LIGHT_SHADOW_TYPE_HARD
	spotLight.Angle = mgl32.DegToRad(30)
	spotLight.Range = 26
	spotLight.Diffuse = mgl32.Vec3{10, 10, 10}
	spotLight.Specular = mgl32.Vec3{1, 1, 1}
	spotLightObject := wengine.NewObject()
	spotLightObject.Translate(mgl32.Vec3{0, 20, -5})
	spotLightObject.Rotate(mgl32.DegToRad(-90), mgl32.Vec3{1, 0, 0})
	spotLightObject.AttachComponent(spotLight)
	spotLightObject.SetEnabled(true)
	scene.RegisterObject("spotLight", spotLightObject)

	cubeMesh := wengine.DefaultCubeMeshAsset()
	context.RegisterAsset("cubeMesh", cubeMesh)

	floorMesh := wengine.DefaultPlaneMeshAsset()
	context.RegisterAsset("floorMesh", floorMesh)

	floorMaterial := &wengine.MeshMaterialAsset{}
	floorMaterial.DiffuseColor = mgl32.Vec4{0.7, 0.7, 0.7, 1}
	context.RegisterAsset("floorMaterial", floorMaterial)

	cubeMaterial := &wengine.MeshMaterialAsset{}
	if *texfile == "" {
		exe, err := os.Executable()
		if err != nil {
			cubeMaterial.DiffuseMapPath = "cube.png"
		} else {
			cubeMaterial.DiffuseMapPath = path.Join(path.Dir(exe), "cube.png")
		}
	} else {
		cubeMaterial.DiffuseMapPath = *texfile
	}
	context.RegisterAsset("cubeMaterial", cubeMaterial)

	cubeObject1 := wengine.NewObject()
	cubeObject1.Translate(mgl32.Vec3{0, 0, -5})
	cubeObject1.SetBehavior(&ControlledRotationBehavior{})
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
	cubeObject2.SetBehavior(&RotationBehavior{axis: mgl32.Vec3{1, 1, 1}})
	cubeObject2.AttachComponent(&wengine.MeshComponent{
		Mesh:          "cubeMesh",
		Material:      "cubeMaterial",
		CastShadow:    true,
		ReceiveShader: true,
	})
	cubeObject2.SetEnabled(true)
	scene.RegisterObject("simpleCube2", cubeObject2)

	floorObject := wengine.NewObject()
	floorObject.Translate(mgl32.Vec3{0, -5, 0})
	floorObject.Scale(mgl32.Vec3{100, 1, 100})
	floorObject.AttachComponent(&wengine.MeshComponent{
		Mesh:          "floorMesh",
		Material:      "floorMaterial",
		CastShadow:    false,
		ReceiveShader: true,
	})
	floorObject.SetEnabled(true)
	scene.RegisterObject("floor", floorObject)
}

type RotationBehavior struct {
	axis mgl32.Vec3
}

func (b *RotationBehavior) Start(bctx *wengine.BehaviorContext) {
}

func (b *RotationBehavior) Update(bctx *wengine.BehaviorContext) {
	bctx.Object.Rotate(0.3*float32(bctx.DeltaTime), b.axis)
}

type ControlledRotationBehavior struct {
}

func (b *ControlledRotationBehavior) Start(bctx *wengine.BehaviorContext) {
}

func (b *ControlledRotationBehavior) Update(bctx *wengine.BehaviorContext) {
	bctx.Object.Rotate(float32(bctx.Context.Input().GetAxis("mouse y")), mgl32.Vec3{1, 0, 0})
	bctx.Object.Rotate(float32(bctx.Context.Input().GetAxis("mouse x")), mgl32.Vec3{0, 1, 0})
}

type CameraBehavior struct {
	spotLight *wengine.LightComponent

	lastCursorMode int
}

func (b *CameraBehavior) Start(bctx *wengine.BehaviorContext) {
	b.lastCursorMode = wengine.CURSOR_MODE_NORMAL
	bctx.Context.Input().SetCursorMode(b.lastCursorMode)

	b.spotLight = bctx.Context.CurrentScene().Objects()["spotLight"].Components()[wengine.COMPO_LIGHT].(*wengine.LightComponent)
}

func (b *CameraBehavior) Update(bctx *wengine.BehaviorContext) {
	if bctx.Context.Input().GetKeyUp("esc") {
		app.Stop()
	}
	if bctx.Context.Input().GetKeyUp("v") {
		if b.lastCursorMode == wengine.CURSOR_MODE_NORMAL {
			bctx.Context.Input().SetCursorMode(wengine.CURSOR_MODE_DISABLED)
			b.lastCursorMode = wengine.CURSOR_MODE_DISABLED
		} else {
			bctx.Context.Input().SetCursorMode(wengine.CURSOR_MODE_NORMAL)
			b.lastCursorMode = wengine.CURSOR_MODE_NORMAL
		}
	}
	b.spotLight.Angle = mgl32.DegToRad(float32(25 + 10*bctx.Context.Input().GetAxis("angle")))
	b.spotLight.Diffuse = mgl32.Vec3{1, 1, 1}.Mul(float32(math.Max(5, 10+30*bctx.Context.Input().GetAxis("intensity"))))
}
