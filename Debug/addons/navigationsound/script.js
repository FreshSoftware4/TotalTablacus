const Addon_Id = "navigationsound";
if (window.Addon == 1) {
	Addons.NavigationSound = {
		Path: await ExtractPath(te, await GetAddonOption(Addon_Id, "Path") || "%SystemRoot%\\Media\\Windows Navigation Start.wav")
	}
	AddEvent("ListViewCreated", function () {
		api.PlaySound(Addons.NavigationSound.Path, null, 3);
	});
} else {
	importScript("addons\\" + Addon_Id + "\\options.js");
}
