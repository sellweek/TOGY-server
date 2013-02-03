"use strict"
$(function() {
	var titleShow = _.template($("#title-show-template").html());
	var titleEdit = _.template($("#title-edit-template").html());
	var descShow = _.template($("#description-show-template").html());
	var descEdit = _.template($("#description-edit-template").html());
	var loading = _.template($("#loading-template").html());

	$("#title-container").on("dblclick", "h1", function(e) {
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

	$("#description-container").on("dblclick", "div", function() {
		$("#description-container").html(loading());
		$.get("/api/presentation/"+presentationKey+"/description", "", function(data) {
			$("#description-container").html(descEdit({markdown: data}));
		});
	});

	$("#description-container").on("submit", "form", function() {
		var text = $("#description-field").val()
		$("#description-container").html(loading());
		$.post("/api/presentation/"+presentationKey+"/description", text, function(data) {
			$("#description-container").html(descShow({text: data}));
		});
		return false;
	});

});