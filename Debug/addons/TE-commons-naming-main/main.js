/*
 * Copyright (c) 2021 istallia
 * Released under the MIT license
 * https://github.com/istallia/TE-commons-naming/blob/main/LICENSE
 */

/* --- 実際にリネーム処理を行う --- */
let FV       = GetFolderView(Ctrl, pt);
let Selected = FV.SelectedItems();
if (Selected && Selected.Count > 0) {
	for (let i = 0; i < Selected.Count; i++) {
		const file_path = Selected.Item(i).Path;
		const file_dir  = fso.GetParentFolderName(file_path);
		const file_name = fso.GetBaseName(file_path);
		const file_ext  = fso.GetExtensionName(file_path);
		const matches   = file_name.match(/^nc\d{1,12}$/);
		if (matches) {
			const title = getCommonsTitle(file_name);
			if (title) fso.moveFile(file_path, file_dir+'\\'+file_name+'_'+replaceForFileName(title)+'.'+file_ext);
		}
	}
}

/* --- インターネット経由でコモンズ素材のタイトル取得 --- */
function getCommonsTitle(id) {
	try {
		const http = new ActiveXObject("Msxml2.ServerXMLHTTP");
		http.open('GET', 'https://commons.nicovideo.jp/material/'+id, false);
		http.send();
		const text    = http.responseText.replace(/\n/g, '');
		const matches = text.match(/<div class="materialHeadTitle">([\s\S]{1,128})<\/div>/);
		if (matches) {
			return decodeHTMLSpecialChars(matches[1]);
		} else {
			return null;
		}
	} catch(e) {
		Addons.DebugConsole.warning(e.message);
		return null;
	}
}

/* --- HTML特殊文字をもとに戻す --- */
function decodeHTMLSpecialChars(text) {
	return text
		.replace(/&nbsp;/g, ' ')
		.replace(/&apos;/g, '\'')
		.replace(/&quot;/g, '"')
		.replace(/&amp;/g, '&')
		.replace(/&lt;/g, '<')
		.replace(/&gt;/g, '>')
		.replace(/&copy;/g, '(C)');
}

/* --- ファイルに使えない文字を置換 --- */
function replaceForFileName(name) {
	return name
		.replace(/\\/g, '＼')
		.replace(/\//g, '／')
		.replace(/:/g, '：')
		.replace(/\*/g, '＊')
		.replace(/\?/g, '？')
		.replace(/\"/g, '”')
		.replace(/\|/g, '｜')
		.replace(/</g, '＜')
		.replace(/>/g, '＞');
}
