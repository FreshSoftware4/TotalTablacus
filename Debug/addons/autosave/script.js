if (window.Addon == 1) {
	setInterval(SaveConfig, (await GetAddonOption("autosave", "Interval") || 5) *  60000);
} else {
	SetTabContents(0, "General", '<table><td><input type="text" name="Interval" size="4" placeholder="5"></td><td><input type="button" value="Default" onclick="document.F.Interval.value=\'\'"></td></tr></table>');
}
