#include "about.h"
#define UNICODE
#include <dwmapi.h>
#include <windows.h>

#include "../res/resource.h"
#include "utils.h"

#define ABOUT_WINDOW_STYLE (WS_OVERLAPPEDWINDOW ^ WS_THICKFRAME ^ WS_MAXIMIZEBOX)

typedef struct AboutWindowData {
    UINT dpi;
    HBITMAP icon_image;
} AboutWindowData;

static LRESULT WINAPI AboutWndProc(HWND hwnd, UINT msg, WPARAM wParam, LPARAM lParam) {
    AboutWindowData *window = (AboutWindowData *)GetWindowLongPtr(hwnd, GWLP_USERDATA);

    // Load icon image PNG
    if (msg == WM_CREATE) {
        window = ((CREATESTRUCTW *)lParam)->lpCreateParams;
        SetWindowLongPtr(hwnd, GWLP_USERDATA, (long long)window);

        window->icon_image = LoadPNGFromResource(L"IMAGE", (wchar_t *)ID_IMAGE_ICON);
        return 0;
    }

    // Handle dpi changes
    if (msg == WM_DPICHANGED) {
        window->dpi = HIWORD(wParam);
        RECT *window_rect = (RECT *)lParam;
        SetWindowPos(hwnd, NULL, window_rect->left, window_rect->top, window_rect->right - window_rect->left, window_rect->bottom - window_rect->top,
                     SWP_NOZORDER | SWP_NOACTIVATE);
        return 0;
    }

    // Paint something nice
    if (msg == WM_ERASEBKGND) {
        return TRUE;
    }
    if (msg == WM_PAINT) {
        PAINTSTRUCT ps;
        HDC hdc = BeginPaint(hwnd, &ps);

        RECT client_rect;
        GetClientRect(hwnd, &client_rect);

        // Create back buffer
        HDC hdc_buffer = CreateCompatibleDC(hdc);
        SetBkMode(hdc_buffer, TRANSPARENT);
        HBITMAP bitmap_buffer = CreateCompatibleBitmap(hdc, client_rect.right, client_rect.bottom);
        SelectObject(hdc_buffer, bitmap_buffer);

        // Draw background color
        HBRUSH brush = CreateSolidBrush(0x0a0a0a);
        RECT rect = {0, 0, client_rect.right, client_rect.bottom};
        FillRect(hdc_buffer, &rect, brush);
        DeleteObject(brush);

        // Draw about icon image
        HDC hdc_image = CreateCompatibleDC(hdc_buffer);
        SelectObject(hdc_image, window->icon_image);
        SetStretchBltMode(hdc_buffer, STRETCH_HALFTONE);
        StretchBlt(hdc_buffer, MulDiv(16, window->dpi, 96), MulDiv(16 + 16, window->dpi, 96), MulDiv(128, window->dpi, 96), MulDiv(128, window->dpi, 96),
                   hdc_image, 0, 0, 256, 256, SRCCOPY);
        DeleteDC(hdc_image);

        // Draw about title
        HFONT title_font = CreateFont(MulDiv(32, window->dpi, 96), 0, 0, 0, FW_NORMAL, FALSE, FALSE, FALSE, ANSI_CHARSET, OUT_DEFAULT_PRECIS,
                                     CLIP_DEFAULT_PRECIS, CLEARTYPE_QUALITY, DEFAULT_PITCH | FF_DONTCARE, L"Segoe UI");
        SelectObject(hdc_buffer, title_font);
        SetTextColor(hdc_buffer, 0xffffff);
        TextOut(hdc_buffer, MulDiv(16 + 128 + 24, window->dpi, 96), MulDiv(32, window->dpi, 96), GetString(ID_STRING_ABOUT_TITLE),
                wcslen(GetString(ID_STRING_ABOUT_TITLE)));
        DeleteObject(title_font);

        // Draw about text
        UINT app_version[4];
        GetAppVersion(app_version);
        wchar_t about_text[512];
        wsprintf(about_text, GetString(ID_STRING_ABOUT_TEXT_FORMAT), app_version[0], app_version[1], app_version[2], app_version[3]);

        HFONT text_font = CreateFont(MulDiv(20, window->dpi, 96), 0, 0, 0, FW_NORMAL, FALSE, FALSE, FALSE, ANSI_CHARSET, OUT_DEFAULT_PRECIS, CLIP_DEFAULT_PRECIS,
                                    CLEARTYPE_QUALITY, DEFAULT_PITCH | FF_DONTCARE, L"Segoe UI");
        SelectObject(hdc_buffer, text_font);
        int y = 32 + 32 + 16;
        wchar_t *c = about_text;
        for (;;) {
            wchar_t *lineStart = c;
            while (*c != '\n' && *c != '\0') c++;
            TextOut(hdc_buffer, MulDiv(16 + 128 + 24, window->dpi, 96), MulDiv(y, window->dpi, 96), lineStart, c - lineStart);
            if (*c == '\0') break;
            c++;
            y += 20 + 8;
        }
        DeleteObject(text_font);

        // Draw and delete back buffer
        BitBlt(hdc, 0, 0, client_rect.right, client_rect.bottom, hdc_buffer, 0, 0, SRCCOPY);
        DeleteObject(bitmap_buffer);
        DeleteDC(hdc_buffer);

        EndPaint(hwnd, &ps);
        return 0;
    }

    // Clean up
    if (msg == WM_DESTROY) {
        DeleteObject(window->icon_image);
        free(window);
        return 0;
    }

    return DefWindowProc(hwnd, msg, wParam, lParam);
}

void OpenAboutWindow(void) {
    // Register window class
    WNDCLASSEX wc = {0};
    wc.cbSize = sizeof(WNDCLASSEX);
    wc.style = CS_HREDRAW | CS_VREDRAW;
    wc.lpfnWndProc = AboutWndProc;
    wc.hInstance = GetModuleHandle(NULL);
    wc.hIcon = (HICON)LoadImage(wc.hInstance, MAKEINTRESOURCE(ID_ICON_APP), IMAGE_ICON, 0, 0, LR_DEFAULTSIZE | LR_DEFAULTCOLOR | LR_SHARED);
    wc.hCursor = LoadCursor(NULL, IDC_ARROW);
    wc.lpszClassName = L"bassiemusic-about";
    wc.hIconSm = (HICON)LoadImage(wc.hInstance, MAKEINTRESOURCE(ID_ICON_APP), IMAGE_ICON, GetSystemMetrics(SM_CXSMICON), GetSystemMetrics(SM_CYSMICON),
                                  LR_DEFAULTCOLOR | LR_SHARED);
    RegisterClassEx(&wc);

    // Create centered window
    AboutWindowData *window = malloc(sizeof(AboutWindowData));
    window->dpi = GetPrimaryDesktopDpi();
    UINT window_width = MulDiv(540, window->dpi, 96);
    UINT window_height = MulDiv(196, window->dpi, 96);
    RECT window_rect;
    window_rect.left = (GetSystemMetrics(SM_CXSCREEN) - window_width) / 2;
    window_rect.top = (GetSystemMetrics(SM_CYSCREEN) - window_height) / 2;
    window_rect.right = window_rect.left + window_width;
    window_rect.bottom = window_rect.top + window_height;
    AdjustWindowRectExForDpi(&window_rect, ABOUT_WINDOW_STYLE, FALSE, 0, window->dpi);
    HWND hwnd = CreateWindowEx(0, wc.lpszClassName, GetString(ID_STRING_ABOUT_TITLE), ABOUT_WINDOW_STYLE, window_rect.left, window_rect.top,
                               window_rect.right - window_rect.left, window_rect.bottom - window_rect.top, HWND_DESKTOP, NULL, wc.hInstance, window);

    // Enable dark window decoration
    BOOL enabled = TRUE;
    if (FAILED(DwmSetWindowAttribute(hwnd, DWMWA_USE_IMMERSIVE_DARK_MODE, &enabled, sizeof(BOOL)))) {
        DwmSetWindowAttribute(hwnd, DWMWA_USE_IMMERSIVE_DARK_MODE_BEFORE_20H1, &enabled, sizeof(BOOL));
    }

    // Show window
    ShowWindow(hwnd, SW_SHOW);
    UpdateWindow(hwnd);
}
