"use strict"
$(function() {
	$(".time").timepicker();
	$("#date").datepicker($.datepicker.regional[ "sk" ]);
	$("#date").datepicker("option", "dateFormat", "yy-mm-dd");
});