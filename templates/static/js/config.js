"use strict"
$(function() {
	$(".timeField").datetimepicker({
		format: "HH:mm",
		pickDate: false,
		language: "sk",

	});

	var setTimeFieldState = function() {
    	if ($("input[name=overrideState]:checked").val() !== "0") {
			$("input:not(input[name=overrideState]):not(input[type=submit])").attr("disabled", "true");
		} else {
			$("input:not(input[name=overrideState]):not(input[type=submit])").removeAttr("disabled");
		}
    }

	setTimeFieldState();

	$("input[name=overrideState]").on("click", function(){
		setTimeFieldState();
	});

	$('.conf-form').on('submit', function() {
        $('input:not(input[name=overrideState]):not(input[type=submit])').removeAttr('disabled');
    });
})