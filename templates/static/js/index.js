"use strict"
$(function() {
	$(".preview").one("dblclick", function(e) {
		var url = $(e.currentTarget).find(".viewer-url").html();
		$(e.currentTarget).html('<iframe src="'+url+'" width="500" height="300"></iframe>')
	});
});

