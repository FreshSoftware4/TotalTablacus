var img = File.openDialog("画像を選択してください");
if (img) {
  var exif = img.exifValues;
  var exifStr = "";
  for (var key in exif) {
    exifStr += key + ": " + exif[key] + "\n";
  }
  System.setClipboard(exifStr);
}
