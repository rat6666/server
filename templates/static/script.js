function addCar(car) {
  appendToUsrTable(car);
}

function flashMessage(msg) {
  $(".flashMsg").remove();
  $(".row").prepend(`
        <div class="col-sm-12"><div class="flashMsg alert alert-success alert-dismissible fade in" role="alert"> <button type="button" class="close" data-dismiss="alert" aria-label="Close"><span aria-hidden="true">×</span></button> <strong>${msg}</strong></div></div>
    `);
}


function appendToUsrTable(car) {
  $("#carTable > tbody:last-child").append(`
        <tr id="car-${car.id}">
            <td class="carData" IDmongo="IDmongo">${car.id}</td>
            '<td class="carData" ID="ID">${car.ID}</td>
            '<td class="carData" name="Model">${car.model}</td>
            '<td class="carData" name="Date">${car.date}</td>
            '<td align="center">
                <button class="btn btn-success form-control" onClick="editCar('${car.id}');" data-toggle="modal" data-target="#myModal")">EDIT</button>
            </td>
            <td align="center">
                <button class="btn btn-danger form-control" onClick="deleteCar('${car.id}');">DELETE</button>
            </td>
        </tr>
    `);
}

function editCar(id) {
  $.getJSON( "/", function(car) {
    $.each( car, function( key, val ) {
      // console.log(val.id);
      if (val.id == id) {
        // console.log(id);
      $(".modal-body").empty().append(`
                <form id="updateCar" action="">
                    <label for="mongoID">mongoID</label>
                    <input class="form-control" type="text" name="id" value="${val.id}"/>
                    <label for="ID">ID</label>
                    <input class="form-control" type="text" name="ID" value="${val.ID}"/>
                    <label for="model">Model</label>
                    <input class="form-control" type="text" name="model" value="${val.model}"/>
                    <label for="date">date</label>
                    <input class="form-control" type="text" name="date" value="${val.date}"/>
            `);
      $(".modal-footer").empty().append(`
                    <button type="button" type="submit" class="btn btn-primary" onClick="updatecar('${id}')">Save changes</button>
                    <button type="button" class="btn btn-default" data-dismiss="modal">Close</button>
                </form>
            `);
    }
  });
  });
}

function updatecar(id) {
  $("#carTable #car-" + id).remove();
  var msg = "Car updated successfully!";
  var car = JSON.stringify($("#updateCar").serializeJSON());
  $.ajax({
    url: "/",
    type: "PUT",
    contentType: false,
    data: car,
    dataType: "json",
  });
  var json_obj = JSON.parse(car);
  addCar(json_obj);
  $(".modal").modal("toggle");
  flashMessage(msg);

}


$("#submitform").click(function (event) {
  var msg = "Car added successfully!";
  var car = JSON.stringify($("#addCar").serializeJSON());
  $.ajax({
    url: "/",
    type: "POST",
    contentType: false,
    data: car,
    dataType: "json",
  });
  var json_obj = JSON.parse(car);
  addCar(json_obj);
  flashMessage(msg);
});


function deleteCar(id) {
  var action = confirm("Are you sure you want to delete this car?");
  var msg = "Car deleted successfully!";
  fetch('/'+ id, {
    method: 'DELETE',
  })
  flashMessage(msg);
  $("#carTable #car-" + id).remove();
}

$(document).ready(function() { 
  var car = {};
  $("#fetch").click(function(event){ 
    $('td').remove();
    $.getJSON( "/", function(emp) {
      $.each( emp, function( key, val ) {
        addCar(val);
      });
    });
	}); 
}); 

