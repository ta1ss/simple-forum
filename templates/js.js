function modifyComment(element) {
  var id = element.getAttribute("data-id");
  var value = element.getAttribute("data-url");
  var commentDiv = document.getElementById("modify-comment-" + id);
  var textarea = commentDiv.querySelector("textarea");
  var postContent = document.getElementById("comment-" + id);
  var postValue = postContent.textContent;
  if (!textarea){
    textarea = document.createElement("textarea");
    textarea.name = "textarea";
    textarea.className = "editarea"
    textarea.value = postValue;
    textarea.required = true
    var submitButton = document.createElement("button");
      submitButton.type = "submit";
      submitButton.textContent = "Save";
      submitButton.className = "save-button"
    var form = document.createElement("form");
      form.method = "post";
      form.action = "/modify-comment";
      commentDiv.appendChild(form);
      form.appendChild(textarea);
      form.appendChild(submitButton);
    var idInput = document.createElement("input");
      idInput.type = "hidden";
      idInput.name = "id";
      idInput.value = id;
      form.appendChild(idInput);
    var url = document.createElement("input");
      url.type = "hidden";
      url.name = "url";
      url.value = value;
      form.appendChild(url);
}
}

function modifyPost(element) {
  var id = element.getAttribute("data-id");
  var value = element.getAttribute("data-url");
  var commentDiv = document.getElementById("modify-post-" + id);
  var textarea = commentDiv.querySelector("textarea");
  var postContent = document.getElementById("post-" + id);
  var postValue = postContent.textContent;
  var titleContent = document.getElementById("title-" + id);
  var titleValue = titleContent.textContent;
  if (!textarea){
    textarea = document.createElement("textarea");
    textarea.name = "postarea";
    textarea.className = "editarea";
    textarea.value= postValue;
    textarea.rows = 10;
    textarea.required = true
    var submitButton = document.createElement("button");
      submitButton.type = "submit";
      submitButton.textContent = "Save";
      submitButton.className = "save-button"
    var textarea2 = document.createElement("textarea");
      textarea2 = document.createElement("textarea");
      textarea2.className = "editarea"
      textarea2.name = "titlearea";
      textarea2.value= titleValue;
      textarea2.rows = 3;
      textarea2.required = true
    var form = document.createElement("form");
      form.method = "post";
      form.action = "/modify-post";
      commentDiv.appendChild(form);
      form.appendChild(textarea2);
      form.appendChild(textarea);
      form.appendChild(submitButton);
    var idInput = document.createElement("input");
      idInput.type = "hidden";
      idInput.name = "id";
      idInput.value = id;
      form.appendChild(idInput);
    var url = document.createElement("input");
      url.type = "hidden";
      url.name = "url";
      url.value = value;
      form.appendChild(url);
  }
}

function validateFileSize() {
  var fileInput = document.getElementById("fileInput");
  var fileSize = fileInput.files[0].size;
  var maxSize = 20*1024*1024; // 20MB
  if (fileSize > maxSize) {
    alert("File size should be less than 20MB");
    fileInput.value = "";
  }
}

function seeNotifications(){
  var notifications = document.getElementById("notifications");
  var value = notifications.getAttribute("data-url");
  console.log(value)
  var div = notifications.querySelector("div")
  if (!div){
    div = document.createElement("div")
    div.className = "notification-div"
    notifications.appendChild(div)
    var ul = document.createElement("ul")
    ul.style.listStyleType = "none";
    div.appendChild(ul);
    fetch('/api/notifications')
      .then(response => {
      if(!response.ok){
        throw new Error("HTTP error "+ response.status);
        }
        return response.json()
        })
        .then(data => {
        data.forEach(function(notification) {
        var li = document.createElement("li");
        li.innerHTML = notification;
        ul.appendChild(li);
      });
    });
    var form = document.createElement("form");
      form.method = "post";
      form.action = "/delete-notifications";
      div.appendChild(form)
    var url = document.createElement("input");
      url.type = "hidden";
      url.name = "url";
      url.value = value;
      form.appendChild(url);
    var submitButton = document.createElement("button");
      submitButton.type = "submit";
      submitButton.className = "clear-button"
      submitButton.textContent = "clear";
      form.appendChild(submitButton)
  } else {
    notifications.removeChild(div)
  }
}

function flagPost(element) {
  var id = element.getAttribute("data-id");
  var flagDiv = document.getElementById("flag-post-" + id);
  var div = flagDiv.querySelector("div")
  if (!div) {
    div = document.createElement("div")
    div.className = "flagpost-div"
    var img = document.createElement("img");
    img.src = "media/flag.png";
    img.alt = "Flag Image";
    img.height = "50";
    img.width = "50";
    div.appendChild(img);
    flagDiv.appendChild(div);

    fetch("/add-flag", {
      method: "POST",
      headers: {
        "Content-Type": "application/json",
      },
      body: JSON.stringify({ id: id, type: "post" }),
    });
  } else {
    flagDiv.removeChild(div)
    fetch("/remove-flag", {
      method: "POST",
      headers: {
        "Content-Type": "application/json",
      },
      body: JSON.stringify({ id: id, type: "post" }),
    });
  }
}

function flagComment(element) {
  var id = element.getAttribute("data-id");
  var flagDiv = document.getElementById("flag-comment-" + id);
  var div = flagDiv.querySelector("div")
  if (!div) {
    div = document.createElement("div")
    div.className = "flagcomment-div"
    var img = document.createElement("img");
    img.src = "media/flag.png";
    img.alt = "Flag Image";
    img.height = "50";
    img.width = "50";
    div.appendChild(img);
    flagDiv.appendChild(div);

    fetch("/add-flag", {
      method: "POST",
      headers: {
        "Content-Type": "application/json",
      },
      body: JSON.stringify({ id: id, type: "comment" }),
    });
  } else {
    flagDiv.removeChild(div)
    fetch("/remove-flag", {
      method: "POST",
      headers: {
        "Content-Type": "application/json",
      },
      body: JSON.stringify({ id: id, type: "comment" }),
    });
  }
}



