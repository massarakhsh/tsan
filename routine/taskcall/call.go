// taskcall - Задача мониторинга входящих звонков.
package taskcall

import (
	"bitbucket.org/961961/tsan/jone"
	"bitbucket.org/961961/tsan/repo"
	"bitbucket.org/961961/tsan/routine"
	"bitbucket.org/shaman/lik"
	"sync"
	"time"
)

//	Декриптор задачи мониторинга
type TaskCall struct {
	routine.Task
	callList   map[int]*sCall
	callSync	sync.Mutex
	callCallingCount	int
	callConnectedCount	int
	callLastUpdate	string
}

//	Дескриптор звонка
type sCall struct {
	Id      int
	At      time.Time
	State   string
	Phone   string
	IP      string
	PhoneTo string
	PinTo   int
	IdBell  lik.IDB
	Used    bool
	Moded   bool
}

//	Указатель дескриптора мониторинга
var Routine		*TaskCall

func GoIt() {
	Routine = &TaskCall{}
	Routine.Name = "calls"
	routine.RegisterTask(Routine)
	go Routine.run()
}

//	Остановить задачу
func CallerStop() {
	Routine.CallerStop()
}

// Опрос состояния звонков для пина pin, результат:
//	- количество звонков в очереди
//	- количество установленных соединений
//	- ID соединения на этом пине
func RequestCall(pin int) (int, int, lik.IDB) {
	return Routine.RequestCall(pin)
}

// Опрос состояния звонка контакта
// результат
func CallerProbe(id lik.IDB) string {
	return Routine.CallerProbe(id)
}

//	Основной процесси задачи
func (it *TaskCall) run() {
	it.callerInit()
	for !it.IsStoping() {
		call,_ := jone.DB.CalculeInt("SELECT COUNT(*) FROM calls WHERE type='incoming' and state='calling'")
		conn,_ := jone.DB.CalculeInt("SELECT COUNT(*) FROM calls WHERE type='incoming' and state='connected'")
		last,_ := jone.DB.CalculeString("SELECT MAX(updated_at) FROM calls WHERE type='incoming' and (state='calling' or state='connected' or state='disconnected')")
		if call != it.callCallingCount || conn != it.callConnectedCount || last != it.callLastUpdate {
			it.callCallingCount = call
			it.callConnectedCount = call
			it.callLastUpdate = last
			it.callerStep()
			time.Sleep(time.Millisecond * 100)
		} else {
			time.Sleep(time.Second * 1)
		}
	}
	it.OnStoped()
}

//	Остановить задачу
func (it *TaskCall) CallerStop() {
	it.OnStoping()
	for nt := 0; nt < 100 && !it.IsStoped(); nt++ {
		time.Sleep(10)
	}
}

// Опрос состояния звонков для пина pin, результат:
//	- количество звонков в очереди
//	- количество установленных соединений
//	- ID соединения на этом пине
func (it *TaskCall) RequestCall(pin int) (int, int, lik.IDB) {
	queue := 0
	connect := 0
	idb := lik.IDB(0)
	it.callSync.Lock()
	for _,scl := range it.callList {
		if scl.State == "calling" {
			queue++
		} else if scl.State == "connected" {
			connect++
			if scl.PinTo == pin {
				idb = scl.IdBell
			}
		}
	}
	it.callSync.Unlock()
	return queue, connect, idb
}

// Опрос состояния звонка контакта
// результат
func (it *TaskCall) CallerProbe(id lik.IDB) string {
	state := ""
	it.callSync.Lock()
	for _,scl := range it.callList {
		if scl.IdBell == id {
			state = scl.State
		}
	}
	it.callSync.Unlock()
	return state
}

//	Инициализация списка звонков
func (it *TaskCall) callerInit() {
	it.callList = make(map[int]*sCall)
	for _,elm := range jone.TableBell.Elms {
		if state := elm.GetString("call/state"); state == "calling" || state == "connected" {
			if idt := elm.GetInt("call/idt"); idt > 0 {
				scl := &sCall{ Id: idt, IdBell: elm.Id }
				if at := elm.GetInt("date"); at > 0 {
					scl.At = time.Unix(int64(at), 0)
				}
				scl.State = state
				if phone := jone.CalculateElmString(elm, "clientid/phone1"); phone != "" {
					scl.Phone = phone
				}
				scl.PinTo = elm.GetInt("pin")
				it.callList[idt] = scl
			}
		}
	}
}

//	Обновление списка звонков
func (it *TaskCall) callerStep() {
	it.callSync.Lock()
	it.readCalls()
	it.scanCalls()
	it.callSync.Unlock()
}

//	Прочитать список звонков
func (it *TaskCall) readCalls() {
	for _,sc := range it.callList {
		sc.Used = false
		sc.Moded = false
	}
	if listcalls := jone.DB.GetListAll("calls"); listcalls != nil {
		for nc := 0; nc < listcalls.Count(); nc++ {
			if call := listcalls.GetSet(nc); call != nil && call.GetString("type") == "incoming" {
				id := call.GetInt("id")
				state := call.GetString("state")
				scl,_ := it.callList[id]
				// new,calling,connected,disconnectod,end
				if scl != nil {
					scl.Used = true
				} else if state == "calling" || state == "connected" {
					scl = &sCall{Id: id, Used: true, Moded: true}
					it.callList[id] = scl
				}
				if scl != nil {
					at := time.Now()
					if dat, err := time.Parse("2006-01-02 15:04:05", call.GetString("updated_at")); err != nil {
						at = time.Now()
					} else if dat.Year() <= 2000 {
						at = time.Now()
					} else {
						at = dat.Add(-3 * time.Hour)
					}
					if at != scl.At {
						scl.At = at
						scl.Moded = true
					}
					if state != scl.State {
						scl.State = state
						scl.Moded = true
					}
					pinto := call.GetInt("request_pin")
					if pinto != scl.PinTo {
						scl.PinTo = pinto
						scl.Moded = true
					}
					if scl.Moded {
						phone := call.GetString("from_number")
						if match := lik.RegExParse(phone, ":(.*)@(.*)"); match != nil {
							phone = match[1]
							scl.IP = match[2]
						}
						scl.Phone = jone.NormalizePhone(phone)
						phoneto := call.GetString("request_number")
						if match := lik.RegExParse(phoneto, ":(.*)@(.*)"); match != nil {
							phoneto = match[1]
						}
						scl.PhoneTo = phoneto
					}
				}
			}
		}
	}
}

//	Сканировать список звонков
func (it *TaskCall) scanCalls() {
	for id,scl := range it.callList {
		if !scl.Used || scl.State == "disconnected" || scl.State == "end" {
			if elm := jone.TableBell.GetElm(scl.IdBell); elm != nil {
				elm.SetValue("", "call/state")
				elm.SetValue(int(scl.At.Unix()), "dateend")
			}
			delete(it.callList, id)
		} else if !scl.Moded {
		} else {
			elm := jone.TableBell.GetElm(scl.IdBell)
			if elm == nil {
				for idb,bel := range jone.TableBell.Elms {
					if bel.GetInt("call/idt") == id {
						scl.IdBell = idb
						elm = bel
						break
					}
				}
			}
			if elm == nil {
				elm = jone.TableBell.CreateElm()
				scl.IdBell = elm.Id
				elm.SetValue(id, "call/idt")
				elm.SetValue(int(scl.At.Unix()), "date")
				if cli := repo.SearchClient(scl.Phone); cli != nil {
					elm.SetValue(cli.Id, "clientid")
				} else {
					elm.SetValue(scl.Phone, "clientid/phone1")
				}
			}
			elm.SetValue(scl.State, "call/state")
			elm.SetValue(scl.IP, "call/ip")
			if scl.State == "connected" {
				elm.SetValue(scl.PinTo, "pin")
			}
		}
	}
}

