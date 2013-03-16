/*
 * jQuery timepicker addon
 * By: Trent Richardson [http://trentrichardson.com]
 * Version 1.1.1
 * Last Modified: 11/07/2012
 *
 * Copyright 2012 Trent Richardson
 * You may use this project under MIT or GPL licenses.
 * http://trentrichardson.com/Impromptu/GPL-LICENSE.txt
 * http://trentrichardson.com/Impromptu/MIT-LICENSE.txt
 */

/*jslint evil: true, white: false, undef: false, nomen: false */

(function(e){function r(){this.regional=[];this.regional[""]={currentText:"Teraz",closeText:"Hotovo",amNames:["AM","A"],pmNames:["PM","P"],timeFormat:"HH:mm",timeSuffix:"",timeOnlyTitle:"Vyberte \u010das",timeText:"\u010cas",hourText:"Hodina",minuteText:"Min\u00fata",secondText:"Sekunda",millisecText:"Milisekunda",timezoneText:"\u010casov\u00e1 z\u00f3na",isRTL:!1};this._defaults={showButtonPanel:!0,timeOnly:!1,showHour:!0,showMinute:!0,showSecond:!1,showMillisec:!1,showTimezone:!1,showTime:!0,stepHour:1,
stepMinute:1,stepSecond:1,stepMillisec:1,hour:0,minute:0,second:0,millisec:0,timezone:null,useLocalTimezone:!1,defaultTimezone:"+0000",hourMin:0,minuteMin:0,secondMin:0,millisecMin:0,hourMax:23,minuteMax:59,secondMax:59,millisecMax:999,minDateTime:null,maxDateTime:null,onSelect:null,hourGrid:0,minuteGrid:0,secondGrid:0,millisecGrid:0,alwaysSetTime:!0,separator:" ",altFieldTimeOnly:!0,altTimeFormat:null,altSeparator:null,altTimeSuffix:null,pickerTimeFormat:null,pickerTimeSuffix:null,showTimepicker:!0,
timezoneIso8601:!1,timezoneList:null,addSliderAccess:!1,sliderAccessArgs:null,controlType:"slider",defaultValue:null,parse:"strict"};e.extend(this._defaults,this.regional[""])}e.ui.timepicker=e.ui.timepicker||{};if(!e.ui.timepicker.version){e.extend(e.ui,{timepicker:{version:"1.1.1"}});e.extend(r.prototype,{$input:null,$altInput:null,$timeObj:null,inst:null,hour_slider:null,minute_slider:null,second_slider:null,millisec_slider:null,timezone_select:null,hour:0,minute:0,second:0,millisec:0,timezone:null,
defaultTimezone:"+0000",hourMinOriginal:null,minuteMinOriginal:null,secondMinOriginal:null,millisecMinOriginal:null,hourMaxOriginal:null,minuteMaxOriginal:null,secondMaxOriginal:null,millisecMaxOriginal:null,ampm:"",formattedDate:"",formattedTime:"",formattedDateTime:"",timezoneList:null,units:["hour","minute","second","millisec"],control:null,setDefaults:function(c){s(this._defaults,c||{});return this},_newInst:function(c,a){var b=new r,d={},f={},g,j;for(g in this._defaults)if(this._defaults.hasOwnProperty(g)){var h=
c.attr("time:"+g);if(h)try{d[g]=eval(h)}catch(l){d[g]=h}}g={beforeShow:function(a,d){if(e.isFunction(b._defaults.evnts.beforeShow))return b._defaults.evnts.beforeShow.call(c[0],a,d,b)},onChangeMonthYear:function(a,d,f){b._updateDateTime(f);e.isFunction(b._defaults.evnts.onChangeMonthYear)&&b._defaults.evnts.onChangeMonthYear.call(c[0],a,d,f,b)},onClose:function(a,d){!0===b.timeDefined&&""!==c.val()&&b._updateDateTime(d);e.isFunction(b._defaults.evnts.onClose)&&b._defaults.evnts.onClose.call(c[0],
a,d,b)}};for(j in g)g.hasOwnProperty(j)&&(f[j]=a[j]||null);b._defaults=e.extend({},this._defaults,d,a,g,{evnts:f,timepicker:b});b.amNames=e.map(b._defaults.amNames,function(a){return a.toUpperCase()});b.pmNames=e.map(b._defaults.pmNames,function(a){return a.toUpperCase()});"string"===typeof b._defaults.controlType?(void 0===e.fn[b._defaults.controlType]&&(b._defaults.controlType="select"),b.control=b._controls[b._defaults.controlType]):b.control=b._defaults.controlType;null===b._defaults.timezoneList&&
(d="-1200 -1100 -1000 -0930 -0900 -0800 -0700 -0600 -0500 -0430 -0400 -0330 -0300 -0200 -0100 +0000 +0100 +0200 +0300 +0330 +0400 +0430 +0500 +0530 +0545 +0600 +0630 +0700 +0800 +0845 +0900 +0930 +1000 +1030 +1100 +1130 +1200 +1245 +1300 +1400".split(" "),b._defaults.timezoneIso8601&&(d=e.map(d,function(a){return"+0000"==a?"Z":a.substring(0,3)+":"+a.substring(3)})),b._defaults.timezoneList=d);b.timezone=b._defaults.timezone;b.hour=b._defaults.hour;b.minute=b._defaults.minute;b.second=b._defaults.second;
b.millisec=b._defaults.millisec;b.ampm="";b.$input=c;a.altField&&(b.$altInput=e(a.altField).css({cursor:"pointer"}).focus(function(){c.trigger("focus")}));if(0===b._defaults.minDate||0===b._defaults.minDateTime)b._defaults.minDate=new Date;if(0===b._defaults.maxDate||0===b._defaults.maxDateTime)b._defaults.maxDate=new Date;void 0!==b._defaults.minDate&&b._defaults.minDate instanceof Date&&(b._defaults.minDateTime=new Date(b._defaults.minDate.getTime()));void 0!==b._defaults.minDateTime&&b._defaults.minDateTime instanceof
Date&&(b._defaults.minDate=new Date(b._defaults.minDateTime.getTime()));void 0!==b._defaults.maxDate&&b._defaults.maxDate instanceof Date&&(b._defaults.maxDateTime=new Date(b._defaults.maxDate.getTime()));void 0!==b._defaults.maxDateTime&&b._defaults.maxDateTime instanceof Date&&(b._defaults.maxDate=new Date(b._defaults.maxDateTime.getTime()));b.$input.bind("focus",function(){b._onFocus()});return b},_addTimePicker:function(c){var a=this.$altInput&&this._defaults.altFieldTimeOnly?this.$input.val()+
" "+this.$altInput.val():this.$input.val();this.timeDefined=this._parseTime(a);this._limitMinMaxDateTime(c,!1);this._injectTimePicker()},_parseTime:function(c,a){this.inst||(this.inst=e.datepicker._getInst(this.$input[0]));if(a||!this._defaults.timeOnly){var b=e.datepicker._get(this.inst,"dateFormat");try{var d=u(b,this._defaults.timeFormat,c,e.datepicker._getFormatConfig(this.inst),this._defaults);if(!d.timeObj)return!1;e.extend(this,d.timeObj)}catch(f){return e.datepicker.log("Error parsing the date/time string: "+
f+"\ndate/time string = "+c+"\ntimeFormat = "+this._defaults.timeFormat+"\ndateFormat = "+b),!1}}else{b=e.datepicker.parseTime(this._defaults.timeFormat,c,this._defaults);if(!b)return!1;e.extend(this,b)}return!0},_injectTimePicker:function(){var c=this.inst.dpDiv,a=this.inst.settings,b=this,d="",f="",g={},j={},h=null;if(0===c.find("div.ui-timepicker-div").length&&a.showTimepicker){for(var h='<div class="ui-timepicker-div'+(a.isRTL?" ui-timepicker-rtl":"")+'"><dl><dt class="ui_tpicker_time_label"'+
(a.showTime?"":' style="display:none;"')+">"+a.timeText+'</dt><dd class="ui_tpicker_time"'+(a.showTime?"":' style="display:none;"')+"></dd>",l=0,m=this.units.length;l<m;l++){d=this.units[l];f=d.substr(0,1).toUpperCase()+d.substr(1);g[d]=parseInt(a[d+"Max"]-(a[d+"Max"]-a[d+"Min"])%a["step"+f],10);j[d]=0;h+='<dt class="ui_tpicker_'+d+'_label"'+(a["show"+f]?"":' style="display:none;"')+">"+a[d+"Text"]+'</dt><dd class="ui_tpicker_'+d+'"><div class="ui_tpicker_'+d+'_slider"'+(a["show"+f]?"":' style="display:none;"')+
"></div>";if(a["show"+f]&&0<a[d+"Grid"]){h+='<div style="padding-left: 1px"><table class="ui-tpicker-grid-label"><tr>';if("hour"==d)for(f=a[d+"Min"];f<=g[d];f+=parseInt(a[d+"Grid"],10)){j[d]++;var k=e.datepicker.formatTime(p(a.pickerTimeFormat||a.timeFormat)?"hht":"HH",{hour:f},a),h=h+('<td data-for="'+d+'">'+k+"</td>")}else for(f=a[d+"Min"];f<=g[d];f+=parseInt(a[d+"Grid"],10))j[d]++,h+='<td data-for="'+d+'">'+(10>f?"0":"")+f+"</td>";h+="</tr></table></div>"}h+="</dd>"}var h=h+('<dt class="ui_tpicker_timezone_label"'+
(a.showTimezone?"":' style="display:none;"')+">"+a.timezoneText+"</dt>"),h=h+('<dd class="ui_tpicker_timezone" '+(a.showTimezone?"":' style="display:none;"')+"></dd>"),n=e(h+"</dl></div>");!0===a.timeOnly&&(n.prepend('<div class="ui-widget-header ui-helper-clearfix ui-corner-all"><div class="ui-datepicker-title">'+a.timeOnlyTitle+"</div></div>"),c.find(".ui-datepicker-header, .ui-datepicker-calendar").hide());l=0;for(m=b.units.length;l<m;l++)d=b.units[l],f=d.substr(0,1).toUpperCase()+d.substr(1),
b[d+"_slider"]=b.control.create(b,n.find(".ui_tpicker_"+d+"_slider"),d,b[d],a[d+"Min"],g[d],a["step"+f]),a["show"+f]&&0<a[d+"Grid"]&&(h=100*j[d]*a[d+"Grid"]/(g[d]-a[d+"Min"]),n.find(".ui_tpicker_"+d+" table").css({width:h+"%",marginLeft:a.isRTL?"0":h/(-2*j[d])+"%",marginRight:a.isRTL?h/(-2*j[d])+"%":"0",borderCollapse:"collapse"}).find("td").click(function(){var a=e(this),c=a.html(),f=parseInt(c.replace(/[^0-9]/g),10),c=c.replace(/[^apm]/ig),a=a.data("for");"hour"==a&&(-1!==c.indexOf("p")&&12>f?f+=
12:-1!==c.indexOf("a")&&12===f&&(f=0));b.control.value(b,b[a+"_slider"],d,f);b._onTimeChange();b._onSelectHandler()}).css({cursor:"pointer",width:100/j[d]+"%",textAlign:"center",overflow:"hidden"}));this.timezone_select=n.find(".ui_tpicker_timezone").append("<select></select>").find("select");e.fn.append.apply(this.timezone_select,e.map(a.timezoneList,function(a){return e("<option />").val("object"==typeof a?a.value:a).text("object"==typeof a?a.label:a)}));"undefined"!=typeof this.timezone&&null!==
this.timezone&&""!==this.timezone?e.timepicker.timeZoneOffsetString(new Date(this.inst.selectedYear,this.inst.selectedMonth,this.inst.selectedDay,12))==this.timezone?q(b):this.timezone_select.val(this.timezone):"undefined"!=typeof this.hour&&null!==this.hour&&""!==this.hour?this.timezone_select.val(a.defaultTimezone):q(b);this.timezone_select.change(function(){b._defaults.useLocalTimezone=!1;b._onTimeChange()});a=c.find(".ui-datepicker-buttonpane");a.length?a.before(n):c.append(n);this.$timeObj=n.find(".ui_tpicker_time");
null!==this.inst&&(c=this.timeDefined,this._onTimeChange(),this.timeDefined=c);if(this._defaults.addSliderAccess){var v=this._defaults.sliderAccessArgs,t=this._defaults.isRTL;v.isRTL=t;setTimeout(function(){if(0===n.find(".ui-slider-access").length){n.find(".ui-slider:visible").sliderAccess(v);var a=n.find(".ui-slider-access:eq(0)").outerWidth(!0);a&&n.find("table:visible").each(function(){var b=e(this),c=b.outerWidth(),d=b.css(t?"marginRight":"marginLeft").toString().replace("%",""),f=c-a,g={width:f,
marginRight:0,marginLeft:0};g[t?"marginRight":"marginLeft"]=d*f/c+"%";b.css(g)})}},10)}}},_limitMinMaxDateTime:function(c,a){var b=this._defaults,d=new Date(c.selectedYear,c.selectedMonth,c.selectedDay);if(this._defaults.showTimepicker){if(null!==e.datepicker._get(c,"minDateTime")&&void 0!==e.datepicker._get(c,"minDateTime")&&d){var f=e.datepicker._get(c,"minDateTime"),g=new Date(f.getFullYear(),f.getMonth(),f.getDate(),0,0,0,0);if(null===this.hourMinOriginal||null===this.minuteMinOriginal||null===
this.secondMinOriginal||null===this.millisecMinOriginal)this.hourMinOriginal=b.hourMin,this.minuteMinOriginal=b.minuteMin,this.secondMinOriginal=b.secondMin,this.millisecMinOriginal=b.millisecMin;c.settings.timeOnly||g.getTime()==d.getTime()?(this._defaults.hourMin=f.getHours(),this.hour<=this._defaults.hourMin?(this.hour=this._defaults.hourMin,this._defaults.minuteMin=f.getMinutes(),this.minute<=this._defaults.minuteMin?(this.minute=this._defaults.minuteMin,this._defaults.secondMin=f.getSeconds(),
this.second<=this._defaults.secondMin?(this.second=this._defaults.secondMin,this._defaults.millisecMin=f.getMilliseconds()):(this.millisec<this._defaults.millisecMin&&(this.millisec=this._defaults.millisecMin),this._defaults.millisecMin=this.millisecMinOriginal)):(this._defaults.secondMin=this.secondMinOriginal,this._defaults.millisecMin=this.millisecMinOriginal)):(this._defaults.minuteMin=this.minuteMinOriginal,this._defaults.secondMin=this.secondMinOriginal,this._defaults.millisecMin=this.millisecMinOriginal)):
(this._defaults.hourMin=this.hourMinOriginal,this._defaults.minuteMin=this.minuteMinOriginal,this._defaults.secondMin=this.secondMinOriginal,this._defaults.millisecMin=this.millisecMinOriginal)}if(null!==e.datepicker._get(c,"maxDateTime")&&void 0!==e.datepicker._get(c,"maxDateTime")&&d){f=e.datepicker._get(c,"maxDateTime");g=new Date(f.getFullYear(),f.getMonth(),f.getDate(),0,0,0,0);if(null===this.hourMaxOriginal||null===this.minuteMaxOriginal||null===this.secondMaxOriginal)this.hourMaxOriginal=b.hourMax,
this.minuteMaxOriginal=b.minuteMax,this.secondMaxOriginal=b.secondMax,this.millisecMaxOriginal=b.millisecMax;c.settings.timeOnly||g.getTime()==d.getTime()?(this._defaults.hourMax=f.getHours(),this.hour>=this._defaults.hourMax?(this.hour=this._defaults.hourMax,this._defaults.minuteMax=f.getMinutes(),this.minute>=this._defaults.minuteMax?(this.minute=this._defaults.minuteMax,this._defaults.secondMax=f.getSeconds()):this.second>=this._defaults.secondMax?(this.second=this._defaults.secondMax,this._defaults.millisecMax=
f.getMilliseconds()):(this.millisec>this._defaults.millisecMax&&(this.millisec=this._defaults.millisecMax),this._defaults.millisecMax=this.millisecMaxOriginal)):(this._defaults.minuteMax=this.minuteMaxOriginal,this._defaults.secondMax=this.secondMaxOriginal,this._defaults.millisecMax=this.millisecMaxOriginal)):(this._defaults.hourMax=this.hourMaxOriginal,this._defaults.minuteMax=this.minuteMaxOriginal,this._defaults.secondMax=this.secondMaxOriginal,this._defaults.millisecMax=this.millisecMaxOriginal)}void 0!==
a&&!0===a&&(b=parseInt(this._defaults.hourMax-(this._defaults.hourMax-this._defaults.hourMin)%this._defaults.stepHour,10),d=parseInt(this._defaults.minuteMax-(this._defaults.minuteMax-this._defaults.minuteMin)%this._defaults.stepMinute,10),f=parseInt(this._defaults.secondMax-(this._defaults.secondMax-this._defaults.secondMin)%this._defaults.stepSecond,10),g=parseInt(this._defaults.millisecMax-(this._defaults.millisecMax-this._defaults.millisecMin)%this._defaults.stepMillisec,10),this.hour_slider&&
(this.control.options(this,this.hour_slider,"hour",{min:this._defaults.hourMin,max:b}),this.control.value(this,this.hour_slider,"hour",this.hour)),this.minute_slider&&(this.control.options(this,this.minute_slider,"minute",{min:this._defaults.minuteMin,max:d}),this.control.value(this,this.minute_slider,"minute",this.minute)),this.second_slider&&(this.control.options(this,this.second_slider,"second",{min:this._defaults.secondMin,max:f}),this.control.value(this,this.second_slider,"second",this.second)),
this.millisec_slider&&(this.control.options(this,this.millisec_slider,"millisec",{min:this._defaults.millisecMin,max:g}),this.control.value(this,this.millisec_slider,"millisec",this.millisec)))}},_onTimeChange:function(){var c=this.hour_slider?this.control.value(this,this.hour_slider,"hour"):!1,a=this.minute_slider?this.control.value(this,this.minute_slider,"minute"):!1,b=this.second_slider?this.control.value(this,this.second_slider,"second"):!1,d=this.millisec_slider?this.control.value(this,this.millisec_slider,
"millisec"):!1,f=this.timezone_select?this.timezone_select.val():!1,g=this._defaults,j=g.pickerTimeFormat||g.timeFormat,h=g.pickerTimeSuffix||g.timeSuffix;"object"==typeof c&&(c=!1);"object"==typeof a&&(a=!1);"object"==typeof b&&(b=!1);"object"==typeof d&&(d=!1);"object"==typeof f&&(f=!1);!1!==c&&(c=parseInt(c,10));!1!==a&&(a=parseInt(a,10));!1!==b&&(b=parseInt(b,10));!1!==d&&(d=parseInt(d,10));var l=g[12>c?"amNames":"pmNames"][0],m=c!=this.hour||a!=this.minute||b!=this.second||d!=this.millisec||
0<this.ampm.length&&12>c!=(-1!==e.inArray(this.ampm.toUpperCase(),this.amNames))||null===this.timezone&&f!=this.defaultTimezone||null!==this.timezone&&f!=this.timezone;m&&(!1!==c&&(this.hour=c),!1!==a&&(this.minute=a),!1!==b&&(this.second=b),!1!==d&&(this.millisec=d),!1!==f&&(this.timezone=f),this.inst||(this.inst=e.datepicker._getInst(this.$input[0])),this._limitMinMaxDateTime(this.inst,!0));p(g.timeFormat)&&(this.ampm=l);this.formattedTime=e.datepicker.formatTime(g.timeFormat,this,g);this.$timeObj&&
(j===g.timeFormat?this.$timeObj.text(this.formattedTime+h):this.$timeObj.text(e.datepicker.formatTime(j,this,g)+h));this.timeDefined=!0;m&&this._updateDateTime()},_onSelectHandler:function(){var c=this._defaults.onSelect||this.inst.settings.onSelect,a=this.$input?this.$input[0]:null;c&&a&&c.apply(a,[this.formattedDateTime,this])},_updateDateTime:function(c){c=this.inst||c;var a=e.datepicker._daylightSavingAdjust(new Date(c.selectedYear,c.selectedMonth,c.selectedDay)),b=e.datepicker._get(c,"dateFormat");
c=e.datepicker._getFormatConfig(c);var d=null!==a&&this.timeDefined,b=this.formattedDate=e.datepicker.formatDate(b,null===a?new Date:a,c);if(!0===this._defaults.timeOnly)b=this.formattedTime;else if(!0!==this._defaults.timeOnly&&(this._defaults.alwaysSetTime||d))b+=this._defaults.separator+this.formattedTime+this._defaults.timeSuffix;this.formattedDateTime=b;if(this._defaults.showTimepicker)if(this.$altInput&&!0===this._defaults.altFieldTimeOnly)this.$altInput.val(this.formattedTime),this.$input.val(this.formattedDate);
else if(this.$altInput){this.$input.val(b);var b="",d=this._defaults.altSeparator?this._defaults.altSeparator:this._defaults.separator,f=this._defaults.altTimeSuffix?this._defaults.altTimeSuffix:this._defaults.timeSuffix;(b=this._defaults.altFormat?e.datepicker.formatDate(this._defaults.altFormat,null===a?new Date:a,c):this.formattedDate)&&(b+=d);b=this._defaults.altTimeFormat?b+(e.datepicker.formatTime(this._defaults.altTimeFormat,this,this._defaults)+f):b+(this.formattedTime+f);this.$altInput.val(b)}else this.$input.val(b);
else this.$input.val(this.formattedDate);this.$input.trigger("change")},_onFocus:function(){if(!this.$input.val()&&this._defaults.defaultValue){this.$input.val(this._defaults.defaultValue);var c=e.datepicker._getInst(this.$input.get(0)),a=e.datepicker._get(c,"timepicker");if(a&&a._defaults.timeOnly&&c.input.val()!=c.lastVal)try{e.datepicker._updateDatepicker(c)}catch(b){e.datepicker.log(b)}}},_controls:{slider:{create:function(c,a,b,d,f,g,j){var h=c._defaults.isRTL;return a.prop("slide",null).slider({orientation:"horizontal",
value:h?-1*d:d,min:h?-1*g:f,max:h?-1*f:g,step:j,slide:function(a,d){c.control.value(c,e(this),b,h?-1*d.value:d.value);c._onTimeChange()},stop:function(){c._onSelectHandler()}})},options:function(c,a,b,d,e){if(c._defaults.isRTL){if("string"==typeof d)return"min"==d||"max"==d?void 0!==e?a.slider(d,-1*e):Math.abs(a.slider(d)):a.slider(d);c=d.min;b=d.max;d.min=d.max=null;void 0!==c&&(d.max=-1*c);void 0!==b&&(d.min=-1*b);return a.slider(d)}return"string"==typeof d&&void 0!==e?a.slider(d,e):a.slider(d)},
value:function(c,a,b,d){return c._defaults.isRTL?void 0!==d?a.slider("value",-1*d):Math.abs(a.slider("value")):void 0!==d?a.slider("value",d):a.slider("value")}},select:{create:function(c,a,b,d,f,g,j){var h='<select class="ui-timepicker-select" data-unit="'+b+'" data-min="'+f+'" data-max="'+g+'" data-step="'+j+'">';for(c._defaults.timeFormat.indexOf("t");f<=g;f+=j)h+='<option value="'+f+'"'+(f==d?" selected":"")+">",h="hour"==b&&p(c._defaults.pickerTimeFormat||c._defaults.timeFormat)?h+e.datepicker.formatTime("hh TT",
{hour:f},c._defaults):"millisec"==b||10<=f?h+f:h+("0"+f.toString()),h+="</option>";h+="</select>";a.children("select").remove();e(h).appendTo(a).change(function(){c._onTimeChange();c._onSelectHandler()});return a},options:function(c,a,b,d,e){b={};var g=a.children("select");if("string"==typeof d){if(void 0===e)return g.data(d);b[d]=e}else b=d;return c.control.create(c,a,g.data("unit"),g.val(),b.min||g.data("min"),b.max||g.data("max"),b.step||g.data("step"))},value:function(c,a,b,d){c=a.children("select");
return void 0!==d?c.val(d):c.val()}}}});e.fn.extend({timepicker:function(c){c=c||{};var a=Array.prototype.slice.call(arguments);"object"==typeof c&&(a[0]=e.extend(c,{timeOnly:!0}));return e(this).each(function(){e.fn.datetimepicker.apply(e(this),a)})},datetimepicker:function(c){c=c||{};var a=arguments;return"string"==typeof c?"getDate"==c?e.fn.datepicker.apply(e(this[0]),a):this.each(function(){var b=e(this);b.datepicker.apply(b,a)}):this.each(function(){var a=e(this);a.datepicker(e.timepicker._newInst(a,
c)._defaults)})}});e.datepicker.parseDateTime=function(c,a,b,d,e){c=u(c,a,b,d,e);c.timeObj&&(a=c.timeObj,c.date.setHours(a.hour,a.minute,a.second,a.millisec));return c.date};e.datepicker.parseTime=function(c,a,b){b=s(s({},e.timepicker._defaults),b||{});var d=function(a,b,c){var d="^"+a.toString().replace(/([hH]{1,2}|mm?|ss?|[tT]{1,2}|[lz]|'.*?')/g,function(a){switch(a.charAt(0).toLowerCase()){case "h":return"(\\d?\\d)";case "m":return"(\\d?\\d)";case "s":return"(\\d?\\d)";case "l":return"(\\d?\\d?\\d)";
case "z":return"(z|[-+]\\d\\d:?\\d\\d|\\S+)?";case "t":a=c.amNames;var b=c.pmNames,d=[];a&&e.merge(d,a);b&&e.merge(d,b);d=e.map(d,function(a){return a.replace(/[.*+?|()\[\]{}\\]/g,"\\$&")});return"("+d.join("|")+")?";default:return"("+a.replace(/\'/g,"").replace(/(\.|\$|\^|\\|\/|\(|\)|\[|\]|\?|\+|\*)/g,function(a){return"\\"+a})+")?"}}).replace(/\s/g,"\\s?")+c.timeSuffix+"$",f=a.toLowerCase().match(/(h{1,2}|m{1,2}|s{1,2}|l{1}|t{1,2}|z|'.*?')/g);a={h:-1,m:-1,s:-1,l:-1,t:-1,z:-1};if(f)for(var g=0;g<
f.length;g++)-1==a[f[g].toString().charAt(0)]&&(a[f[g].toString().charAt(0)]=g+1);f="";d=b.match(RegExp(d,"i"));b={hour:0,minute:0,second:0,millisec:0};if(d){-1!==a.t&&(void 0===d[a.t]||0===d[a.t].length?(f="",b.ampm=""):(f=-1!==e.inArray(d[a.t].toUpperCase(),c.amNames)?"AM":"PM",b.ampm=c["AM"==f?"amNames":"pmNames"][0]));-1!==a.h&&(b.hour="AM"==f&&"12"==d[a.h]?0:"PM"==f&&"12"!=d[a.h]?parseInt(d[a.h],10)+12:Number(d[a.h]));-1!==a.m&&(b.minute=Number(d[a.m]));-1!==a.s&&(b.second=Number(d[a.s]));-1!==
a.l&&(b.millisec=Number(d[a.l]));if(-1!==a.z&&void 0!==d[a.z]){a=d[a.z].toUpperCase();switch(a.length){case 1:a=c.timezoneIso8601?"Z":"+0000";break;case 5:c.timezoneIso8601&&(a="0000"==a.substring(1)?"Z":a.substring(0,3)+":"+a.substring(3));break;case 6:c.timezoneIso8601?"00:00"==a.substring(1)&&(a="Z"):a="Z"==a||"00:00"==a.substring(1)?"+0000":a.replace(/:/,"")}b.timezone=a}return b}return!1};if("function"===typeof b.parse)return b.parse(c,a,b);if("loose"===b.parse){var f;a:{try{var g=new Date("2012-01-01 "+
a);f={hour:g.getHours(),minutes:g.getMinutes(),seconds:g.getSeconds(),millisec:g.getMilliseconds(),timezone:e.timepicker.timeZoneOffsetString(g)};break a}catch(j){try{f=d(c,a,b);break a}catch(h){e.datepicker.log("Unable to parse \ntimeString: "+a+"\ntimeFormat: "+c)}}f=!1}return f}return d(c,a,b)};e.datepicker.formatTime=function(c,a,b){b=b||{};b=e.extend({},e.timepicker._defaults,b);a=e.extend({hour:0,minute:0,second:0,millisec:0,timezone:"+0000"},a);var d=b.amNames[0],f=parseInt(a.hour,10);11<f&&
(d=b.pmNames[0]);c=c.replace(/(?:HH?|hh?|mm?|ss?|[tT]{1,2}|[lz]|('.*?'|".*?"))/g,function(c){switch(c){case "HH":return("0"+f).slice(-2);case "H":return f;case "hh":return("0"+w(f)).slice(-2);case "h":return w(f);case "mm":return("0"+a.minute).slice(-2);case "m":return a.minute;case "ss":return("0"+a.second).slice(-2);case "s":return a.second;case "l":return("00"+a.millisec).slice(-3);case "z":return null===a.timezone?b.defaultTimezone:a.timezone;case "T":return d.charAt(0).toUpperCase();case "TT":return d.toUpperCase();
case "t":return d.charAt(0).toLowerCase();case "tt":return d.toLowerCase();default:return c.replace(/\'/g,"")||"'"}});return c=e.trim(c)};e.datepicker._base_selectDate=e.datepicker._selectDate;e.datepicker._selectDate=function(c,a){var b=this._getInst(e(c)[0]),d=this._get(b,"timepicker");d?(d._limitMinMaxDateTime(b,!0),b.inline=b.stay_open=!0,this._base_selectDate(c,a),b.inline=b.stay_open=!1,this._notifyChange(b),this._updateDatepicker(b)):this._base_selectDate(c,a)};e.datepicker._base_updateDatepicker=
e.datepicker._updateDatepicker;e.datepicker._updateDatepicker=function(c){var a=c.input[0];if(!e.datepicker._curInst||!(e.datepicker._curInst!=c&&e.datepicker._datepickerShowing&&e.datepicker._lastInput!=a))if("boolean"!==typeof c.stay_open||!1===c.stay_open)if(this._base_updateDatepicker(c),a=this._get(c,"timepicker"))a._addTimePicker(c),a._defaults.useLocalTimezone&&(q(a,new Date(c.selectedYear,c.selectedMonth,c.selectedDay,12)),a._onTimeChange())};e.datepicker._base_doKeyPress=e.datepicker._doKeyPress;
e.datepicker._doKeyPress=function(c){var a=e.datepicker._getInst(c.target),b=e.datepicker._get(a,"timepicker");if(b&&e.datepicker._get(a,"constrainInput")){var d=p(b._defaults.timeFormat),a=e.datepicker._possibleChars(e.datepicker._get(a,"dateFormat")),b=b._defaults.timeFormat.toString().replace(/[hms]/g,"").replace(/TT/g,d?"APM":"").replace(/Tt/g,d?"AaPpMm":"").replace(/tT/g,d?"AaPpMm":"").replace(/T/g,d?"AP":"").replace(/tt/g,d?"apm":"").replace(/t/g,d?"ap":"")+" "+b._defaults.separator+b._defaults.timeSuffix+
(b._defaults.showTimezone?b._defaults.timezoneList.join(""):"")+b._defaults.amNames.join("")+b._defaults.pmNames.join("")+a,d=String.fromCharCode(void 0===c.charCode?c.keyCode:c.charCode);return c.ctrlKey||" ">d||!a||-1<b.indexOf(d)}return e.datepicker._base_doKeyPress(c)};e.datepicker._base_updateAlternate=e.datepicker._updateAlternate;e.datepicker._updateAlternate=function(c){var a=this._get(c,"timepicker");if(a){var b=a._defaults.altField;if(b){var d=this._getDate(c);c=e.datepicker._getFormatConfig(c);
var f,g=a._defaults.altSeparator?a._defaults.altSeparator:a._defaults.separator;f=a._defaults.altTimeSuffix?a._defaults.altTimeSuffix:a._defaults.timeSuffix;f=""+(e.datepicker.formatTime(null!==a._defaults.altTimeFormat?a._defaults.altTimeFormat:a._defaults.timeFormat,a,a._defaults)+f);!a._defaults.timeOnly&&!a._defaults.altFieldTimeOnly&&(f=a._defaults.altFormat?e.datepicker.formatDate(a._defaults.altFormat,null===d?new Date:d,c)+g+f:a.formattedDate+g+f);e(b).val(f)}}else e.datepicker._base_updateAlternate(c)};
e.datepicker._base_doKeyUp=e.datepicker._doKeyUp;e.datepicker._doKeyUp=function(c){var a=e.datepicker._getInst(c.target),b=e.datepicker._get(a,"timepicker");if(b&&b._defaults.timeOnly&&a.input.val()!=a.lastVal)try{e.datepicker._updateDatepicker(a)}catch(d){e.datepicker.log(d)}return e.datepicker._base_doKeyUp(c)};e.datepicker._base_gotoToday=e.datepicker._gotoToday;e.datepicker._gotoToday=function(c){var a=this._getInst(e(c)[0]),b=a.dpDiv;this._base_gotoToday(c);c=this._get(a,"timepicker");q(c);this._setTime(a,
new Date);e(".ui-datepicker-today",b).click()};e.datepicker._disableTimepickerDatepicker=function(c){var a=this._getInst(c);if(a){var b=this._get(a,"timepicker");e(c).datepicker("getDate");b&&(b._defaults.showTimepicker=!1,b._updateDateTime(a))}};e.datepicker._enableTimepickerDatepicker=function(c){var a=this._getInst(c);if(a){var b=this._get(a,"timepicker");e(c).datepicker("getDate");b&&(b._defaults.showTimepicker=!0,b._addTimePicker(a),b._updateDateTime(a))}};e.datepicker._setTime=function(c,a){var b=
this._get(c,"timepicker");if(b){var d=b._defaults;b.hour=a?a.getHours():d.hour;b.minute=a?a.getMinutes():d.minute;b.second=a?a.getSeconds():d.second;b.millisec=a?a.getMilliseconds():d.millisec;b._limitMinMaxDateTime(c,!0);b._onTimeChange();b._updateDateTime(c)}};e.datepicker._setTimeDatepicker=function(c,a,b){if(c=this._getInst(c)){var d=this._get(c,"timepicker");d&&(this._setDateFromField(c),a&&("string"==typeof a?(d._parseTime(a,b),a=new Date,a.setHours(d.hour,d.minute,d.second,d.millisec)):a=new Date(a.getTime()),
"Invalid Date"==a.toString()&&(a=void 0),this._setTime(c,a)))}};e.datepicker._base_setDateDatepicker=e.datepicker._setDateDatepicker;e.datepicker._setDateDatepicker=function(c,a){var b=this._getInst(c);if(b){var d=a instanceof Date?new Date(a.getTime()):a;this._updateDatepicker(b);this._base_setDateDatepicker.apply(this,arguments);this._setTimeDatepicker(c,d,!0)}};e.datepicker._base_getDateDatepicker=e.datepicker._getDateDatepicker;e.datepicker._getDateDatepicker=function(c,a){var b=this._getInst(c);
if(b){var d=this._get(b,"timepicker");return d?(void 0===b.lastVal&&this._setDateFromField(b,a),(b=this._getDate(b))&&d._parseTime(e(c).val(),d.timeOnly)&&b.setHours(d.hour,d.minute,d.second,d.millisec),b):this._base_getDateDatepicker(c,a)}};e.datepicker._base_parseDate=e.datepicker.parseDate;e.datepicker.parseDate=function(c,a,b){var d;try{d=this._base_parseDate(c,a,b)}catch(f){d=this._base_parseDate(c,a.substring(0,a.length-(f.length-f.indexOf(":")-2)),b),e.datepicker.log("Error parsing the date string: "+
f+"\ndate string = "+a+"\ndate format = "+c)}return d};e.datepicker._base_formatDate=e.datepicker._formatDate;e.datepicker._formatDate=function(c){var a=this._get(c,"timepicker");return a?(a._updateDateTime(c),a.$input.val()):this._base_formatDate(c)};e.datepicker._base_optionDatepicker=e.datepicker._optionDatepicker;e.datepicker._optionDatepicker=function(c,a,b){var d=this._getInst(c),f;if(!d)return null;if(d=this._get(d,"timepicker")){var g=null,j=null,h=null,l=d._defaults.evnts,m={},k;if("string"==
typeof a)if("minDate"===a||"minDateTime"===a)g=b;else if("maxDate"===a||"maxDateTime"===a)j=b;else if("onSelect"===a)h=b;else{if(l.hasOwnProperty(a)){if("undefined"===typeof b)return l[a];m[a]=b;f={}}}else if("object"==typeof a)for(k in a.minDate?g=a.minDate:a.minDateTime?g=a.minDateTime:a.maxDate?j=a.maxDate:a.maxDateTime&&(j=a.maxDateTime),l)l.hasOwnProperty(k)&&a[k]&&(m[k]=a[k]);for(k in m)m.hasOwnProperty(k)&&(l[k]=m[k],f||(f=e.extend({},a)),delete f[k]);if(k=f)a:{k=f;for(var n in k)if(k.hasOwnProperty(k)){k=
!1;break a}k=!0}if(k)return;g?(g=0===g?new Date:new Date(g),d._defaults.minDate=g,d._defaults.minDateTime=g):j?(j=0===j?new Date:new Date(j),d._defaults.maxDate=j,d._defaults.maxDateTime=j):h&&(d._defaults.onSelect=h)}return void 0===b?this._base_optionDatepicker.call(e.datepicker,c,a):this._base_optionDatepicker.call(e.datepicker,c,f||a,b)};var s=function(c,a){e.extend(c,a);for(var b in a)if(null===a[b]||void 0===a[b])c[b]=a[b];return c},p=function(c){return-1!==c.indexOf("t")&&-1!==c.indexOf("h")},
w=function(c){12<c&&(c-=12);0==c&&(c=12);return String(c)},u=function(c,a,b,d,f){var g;a:{try{var j=f&&f.separator?f.separator:e.timepicker._defaults.separator,h=(f&&f.timeFormat?f.timeFormat:e.timepicker._defaults.timeFormat).split(j).length,l=b.split(j),m=l.length;if(1<m){g=[l.splice(0,m-h).join(j),l.splice(0,h).join(j)];break a}}catch(k){if(e.datepicker.log("Could not split the date from the time. Please check the following datetimepicker options\nthrown error: "+k+"\ndateTimeString"+b+"\ndateFormat = "+
c+"\nseparator = "+f.separator+"\ntimeFormat = "+f.timeFormat),0<=k.indexOf(":")){g=b.length-(k.length-k.indexOf(":")-2);b.substring(g);g=[e.trim(b.substring(0,g)),e.trim(b.substring(g))];break a}else throw k;}g=[b,""]}c=e.datepicker._base_parseDate(c,g[0],d);if(""!==g[1]){a=e.datepicker.parseTime(a,g[1],f);if(null===a)throw"Wrong time format";return{date:c,timeObj:a}}return{date:c}},q=function(c,a){if(c&&c.timezone_select){c._defaults.useLocalTimezone=!0;var b=e.timepicker.timeZoneOffsetString("undefined"!==
typeof a?a:new Date);c._defaults.timezoneIso8601&&(b=b.substring(0,3)+":"+b.substring(3));c.timezone_select.val(b)}};e.timepicker=new r;e.timepicker.timeZoneOffsetString=function(c){c=-1*c.getTimezoneOffset();var a=c%60;return(0<=c?"+":"-")+("0"+(101*((c-a)/60)).toString()).substr(-2)+("0"+(101*a).toString()).substr(-2)};e.timepicker.timeRange=function(c,a,b){return e.timepicker.handleRange("timepicker",c,a,b)};e.timepicker.dateTimeRange=function(c,a,b){e.timepicker.dateRange(c,a,b,"datetimepicker")};
e.timepicker.dateRange=function(c,a,b,d){e.timepicker.handleRange(d||"datepicker",c,a,b)};e.timepicker.handleRange=function(c,a,b,d){function f(c,d,e){d.val()&&new Date(a.val())>new Date(b.val())&&d.val(e)}function g(a,b,d){e(a).val()&&(a=e(a)[c].call(e(a),"getDate"),a.getTime&&e(b)[c].call(e(b),"option",d,a))}e.fn[c].call(a,e.extend({onClose:function(a){f(this,b,a)},onSelect:function(){g(this,b,"minDate")}},d,d.start));e.fn[c].call(b,e.extend({onClose:function(b){f(this,a,b)},onSelect:function(){g(this,
a,"maxDate")}},d,d.end));"timepicker"!=c&&d.reformat&&e([a,b]).each(function(){var a=e(this)[c].call(e(this),"option","dateFormat"),b=new Date(e(this).val());e(this).val()&&b&&e(this).val(e.datepicker.formatDate(a,b))});f(a,b,a.val());g(a,b,"minDate");g(b,a,"maxDate");return e([a.get(0),b.get(0)])};e.timepicker.version="1.1.1"}})(jQuery);