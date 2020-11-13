var script_second = 0;
var locate_menu = null;
var UploadFiles = null;

function script_start() {
    pool_step.push(script_step);
    if (!lik_trust) lik_set_trust("");
    if (lik_trust) lik_set_marshal(1000, "/marshal");
}

function script_step() {
    script_showtime();
    script_redraw();
}

function script_showtime() {
    if (script_second!=tick_second) {
        script_second = tick_second;
        var elm = jQuery('#srvtime');
        if (elm.size()>0) {
            var text = "";
            var ok = true;
            if (tick_total - tick_answer < 10000) {
                var tt = tick_server - tick_shift_minute * 60;
                text = build_showtime(tt, true);
                ok = true;
            } else if (tick_total - tick_answer < 300000) {
                var tt = Math.floor((tick_total - tick_answer) / 1000);
                text = "нет связи " + build_showtime(tt, false);
                ok = false;
            } else {
                lik_stop();
                text = "<b>СИСТЕМА ОСТАНОВЛЕНА</b>";
                ok = false;
            }
            if (ok && elm.hasClass("srvoff")) {
                elm.removeClass("srvoff");
            } else if (!ok && !elm.hasClass("srvoff")) {
                elm.addClass("srvoff");
            }
            elm.html(text);
        }
        monitor_bell();
    }
}

function build_showtime(tt, ok) {
    var ts = tt % 60;
    tt = (tt - ts) / 60;
    var tm = tt % 60;
    tt = (tt - tm) / 60;
    var th = tt % 24;
    tt = (tt - th) / 24;
    var text = "";
    if (ok) {
        text += (th >= 10) ? "" + th : "0" + th;
        text += (ts & 1) ? ":" : " ";
    }
    text += (tm >= 10) ? tm : "0" + tm;
    if (!ok) {
        text += (ts >= 10) ? ":" + ts : ":0" + ts;
    }
    return text;
}

function script_redraw() {
    let rdr = jQuery('[redraw]');
    if (rdr.size() > 0) {
        rdr.each(function (idx, item) {
            let elm = jQuery(item);
            let redraw = elm.attr('redraw');
            elm.removeAttr('redraw');
            if (redraw in window) {
                setTimeout(function () {
                    window[redraw](elm);
                }, 100);
            }
        });
    }
}

var last_bellform = 0;

function monitor_bell() {
    var bell = '0';
    var elm = jQuery('#belledit');
    if (elm.size()>0) {
        bell = elm.text();
        if (bell != last_bellform) {
            front_get("/bell/bellform/"+bell);
        }
    }
    last_bellform = bell;
}

function front_get(cmd) {
    get_data_part("/front" + cmd);
}

function front_post(cmd, data) {
    post_data_part("/front" + cmd, data);
}

function front_proc(cmd, proc, parm) {
    get_data_proc("/front" + cmd, proc, parm);
}

function front_post_proc(cmd, data, proc, parm) {
    post_data_proc("/front" + cmd, data, proc, parm);
}

function path_trio(main,zone,cmd) {
    let path = "";
    if (main) path += "/" + main;
    if (zone) path += "/" + zone;
    if (cmd) path += "/" + cmd;
    if (!path) path = "/";
    return path;
}

function frio_get(main,zone,cmd) {
    front_get(path_trio(main, zone, cmd));
}

function frio_post(main,zone,cmd, data) {
    front_post(path_trio(main, zone, cmd), data);
}

function frio_proc(main,zone,cmd, proc, parm) {
    front_proc(path_trio(main, zone, cmd), proc, parm);
}

function click_exit(main) {
    window.close();
}

function click_print(main) {
    fancy_bind_ctrl("print","print","_show");
}

function click_project(main) {
    fancy_bind_rows("project","project","_show");
}

function click_role() {
    fancy_bind_ctrl("command","role","_show");
}

function click_segment() {
    fancy_bind_ctrl("command","segment","_show");
}

function click_menu(main, cmd) {
    frio_get("menu", main, cmd);
}

function openbox(id){
    $("#"+id).fadeIn(); //плавное появление блока
}

function closebox(id){
    $("#"+id).fadeOut(); //плавное исчезание блока
}

function grid_command_segment(field, value, oldValue) {
    frio_get("all","segment", value);
}

function grid_command_realty(field, value, oldValue) {
    frio_get("all","realty", value);
}

function grid_command_filter(field, value, oldValue) {
    frio_get("all","filter", value);
}

function grid_command_locate(field, value, oldValue) {
    frio_get("all","locate", value);
}

function grid_command_status(field, value, oldValue) {
    frio_get("all","status", value);
}

function media_store(zone) {
    front_get('/files/' + zone + '/store');
}

function media_cancel(zone) {
    front_get('/files/' + zone + '/cancel');
}

function tunetree_switch(main, path, dir) {
    front_get("/" + main + "/tree/switch/"+path+"/"+dir);
}

function tunetree_select(main, path) {
    front_get("/" + main + "/tree/select/"+path);
}

function export_test(obj) {
    var grid = fancy_seekform(obj);
    front_get("/" + grid.likMain + "/" + grid.likZone + "/test");
}

function dump(params) {
    var text = "";
    for (var key in params) {
        text += key + ": " + params[key] + "\n";
    }
    alert(text);
}

function bind_offer(idoff) {
    fancy_window_destroy();
    lik_window_part("/show"+idoff+"?_tp=1");
    window.location.reload();
}

function bind_cabinet(idmember) {
    lik_go_part("/member" + idmember+ "/membercard"+idmember);
}

function offer_goto(idoff) {
    fancy_window_destroy();
    front_get("/all/gooffer/"+idoff);
}

function go_window_offer(mode, idoff) {
    fancy_window_destroy();
    //lik_window_part("/"+mode+idoffer+"/offershow"+idoffer+"?_tp=1");
    lik_window_part("/offershow"+idoffer+"?_tp=1");
}

function go_window_cabinet(idmember) {
    lik_window_part("/member" + idmember+ "/membercard"+idmember+"?_tp=1");
}

function member_enter_cabinet() {
    fancy_window_destroy();
    front_get("/all/gomember");
}

function member_create_cabinet() {
    fancy_window_destroy();
    front_get("/all/newmember");
}

function bell_goto(idbell) {
    fancy_window_destroy();
    front_get("/all/contact/gobell/"+idbell);
}

function bell_bind(idbell) {
    fancy_window_destroy();
    lik_window_part("/bell"+idbell+"/_show?_tp=1");
}

function goto_file(path) {
    lik_window_part(path+"?_tp=1");
}

function bell_accept(wrt) {
    var data = (wrt) ? form_collect_data() : "";
    front_post("/bell/all/bellaccept", data);
    setTimeout(function() { form_reshow('bell','all'); }, 500);
}

function bell_cancel(wrt) {
    let sel = jQuery("#why_cancel");
    let why = sel.val();
    fancy_window_destroy();
    front_get("/bell/all/bellcancel/"+why);
}

function bell_newoffer(wrt) {
    var data = (wrt) ? form_collect_data() : "";
    front_post("/bell/all/offercreate", data);
    fancy_window_destroy();
}

function form_collect_data() {
    return (FancyWindowForm) ? fancy_collect_edit(FancyWindowForm) : null;
}

function print_as(what) {
    fancy_window_destroy();
    front_get("/print/" + what);
}

function choose_command(what) {
    fancy_window_destroy();
    front_get("/command/" + what);
}

//////////////// Filters

var symbolValue = "";

function symbol_value(grid, value) {
    symbolValue = value;
}

function filter_clear(obj) {
    var form = fancy_seekform(obj);
    frio_get(form.likMain,"filterclear", "");
}

function filter_save(obj) {
    var form = fancy_seekform(obj);
    frio_get(form.likMain,"filtersave","");
}

function filter_saveas(obj) {
    var form = fancy_seekform(obj);
    frio_get(form.likMain,"filtersaveas", string_to_XS(symbolValue));
}

function filter_delete(obj) {
    var form = fancy_seekform(obj);
    frio_get(form.likMain,"filterdelete","");
}

function bell_set_change(obj) {
    let form = fancy_seekform(obj);
    form.get("phone1");
}

function form_reshow(main, zone) {
    fancy_window_destroy();
    fancy_bind_ctrl(main, zone,"_show");
}

function change_pin(pin) {
    var inpin = prompt("Укажите внутренний номер телефона. 0 - нет", pin);
    if (inpin !== null && inpin != pin) {
        front_get("/bell/setpin/" + string_to_XS(inpin));
    }
}

function set_status_offer(item, stat) {
    let sta = stat;
    let txt = stat;
    if (match = /(.*?),(.*)/.exec(stat)) {
        sta = match[1];
        txt = match[2];
    }
    let why = prompt("Установить \""+txt+ "\", можете уточнить: ", "");
    if (why !== null) {
        let button = Fancy.getWidget('topofferstatus');
        button.menu.hide();
        front_get("/offerstaff/status/setstatus/" + sta + "/" + string_to_XS(why));
    }
}

function choose_offer_pair(id) {
    front_get("/offerdeal/manage/choose/" + id);
}

function ava_control(cmd) {
    if (cmd == "delete") {
        if (!confirm("Действительно удалить эту фотографию?")) return;
    }
    front_get("/membercard/manage/" + cmd);
}

function system_clear_base() {
    if (confirm("Удалить все объекты, вы уверены?")) {
        front_get("/system/clearbase");
    }
}

function deal_manage_delete() {
    if (confirm("Удалить эту сделку, вы уверены?")) {
        front_get("/offerdeal/manage/delete");
    }
}

function change_id_new() {
    if (confirm("Изменить ID, вы уверены?")) {
        front_get("/offerstaff/editor/redoid");
    }
}

function go_godoc() {
    window.open("http://rltweb.ru:8090", '_blank');
}

function ask_message_delete(obj) {
    if (confirm("Удалить это сообщение, вы уверены?")) {
        var grid = fancy_seekform(obj);
        front_get("/" + grid.likMain + "/" + grid.likZone + "/delete");
    }
}

function new_message_append(obj) {
    var grid = fancy_seekform(obj);
    var text = prompt("Новое сообщение", "");
    var data = "text=" + string_to_XS(text);
    front_post("/" + grid.likMain + "/" + grid.likZone + "/append", data);
}

