#pragma once

#define UNICODE
#include <windows.h>

#define WS_EX_NOREDIRECTIONBITMAP 0x00200000L

#define GET_X_LPARAM(lParam) (int)(short)LOWORD(lParam)
#define GET_Y_LPARAM(lParam) (int)(short)HIWORD(lParam)

#define DWMWA_USE_IMMERSIVE_DARK_MODE_BEFORE_20H1 19
#define DWMWA_USE_IMMERSIVE_DARK_MODE 20

UINT GetPrimaryDesktopDpi(void);

BOOL AdjustWindowRectExForDpi(RECT *lpRect, DWORD dwStyle, BOOL bMenu, DWORD dwExStyle, UINT dpi);

wchar_t *GetString(UINT id);

HBITMAP LoadPNGFromResource(wchar_t *type, wchar_t *name);

void GetAppVersion(UINT *version);
