package main

import (
	"fmt"
	"log"
	"regexp"
	"runtime"
	"strings"

	"github.com/jroimartin/gocui"
	runewidth "github.com/mattn/go-runewidth"
)

var (
	jenis string
	max   int
)

func main() {
	CheckFolder()
	err := FindTmp()
	if err != nil {
		log.Panicln(err)
	}
	g, err := gocui.NewGui(gocui.OutputNormal)
	if err != nil {
		log.Panicln(err)
	}
	defer g.Close()
	g.SetManagerFunc(layout)
	if runtime.GOOS == "windows" && runewidth.IsEastAsian() {
		g.ASCII = true
	}
	g.Mouse = true
	if err := keybindings(g); err != nil {
		log.Panicln(err)
	}
	if err := g.MainLoop(); err != nil && err != gocui.ErrQuit {
		log.Panicln(err)
	}
}

//layout stuff
func layout(g *gocui.Gui) error {
	g.Cursor = false
	maxX, maxY := g.Size()
	g.DeleteView("BtnPrint")
	if v, err := g.SetView("Title", maxX/2-6, 0, maxX/2+5, 2); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		v.Frame = false
		fmt.Fprintln(v, "AUTO SURAT")

	}
	if _, err := g.SetView("Background", 0, 0, maxX-1, maxY-1); err != nil {
	}
	g.SetViewOnBottom("Background")
	if v, err := g.SetView("Menu", maxX/2-10, 5, maxX/2+22, 11); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		v.Frame = false
		v.Highlight = true
		v.SelBgColor = gocui.ColorGreen
		v.SelFgColor = gocui.ColorBlack
		fmt.Fprintln(v, "1.Isi Data Surat")
		fmt.Fprintln(v, "2.Check Template Surat")
		fmt.Fprintln(v, "3.Exit")
	}

	return nil
}

func layoutIns(g *gocui.Gui) error {
	g.Cursor = true
	maxX, maxY := g.Size()
	if v, err := g.SetView("ListJenis", 0, 0, maxX*1/5, maxY-1); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		v.Highlight = true
		v.SelBgColor = gocui.ColorGreen
		v.SelFgColor = gocui.ColorBlack
		for i, s := range TmpNumbering {
			fmt.Fprintln(v, fmt.Sprintf("%v. %v", i+1, s))
		}
		fmt.Fprintln(v, "   BACK")
	}
	if _, err := g.SetView("FormBlank", maxX*1/5, 0, maxX-1, maxY-1); err != nil {
	}
	g.SetViewOnBottom("FormBlank")
	return nil
}

func InputForm(g *gocui.Gui, v *gocui.View) error {
	var l string
	var err error
	maxX, maxY := g.Size()
	if _, err := g.SetCurrentView(v.Name()); err != nil {
		return err
	}

	_, cy := v.Cursor()
	if l, err = v.Line(cy); err != nil {
		l = ""
	}
	g.DeleteView("Form")
	if l == "   BACK" {
		g.DeleteView("ListTmp")
		err := ChangeView(g, layout)
		if err != nil {
			return err
		}
		return nil
	}

	num := regexp.MustCompile("[0-9]+. ")
	l = strings.Replace(l, num.FindString(l), "", -1)
	if v, err := g.SetView("Form", maxX*1/5, 0, maxX-1, maxY-3); err != nil {
		v.Editable = true
		v.Title = "Masukan Data"
		val, ok := Tmp[l]
		if ok {
			max = FindLongestStr(val.Items)
			jenis = val.Name
			for _, t := range val.Items {
				_, ok := CheckVarTgl(t)
				if !ok {
					TtkDua := strings.Repeat(" ", max-len(t)) + " : "
					fmt.Fprintln(v, VarName(t)+TtkDua)
				}
			}
			setCurrentViewOnTop(g, "Form")
		}
	}
	if l != "" {
		if print, err := g.SetView("BtnPrint", maxX*1/5, maxY-3, maxX-1, maxY-1); err != nil {
			print.Highlight = true
			print.SelBgColor = gocui.ColorBlack
			print.SelFgColor = gocui.ColorGreen
			fmt.Fprintln(print, "                                            Print")
		}
	}
	return nil
}

func layoutTmp(g *gocui.Gui) error {
	maxX, maxY := g.Size()
	if v, err := g.SetView("ListTmp", 0, 0, maxX*1/4, maxY-1); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		v.Title = "List Template"
		v.Highlight = true
		v.SelBgColor = gocui.ColorGreen
		v.SelFgColor = gocui.ColorBlack
		for i, s := range TmpNumbering {
			fmt.Fprintln(v, fmt.Sprintf("%v. %v", i+1, s))
		}
		fmt.Fprintln(v, "   BACK")
	}
	if _, err := g.SetView("FormBlank", maxX*1/4, 0, maxX-1, maxY-1); err != nil {
	}
	if print, err := g.SetView("BtnUpdate", maxX*1/4, maxY-3, maxX-1, maxY-1); err != nil {
		print.Highlight = true
		print.SelBgColor = gocui.ColorBlack
		print.SelFgColor = gocui.ColorGreen
		fmt.Fprintln(print, "                                            Update")
	}
	return nil
}

func TmpForm(g *gocui.Gui, v *gocui.View) error {
	var l string
	var err error
	maxX, maxY := g.Size()
	if _, err := g.SetCurrentView(v.Name()); err != nil {
		return err
	}

	_, cy := v.Cursor()
	if l, err = v.Line(cy); err != nil {
		l = ""
	}
	g.DeleteView("VarForm")
	if l == "   BACK" {
		g.DeleteView("ListTmp")
		err := ChangeView(g, layout)
		if err != nil {
			return err
		}
		return nil
	}
	num := regexp.MustCompile("[0-9]+. ")
	l = strings.Replace(l, num.FindString(l), "", -1)
	val, ok := Tmp[l]
	if ok {
		if v, err := g.SetView("VarForm", maxX*1/4, 0, maxX-1, maxY-3); err != nil {
			if err != gocui.ErrUnknownView {
				return err
			}
			v.Title = "Variable Template"
			for _, t := range val.Items {
				fmt.Fprintln(v, "-"+VarName(t))
			}
			fmt.Fprintln(v, fmt.Sprintf("\nTotal Variable : %v", len(val.Items)))
		}
		// g.DeleteView("VarName")
		// if v, err := g.SetView("VarName", maxX*1/4+30, 0, maxX-1, maxY-3); err != nil {
		// 	v.Title = "Variable Name"
		// 	for _, t := range val.Items {
		// 		fmt.Fprintln(v, "-"+VarName(t))
		// 	}
		// }
	}
	return nil
}

//fungsi dan utilities
func keybindings(g *gocui.Gui) error {
	if err := g.SetKeybinding("", gocui.KeyCtrlC, gocui.ModNone, quit); err != nil {
		return err
	}

	if err := g.SetKeybinding("Menu", gocui.MouseLeft, gocui.ModNone, ChangeIns); err != nil {
		return err
	}

	if err := g.SetKeybinding("ListJenis", gocui.MouseLeft, gocui.ModNone, InputForm); err != nil {
		return err
	}
	if err := g.SetKeybinding("ListTmp", gocui.MouseLeft, gocui.ModNone, TmpForm); err != nil {
		return err
	}
	if err := g.SetKeybinding("BtnPrint", gocui.MouseLeft, gocui.ModNone, ToDocx); err != nil {
		return err
	}
	if err := g.SetKeybinding("BtnUpdate", gocui.MouseLeft, gocui.ModNone, UpdateTmp); err != nil {
		return err
	}
	if err := g.SetKeybinding("Form", gocui.KeyTab, gocui.ModNone, Tabs); err != nil {
		return err
	}

	return nil
}

func UpdateTmp(g *gocui.Gui, v *gocui.View) error {
	err := FindTmp()
	if err != nil {
		return err
	}
	v, _ = g.SetCurrentView("ListTmp")
	TmpForm(g, v)
	// g.DeleteView("ListTmp")
	// g.Update(layoutTmp)
	return nil
}

func Status(g *gocui.Gui, v *gocui.View, s string) {
	g.DeleteView("Status")
	maxX, maxY := g.Size()
	if v, err := g.SetView("Status", 0, maxY-3, maxX*1/5, maxY-1); err != nil {
		fmt.Fprintln(v, s)
	}
}

func Tabs(g *gocui.Gui, v *gocui.View) error {
	_, y := v.Cursor()
	err := v.SetCursor(max-1, y+1)
	if err != nil {
		return err
	}
	return nil
}

func ToDocx(g *gocui.Gui, v *gocui.View) error {
	v, err := g.SetCurrentView("Form")
	if err != nil {
		return err
	}
	err = SuratSingle(v.Buffer(), jenis)
	if err != nil {
		return err
	}
	g.DeleteView("Form")
	v, _ = g.SetCurrentView("ListJenis")
	InputForm(g, v)
	return nil
}

func ChangeIns(g *gocui.Gui, v *gocui.View) error {
	var l string
	var err error

	if _, err := g.SetCurrentView(v.Name()); err != nil {
		return err
	}

	_, cy := v.Cursor()
	if l, err = v.Line(cy); err != nil {
		l = ""
	}
	switch l {
	case "1.Isi Data Surat":
		err := ChangeView(g, layoutIns)
		if err != nil {
			return err
		}
		break
	case "2.Check Template Surat":
		err := ChangeView(g, layoutTmp)
		if err != nil {
			return err
		}
		break
	case "3.Exit":
		return gocui.ErrQuit
	}
	return nil
}

func ChangeView(g *gocui.Gui, manager func(*gocui.Gui) error) error {
	g.DeleteView(g.CurrentView().Name())
	g.SetManagerFunc(manager)
	if err := keybindings(g); err != nil {
		return err
	}
	return nil
}

func quit(g *gocui.Gui, v *gocui.View) error {
	return gocui.ErrQuit
}

func setCurrentViewOnTop(g *gocui.Gui, name string) (*gocui.View, error) {
	if _, err := g.SetCurrentView(name); err != nil {
		return nil, err
	}
	return g.SetViewOnTop(name)
}
