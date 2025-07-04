const Addon_Id = "recentlyclosedtabs";
const Default = "ToolBar2Left";
let item = GetAddonElement(Addon_Id);
if (!item.getAttribute("Set")) {
	item.setAttribute("MenuExec", 1);
	item.setAttribute("Menu", "Tabs");
	item.setAttribute("MenuPos", -1);
}
if (window.Addon == 1) {
	Addons.RecentlyClosedTabs = {
		Exec: async function (Ctrl, pt) {
			const FV = await GetFolderView(Ctrl, pt);
			FV.Focus();
			if (!pt) {
				pt = await api.Memory("POINT");
				await api.GetCursorPos(pt);
			}
			Common.RecentlyClosedTabs.nCommand = 1;
			const hMenu = await Sync.RecentlyClosedTabs.CreateMenu();
			await Sync.RecentlyClosedTabs.MenuCommand(Ctrl, pt, "", await api.TrackPopupMenuEx(hMenu, TPM_RIGHTBUTTON | TPM_RETURNCMD, await pt.x, await pt.y, ui_.hwnd, null, null), hMenu);
			api.DestroyMenu(hMenu);
			return S_OK;
		}
	};
	Common.RecentlyClosedTabs = await api.CreateObject("Object");
	Common.RecentlyClosedTabs.nCommand = 1;
	//Menu
	if (item.getAttribute("MenuExec")) {
		Common.RecentlyClosedTabs.strMenu = item.getAttribute("Menu");
		Common.RecentlyClosedTabs.strName = item.getAttribute("MenuName") || await GetAddonInfo(Addon_Id).Name;
		Common.RecentlyClosedTabs.nPos = GetNum(item.getAttribute("MenuPos"));
	}
	//Key
	if (item.getAttribute("KeyExec")) {
		SetKeyExec(item.getAttribute("KeyOn"), item.getAttribute("Key"), Addons.RecentlyClosedTabs.Exec, "Async");
	}
	//Mouse
	if (item.getAttribute("MouseExec")) {
		SetGestureExec(item.getAttribute("MouseOn"), item.getAttribute("Mouse"), Addons.RecentlyClosedTabs.Exec, "Async");
	}

	AddEvent("Layout", async function () {
		SetAddon(Addon_Id, Default, ['<span class="button" onclick="SyncExec(Addons.RecentlyClosedTabs.Exec, this, 9)" onmouseover="MouseOver(this)" onmouseout="MouseOut()">', await GetImgTag({
			title: item.getAttribute("MenuName") || await GetAddonInfo(Addon_Id).Name,
			src: item.getAttribute("Icon") || "icon:dsuiext.dll,0"
		}, GetIconSizeEx(item)), '</span>']);
		delete item;
	});

	AddTypeEx("Add-ons", "Recently closed tabs", Addons.RecentlyClosedTabs.Exec);

	$.importScript("addons\\" + Addon_Id + "\\sync.js");
} else {
	EnableInner();
}
