package widget

import (
	"fmt"
	"testing"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/internal/cache"
	"fyne.io/fyne/v2/internal/widget"
	"fyne.io/fyne/v2/test"
	"fyne.io/fyne/v2/theme"

	"github.com/stretchr/testify/assert"
)

func TestButton_Style(t *testing.T) {
	button := NewButton("Test", nil)
	bg := test.WidgetRenderer(button).(*buttonRenderer).buttonColor()

	button.Importance = HighImportance
	assert.NotEqual(t, bg, test.WidgetRenderer(button).(*buttonRenderer).buttonColor())
}

func TestButton_DisabledColor(t *testing.T) {
	button := NewButton("Test", nil)
	bg := test.WidgetRenderer(button).(*buttonRenderer).buttonColor()
	button.Importance = MediumImportance
	assert.Equal(t, bg, theme.ButtonColor())

	button.Disable()
	bg = test.WidgetRenderer(button).(*buttonRenderer).buttonColor()
	assert.Equal(t, bg, theme.DisabledButtonColor())
}

func TestButton_DisabledIcon(t *testing.T) {
	button := NewButtonWithIcon("Test", theme.CancelIcon(), nil)
	render := test.WidgetRenderer(button).(*buttonRenderer)
	assert.Equal(t, render.icon.Resource.Name(), theme.CancelIcon().Name())

	button.Disable()
	assert.Equal(t, render.icon.Resource.Name(), fmt.Sprintf("disabled_%v", theme.CancelIcon().Name()))

	button.Enable()
	assert.Equal(t, render.icon.Resource.Name(), theme.CancelIcon().Name())
}

func TestButton_DisabledIconChangeUsingSetIcon(t *testing.T) {
	button := NewButtonWithIcon("Test", theme.CancelIcon(), nil)
	render := test.WidgetRenderer(button).(*buttonRenderer)
	assert.Equal(t, render.icon.Resource.Name(), theme.CancelIcon().Name())

	// assert we are using the disabled original icon
	button.Disable()
	assert.Equal(t, render.icon.Resource.Name(), fmt.Sprintf("disabled_%v", theme.CancelIcon().Name()))

	// re-enable, then change the icon
	button.Enable()
	button.SetIcon(theme.SearchIcon())
	assert.Equal(t, render.icon.Resource.Name(), theme.SearchIcon().Name())

	// assert we are using the disabled new icon
	button.Disable()
	assert.Equal(t, render.icon.Resource.Name(), fmt.Sprintf("disabled_%v", theme.SearchIcon().Name()))

}

func TestButton_DisabledIconChangedDirectly(t *testing.T) {
	button := NewButtonWithIcon("Test", theme.CancelIcon(), nil)
	render := test.WidgetRenderer(button).(*buttonRenderer)
	assert.Equal(t, render.icon.Resource.Name(), theme.CancelIcon().Name())

	// assert we are using the disabled original icon
	button.Disable()
	assert.Equal(t, render.icon.Resource.Name(), fmt.Sprintf("disabled_%v", theme.CancelIcon().Name()))

	// re-enable, then change the icon
	button.Enable()
	button.Icon = theme.SearchIcon()
	render.Refresh()
	assert.Equal(t, render.icon.Resource.Name(), theme.SearchIcon().Name())

	// assert we are using the disabled new icon
	button.Disable()
	assert.Equal(t, render.icon.Resource.Name(), fmt.Sprintf("disabled_%v", theme.SearchIcon().Name()))

}

func TestButtonRenderer_Layout(t *testing.T) {
	button := NewButtonWithIcon("Test", theme.CancelIcon(), nil)
	render := test.WidgetRenderer(button).(*buttonRenderer)
	render.Layout(render.MinSize())

	assert.True(t, render.icon.Position().X < render.label.Position().X)
	assert.Equal(t, theme.Padding()*3, render.icon.Position().X)
	assert.Equal(t, theme.Padding()*3, render.MinSize().Width-render.label.Position().X-render.label.Size().Width)
}

func TestButtonRenderer_Layout_Stretch(t *testing.T) {
	button := NewButtonWithIcon("Test", theme.CancelIcon(), nil)
	button.Resize(button.MinSize().Add(fyne.NewSize(100, 100)))
	render := test.WidgetRenderer(button).(*buttonRenderer)

	textHeight := render.label.MinSize().Height
	minIconHeight := fyne.Max(theme.IconInlineSize(), textHeight)
	assert.Equal(t, 50+theme.Padding()*3, render.icon.Position().X, "icon x")
	assert.Equal(t, 50+theme.Padding()*2, render.icon.Position().Y, "icon y")
	assert.Equal(t, theme.IconInlineSize(), render.icon.Size().Width, "icon width")
	assert.Equal(t, minIconHeight, render.icon.Size().Height, "icon height")
	assert.Equal(t, 50+theme.Padding()*4+theme.IconInlineSize(), render.label.Position().X, "label x")
	assert.Equal(t, 50+theme.Padding()*2, render.label.Position().Y, "label y")
	assert.Equal(t, render.label.MinSize(), render.label.Size(), "label size")
}

func TestButtonRenderer_Layout_NoText(t *testing.T) {
	button := NewButtonWithIcon("", theme.CancelIcon(), nil)
	render := test.WidgetRenderer(button).(*buttonRenderer)

	button.Resize(fyne.NewSize(100, 100))

	assert.Equal(t, 50-theme.IconInlineSize()/2, render.icon.Position().X)
	assert.Equal(t, 50-theme.IconInlineSize()/2, render.icon.Position().Y)
}

func TestButton_Shadow(t *testing.T) {
	{
		button := NewButton("Test", func() {})
		shadowFound := false
		for _, o := range test.LaidOutObjects(button) {
			if _, ok := o.(*widget.Shadow); ok {
				shadowFound = true
			}
		}
		if !shadowFound {
			assert.Fail(t, "button should cast a shadow")
		}
	}
	{
		button := NewButton("Test", func() {})
		button.Importance = LowImportance
		for _, o := range test.LaidOutObjects(button) {
			if _, ok := o.(*widget.Shadow); ok {
				assert.Fail(t, "button with LowImportance should not create a shadow")
			}
		}
	}
}

func TestButtonRenderer_ApplyTheme(t *testing.T) {
	button := &Button{}
	render := test.WidgetRenderer(button).(*buttonRenderer)

	textSize := render.label.TextSize
	customTextSize := textSize
	test.WithTestTheme(t, func() {
		render.applyTheme()
		customTextSize = render.label.TextSize
	})

	assert.NotEqual(t, textSize, customTextSize)
}

func TestButtonRenderer_TapAnimation(t *testing.T) {
	test.NewApp()
	defer test.NewApp()

	button := NewButton("Hi", func() {})
	w := test.NewWindow(button)
	defer w.Close()
	w.Resize(fyne.NewSize(50, 50).Add(fyne.NewSize(10, 10)))
	button.Resize(fyne.NewSize(50, 50))

	test.ApplyTheme(t, test.NewTheme())
	button.Refresh()

	render1 := test.WidgetRenderer(button).(*buttonRenderer)
	test.Tap(button)
	button.tapAnim.Tick(0.5)
	test.AssertImageMatches(t, "button/tap_animation.png", w.Canvas().Capture())

	cache.DestroyRenderer(button)
	button.Refresh()

	render2 := test.WidgetRenderer(button).(*buttonRenderer)

	assert.NotEqual(t, render1, render2)

	test.Tap(button)
	button.tapAnim.Tick(0.5)
	test.AssertImageMatches(t, "button/tap_animation.png", w.Canvas().Capture())
}
