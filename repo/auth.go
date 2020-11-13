package repo

import (
	"github.com/massarakhsh/tsan/jone"
	"github.com/massarakhsh/lik"
)

//	Проверка на выполенный логин
func (rule *DataRule) IsLogin() bool {
	return rule.ItSession.IdMember > 0
}

//	Проверка на режим диспетчера
func (rule *DataRule) IAmDispatch() bool {
	return rule.IAmRole(jone.ItDispatch)
}

//	Проверка на режим риэлтора
func (rule *DataRule) IAmRealtor() bool {
	return rule.IAmRole(jone.ItRealtor)
}

//	Проверка на режим менеджера
func (rule *DataRule) IAmManager() bool {
	return rule.IAmRole(jone.ItManager)
}

//	Проверка на режим менеджера по рекламе
func (rule *DataRule) IAmAdvert() bool {
	return rule.IAmRole(jone.ItAdvert)
}

//	Проверка на режим администратора
func (rule *DataRule) IAmAdmin() bool {
	return rule.IAmRole(jone.ItAdmin)
}

//	Проверка на режим
func (rule *DataRule) IAmRole(role string) bool {
	return rule.IsLogin() && rule.ItSession.RoleMember == role
}

//	Проверка на Шамана
func (rule *DataRule) IAmShaman() bool {
	ok := false
	if member := rule.GetMember(); member != nil {
		if member.GetString("login") == "shaman" {
			ok = true
		}
	}
	return ok
}

//	Проверка на допуск роли
func (rule *DataRule) ICanRole(role string) bool {
	ok := false
	if member := rule.GetMember(); member != nil {
		roles := jone.CalculateElmString(member, "role")
		if lik.RegExCompare(roles, role) {
			ok = true
		} else if lik.RegExCompare(roles, jone.ItAdmin) {
			ok = true
		} else if lik.RegExCompare(roles, jone.ItManager) {
			if role == jone.ItRealtor || role == jone.ItDispatch {
				ok = true
			}
		} else if lik.RegExCompare(roles, jone.ItRealtor) {
			if role == jone.ItDispatch {
				ok = true
			}
		}
	}
	return ok
}

//	Проверка на допуск к сегменту
func (rule *DataRule) ICanDo(segment string) bool {
	ok := false
	if segment == jone.DoCall {
		ok = true
	} else if segment == jone.DoTune && rule.IAmAdmin() {
		ok = true
	} else if member := rule.GetMember(); member != nil {
		ok = member.GetBool("do_" + segment)
	}
	return ok
}

//	ПОдключиться к сессии в режиме автологина
func (rule *DataRule) BindSession() bool {
	if !rule.IsLogin() && !rule.ItPage.Session.ExitLogin {
		if idop := lik.StrToInt(rule.GetCookie("auto_login_id")); idop > 0 {
			rule.SessionLogin(lik.IDB(idop), 0, true)
		}
	}
	return rule.IsLogin()
}

//	Войти в сессию с идентификатором idmem, выйти если 0
func (rule *DataRule) SessionLogin(idmem lik.IDB, pin int, store bool) {
	if idmem != 0 && rule.ItSession.IdMember == idmem {
	} else if member := jone.TableMember.GetElm(idmem); member != nil {
		if rule.ItSession.IdMember != 0 {
			rule.ClearSession()
		}
		rule.ItSession.IdMember = idmem
		roles := jone.CalculateElmString(member,"role")
		if roles == "" {
			roles = jone.ItDispatch
			jone.SetElmValue(member, roles, "role")
		}
		rolit := jone.CalculateElmString(member,"rolit")
		if rolit != "" && !rule.ICanRole(rolit) {
			rolit = ""
		}
		if rolit == "" {
			if match := lik.RegExParse(roles, "^([\\,]+)"); match != nil {
				rolit = match[1]
			} else {
				rolit = roles
			}
		}
		rule.ItSession.RoleMember = rolit
		jone.SetElmValue(member, rolit, "rolit")
		if pin != 0 {
			rule.SetPinMember(pin)
		}
		rule.SayWarning("login " + rule.GetMember().GetString("login"))
	} else {
		idmem = 0
		rule.ItSession.IdMember = 0
		rule.ItSession.RoleMember = ""
		rule.SetPagePart(0, "cabinet")
	}
	if store {
		rule.SetCookie(lik.IntToStr(int(idmem)),"auto_login_id")
	}
}

//	Установить внутренний номер телефона
func (rule *DataRule) SetPinMember(pin int) {
	for id,mem := range jone.TableMember.Elms {
		if id != rule.ItSession.IdMember {
			if pin == mem.GetInt("pin") {
				mem.SetValue(0, "pin")
			}
		}
	}
	rule.GetMember().SetValue(pin, "pin")
}

