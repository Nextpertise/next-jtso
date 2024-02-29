let eventSource;
const browseButton = document.getElementById("browse");
const resultDiv = document.getElementById("result");
const modal = document.getElementById("modalcore")
modal.style.scrollBehavior = 'smooth';

$(document).ready(function () {

  $('#searching').on('input', function () {
    var searchString = $(this).val();
    $('#jstree').jstree('search', searchString);
  });

  // Collapse button click event
  $('#collapse').on('click', function () {
    $('#result').jstree('close_all');
  });

  // Expand button click event
  $('#expand').on('click', function () {
    $('#result').jstree('open_all');
  });

  $("#result").jstree({
    "core": {
      "data": []
    },
    "themes": {
      "theme": "default",
      "dots": true,
      "icons": true
    },
    "plugins": ["state", "sort", "search"]
  });
});


browseButton.addEventListener("click", function () {

  var p = document.getElementById("pathName").value.trim();
  var m = document.getElementById("merge").checked;
  var r = document.getElementById("router").value.trim();

  var dataToSend = {
    "shortname": r,
    "xpath": p,
    "merge": m
  };
  fetch("/searchxpath", {
      method: "POST",
      headers: {
        "Content-Type": "application/json",
      },
      body: JSON.stringify(dataToSend),
    })
    .then(response => response.json())
    .then(data => {
      browseButton.disabled = true;


      // Start the EventSource for streaming
      eventSource = new EventSource("/stream");
      modal.innerHTML = '';
      $('#logs').modal('show');
      eventSource.onmessage = function (event) {
        const data = JSON.parse(event.data);
        appendContent(data.msg);
        scrollToBottom()
        if (data.status == "END") {

          eventSource.close();
          browseButton.disabled = false;
          $('#result').jstree(true).settings.core.data = JSON.parse(data.payload);
          $('#result').jstree(true).refresh();
          $('#logs').modal('hide');
          alertify.success('Here the results!')
        }
      };

      eventSource.onerror = function (event) {
        alertify.alert("JSTO...", "Unexpected error: " + event);
        browseButton.disabled = false;
        $('#logs').modal('hide');
        eventSource.close();
      };
    })
    .catch(error => {
      alertify.alert("JSTO...", "Unexpected error: " + error);
      browseButton.disabled = false;
      $('#logs').modal('hide');
    });
});

// Function to append new content
function appendContent(text) {
  var newElement = document.createElement('div');
  newElement.innerHTML = text;
  modal.appendChild(newElement);
}

// Function to scroll to the bottom with smooth scrolling
function scrollToBottom() {
  modal.scrollTop = modal.scrollHeight;
}