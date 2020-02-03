package main

import (
	"encoding/json"
	"os"
	"runtime"
	"strconv"

	ketchup "./ketchup"
	mustard "./mustard"
	sauce "./sauce"

	"github.com/fogleman/gg"
	"github.com/go-gl/gl/v4.1-core/gl"
	"github.com/go-gl/glfw/v3.3/glfw"
)

func main() {
	runtime.LockOSThread()
	glfw.Init()
	gl.Init()

	mustard.SetGLFWHints()

	url := os.Args[1]
	resource := sauce.GetResource(url)
	htmlString := string(resource.Body)
	parsedDocument := ketchup.ParseDocument(htmlString)
	res, _ := json.MarshalIndent(parsedDocument.Children[0], "", "--")

	app := mustard.CreateNewApp("THDWB")
	window := mustard.CreateNewWindow("THDWB", 600, 600)
	rootFrame := mustard.CreateFrame(mustard.HorizontalFrame)

	appBar := mustard.CreateFrame(mustard.VerticalFrame)

	titleBar := mustard.CreateLabelWidget("THDWB - nil")
	titleBar.SetFontColor("#fff")

	logo := mustard.CreateImageWidget("logo.png")
	logo.SetWidth(20)

	appBar.SetHeight(28)
	appBar.AttachWidget(logo)
	appBar.AttachWidget(titleBar)
	appBar.SetBackgroundColor("#5f6368")

	rootFrame.AttachWidget(appBar)

	viewPort := mustard.CreateContextWidget(func(ctx *gg.Context) {
		ctx.SetHexColor("#000000")
		ctx.LoadFontFace("roboto.ttf", 12)
		ctx.DrawStringWrapped(string(res), float64(0+12/4), float64(0+12*2/2), float64(0), float64(0), float64(ctx.Width()), 2, gg.AlignLeft)
		ctx.Fill()
	})

	rootFrame.AttachWidget(viewPort)

	statusBar := mustard.CreateFrame(mustard.HorizontalFrame)
	statusBar.SetBackgroundColor("#babcbe")
	statusBar.SetHeight(20)

	statusLabel := mustard.CreateLabelWidget("Processed Events:")
	statusLabel.SetFontSize(16)
	frameEvents := 0

	rootFrame.AttachWidget(statusBar)
	statusBar.AttachWidget(statusLabel)

	window.SetRootFrame(rootFrame)

	app.AddWindow(window)

	window.Show()
	app.Run(func() {
		frameEvents++
		width, height := window.GetSize()
		statusLabel.SetContent("Processed Events: " + strconv.Itoa(frameEvents) + "; Resolution: " + strconv.Itoa(width) + "X" + strconv.Itoa(height))
	})

}