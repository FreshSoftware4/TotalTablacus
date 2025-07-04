SetGestureExec("List", "1", function (Ctrl, pt) {
	if (Ctrl.hwndList && Ctrl.FolderFlags & FWF_SINGLECLICKACTIVATE) {
		const ptc = pt.Clone();
		api.ScreenToClient(Ctrl.hwndList, ptc);
		const ht = api.Memory("LVHITTESTINFO");
		ht.pt = ptc;
		api.SendMessage(Ctrl.hwndList, LVM_HITTEST, 0, ht);
		if (ht.flags & (LVHT_ONITEMICON | LVHT_ONITEMLABEL)) {
			const b = api.GetSystemMetrics(SM_SWAPBUTTON);
			api.mouse_event(b ? MOUSEEVENTF_RIGHTDOWN : MOUSEEVENTF_LEFTDOWN, 0, 0, 0, 0);
			api.mouse_event(b ? MOUSEEVENTF_RIGHTUP : MOUSEEVENTF_LEFTUP, 0, 0, 0, 0);
			return S_OK;
		}
	}
}, "Func", true);
