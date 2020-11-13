var YandexMain = null;
var YandexZone = null;
var YandexId = null;
var YandexIsEdit = false;
var YandexMap = null;
var YandexCenter = null;
var YandexZoom = 0;
var YandexPoint = null;

function map_bind_zone(main, zone, id) {
    YandexMain = main;
    YandexZone = zone;
    YandexId = id;
    frio_proc(main, zone,"showmap", map_answer_bind, id);
}

function map_redraw() {
    frio_proc(YandexMain, YandexZone,"showmap", map_answer_bind, YandexId);
}

function map_answer_bind(id, lika) {
    YandexIsEdit = ('isedit' in lika) ? lika.isedit : false;
    if (lika.isfull) {
        YandexCenter = ('center' in lika) ? lika.center : null;
        YandexZoom = ('zoom' in lika) ? lika.zoom : null;
    }
    let points = ('points' in lika) ? lika.points : null;
    YandexPoint = null;
    jQuery('#'+id).html('');
    YandexMap = new ymaps.Map(id, { center: YandexCenter, zoom: YandexZoom });
    YandexMap.controls.remove('geolocationControl');
    YandexMap.controls.remove('trafficControl');
    if (points) {
        setTimeout(function(){ map_set_points(points); }, 250);
    }
    YandexMap.events.add('boundschange', map_bounds_change);
    if (YandexIsEdit) {
        YandexMap.events.add('contextmenu', map_mouse_right);
    }
}

function map_setedit() {
    if (!YandexIsEdit) {
        frio_proc(YandexMain,YandexZone, "toedit", map_answer_bind, YandexId);
    }
}

function map_setclear() {
    if (YandexIsEdit && YandexPoint) {
        YandexMap.geoObjects.remove(YandexPoint);
        YandexPoint = null;
    }
}

function map_setcancel() {
    if (YandexIsEdit) {
        frio_proc(YandexMain,YandexZone,"tocancel", map_answer_bind, YandexId);
    }
}

function map_setwrite() {
    if (YandexIsEdit) {
        let data = "up_zoom=" + string_to_XS(YandexZoom);
        data += "&up_centerx=" + string_to_XS(YandexCenter[0]) + "&up_centery=" + string_to_XS(YandexCenter[1]);
        if (YandexPoint) {
            let points = YandexPoint.geometry._coordinates;
            data += "&up_points=" + string_to_XS(points.join(','));
        }
        frio_post(YandexMain,YandexZone,"write", data);
    }
}

function map_snap() {
    html2canvas(document.getElementById(YandexId)).then(function(canvas) {
        var img = canvas.toDataURL()
        window.open(img);
        //var dt64 = canvas.toDataURL('image/jpeg').replace('image/jpeg', 'image/octet-stream');
        //var data = "up_snap=" + string_to_XS(dt64);
        //frio_post(YandexMain,YandexZone, "snap", data);
    });
}

function map_set_points(points) {
    if (YandexPoint) {
        YandexPoint.geometry.setCoordinates(points);
    } else {
        YandexPoint = new ymaps.GeoObject({
            geometry: {
                type: "Point",
                coordinates: points
            },
            properties: {
                iconContent: (YandexIsEdit) ? "Перемещайте" : "Здесь"
            }
        }, {
            preset: 'islands#blackStretchyIcon',
            draggable: YandexIsEdit
        });
        YandexMap.geoObjects.add(YandexPoint);
    }
}

function map_mouse_right(e) {
    let points = e.get('coords');
    map_set_points(points);
}

function map_bounds_change(e) {
    YandexCenter = e.get('newCenter');
    YandexZoom = e.get('newZoom');
}

