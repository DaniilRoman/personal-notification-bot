function showHide(id) {
  var x = document.getElementById(id);
  if (x.style.display === "none") {
    x.style.display = "block";
  } else {
    x.style.display = "none";
  }
}

function trackClick(url) {
  const websiteUrl = new URL(url);
  const hostname = websiteUrl.hostname
  const appScriptId = "{{.AppScriptId}}"
  const appscriptEndpoint = `https://script.google.com/macros/s/${appScriptId}/exec?url="${hostname}"`
  doGet(appscriptEndpoint)
}

function doGet(url) {
  const xhr = new XMLHttpRequest();
  xhr.open('GET', url);

  xhr.onerror = function(error) {
    console.error('Error:', error);
  };

  xhr.send();
}