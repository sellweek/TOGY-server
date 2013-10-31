"use strict"
$(function() {
	$(".time").datetimepicker({
		format: "HH:mm",
		pickDate: false,
		language: "sk"
	});
	$("#date").datetimepicker({
		language:"sk",
		pickTime: false
	});
});