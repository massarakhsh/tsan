var FancyWindowForm = null;
var FancyCurrentGrid = null;
var FancyIntoSelect = false;
var FancyInitialized = false;

function fancy_item_find(main) {
    if (!main) return FancyWindowForm;
    jQuery('[likFancy]').each(function(idx,item) {
        var elm = $(item);
        if (elm.attr('likMain') == main) {
            return elm.attr(likFancy);
        }
    });
    return null;
}

function fancy_item_set(main, fancy) {
    if (!main) FancyWindowForm = fancy;
    else {
    }
}

function fancy_redraw_grid(elm) {
    let id = elm.attr('id');
    let remain = elm.attr('remain');
    let rezone = elm.attr('rezone');
    front_proc("/"+remain+"/"+rezone + "/showgrid", fancy_answer_draw, elm);
}

function fancy_redraw_form(elm) {
    let id = elm.attr('id');
    let remain = elm.attr('remain');
    let rezone = elm.attr('rezone');
    front_proc("/"+remain+"/"+rezone + "/showform", fancy_answer_draw, elm);
}

function fancy_answer_draw(elm, lika) {
    if (!lika) {
        //alert("FAD: "+elm.attr('remain')+"/"+elm.attr('rezone')+"/"+elm.attr('id'));
    } else if ('fancy' in lika) {
        let fancy = lika.fancy;
        if (!('id' in fancy)) fancy.id = elm.attr('id');
        fancy_build_bind(fancy);
    }
}

function fancy_trio_form(parm) {
    if (match = /(.*?)_(.*?)(_.*)/.exec(parm)) {
        var main = match[1];
        var zone = match[2];
        var id = match[3];
        fancy_bind_ctrl(main, zone, id);
    }
}

function fancy_bind_ctrl(ctrl,zone,id) {
    var path = "/"+ctrl+"/"+zone + "/showform/"+id;
    front_proc(path, fancy_answer_bind, id);
}

function fancy_bind_rows(ctrl,zone,id) {
    var path = "/"+ctrl+"/"+zone + "/showgrid/"+id;
    front_proc(path, fancy_answer_bind, id);
}

function fancy_answer_bind(path, lika) {
    if (!lika) {
        //alert("FAB: "+path);
    } else if ('fancy' in lika) {
        let id = path;
        let fancy = lika.fancy;
        if (!('id' in fancy)) fancy.id = id;
        fancy_build_bind(fancy);
    }
}

function fancy_build_bind(instant) {
    let iswindow = (!instant.id || instant.id.length==0 || instant.id[0]=='_');
    fancy_prepare(instant);
    if ('license' in instant) FancyGrid.LICENSE = [ instant.license ];
    //if (!('theme' in instant)) instant.theme = "material";
    instant.nativeScroller = true;
    instant.trackOver = true;
    instant.columnTrackOver = false;
    instant.rowTrackOver = true;
    instant.stateful = false;
    instant.i18n = "ru"
    if (!('defaults' in instant)) instant.defaults = {};
    if (!('events' in instant)) instant.events = [];
    if (!('type' in instant.defaults)) instant.defaults.type = "string";
    if (instant.isform) {
        instant.events.push({ init: function () { fancy_form_init(this);}});
    } else {
        instant.events.push({ init: function () { fancy_grid_init(this);}});
        instant.events.push({ load: function () { fancy_grid_load(this);}});
        if (!('selModel' in instant)) instant.selModel = "row";
        if (!('align' in instant.defaults)) instant.defaults.align = "center";
        if (!('width' in instant.defaults)) instant.defaults.width = 96;
        if (!('resizable' in instant.defaults)) instant.defaults.resizable = true;
        if (!('sortable' in instant.defaults)) instant.defaults.sortable = true;
    }
    if (iswindow) {
        fancy_window_destroy();
        instant.draggable = true;
        instant.window = true;
    } else {
        instant.renderTo = instant.id;
        let elm = jQuery("#" + instant.id);
        if (elm.size()>0) {
            let itm = elm.attr('likFancy');
            if (itm) {
                itm.destroy();
                elm.removeAttr('likFancy')
            }
            elm.html("");
        }
        //jQuery("#" + instant.id).replaceWith("<div id=" + instant.id + "></div>");
    }
    FancyIntoSelect = false;
    FancyInitialized = false;
    var item;
    if (instant.isform) {
        item = new FancyForm(instant);
    } else {
        item = new FancyGrid(instant);
    }
    if (!iswindow) {
        jQuery("#" + instant.id).attr('likFancy', item);
    }
    if (iswindow) {
        FancyWindowForm = item;
    } else if (instant.isform) {
    } else if ('likFull' in instant) {
        FancyCurrentGrid = item;
    } else {
        FancyCurrentGrid = null;
    }
    item.likFase = 0;
    if ('id' in instant) item.likId = instant.id;
    if ('likLeft' in instant) item.likLeft = instant.likLeft;
    if ('likTop' in instant) item.likTop = instant.likTop;
    if ('likFull' in instant) item.likFull = instant.likFull;
    if ('likMain' in instant) item.likMain = instant.likMain;
    if ('likPart' in instant) item.likPart = instant.likPart;
    if ('likZone' in instant) item.likZone = instant.likZone;
    if ('likLock' in instant) item.likLock = instant.likLock;
    if ('likSelect' in instant) item.likSelect = instant.likSelect;
}

function fancy_window_destroy() {
    if (FancyWindowForm) {
        FancyWindowForm.destroy();
        FancyWindowForm = null;
    }
}

function fancy_prepare(data) {
    if (data !== null && typeof(data) == 'object') {
        for (var key in data) {
            let value = data[key];
            if (key == "*") {
                delete data[key];
                key = "";
                data[key] = value;
            }
            if (typeof(value) == 'string') {
                var match;
                if (match = /^function_(.+)\((.*)\)/.exec(value)) {
                    let func = match[1];
                    let parm = match[2];
                    if (func in window) {
                        data[key] = function () {
                            window[func](this, parm);
                        };
                    } else {
                        data[key] = fancy_nothing;
                    }
                } else if (match = /^function_(.+)/.exec(value)) {
                    let func = match[1];
                    if (func in window) {
                        data[key] = window[func];
                    } else {
                        data[key] = fancy_nothing;
                    }
                }
            } else if (value !== null && typeof(value) == 'object') {
                fancy_prepare(data[key]);
            }
        }
    }
}

function fancy_seekform(obj) {
    var form = null;
    if ('likZone' in obj) {
        form = obj;
    } else if (obj.scope) {
        if (!form && 'likZone' in obj.scope) {
            form = obj.scope;
        }
        if (!form && obj.scope._tbar && obj.scope._tbar.scope) {
            if ('likZone' in obj.scope._tbar.scope) {
                form = obj.scope._tbar.scope;
            }
        }
        if (!form && obj.scope.items && obj.scope.items.length > 0) {
            if ('likZone' in obj.scope.items[0]) {
                form = obj.scope.items[0];
            }
        }
    } else if (obj.events) {
        if (obj.events[0].scope) {
            let sco = obj.events[0].scope;
            if ('likZone' in sco) {
                form = sco;
            } else if (sco.scope) {
                if ('likZone' in sco.scope) {
                    form = sco.scope;
                }
            }
        }
    }
    if (form == null) {
        form = FancyCurrentGrid;
    }
    return form
}

function fancy_nothing() {
}

function fancy_form_init(form) {
    setTimeout( function() {
        if ('likLeft' in form && 'likTop' in form) {
            form.showAt(form.likLeft, form.likTop);
        } else {
            form.show();
        }
    }, 100);
}

function fancy_grid_init(grid) {
    grid.show();
    setTimeout(function(){
        grid.load();
        FancyInitialized = true;
        //fancy_grid_seek(grid);
    }, 100);
}

function fancy_before_request(o) {
    //if (FancyInitialized) o.params.up_reality = 1;
    //o.params.up_filter = fancy_build_filter(FancyCurrentGrid);
    //o.params.up_sort = fancy_build_sorter(FancyCurrentGrid);
    return o;
}

function fancy_after_request(o) {
    if (FancyCurrentGrid && ('items' in o.response)) {
        fancy_prepare(o.response.items);
        if ('likSelect' in o.response) {
            FancyCurrentGrid.likSelect = o.response.likSelect
        }
        //fancy_grid_seek(FancyCurrentGrid);
    }
    return o;
}

function fancy_grid_update(cmd) {
    if (FancyCurrentGrid) {
        FancyCurrentGrid.waitingForFilters = true;
        FancyCurrentGrid.addFilter('id', 0, '>');
        FancyCurrentGrid.clearFilter('id');
        FancyCurrentGrid.updateFilters();
    }
}

function fancy_grid_seek(grid) {
    setTimeout( function() {
        //var index = ('likSelect' in grid) ? grid.likSelect : -1;
        //fancy_grid_rowselect(grid, index,null);
        //grid.selectRow(4);
    }, 100);
}

function fancy_grid_load(grid) {
    setTimeout(function(){
        var index = ('likSelect' in grid) ? grid.likSelect : -1;
        fancy_grid_rowselect(grid, index,null);
    }, 100);
}

function fancy_grid_rowselect(grid, rowIndex, dataItem) {
    if (grid.tbar) {
        for (var nb = 0; nb < grid.tbar.length; nb++) {
            var cls = grid.tbar[nb].imageCls;
            if (!cls) {
            } else if (/(bell|add)$/.test(cls)) {
                grid.tbar[nb].enable();
            } else if (/(show|mod|del|member)$/.test(cls)) {
                if (rowIndex >= 0) grid.tbar[nb].enable();
                else grid.tbar[nb].disable();
            }
        }
    }
    if (FancyIntoSelect) {
    } else if (dataItem) {
        var zone = ('likZone' in grid) ? grid.likZone : '';
        var id = ('id' in dataItem) ? dataItem.id: '';
        var parm = "/"+grid.likMain+"/"+grid.likZone + "/rowselect/" + id;
        grid.likSelect = id;
        front_get(parm);
    } else if (rowIndex >= 0) {
        FancyIntoSelect = true;
        grid.selectRow(rowIndex);
    }
    FancyIntoSelect = false;
    fancy_clip_scan(grid);
}

function fancy_grid_rowdeselect(grid, rowIndex, dataItem) {
    fancy_clip_scan(grid);
}

function fancy_grid_mark(grid, obj) {
    var cmd = obj.column.index;
    if (cmd == "ready" || cmd == "use" ||
        cmd == "mark" || cmd == "enable" ||
        (match = /^tag_(.+)/.exec(cmd))) {
        let main = ('likMain' in grid) ? grid.likMain : '';
        let zone = ('likZone' in grid) ? grid.likZone : '';
        var id = obj.id;
        var val = (obj.value) ? 0 : 1;
        front_get("/" + main + "/" + zone + "/" + cmd + "/" + id + "/" + val);
    }
}

var listContextGrid = null;
var listContextItem = null;

function fancy_grid_rclick(grid,o) {
    listContextGrid = grid;
    listContextItem = o.item;
    var rowIndex = grid.getRowById(listContextItem.data.id);
    if (rowIndex === 0 || rowIndex > 0) {
        grid.selectRow(rowIndex);
    }
}

function fancy_grid_context(obj,item) {
    var grid = fancy_seekform(obj);
    if (!grid) grid = listContextGrid;
    if ('part' in item && item.part && listContextItem) {
        var id = listContextItem.data.id;
        var part = item.part;
        if (part == "mark") {
            listContextItem.data.mark = 1 - listContextItem.data.mark;
            front_get("/"+ grid.likMain+"/"+grid.likZone+ "/" + part + "/" + id + "/" + listContextItem.data.mark);
            grid.update();
        } else {
            front_get("/"+ grid.likMain+"/"+grid.likZone+ "/" + part + "/" + id);
        }
    }
}

function fancy_goto_form(parm, lika) {
    fancy_trio_form(parm);
}

function fancy_grid_cmd(obj) {
    var grid = fancy_seekform(obj);
    front_get("/" + grid.likMain + "/" + grid.likZone + "/" + obj.cmd);
}

function fancy_grid_show(obj) {
    var form = fancy_seekform(obj);
    var zone = ('likZone' in form) ? form.likZone : '';
    fancy_bind_ctrl(form.likMain, zone, "_show");
}

function fancy_grid_bell(obj) {
    var form = fancy_seekform(obj);
    var zone = ('likZone' in form) ? form.likZone : '';
    fancy_bind_ctrl(form.likMain, zone, "_bell");
}

function fancy_grid_mail(obj) {
    var form = fancy_seekform(obj);
    var zone = ('likZone' in form) ? form.likZone : '';
    fancy_bind_ctrl(form.likMain, zone, "_mail");
}

function fancy_grid_add(obj) {
    var form = fancy_seekform(obj);
    var zone = ('likZone' in form) ? form.likZone : '';
    fancy_bind_ctrl(form.likMain, zone, "_add");
}

function fancy_grid_mod(obj) {
    var form = fancy_seekform(obj);
    var main = ('likMain' in form) ? form.likMain : '';
    var zone = ('likZone' in form) ? form.likZone : '';
    front_proc("/"+main+"/"+zone + "/settab/" + form.activeTab, fancy_goto_form, main+"_"+zone+"_mod");
}

function fancy_grid_del(obj) {
    var form = fancy_seekform(obj);
    var zone = ('likZone' in form) ? form.likZone : '';
    fancy_bind_ctrl(form.likMain, zone, "_del");
}

function fancy_grid_append(obj) {
    var form = fancy_seekform(obj);
    front_get("/"+form.likMain+"/"+form.likZone+"/append");
}

function fancy_grid_dblclick(grid, obj) {
    if (FancyWindowForm) {
        let part = ('likPart' in grid) ? grid.likPart : "";
        let id = ('id' in obj) ? obj.id : "";
        let main = ('likMain' in FancyWindowForm) ? FancyWindowForm.likMain : '';
        let zone = ('likZone' in FancyWindowForm) ? FancyWindowForm.likZone : '';
        let mod = ('likId' in FancyWindowForm) ? FancyWindowForm.likId : '';
        //alert(id+","+main+","+zone+","+part+","+mod);
        if (part == "offer" && (mod == "_add" || mod == "_mod" || mod == "_edit")) {
            fancy_offer_probe(main, id);
        }
    //} else if ('pathopen' in obj.item) {
    //    var url = lik_build_url(obj.item.pathopen);
    //    window.open(url, '_blank');
    //} else if ('pathopen' in obj.data) {
    //    var url = lik_build_url(obj.data.pathopen);
    //    window.open(url, '_blank');
    } else {
        front_get("/"+grid.likMain+"/"+grid.likZone + "/toenter");
    }
}

function fancy_grid_drag(grid, params) {
    var it = params[0];
    var id = it.id;
    setTimeout(function () {
        fancy_grid_drop(grid, id);
    }, 500);
}

function fancy_grid_drop(grid, id) {
    var data = grid.getData();
    let main = ('likMain' in grid) ? grid.likMain : '';
    let zone = ('likZone' in grid) ? grid.likZone : '';
    for (var nr = 0; nr < data.length; nr++) {
        var row = data[nr];
        if (id == row.id) {
            front_get("/" + main + "/" + zone + "/order/" + id + "/" + nr);
            break;
        }
    }
}

function fancy_col_size(grid, params) {
    var form = fancy_seekform(grid);
    var zone = ('likZone' in form) ? form.likZone : '';
    var index = params.column.index;
    front_get("/"+form.likMain+"/"+zone+"/colsize/" + index + "/" + params.width);
}

function fancy_col_drag(grid, params) {
    var grid = fancy_seekform(grid);
    var index = params.column.index;
    for (nc=0; nc<grid.columns.length; nc++) {
        if (grid.columns[nc].index == index) {
            front_get("/"+grid.likMain+"/"+grid.likZone+"/coldrag/" + index + "/" + nc);
            break;
        }
    }
}

function fancy_column_show(item, parm) {
    if (FancyCurrentGrid) {
        let sw = ""
        var imgEl = item.el.select('.fancy-menu-item-image');
        if (imgEl.hasClass("imgplus")) {
            sw = "/0";
            imgEl.removeCls('imgplus');
            imgEl.addCls('imgno');
        } else if (imgEl.hasClass("imgno")) {
            sw = "/1";
            imgEl.removeCls('imgno');
            imgEl.addCls('imgplus');
        }
        if (sw) {
            //let button = Fancy.getWidget('topcolumns');
            //button.menu.hide();
            front_get("/" + FancyCurrentGrid.likMain + "/" + FancyCurrentGrid.likZone + "/colshow/" + parm + sw);
        }
    }
}

function fcs() {
    alert("fcs");
}

/////////////// Form

function fancy_form_toshow(obj) {
    fancy_form_tomode(obj, "show");
}
function fancy_form_tocreate(obj) {
    fancy_form_tomode(obj, "add");
}
function fancy_form_toedit(obj) {
    fancy_form_tomode(obj, "mod");
}
function fancy_form_todelete(obj) {
    fancy_form_tomode(obj, "del");
}
function fancy_form_tomode(obj, mode) {
    let form = fancy_seekform(obj);
    var main = ('likMain' in form) ? form.likMain : '';
    let zone = ('likZone' in form) ? form.likZone : '';
    let tab = form.activeTab;
    front_proc("/"+main+"/"+zone + "/tab/" + tab, fancy_trio_form, main+"_"+zone+"_"+mode);
}

function fancy_form_cancel(obj) {
    var form = fancy_seekform(obj);
    front_get("/"+form.likMain+"/"+form.likZone + "/cancel");
    fancy_window_destroy();
}

function fancy_form_write(obj) {
    var form = fancy_seekform(obj);
    var data = fancy_collect_edit(form);
    front_post("/"+form.likMain+"/"+form.likZone + "/write/" + form.activeTab, data);
    fancy_window_destroy();
}

///////////// Editor

function fancy_edit_start(obj) {
    var form = fancy_seekform(obj);
    var zone = ('likZone' in form) ? form.likZone : '';
    front_get("/"+form.likMain+"/"+zone + "/edit/" + form.activeTab);
}

function fancy_edit_write(obj) {
    var form = fancy_seekform(obj);
    var zone = ('likZone' in form) ? form.likZone : '';
    var data = fancy_collect_edit(form);
    front_post("/"+form.likMain+"/"+zone + "/write/" + form.activeTab, data);
}

function fancy_edit_cancel(obj) {
    var form = fancy_seekform(obj);
    var zone = ('likZone' in form) ? form.likZone : '';
    front_get("/"+form.likMain+"/"+zone + "/cancel/" + form.activeTab);
}

////////////////////

function fancy_real_delete(obj) {
    var form = fancy_seekform(obj);
    var zone = ('likZone' in form) ? form.likZone : '';
    front_get("/"+form.likMain+"/"+zone + "/delete");
    fancy_window_destroy();
}

function fancy_clip_copy(obj) {
    fancy_clip_call(obj,"copy");
}
function fancy_clip_cut(obj) {
    fancy_clip_call(obj,"cut");
}
function fancy_clip_paste(obj) {
    fancy_clip_call(obj,"paste");
}
function fancy_clip_clear(obj) {
    fancy_clip_call(obj,"clear");
}

function fancy_clip_call(obj, fun) {
    let form = fancy_seekform(obj);
    front_get("/"+form.likMain+"/"+form.likZone + "/clip" + fun);
}

function fancy_clip_scan(form) {
}

function fancy_collect_edit(form) {
    let data = "likid=" + form.likId;
    let values = form.get();
    if (values) {
        for (var key in values) {
            data += "&up_" + key + "=" + string_to_XS(values[key]);
        }
    }
    return data;
}

/////////////////////////////

function fancy_bell_probe(value) {
    value = value.replace(/\D/g, '');
    let main = (FancyWindowForm) ? FancyWindowForm.likMain : null;
    front_proc("/"+main+"/bell/phonesearch/"+value, fancy_bell_answer, null);
}

function fancy_bell_answer(over, lika) {
    var code = '';
    if ('ok' in lika && lika.ok) {
        code = lika.namely+" "+lika.paterly+" "+lika.family;
    }
    if (FancyWindowForm) {
        var item = FancyWindowForm.getItem('search');
        if (item) {
            FancyWindowForm.set('search', code);
        }
    }
}

function fancy_phone_probe(value) {
    value = value.replace(/\D/g, '');
    let main = (FancyWindowForm) ? FancyWindowForm.likMain : null;
    front_proc("/"+main+"/bell/phonesearch/"+value, fancy_phone_answer, null);
}

function fancy_phone_answer(over, lika) {
    fancy_phone_fill(lika,false);
}

function fancy_phone_fill(lika, force) {
    let ok = (lika && 'ok' in lika) ? lika.ok : 0;
    let id = (lika && 'id' in lika) ? lika.id : 0;
    let namely = (lika && 'namely' in lika) ? lika.namely : '';
    let paterly = (lika && 'paterly' in lika) ? lika.paterly : '';
    let family = (lika && 'family' in lika) ? lika.family : '';
    let div = jQuery('#phone_search');
    let code = "";
    if (ok) {
        if (id) code = "Клиент №" + id + ": ";
        else code = "Контакт: "
        code += namely + " " + paterly + " " + family;
    } else {
        code = "Телефон не найден"
    }
    code += ". [<a href='#' onclick='fancy_phone_fill(null,true)'>Стереть</a>]";
    div.html(code);
    if (ok || force) {
        FancyWindowForm.set('u_clientid', id);
        FancyWindowForm.set('s_clientid__namely', namely);
        FancyWindowForm.set('s_clientid__paterly', paterly);
        FancyWindowForm.set('s_clientid__family', family);
    }
}

function fancy_offer_probe(main, id) {
    front_proc("/"+main+"/offersearch/"+id, fancy_offer_answer, null);
}

function fancy_offer_answer(over, lika) {
    let id = (lika && 'id' in lika) ? lika.id : 0;
    let target = (lika && 'target' in lika) ? lika.target : '';
    let segment = (lika && 'segment' in lika) ? lika.segment : '';
    let realty = (lika && 'realty' in lika) ? lika.realty : '';
    let subcity = (lika && 'subcity' in lika) ? lika.subcity : '';
    let rooms = (lika && 'rooms' in lika) ? lika.rooms : '';
    let memberid = (lika && 'memberid' in lika) ? lika.memberid : '';
    let div = jQuery('#offer_search');
    if (id > 0) {
        let code = "Объект по заявке №" + id;
        code += ". [<a href='#' onclick='fancy_offer_answer(null,null)'>Стереть</a>]";
        div.html(code);
    } else {
        div.html("...");
    }
    if (FancyWindowForm) {
        FancyWindowForm.set('u_targetid', id);
        FancyWindowForm.set('c_target', target);
        FancyWindowForm.set('c_segment', segment);
        FancyWindowForm.set('c_realty', realty);
        FancyWindowForm.set('c_address__subcity', subcity);
        FancyWindowForm.set('c_define__rooms', rooms);
        FancyWindowForm.set('u_memberid', memberid);
    }
}

function fancy_build_filter(grid) {
    if (!grid) return null;
    if (!('filter' in grid)) return null;
    let filter = grid.filter;
    if (!('filters' in filter)) return null;
    let filters = filter.filters;
    let search = "";
    if (grid.tbar) {
        for (var nt = 0; nt < grid.tbar.length; nt++) {
            var top = grid.tbar[nt];
            if (top.type == "search" && ('acceptedValue' in top)) {
                search = grid.tbar[nt].acceptedValue;
                break;
            }
        }
    }
    let condition = [];
    for (var key in filters) {
        let flt = filters[key];
        for (var op in flt) {
            let flop = flt[op];
            let val = "";
            if (op == "|") {
                for (var ch in flop) {
                    if (val != "") val += ",";
                    val += ch;
                }
            } else if (flop &&
                (typeof(flop)==="string" || typeof(flop)==="number") &&
                (op != '*' || flop != search)) {
                val = flop;
            }
            if (val != "") {
                condition.push(key+"/"+op+"/"+val);
            }
        }
    }
    let result = string_to_XS(search);
    for (n1=condition.length-2; n1>=0; n1--) {
        for (n2=0; n2<n1; n2++) {
            let c1 = condition[n2];
            let c2 = condition[n2+1];
            if (c1>c2) {
                condition[n2] = c2;
                condition[n2+1] = c1;
            }
        }
    }
    for (n=0; n<condition.length; n++) {
        result += "/" + string_to_XS(condition[n]);
    }
    return result;
}

function fancy_build_sorter(grid) {
    if (!grid) return null;
    if (!('store' in grid)) return null;
    let store = grid.store;
    if (!('sorters' in store)) return null;
    let sorters = store.sorters;
    if (!sorters || sorters.length==0) return null;
    let result = sorters[0].key;
    if (!result) return null;
    if (sorters[0].dir.toLowerCase() == "desc") result += "_r";
    return result;
}
