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
  fetch(url, { mode: 'no-cors'})
  .then(function(response) {
    console.log(response);
  })
  .catch(function(err) {
    console.log('Fetch Error: ', err);
  });
}