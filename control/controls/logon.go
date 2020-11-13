package controls

import (
	"bitbucket.org/961961/tsan/control"
	"bitbucket.org/961961/tsan/fancy"
	"bitbucket.org/961961/tsan/jone"
	"bitbucket.org/961961/tsan/repo"
	"bitbucket.org/961961/tsan/show"
	"bitbucket.org/shaman/lik"
	"bitbucket.org/shaman/lik/likbase"
	"bitbucket.org/shaman/lik/likdom"
	"strings"
)

//	Дескриптор огна входа в систему
type LogonControl struct {
	control.DataControl
	fancy.DataFancy
	Diagnosis string
	ProLogin  string
}

//	Обработчик событий окна
type dealLogonExecute struct {
	It	*LogonControl
}
func (it *dealLogonExecute) Run(rule *repo.DataRule, cmd string, data lik.Seter) {
	it.It.DoExecute(rule, cmd, data)
}

//	Конструктор дескриптора окна
func BuildLogon(rule *repo.DataRule, id lik.IDB) *LogonControl {
	it := &LogonControl{}
	it.Mode = "logon"
	it.FancyInitialize("logon","member","logon")
	it.ItExecute = &dealLogonExecute{it}
	return it
}

//	Обработка маршализации
func (it *LogonControl) DoMarshal(rule *repo.DataRule) {
}

//	Обраюотка событий в окне
func (it *LogonControl) DoExecute(rule *repo.DataRule, cmd string, data lik.Seter) {
	if cmd == "all" || cmd == "logon" {
		it.DoExecute(rule, rule.Shift(), data)
	} else if cmd == "showform" {
		it.cmdRefresh(rule)
	} else if cmd == "logout" {
		rule.SessionLogin(0, 0, true)
		rule.ItSession.ExitLogin = true
		rule.OnChangeData()
	} else if cmd == "write" {
		it.cmdProbeLogin(rule, data)
		if rule.IsLogin() {
			rule.SetGoPart("/bell")
		}
	} else if lik.RegExCompare(cmd, "^("+jone.ItDispatch+"|"+jone.ItRealtor+"|"+jone.ItManager+"|"+jone.ItAdmin+")$") {
		if rule.ICanRole(cmd) {
			rule.ItSession.RoleMember = cmd
			rule.SetGoPart("/bell")
		}
	}
}

//	Отображение окна
func (it *LogonControl) BuildShow(rule *repo.DataRule, sx int, sy int) likdom.Domer {
	it.SetSize(sx, sy)
	var div likdom.Domer
	if !rule.IsLogin() {
		div = it.BuildDataLogin(rule)
	} else if rule.IsTechno {
		div = it.BuildDataTechno(rule)
	}
	return div
}

//	Отображение входа в систему
func (it *LogonControl) BuildDataLogin(rule *repo.DataRule) likdom.Domer {
	div := likdom.BuildDivClassId("roll_data", "logon_data")
	td := div.BuildTableClass("fill").BuildTrTd("align=center", "valign=middle")
	divin := td.BuildDiv("align=center", "id=div_login")
	divin.AppendItem(show.BuildFancyForm(it.Main,"logon"))
	return div
}

//	Отображение технологического режима
func (it *LogonControl) BuildDataTechno(rule *repo.DataRule) likdom.Domer {
	div := likdom.BuildDivClassId("roll_data", "logon_data")
	divin := div.BuildTableClass("fill").BuildTrTd().BuildDiv("align=center")
	divin.BuildString("ТЕХНОЛОГИЧЕСКИЙ РЕЖИМ<br>")
	divin.BuildString("Идут внутренние работы, извините за временные неудобства<br>")
	return div
}

//	Обновить окно
func (it *LogonControl) cmdRefresh(rule *repo.DataRule) {
	login := ""
	password := ""
	pin := 0;
	store := 0
	if rule.ItSession.WaitLogin && !rule.ItSession.ExitLogin {
		if idop := lik.StrToInt(rule.GetContext("auto_login_id")); idop > 0 {
			if operator := jone.TableMember.GetElm(lik.IDB(idop)); operator != nil {
				login = jone.CalculateElmString(operator,"login")
				password = "********"
				pin = jone.CalculateElmInt(operator,"pin")
				store = 1
			}
		}
	}
	if login == "" {
		login = it.ProLogin
	}
	it.FormClear()
	it.SetTitle(rule,fancy.FunShow,"Вход в систему")
	it.Form.Items.AddItems(
		lik.BuildSet("label=Логин", "name=login", "value", login),
		lik.BuildSet("label=Пароль", "name=password", "type=password", "value", password),
		lik.BuildSet("label=Гарнитура", "name=pin", "format=number", "value", pin),
		lik.BuildSet("label=Запомнить", "name=store", "type=checkbox", "value", store > 0),
	)
	if it.Diagnosis != "" {
		it.Form.Items.AddItems(
			lik.BuildSet("type=html", "value", "<font color=#a00 weight=bold>" + it.Diagnosis + "</font>"),
		)
	}
	it.AddTitleToolText(rule,"Войти в систему","function_fancy_form_write")
	it.AddTitleToolText(rule,"Отменить","function_fancy_form_cancel")
	it.ShowForm(rule)
}

//	Проверка возможности входа
func (it *LogonControl) cmdProbeLogin(rule *repo.DataRule, data lik.Seter) {
	members := []*likbase.ItElm{}
	var admin *likbase.ItElm
	for _, elm := range(jone.TableMember.Elms) {
		if jone.CalculateElmString(elm,"login") == "admin" {
			admin = elm
		}
		members = append(members, elm)
	}
	if admin == nil {
		admin = jone.TableMember.CreateElm()
		jone.SetElmValue(admin,"Администратор","family")
		jone.SetElmValue(admin,"admin","login")
		jone.SetElmValue(admin, lik.GetMD5Hash("admin"),"password")
		jone.SetElmValue(admin, jone.ItAdmin,"role")
	}
	jone.SetElmValue(admin, jone.ItAdmin,"role")
	it.Diagnosis = ""
	it.ProLogin = lik.StringFromXS(data.GetString("login"))
	passtext := lik.StringFromXS(data.GetString("password"))
	password := lik.GetMD5Hash(passtext)
	pin := lik.StrToInt(lik.StringFromXS(data.GetString("pin")))
	str := lik.StringFromXS(data.GetString("store"))
	store := (str == "true")
	islogin := false
	for _, member := range (members) {
		if strings.ToLower(jone.CalculateElmString(member,"login")) == strings.ToLower(it.ProLogin) {
			islogin = true
			if jone.CalculateElmString(member,"password") == password {
				rule.SessionLogin(member.Id, pin, store)
				return
			}
		}
	}
	if !islogin {
		it.Diagnosis = "Такой логин не существует"
	} else {
		it.Diagnosis = "Неверный пароль"
	}
	rule.OnChangeData()
}

