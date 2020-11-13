//	Система РИЭЛТОР.
package main

import (
	"github.com/massarakhsh/tsan/front"
	"bitbucket.org/961961/tsan/jone"
	"bitbucket.org/961961/tsan/one"
	"bitbucket.org/961961/tsan/repo"
	"bitbucket.org/961961/tsan/routine/taskcall"
	"bitbucket.org/961961/tsan/routine/taskexport/avito"
	"bitbucket.org/961961/tsan/routine/taskexport/cian"
	"bitbucket.org/961961/tsan/routine/taskexport/domclick"
	"bitbucket.org/961961/tsan/routine/taskexport/tsan"
	"bitbucket.org/961961/tsan/routine/taskexport/yandex"
	"bitbucket.org/961961/tsan/show"
	"bitbucket.org/shaman/lik"
	"bitbucket.org/shaman/lik/likapi"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"
)

//	Имя PID-файла
const PidFile = "var/tsan.pid"

//	Получить и обработать список аргументов
func getArgs() bool {
	args, ok := lik.GetArgs(os.Args[1:])
	if val := args.GetInt("port"); val > 0 {
		repo.HostPort = val
	}
	if val := args.GetString("serv"); val != "" {
		repo.HostServ = val
	}
	if val := args.GetString("base"); val != "" {
		repo.HostBase = val
	}
	if val := args.GetString("user"); val != "" {
		repo.HostUser = val
	}
	if val := args.GetString("pass"); val != "" {
		repo.HostPass = val
	}
	if val := args.GetString("signal"); val != "" {
		repo.HostSignal = val
	}
	if len(repo.HostBase) <= 0 {
		fmt.Println("HostBase name must be present")
		ok = false
	}
	if !ok {
		fmt.Println("Usage: tsan [-key val | --key=val]...")
		fmt.Println("signal  - signal to process: T(erminate), (sto)P, C(ontinue)")
		fmt.Println("port    - port value (80)")
		fmt.Println("serv    - Database server")
		fmt.Println("base    - Database name")
		fmt.Println("user    - Database user")
		fmt.Println("pass    - Database pass")
	}
	return ok
}

//	Роутер запросов
func router(w http.ResponseWriter, r *http.Request) {
	uri := r.RequestURI
	isfront := lik.RegExCompare(uri, "front")
	ismarshal := lik.RegExCompare(uri, "marshal")
	isfile := lik.RegExCompare(uri, "\\.(js|css|htm|html|ico|gif|png|jpg|jpeg|pdf|csv|xml|doc|docx|xls|xlsx)")
	if isfile && likapi.ProbeRouteFile(w, r, uri) {
		return
	}
	var page *repo.DataPage
	if sp := lik.StrToInt(likapi.GetParm(r, "_sp")); sp > 0 {
		if pager := likapi.FindPage(sp); pager != nil {
			page = pager.(repo.DataPager).GetItPage()
		}
	}
	if match := lik.RegExParse(uri, "([^/]+\\.pix)"); match != nil {
		img := show.BuildPixRast(match[1])
		likapi.RouteRast(w, 200, img)
		return
	}
	if page == nil {
		if isfile {
			return
		}
		if time.Now().Sub(repo.TimeStart) > time.Second*15 && ismarshal {
			return
		}
		page = repo.StartPage(uri)
	} else if lik.StrToInt(likapi.GetParm(r, "_tp")) > 0 {
		page = repo.ClonePage(page)
	}
	rule := repo.BindRule(page)
	rule.LoadRequest(r)
	var response lik.Seter
	if isfront {
		rule.Shift()
		response = front.FrontExecute(rule)
	} else if ismarshal {
		rule.Shift()
		response = front.MarshalExecute(rule)
	}
	if rule.IsJson {
		likapi.RouteCookies(w, rule.GetAllCookies())
		likapi.RouteJson(w, 200, response, rule.ResultFormat)
	} else if rule.IsXml {
		likapi.RouteXml(w, 200, response.GetList("xml"))
	} else if rule.IsCsv {
		likapi.RouteCsv(w, 200, response.GetList("csv"), "^")
	}
	if !rule.IsJson && !rule.IsXml && !rule.IsCsv {
		rc, html := front.FrontPage(rule)
		likapi.RouteCookies(w, rule.GetAllCookies())
		likapi.RouteHtml(w, rc, html.ToString())
	}
}

//	Основная точка входа
func main() {
	one.InitializeABC()
	pidgo := getActiveProcess()
	repo.HostPid = syscall.Getpid()
	repo.TimeStart = time.Now()
	if fi, err := os.Lstat("tsan.bin"); err == nil {
		repo.TimeBin = fi.ModTime()
	}
	lik.SetLevelInf()
	lik.SayError("System started")
	if !getArgs() {
		return
	}
	cmd := ""
	if repo.HostSignal != "" {
		cmd = strings.ToUpper(repo.HostSignal[:1])
	}
	if pidgo > 0 {
		if prc, err := os.FindProcess(pidgo); prc != nil && err == nil {
			if cmd == "S" {
				prc.Signal(syscall.Signal(23))
				lik.SayWarning("Send stop process")
			} else if cmd == "C" {
				prc.Signal(syscall.Signal(25))
				lik.SayWarning("Send continue process")
			} else {
				prc.Signal(syscall.SIGTERM)
				lik.SayWarning("Send term process")
			}
		} else {
			lik.SayWarning("Proceess not found")
		}
		setActiveProcess(0)
	}
	if cmd != "" {
		return
	}

	setActiveProcess(repo.HostPid)
	repo.ChanSignal = make(chan os.Signal, 1)
	signal.Notify(repo.ChanSignal, syscall.SIGKILL, syscall.SIGTERM, syscall.Signal(23), syscall.Signal(25))
	go waitSignal()

	jone.GoIt(repo.HostServ, repo.HostBase, repo.HostUser, repo.HostPass)
	one.GoIt(repo.HostServ, repo.HostBase, repo.HostUser, repo.HostPass)
	repo.GoIt(repo.HostServ, repo.HostBase, repo.HostUser, repo.HostPass)
	time.Sleep(250 * time.Millisecond)
	avito.GoIt()
	yandex.GoIt()
	cian.GoIt()
	tsan.GoIt()
	domclick.GoIt()
	taskcall.GoIt()

	http.HandleFunc("/", router)
	if err := http.ListenAndServe(":"+fmt.Sprint(repo.HostPort), nil); err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}

//	Получить код активного процесса из PID-файла
func getActiveProcess() int {
	pid := 0
	if data, err := ioutil.ReadFile(PidFile); err == nil {
		pid = lik.StrToInt(string(data))
	}
	return pid
}

//	Записать код процесса в PID-файл
func setActiveProcess(pid int) {
	var data []byte
	if pid > 0 {
		data = []byte(lik.IntToStr(pid))
	}
	ioutil.WriteFile(PidFile, data, 0777)
}

//	Процесс ожидания и обработки сигналов
func waitSignal() {
	for {
		signal := <-repo.ChanSignal
		if signal == syscall.Signal(23) {
			repo.ToPause = true
		} else if signal == syscall.Signal(25) {
			repo.ToPause = false
		} else {
			jone.StopBase()
			time.Sleep(1000 * time.Millisecond)
			os.Exit(3)
		}
	}
}
