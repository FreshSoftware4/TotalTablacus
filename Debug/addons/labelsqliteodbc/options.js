const s = [];
if (await MainWindow.Sync.Label && await MainWindow.Sync.LabelSQLiteOdbc) {
	s.push('<button onclick="Addons.LabelSQLiteOdbc.Import()">Load</button><br><button onclick="Addons.LabelSQLiteOdbc.Export()">Save</button><br><br>');
}
s.push('<button title="http://www.ch-werner.de/sqliteodbc/" onclick="wsh.Run(this.title)">Get SQLite ODBC Driver...</button>');
SetTabContents(0, "", s);

Addons.LabelSQLiteOdbc = {
	Import: async function () {
		const commdlg = await api.CreateObject("CommonDialog");
		commdlg.InitDir = BuildPath(ui_.DataFolder, "config")
		commdlg.Filter = await MakeCommDlgFilter("*.tsv");
		commdlg.Flags = OFN_FILEMUSTEXIST;
		if (await commdlg.ShowOpen()) {
			MainWindow.Sync.LabelSQLiteOdbc.Load(await commdlg.FileName);
		}
	},

	Export: async function () {
		const commdlg = await api.CreateObject("CommonDialog");
		commdlg.InitDir = BuildPath(ui_.DataFolder, "config")
		commdlg.Filter = await MakeCommDlgFilter("*.tsv");
		commdlg.DefExt = "tsv";
		commdlg.Flags = OFN_OVERWRITEPROMPT;
		if (await commdlg.ShowSave()) {
			MainWindow.Sync.LabelSQLiteOdbc.Save(await commdlg.FileName);
		}
	}
}
