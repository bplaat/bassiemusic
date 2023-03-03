#pragma once

#define UNICODE
#include <d2d1.h>
#include <windows.h>

#define WS_EX_NOREDIRECTIONBITMAP 0x00200000L

#define GET_X_LPARAM(lParam) (int)(short)LOWORD(lParam)
#define GET_Y_LPARAM(lParam) (int)(short)HIWORD(lParam)

#define DWMWA_USE_IMMERSIVE_DARK_MODE_BEFORE_20H1 19
#define DWMWA_USE_IMMERSIVE_DARK_MODE 20

int wprintf(const wchar_t *format, ...);

UINT GetPrimaryDesktopDpi(void);

BOOL AdjustWindowRectExForDpi(RECT *lpRect, DWORD dwStyle, BOOL bMenu, DWORD dwExStyle, UINT dpi);

wchar_t *GetString(UINT id);

HBITMAP LoadPNGFromResource(wchar_t *type, wchar_t *name);

void GetAppVersion(UINT *version);

typedef struct CanvasRect {
    float x;
    float y;
    float width;
    float height;
} CanvasRect;

typedef UINT CanvasColor;

void Direct2d_FillRect(ID2D1RenderTarget *render_target, CanvasRect *rect, CanvasColor color);

void Direct2d_FillPath(ID2D1Factory *d2d_factory, ID2D1RenderTarget *render_target, CanvasRect *rect, int viewport_width, int viewport_height, char *path,
                       CanvasColor color);
