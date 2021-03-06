package bun

import (
	"fmt"

	assets "thdwb/assets"
	gg "thdwb/gg"
	structs "thdwb/structs"
)

func RenderDocument(ctx *gg.Context, document *structs.HTMLDocument) {
	//tree.Children[0] is head
	body := document.RootElement.Children[1]

	document.RootElement.Style.Width = float64(ctx.Width())
	document.RootElement.Style.Height = float64(ctx.Height())

	layoutDOM(ctx, body, 0)
}

func getNodeContent(NodeDOM *structs.NodeDOM) string {
	return NodeDOM.Content
}

func getElementName(NodeDOM *structs.NodeDOM) string {
	return NodeDOM.Element
}

func getNodeChildren(NodeDOM *structs.NodeDOM) []*structs.NodeDOM {
	return NodeDOM.Children
}

func walkDOM(TreeDOM *structs.NodeDOM, d string) {
	fmt.Println(d, getElementName(TreeDOM))
	nodeChildren := getNodeChildren(TreeDOM)

	for i := 0; i < len(nodeChildren); i++ {
		walkDOM(nodeChildren[i], d+"-")
	}
}

func layoutDOM(ctx *gg.Context, node *structs.NodeDOM, childIdx int) {
	nodeChildren := getNodeChildren(node)

	calculateNode(ctx, node, childIdx)

	for i := 0; i < len(nodeChildren); i++ {
		layoutDOM(ctx, nodeChildren[i], i)
		node.Style.Height += nodeChildren[i].Style.Height
	}

	paintNode(ctx, node)
}

func paintNode(ctx *gg.Context, node *structs.NodeDOM) {
	switch node.Style.Display {
	case "block":
		paintBlockElement(ctx, node)
	case "inline":
		paintInlineElement(ctx, node)
	case "list-item":
		paintListItemElement(ctx, node)
	}
}

func calculateNode(ctx *gg.Context, node *structs.NodeDOM, postion int) {
	switch node.Style.Display {
	case "block":
		calculateBlockLayout(ctx, node, postion)
	case "inline":
		calculateInlineLayout(ctx, node, postion)
	case "list-item":
		calculateListItemLayout(ctx, node, postion)
	}
}

func paintBlockElement(ctx *gg.Context, node *structs.NodeDOM) {
	ctx.DrawRectangle(node.Style.Left, node.Style.Top, node.Style.Width, node.Style.Height)
	ctx.SetRGBA(node.Style.BackgroundColor.R, node.Style.BackgroundColor.G, node.Style.BackgroundColor.B, node.Style.BackgroundColor.A)
	ctx.Fill()

	ctx.SetRGBA(node.Style.Color.R, node.Style.Color.G, node.Style.Color.B, node.Style.Color.A)
	ctx.LoadAssetFont(assets.SansSerif(), node.Style.FontSize)
	ctx.DrawStringWrapped(node.Content, node.Style.Left, node.Style.Top+1, 0, 0, node.Style.Width, 1.5, gg.AlignLeft)
	ctx.Fill()
}

func calculateBlockLayout(ctx *gg.Context, node *structs.NodeDOM, childIdx int) {
	if node.Style.Width == 0 {
		node.Style.Width = node.Parent.Style.Width
	}

	if node.Style.Height == 0 {
		ctx.LoadAssetFont(assets.SansSerif(), node.Style.FontSize)
		node.Style.Height = ctx.MeasureStringWrapped(node.Content, node.Style.Width, 1.5) + 2 + ctx.FontHeight()*.5
	}

	if childIdx > 0 {
		prev := node.Parent.Children[childIdx-1]

		if prev.Style.Display != "inline" {
			node.Style.Top = prev.Style.Top + prev.Style.Height
		} else {
			node.Style.Top = prev.Style.Top
		}
	} else {
		node.Style.Top = node.Parent.Style.Top
	}
}

func paintListItemElement(ctx *gg.Context, node *structs.NodeDOM) {
	ctx.DrawRectangle(node.Style.Left, node.Style.Top, node.Style.Width, node.Style.Height)
	ctx.SetRGBA(node.Style.BackgroundColor.R, node.Style.BackgroundColor.G, node.Style.BackgroundColor.B, node.Style.BackgroundColor.A)
	ctx.Fill()

	ctx.DrawCircle(node.Style.Left+15, node.Style.Top+node.Style.FontSize/2, 3)
	ctx.SetRGBA(node.Style.Color.R, node.Style.Color.G, node.Style.Color.B, node.Style.Color.A)
	ctx.LoadAssetFont(assets.SansSerif(), node.Style.FontSize)
	ctx.DrawStringWrapped(node.Content, node.Style.Left+30, node.Style.Top+1, 0, 0, node.Style.Width, 1.5, gg.AlignLeft)
	ctx.Fill()
}

func calculateListItemLayout(ctx *gg.Context, node *structs.NodeDOM, childIdx int) {
	if node.Style.Width == 0 {
		node.Style.Width = node.Parent.Style.Width
	}

	if node.Style.Height == 0 && len(node.Content) > 0 {
		ctx.LoadAssetFont(assets.SansSerif(), node.Style.FontSize)
		node.Style.Height = ctx.MeasureStringWrapped(node.Content, node.Style.Width, 1.5) + 2 + ctx.FontHeight()*.5
	}

	if childIdx > 0 {
		prev := node.Parent.Children[childIdx-1]
		node.Style.Top = prev.Style.Top + prev.Style.Height
	} else {
		node.Style.Top = node.Parent.Style.Top
	}
}

func paintInlineElement(ctx *gg.Context, node *structs.NodeDOM) {
	ctx.DrawRectangle(node.Style.Left, node.Style.Top, node.Style.Width, node.Style.Height)
	ctx.SetRGBA(node.Style.BackgroundColor.R, node.Style.BackgroundColor.G, node.Style.BackgroundColor.B, node.Style.BackgroundColor.A)
	ctx.Fill()

	ctx.SetRGBA(node.Style.Color.R, node.Style.Color.G, node.Style.Color.B, node.Style.Color.A)
	ctx.LoadAssetFont(assets.SansSerif(), node.Style.FontSize)
	ctx.DrawStringWrapped(node.Content, node.Style.Left, node.Style.Top, 0, 0, node.Style.Width, 1.5, gg.AlignLeft)
	ctx.Fill()
}

func calculateInlineLayout(ctx *gg.Context, node *structs.NodeDOM, childIdx int) {
	ctx.LoadAssetFont(assets.SansSerif(), node.Style.FontSize)

	if childIdx > 0 && node.Parent.Children[childIdx-1] != nil {
		prev := node.Parent.Children[childIdx-1]
		if prev.Style.Display == "inline" {
			node.Style.Top = prev.Style.Top
		} else {
			node.Style.Top = prev.Style.Top + prev.Style.Height
		}
	} else {
		node.Style.Top = node.Parent.Style.Top
	}

	if childIdx > 0 && node.Parent.Children[childIdx-1] != nil {
		prev := node.Parent.Children[childIdx-1]
		if prev.Style.Display == "inline" {
			node.Style.Left = prev.Style.Left + prev.Style.Width
		}
	} else {
		node.Style.Left = node.Parent.Style.Left
	}

	node.Style.Width, node.Style.Height = ctx.MeasureMultilineString(node.Content, 1.5)
	node.Style.Height++
}

func GetPageTitle(TreeDOM *structs.NodeDOM) string {
	nodeChildren := getNodeChildren(TreeDOM)
	pageTitle := "Sem Titulo"

	if getElementName(TreeDOM) == "title" {
		return getNodeContent(TreeDOM)
	}

	for i := 0; i < len(nodeChildren); i++ {
		nPageTitle := GetPageTitle(nodeChildren[i])

		if nPageTitle != "Sem Titulo" {
			pageTitle = nPageTitle
		}
	}

	return pageTitle
}
