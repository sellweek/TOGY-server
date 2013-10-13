"use strict"
$(function() {
	$.datepicker.setDefaults($.datepicker.regional["sk"])
	$(".time").timepicker();
	$("#date").datepicker({
		"dateFormat": "yy-mm-dd"
	});
	$("#date").datepicker("setDate", $("#date").val())
});