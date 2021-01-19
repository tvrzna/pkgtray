package main

import (
	"log"
	"os"
	"strconv"
	"time"

	"github.com/mattn/go-gtk/gtk"
)

type getUpdatesFce func() int

type trayIcon struct {
	statusIcon *gtk.StatusIcon
	conf       *config
	getUpdates getUpdatesFce
}

func createTrayIcon(conf *config, refresh getUpdatesFce) *trayIcon {
	gtk.Init(nil)

	trayIcon := &trayIcon{gtk.NewStatusIcon(), conf, refresh}
	trayIcon.init()

	return trayIcon
}

func (t *trayIcon) start() {
	gtk.Main()
}

func (t *trayIcon) init() {
	t.statusIcon.Connect("popup-menu", t.makeMenu)
	t.statusIcon.SetVisible(true)
	t.updateIcon(0)
}

func (t *trayIcon) updateIcon(updates int) {
	var msg string
	if updates > 0 {
		t.statusIcon.SetFromStock(gtk.STOCK_REFRESH)
		msg = strconv.Itoa(updates) + " updates available."
	} else if updates == 0 {
		t.statusIcon.SetFromStock(gtk.STOCK_OK)
		msg = "System is up to date."
	} else {
		t.statusIcon.SetFromStock(gtk.STOCK_CLOSE)
		msg = "Checking for updates failed."
	}

	t.statusIcon.SetTooltipText(msg)
	t.statusIcon.SetTitle(msg)

	log.Println(msg)
}

func (t *trayIcon) makeMenu() {
	menu := gtk.NewMenu()

	refresh := gtk.NewImageMenuItemFromStock(gtk.STOCK_REFRESH, &gtk.AccelGroup{})
	refresh.Connect("activate", t.doRefresh, "Refresh")
	refresh.SetLabel("Refresh")
	refresh.Show()
	menu.Append(refresh)

	about := gtk.NewImageMenuItemFromStock(gtk.STOCK_INFO, &gtk.AccelGroup{})
	about.Connect("activate", t.doAbout, "About")
	about.SetLabel("About")
	about.Show()
	menu.Append(about)

	close := gtk.NewImageMenuItemFromStock(gtk.STOCK_CLOSE, &gtk.AccelGroup{})
	close.Connect("activate", t.doClose, "Quit")
	close.SetLabel("Quit")
	close.Show()
	menu.Append(close)

	menu.Popup(nil, nil, gtk.StatusIconPositionMenu, t.statusIcon, 1, gtk.GetCurrentEventTime())
}

func (t *trayIcon) doRefresh() {
	log.Print("refresh")
	updates := t.getUpdates()
	t.updateIcon(updates)
}

func (t *trayIcon) doAbout() {
	msg := gtk.NewMessageDialog(nil,
		gtk.DIALOG_MODAL|gtk.DIALOG_DESTROY_WITH_PARENT,
		gtk.MESSAGE_INFO,
		gtk.BUTTONS_OK,
		getVersion())
	msg.SetTitle("About pkgtray")
	msg.Run()
	msg.Destroy()
}

func (t *trayIcon) doClose() {
	log.Print("close")
	gtk.MainQuit()
	os.Exit(0)
}

func (t *trayIcon) processRefreshLoop() {
	go func() {
		for true {
			go t.doRefresh()
			time.Sleep(time.Duration(t.conf.intervalSeconds) * time.Second)
		}
	}()
}
