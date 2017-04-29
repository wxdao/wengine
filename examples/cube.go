package main

import (
	"flag"
	"github.com/go-gl/mathgl/mgl32"
	"github.com/wxdao/wengine"
	_ "github.com/wxdao/wengine/opengl"
	"log"
	"os"
	"path"
	"runtime/pprof"
)

var cpuprofile = flag.String("cpuprofile", "", "write cpu profile `file`")
var texfile = flag.String("tex", "", "texture file for cube `file`")

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
	dirLight.Diffuse = mgl32.Vec3{0.5, 0.5, 0.5}
	dirLight.Specular = mgl32.Vec3{0.1, 0.1, 0.1}
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
	pointLightObject.Translate(mgl32.Vec3{5, 8, -5})
	pointLightObject.AttachComponent(pointLight)
	pointLightObject.SetEnabled(false)
	scene.RegisterObject("pointLight", pointLightObject)

	spotLight := &wengine.LightComponent{}
	spotLight.LightSource = wengine.LIGHT_SOURCE_SPOT
	spotLight.ShadowType = wengine.LIGHT_SHADOW_TYPE_HARD
	spotLight.Angle = mgl32.DegToRad(30)
	spotLight.Range = 25
	spotLight.Diffuse = mgl32.Vec3{10, 10, 10}
	spotLight.Specular = mgl32.Vec3{1, 1, 1}
	spotLightObject := wengine.NewObject()
	spotLightObject.Translate(mgl32.Vec3{0, 20, -5})
	spotLightObject.Rotate(mgl32.DegToRad(-90), mgl32.Vec3{1, 0, 0})
	spotLightObject.AttachComponent(spotLight)
	spotLightObject.SetEnabled(true)
	scene.RegisterObject("spotLight", spotLightObject)

	cubeMesh := wengine.DefaultCubeAsset()
	context.RegisterAsset("cubeMesh", cubeMesh)

	cubeColorMaterial := &wengine.MaterialAsset{}
	cubeColorMaterial.DiffuseColor = mgl32.Vec4{0.7, 0.7, 0.7, 1}
	context.RegisterAsset("cubeColorMaterial", cubeColorMaterial)

	cubeMaterial := &wengine.MaterialAsset{}
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
