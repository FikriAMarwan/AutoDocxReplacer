package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"time"

	"github.com/goodsign/monday"
	"github.com/nguyenthenguyen/docx"
)

type TmpDoc struct {
	Name  string
	Count int
	Items []string
}

var (
	Tmp          map[string]TmpDoc
	TmpNumbering []string
)

func GetData(path string) (string, error) {
	info, err := os.Stat(path)
	if err != nil {
		if os.IsNotExist(err) {
			return "", err
		}
	}
	if filepath.Ext(path) != ".txt" {
		return "", fmt.Errorf("txt")
	}

	if info.IsDir() {
		return "", fmt.Errorf("Folder")
	}
	file, err := os.Open(path)
	if err != nil {
		return "", err
	}
	defer file.Close()
	result, err := ioutil.ReadAll(file)
	if err != nil {
		return "", err
	}
	return string(result), nil
}

func FindTmp() error {
	TmpNumbering = []string{}
	Tmp = make(map[string]TmpDoc)
	item := regexp.MustCompile("{{[a-z A-Z 0-9 _ . , / -]+}}")
	filepath.Walk("template/", func(path string, info os.FileInfo, err error) error {
		if !info.IsDir() && !strings.Contains(path, "~$") && strings.ToLower(filepath.Ext(path)) == ".docx" {
			var t TmpDoc
			r, _ := docx.ReadDocxFile(path)
			t.Name = strings.Replace(info.Name(), ".docx", "", -1)
			t.Items = item.FindAllString(r.Editable().GetContent(), -1)
			t.Count = len(t.Items)
			TmpNumbering = append(TmpNumbering, t.Name)
			Tmp[strings.Replace(info.Name(), ".docx", "", -1)] = t
			r.Close()
		}
		return nil
	})
	return nil
}

func SuratMulti(data string) error {
	list, err := FindContent(data)
	if err != nil {
		return err
	}
	if len(list) == 0 {
		return fmt.Errorf("Tidak Dapat Menemukan Data")
	}
	for _, v := range list {
		val, ok := Tmp[v[0]]
		if ok {
			r, err := docx.ReadDocxFile("template/" + v[0] + ".docx")
			if err != nil {
				return err
			}
			defer r.Close()
			doc := r.Editable()
			i := 1
			for _, val := range val.Items {
				wkt, ok := CheckVarTgl(val)
				if ok {
					doc.Replace(val, wkt, -1)
					continue
				}
				if strings.ToLower(val) == "{{no surat}}" {
					v[i] = v[i] + " / " + time.Now().Format("2006")
				}
				if strings.ToLower(val) == "{{jenis kelamin}}" {
					if strings.ToLower(v[i]) == "l" {
						v[i] = "Laki-Laki"
					}
					if strings.ToLower(v[i]) == "p" {
						v[i] = "Perempuan"
					}
				}
				if strings.ToLower(val) == "{{status kawin}}" {
					if strings.ToLower(v[i]) == "bk" {
						v[i] = "Belum Kawin"
					}
					if strings.ToLower(v[i]) == "k" {
						v[i] = "Kawin"
					}
				}
				doc.Replace(val, v[i], -1)
				i++
			}
			err = doc.WriteToFile(`hasil\` + v[0] + `_` + v[2] + "_" + time.Now().Format("02-01-2006") + `.docx`)
		}
	}
	return nil
}

func SuratSingle(data string, jenis string) error {
	var item []string
	item, err := FindContentSingle(data)
	if err != nil {
		return err
	}
	if len(item) == 0 {
		return fmt.Errorf("Tidak Dapat Menemukan Data")
	}
	val, ok := Tmp[jenis]
	if ok {
		r, err := docx.ReadDocxFile("template/" + jenis + ".docx")
		if err != nil {
			return err
		}
		defer r.Close()
		doc := r.Editable()
		i := 0
		for _, val := range val.Items {
			wkt, ok := CheckVarTgl(val)
			if ok {
				doc.Replace(val, wkt, -1)
				continue
			}
			if strings.ToLower(val) == "{{no surat}}" {
				item[i] = item[i] + " / " + time.Now().Format("2006")
			}
			if strings.ToLower(val) == "{{jenis kelamin}}" {
				if strings.ToLower(item[i]) == "l" {
					item[i] = "Laki-Laki"
				}
				if strings.ToLower(item[i]) == "p" {
					item[i] = "Perempuan"
				}
			}
			if strings.ToLower(val) == "{{status kawin}}" {
				if strings.ToLower(item[i]) == "bk" {
					item[i] = "Belum Kawin"
				}
				if strings.ToLower(item[i]) == "k" {
					item[i] = "Kawin"
				}
			}
			doc.Replace(val, item[i], -1)
			i++
		}
		err = doc.WriteToFile(`hasil\` + jenis + `_` + item[1] + "_" + time.Now().Format("02-01-2006") + `.docx`)
	}
	return nil
}

//Cari Konten Data Multiple
func FindContent(data string) ([][]string, error) {
	var list [][]string
	tmp := regexp.MustCompile(`: *([a-z A-Z 0-9 / . , -]+)`)
	s := strings.Split(data, "##")
	for _, v := range s {
		var item []string
		t := tmp.FindAllStringSubmatch(v, -1)
		for _, val := range t {
			item = append(item, val[1])
		}
		if len(item) > 0 {
			list = append(list, item)
		}
	}
	return list, nil
}

//Cari Konten Data Single
func FindContentSingle(data string) ([]string, error) {
	var item []string
	tmp := regexp.MustCompile(`: *([a-z A-Z 0-9 / . , -]+)`)
	t := tmp.FindAllStringSubmatch(data, -1)
	for _, val := range t {
		item = append(item, val[1])
	}
	return item, nil
}

//cek folder klo ngak da bikin
func CheckFolder() {
	if _, err := os.Stat("./hasil"); os.IsNotExist(err) {
		os.Mkdir("./hasil", os.ModeDir)
	}
	if _, err := os.Stat("./template"); os.IsNotExist(err) {
		os.Mkdir("./template", os.ModeDir)
	}
}

//cek template di folder template klo tidak ada warning untuk masukin
func CheckTemplate() string {
	var notexist string
	tmpList := []string{"SKU", "SKCK"}
	for _, v := range tmpList {
		if _, err := os.Stat("./template/" + v + ".docx"); os.IsNotExist(err) {
			notexist = notexist + v + "\n"
		}
	}
	return notexist
}

//Translate Variable strings.ReplaceAll(v, "_", " ")
func VarName(v string) string {
	v = strings.ReplaceAll(v, "{{", "")
	v = strings.ReplaceAll(v, "}}", "")
	return strings.Title(strings.ToLower(v))
}

func FindLongestStr(src []string) int {
	max := 0
	for _, v := range src {
		_, ok := CheckVarTgl(v)
		if len(v) > max && !ok {
			max = len(v)
		}
	}
	return max
}

func CheckVarTgl(v string) (string, bool) {
	switch v {
	case "{{d-m-y}}":
		return time.Now().Format("02 01 2006"), true
	case "{{d-month-y}}":
		return monday.Format(time.Now(), "2 January 2006", monday.LocaleIdID), true
	case "{{day-month-y}}":
		return monday.Format(time.Now(), "Monday, 2 January 2006", monday.LocaleIdID), true
	case "{{d-mon-y}}":
		return monday.Format(time.Now(), "2 Jan 2006", monday.LocaleIdID), true
	case "{{dd-mon-y}}":
		return monday.Format(time.Now(), "Mon, 2 Jan 2006", monday.LocaleIdID), true
	case "{{dd-month-y}}":
		return monday.Format(time.Now(), "Mon, 2 Januari 2006", monday.LocaleIdID), true
	}
	return v, false
}
