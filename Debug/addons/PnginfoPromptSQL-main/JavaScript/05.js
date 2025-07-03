var fso = new ActiveXObject("Scripting.FileSystemObject");
var System = WScript.CreateObject("WScript.Shell");
var imgFile = fso.openDialog();
if (imgFile) {
  var exif = imgFile.exifValues;
  var exifStr = "";
  for (var key in exif) {
    exifStr += key + ": " + exif[key] + "\n";
  }
  System
}