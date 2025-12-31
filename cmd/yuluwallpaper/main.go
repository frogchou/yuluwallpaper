package main

import (
	_ "embed"
	"errors"
	"image/color"
	"log"
	"os"
	"path/filepath"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/driver/desktop"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"

	wallapp "yuluwallpaper/internal/app"
	"yuluwallpaper/internal/autostart"
	"yuluwallpaper/internal/config"
	"yuluwallpaper/internal/logger"
)

const appID = "com.yulu.wallpaper"
const appDisplayName = "Yulu Wallpaper"

//go:embed assets/fonts/NotoSansCJKsc-Regular.otf
var appFontData []byte

var appFontResource = fyne.NewStaticResource("NotoSansCJKsc-Regular.otf", appFontData)

func main() {
	logPath, err := config.LogPath()
	if err == nil {
		_ = logger.Init(logPath)
	}
	defer logger.Close()

	cfg, err := config.Load()
	if err != nil {
		log.Printf("config load failed: %v", err)
		cfg = config.Default()
	}

	assetsDir, err := config.AssetsDir()
	if err != nil {
		log.Printf("assets dir failed: %v", err)
		assetsDir = os.TempDir()
	}

	service := wallapp.NewService(cfg, assetsDir)
	go service.Run()

	fyneApp := app.NewWithID(appID)
	fyneApp.Settings().SetTheme(appTheme{})
	settingsUI := newSettingsUI(fyneApp, &cfg, logPath, func(newCfg config.Config) {
		cfg = newCfg
		service.UpdateConfig(newCfg)
	})

	menu := fyne.NewMenu("",
		fyne.NewMenuItem("设置", func() {
			settingsUI.ApplyConfig(cfg)
			settingsUI.Show()
		}),
		fyne.NewMenuItem("立即刷新", func() {
			service.RequestRefresh()
		}),
		fyne.NewMenuItemSeparator(),
		newQuitMenuItem(func() {
			service.Stop()
			fyneApp.Quit()
		}),
	)

	if desktopApp, ok := fyneApp.(desktop.App); ok {
		desktopApp.SetSystemTrayMenu(menu)
		desktopApp.SetSystemTrayIcon(appIconResource())
	}

	fyneApp.Run()
}

type settingsUI struct {
	window         fyne.Window
	intervalSelect *widget.Select
	layoutSelect   *widget.Select
	autoStartCheck *widget.Check

	labelToMinutes map[string]int
	logPath        string

	onApply    func(config.Config)
	currentCfg *config.Config
}

func newSettingsUI(fyneApp fyne.App, cfg *config.Config, logPath string, onApply func(config.Config)) *settingsUI {
	ui := &settingsUI{
		window:     fyneApp.NewWindow("壁纸设置"),
		logPath:    logPath,
		onApply:    onApply,
		currentCfg: cfg,
	}

	intervalOptions := config.IntervalOptions()
	labels := make([]string, 0, len(intervalOptions))
	labelToMinutes := make(map[string]int, len(intervalOptions))
	for _, opt := range intervalOptions {
		labels = append(labels, opt.Label)
		labelToMinutes[opt.Label] = opt.Minutes
	}

	ui.intervalSelect = widget.NewSelect(labels, nil)
	ui.layoutSelect = widget.NewSelect([]string{"平铺", "拉伸", "适应", "填充", "居中"}, nil)
	ui.autoStartCheck = widget.NewCheck("开机自启动", nil)

	form := &widget.Form{
		Items: []*widget.FormItem{
			{Text: "更换周期", Widget: ui.intervalSelect},
			{Text: "桌面布局", Widget: ui.layoutSelect},
		},
	}

	title := canvas.NewText("余露壁纸", theme.ForegroundColor())
	title.TextStyle = fyne.TextStyle{Bold: true}
	title.TextSize = 20
	subtitle := canvas.NewText("风起时更换，心安处久居。", theme.DisabledColor())
	subtitle.TextStyle = fyne.TextStyle{Italic: true}
	subtitle.TextSize = 12
	header := container.NewVBox(title, subtitle)

	formCard := widget.NewCard("基础设置", "让桌面在时光里悄然更迭", form)
	autoCard := widget.NewCard("启动方式", "静默守候，需要时即现", container.NewVBox(ui.autoStartCheck))

	saveBtn := widget.NewButton("保存", func() {
		newCfg, err := ui.configFromInputs()
		if err != nil {
			dialog.ShowError(err, ui.window)
			return
		}
		if err := ui.applyAutoStart(*ui.currentCfg, newCfg); err != nil {
			dialog.ShowError(err, ui.window)
			return
		}
		if err := config.Save(newCfg); err != nil {
			dialog.ShowError(err, ui.window)
			return
		}
		ui.onApply(newCfg)
		ui.window.Hide()
	})
	cancelBtn := widget.NewButton("取消", func() {
		ui.window.Hide()
	})
	buttons := container.NewHBox(saveBtn, cancelBtn)

	content := container.NewVBox(
		header,
		widget.NewSeparator(),
		formCard,
		autoCard,
		buttons,
	)
	ui.window.SetContent(container.NewPadded(content))
	ui.window.Resize(fyne.NewSize(420, 320))
	ui.window.SetCloseIntercept(func() {
		ui.window.Hide()
	})
	ui.labelToMinutes = labelToMinutes
	ui.window.SetIcon(appIconResource())

	ui.ApplyConfig(*cfg)
	return ui
}

func (ui *settingsUI) ApplyConfig(cfg config.Config) {
	label := config.IntervalLabel(cfg.IntervalMinutes)
	if label != "" {
		ui.intervalSelect.SetSelected(label)
	}

	switch cfg.Layout {
	case config.LayoutTile:
		ui.layoutSelect.SetSelected("平铺")
	case config.LayoutFit:
		ui.layoutSelect.SetSelected("适应")
	case config.LayoutFill:
		ui.layoutSelect.SetSelected("填充")
	case config.LayoutCenter:
		ui.layoutSelect.SetSelected("居中")
	default:
		ui.layoutSelect.SetSelected("拉伸")
	}

	ui.autoStartCheck.SetChecked(cfg.AutoStart)
}

func (ui *settingsUI) Show() {
	ui.window.Show()
	ui.window.RequestFocus()
}

func (ui *settingsUI) configFromInputs() (config.Config, error) {
	minutes, ok := ui.labelToMinutes[ui.intervalSelect.Selected]
	if !ok {
		return config.Config{}, errors.New("请选择更换周期")
	}

	layout := config.LayoutStretch
	switch ui.layoutSelect.Selected {
	case "平铺":
		layout = config.LayoutTile
	case "适应":
		layout = config.LayoutFit
	case "填充":
		layout = config.LayoutFill
	case "居中":
		layout = config.LayoutCenter
	}

	return config.Config{
		IntervalMinutes: minutes,
		Layout:          layout,
		AutoStart:       ui.autoStartCheck.Checked,
	}, nil
}

func (ui *settingsUI) applyAutoStart(current, next config.Config) error {
	if current.AutoStart == next.AutoStart {
		return nil
	}
	if next.AutoStart {
		exe, err := os.Executable()
		if err != nil {
			return err
		}
		return autostart.Enable(appDisplayName, exe, ui.logPath)
	}
	return autostart.Disable(appDisplayName)
}

func appIconResource() fyne.Resource {
	if data, err := readTrayIcon(); err == nil {
		return fyne.NewStaticResource("favicon.png", data)
	}
	return theme.FyneLogo()
}

func readTrayIcon() ([]byte, error) {
	paths := []string{"favicon.png"}
	if exePath, err := os.Executable(); err == nil {
		paths = append([]string{filepath.Join(filepath.Dir(exePath), "favicon.png")}, paths...)
	}
	for _, path := range paths {
		data, err := os.ReadFile(path)
		if err == nil && len(data) > 0 {
			return data, nil
		}
	}
	return nil, errors.New("tray icon not found")
}

func newQuitMenuItem(action func()) *fyne.MenuItem {
	item := fyne.NewMenuItem("退出", action)
	item.IsQuit = true
	return item
}

type appTheme struct{}

func (appTheme) Color(name fyne.ThemeColorName, variant fyne.ThemeVariant) color.Color {
	return theme.DefaultTheme().Color(name, variant)
}

func (appTheme) Font(style fyne.TextStyle) fyne.Resource {
	if len(appFontData) == 0 {
		return theme.DefaultTheme().Font(style)
	}
	return appFontResource
}

func (appTheme) Icon(name fyne.ThemeIconName) fyne.Resource {
	return theme.DefaultTheme().Icon(name)
}

func (appTheme) Size(name fyne.ThemeSizeName) float32 {
	return theme.DefaultTheme().Size(name)
}
