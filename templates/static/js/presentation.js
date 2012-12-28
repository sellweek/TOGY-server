"use strict"
$(function() {
	var titleShow = _.template($("#title-show-template").html());
	var titleEdit = _.template($("#title-edit-template").html());
	var descShow = _.template($("#description-show-template").html());
	var descEdit = _.template($("#description-edit-template").html());
	var loading = _.template($("#loading-template").html());

	$("#title").on("dblclick", function(e) {
		$("#title-container").html(titleEdit({title: $(e.currentTarget).html()}));
	});

	$("#title-container").on("submit", "form", function() {
		var title = $("#title-field").val();
		$("#title-container").html(loading());
		$.post("/api/presentation/"+presentationKey+"/name", title, function() {
			$("#title-container").html(titleShow({title: title}));
		});
		return false;
	});

});