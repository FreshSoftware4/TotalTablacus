const Addon_Id = "retouch";
const Default = "ToolBar2Left";
const item = await GetAddonElement(Addon_Id);
if (!item.getAttribute("Set")) {
	item.setAttribute("MenuExec", 1);
	item.setAttribute("Menu", "Edit");
	item.setAttribute("MenuPos", -1);

	item.setAttribute("KeyOn", "List");

	item.setAttribute("MouseOn", "List");
}
if (window.Addon == 1) {
	AddEvent("Layout", async function () {
		SetAddon(Addon_Id, Default, ['<span class="button" id="RetouchButton" onclick="SyncExec(Sync.Retouch.Exec, this)" onmouseover="MouseOver(this)" onmouseout="MouseOut()">', await GetImgTag({
			title: item.getAttribute("MenuName") || await GetAddonInfo(Addon_Id).Name,
			src: item.getAttribute("Icon") || "icon:shell32.dll,141"
		}, GetIconSizeEx(item)), '</span>']);
	});

	$.importScript("addons\\" + Addon_Id + "\\sync.js");
} else {
	EnableInner();
}
