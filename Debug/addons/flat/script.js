const Addon_Id = "flat";
const Default = "None";
const item = await GetAddonElement(Addon_Id);
if (!item.getAttribute("Set")) {
	item.setAttribute("MenuExec", -1);
	item.setAttribute("Menu", "Tabs");
	item.setAttribute("MenuPos", -1);
}
if (window.Addon == 1) {
	AddEvent("Layout", async function () {
		await SetAddon(Addon_Id, Default, ['<span class="button" onclick="SyncExec(Sync.Flat.Exec, this);" onmouseover="MouseOver(this)" onmouseout="MouseOut()">', await GetImgTag({
			title: item.getAttribute("MenuName") || await GetText("Flat"),
			src: item.getAttribute("Icon") || ("string" === ui_.MiscIcon[Addon_Id] ? ui_.MiscIcon[Addon_Id] : await GetMiscIcon(Addon_Id)) || "icon:browser,15"
		}, GetIconSizeEx(item)), '</span>']);
	});
	$.importScript("addons\\" + Addon_Id + "\\sync.js");
} else {
	SetTabContents(0, "General", '<label><input type="checkbox" name="Arc">@sdcpl.dll,-132[Compressed Files]</label>');
	EnableInner();
}
