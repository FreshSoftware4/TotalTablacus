
var File = new ActiveXObject("Scripting.FileSystemObject");

var imgFile = File.openDialog();
if (imgFile) {
  var exif = imgFile.exifValues;
  var exifStr = "";
  for (var key in exif) {
    exifStr += key + ": " + exif[key] + "\n";
  }
  System.setClipboard(exifStr);
}