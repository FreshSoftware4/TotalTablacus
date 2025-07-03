// 画像ファイルをData URL形式に変換する
var imgUrl = window.URL.createObjectURL(imgFile);

// Exifライブラリを使用してExifデータを取得する
var exifData = EXIF.readFromBinaryFile(imgUrl);

// SQLiteデータベースにExifデータを登録する
db.transaction(function(tx) {
  tx.executeSql("INSERT INTO exif (data) VALUES (?)", [exifData]);
});
