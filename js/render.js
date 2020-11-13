function fancy_render_number_n(o) { return fancy_render_number_(o,'n'); }
function fancy_render_number_m(o) { return fancy_render_number_(o,'m'); }
function fancy_render_number_i(o) { return fancy_render_number_(o,'i'); }
function fancy_render_number_(o,fun) {
    var stl = { };
    var val = ''+o.value;
    val = val.replace(/ /g,"");
    val = val.replace(/\+/g,"");
    val = val.replace(/\,/g,".");
    var valpic = val.replace(/[^0-9\-\.]/g,"");
    if (valpic != val) {
        stl.color = '#d00000';
    } else if (val) {
        var flt = '';
        var match;
        if (match=/^(.*)[,\.](.*)/.exec(val)) {
            val = match[1];
            flt = match[2];
        }
        var res = flt;
        if (flt.length > 0) res = '.'+res;
        var len = val.length;
        for (var pos=0; pos<len; pos++) {
            if (pos>0 && pos%3==0) { res = ' '+res; }
            res = val.substr(len-1-pos,1) + res;
        }
        while (fun=='i' && res.length<3) res = '0'+res;
        res = res.replace(/ /g,"&nbsp;");
        o.value = res;
    }
    o.style = stl;
    return o;
}

function fancy_render_datetime(o) {
    if (match = /(\d+)\D+(\d+)\D+(\d+)\D+(\d+)\D+(\d+)/.exec(o.value)) {
        o.value = match[3]+"/"+match[2]+"/"+match[1]+" "+match[4]+":"+match[5];
    } else if (match = /(\d+)\D+(\d+)\D+(\d+)/.exec(o.value)) {
        o.value = match[3]+"/"+match[2]+"/"+match[1];
    }
    return o;
}

function fancy_render_ymd(o) {
    if (match = /(\d\d\d\d)(\D+)(\d\d)(\D+)(\d\d)(.*)/.exec(o.value)) {
        o.value = match[5]+match[2]+match[3]+match[4]+match[1]+match[6];
    }
    return o;
}

function fancy_render_phone(o) {
    let phone = fancy_input_telephone(o.value);
    if (phone.length >= 5 && (phone[0] != "+" || phone[1] != "7")) {
        phone = "+7" + phone;
    }
    o.value = phone;
    return o;
}

function fancy_input_number(value) {
    value = value.toString().replace(/[^0-9\.\,]/g, '').replace(',','.');
    if (value.indexOf('.')<0) {
        if (value.length === 0) {
            value = '';
        } else {
            var res = '';
            while (value.length >= 3) {
                if (res.length>0) res = ' '+res;
                res = value.substr(value.length-3,3) + res;
                value = value.substr(0, value.length-3);
            }
            if (value.length > 0) {
                if (res.length>0) res = ' '+res;
                res = value + res;
            }
            value = res;
        }
    }
    return value;
}

function fancy_input_telephone(value) {
    value = " "+value;
    value = value.replace(/\D/g, '');
    if (value.length > 10 && (value.substr(0,1)=='7' || value.substr(0,1)=='8')) {
        value = value.substr(1);
    }
    var res = '';
    for (p=0; p<value.length; p++) {
        if (p==0) res += "(";
        else if (p==3) res += ")";
        else if (p==6 || p==8) res += "-";
        else if (p==10) res += " ";
        res += value.substr(p,1);
    }
    return res;
}

function fancy_bell_control(value) {
    setTimeout(function() { fancy_bell_probe(' '+value); }, 100 )
    return fancy_input_telephone(value)
}

function fancy_phone_control(value) {
    setTimeout(function() { fancy_phone_probe(' '+value); }, 100 )
    return fancy_input_telephone(value)
}

function fancy_input_time(value) {
    value = value.replace(/\D/g, '');
    var res = '';
    for (var p = 0; p < 4 && p < value.length; p++) {
        if (p == 2) res += ":";
        res += value.substr(p, 1);
    }
    return res;
}

