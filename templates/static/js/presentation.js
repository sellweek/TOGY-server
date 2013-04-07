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

	$("input[name=datetime]").datetimepicker({
		stepMinute: 5
	});
	$("input[name=datetime]").datetimepicker($.datepicker.regional[ "sk" ]);
	$("input[name=datetime]").datetimepicker("option", "dateFormat", "yy-mm-dd");

	$("#schedule-activation").on("submit", function() {
		$.post("/api/presentation/"+presentationKey+"/schedule", $("input[name=datetime]").datetimepicker("getDate").toString(), function(data) {
			if (data==="") {
				$("#schedule-activation-container").html("Aktivácia naplánovaná");
				$("#schedule-activation-container").addClass("alert alert-success");
			} else {
				$("#schedule-activation-container").html("Chyba: "+data);
				$("#schedule-activation-container").addClass("alert alert-error");
			}
		});
		return false;
	})
});