function doGet(e) {
  var blogUrl = JSON.parse(e.parameters.url);
  var spreadsheet = SpreadsheetApp.getActiveSheet();

  var tf = spreadsheet.createTextFinder(blogUrl);
  var all = tf.findAll();

  if (all.length == 0) {
    spreadsheet.appendRow([blogUrl,1]);
  } else {
    if (all.length > 1) {
      Logger.log("Blog url: %s has more then 1 matches", blogUrl); 
    }
    for (var i = 0; i < all.length; i++) {
      const clickCell = "B"+all[i].getRow();
      const clickCount = spreadsheet.getRange(clickCell).getValue();
      Logger.log('Value: %s for range: %s', clickCount, clickCell);
      spreadsheet.getRange(clickCell).setValue(clickCount+1);
    }
  }
}
