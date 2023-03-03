#define UNICODE
#include <dwmapi.h>
#include <objbase.h>
#include <shlobj.h>
#include <windows.h>
#define COBJMACROS
#include <WinHttp.h>
#include <d2d1.h>
#include <d3d11.h>

#include "../res/resource.h"
#include "WebView2.h"
#include "about.h"
#include "dcomp.h"
#include "update.h"
#include "utils.h"

#define ID_MENU_ABOUT 2
#define WINDOW_STYLE (WS_POPUP | WS_THICKFRAME | WS_CAPTION | WS_SYSMENU | WS_MAXIMIZEBOX | WS_MINIMIZEBOX)
#define WM_SHOW_UPDATE_WINDOW (WM_USER + 1)
#define TITLEBAR_BUTTON_WIDTH 46
#define TITLEBAR_BUTTON_HEIGHT_DESKTOP 30
#define TITLEBAR_BUTTON_HEIGHT_MOBILE 52

HWND window_hwnd;
UINT window_dpi;
ID2D1Factory *d2d_factory = NULL;
IDCompositionDevice *composition_device = NULL;
IDCompositionVisual *webview_visual = NULL;
ICoreWebView2CompositionController *composition_controller = NULL;
ICoreWebView2Controller *controller = NULL;
ICoreWebView2 *webview = NULL;
IDCompositionVirtualSurface *titlebar_surface;
BOOL minimize_hover = FALSE;
BOOL maximize_hover = FALSE;
BOOL close_hover = FALSE;

typedef HRESULT(STDMETHODCALLTYPE *_CreateCoreWebView2EnvironmentWithOptions)(
    PCWSTR browserExecutableFolder, PCWSTR userDataFolder, ICoreWebView2EnvironmentOptions *environmentOptions,
    ICoreWebView2CreateCoreWebView2EnvironmentCompletedHandler *environmentCreatedHandler);

void FatalError(wchar_t *message) {
    MessageBox(HWND_DESKTOP, message, L"BassieMusic Error", MB_OK | MB_ICONSTOP);
    ExitProcess(1);
}

void ResizeBrowser(HWND hwnd) {
    if (controller == NULL) return;
    RECT window_rect;
    GetClientRect(hwnd, &window_rect);
    ICoreWebView2Controller_put_Bounds(controller, window_rect);
}

void WindowToggleMaximize(void) {
    WINDOWPLACEMENT placement;
    GetWindowPlacement(window_hwnd, &placement);
    if (placement.showCmd == SW_MAXIMIZE) {
        ShowWindow(window_hwnd, SW_RESTORE);
    } else {
        ShowWindow(window_hwnd, SW_MAXIMIZE);
    }
}

void WindowClose(void) {
    wchar_t windowStatePath[MAX_PATH];
    SHGetFolderPath(NULL, CSIDL_LOCAL_APPDATA, NULL, 0, windowStatePath);
    wcscat(windowStatePath, L"\\BassieMusic\\window");
    HANDLE windowStateFile = CreateFile(windowStatePath, GENERIC_WRITE, 0, NULL, CREATE_ALWAYS, FILE_ATTRIBUTE_NORMAL, NULL);
    WINDOWPLACEMENT windowState;
    GetWindowPlacement(window_hwnd, &windowState);
    WriteFile(windowStateFile, &windowState, sizeof(WINDOWPLACEMENT), NULL, NULL);
    CloseHandle(windowStateFile);
    DestroyWindow(window_hwnd);
}

BOOL HandleSystemKeyDown(UINT key) {
    if (key == VK_F4) {
        WindowClose();
        return TRUE;
    }
    return FALSE;
}

void CheckForUpdates(void) {
    // Get current app version
    UINT app_version[4];
    GetAppVersion(app_version);
    wchar_t current_version[255];
    wsprintf(current_version, L"%d.%d.%d.%d", app_version[0], app_version[1], app_version[2], app_version[3]);

    // Fetch latest version from API
    HINTERNET hSession = WinHttpOpen(L"WinHTTP Example/1.0", WINHTTP_ACCESS_TYPE_DEFAULT_PROXY, WINHTTP_NO_PROXY_NAME, WINHTTP_NO_PROXY_BYPASS, 0);
    HINTERNET hConnect = WinHttpConnect(hSession, L"bassiemusic-api.plaatsoft.nl", INTERNET_DEFAULT_HTTPS_PORT, 0);
    HINTERNET hRequest =
        WinHttpOpenRequest(hConnect, L"GET", L"/apps/windows/version", NULL, WINHTTP_NO_REFERER, WINHTTP_DEFAULT_ACCEPT_TYPES, WINHTTP_FLAG_SECURE);
    if (!WinHttpSendRequest(hRequest, WINHTTP_NO_ADDITIONAL_HEADERS, 0, WINHTTP_NO_REQUEST_DATA, 0, 0, 0)) {
        WinHttpCloseHandle(hRequest);
        WinHttpCloseHandle(hConnect);
        WinHttpCloseHandle(hSession);
        return;
    }
    if (!WinHttpReceiveResponse(hRequest, NULL)) {
        WinHttpCloseHandle(hRequest);
        WinHttpCloseHandle(hConnect);
        WinHttpCloseHandle(hSession);
        return;
    }

    // Read http response and covert to UTF-16
    DWORD latest_version_size;
    WinHttpQueryDataAvailable(hRequest, &latest_version_size);
    char latest_version[255];
    WinHttpReadData(hRequest, (void *)latest_version, latest_version_size, NULL);
    wchar_t latest_version_wide[255];
    MultiByteToWideChar(CP_ACP, 0, latest_version, -1, latest_version_wide, latest_version_size);

    // Compare version when diffrent send show update window message
    if (wcscmp(current_version, latest_version_wide)) {
        SendMessage(window_hwnd, WM_SHOW_UPDATE_WINDOW, 0, 0);
    }

    WinHttpCloseHandle(hRequest);
    WinHttpCloseHandle(hConnect);
    WinHttpCloseHandle(hSession);
}

// Default IUnknown method wrappers
HRESULT STDMETHODCALLTYPE Unknown_QueryInterface(IUnknown *This, REFIID riid, void **ppvObject) { return E_NOINTERFACE; }

ULONG STDMETHODCALLTYPE Unknown_AddRef(IUnknown *This) { return E_NOTIMPL; }

ULONG STDMETHODCALLTYPE Unknown_Release(IUnknown *This) { return E_NOTIMPL; }

// Forward interface reference
ICoreWebView2CreateCoreWebView2EnvironmentCompletedHandlerVtbl EnvironmentCompletedHandlerVtbl;
ICoreWebView2CreateCoreWebView2CompositionControllerCompletedHandlerVtbl ControllerCompletedHandlerVtbl;
ICoreWebView2NewWindowRequestedEventHandlerVtbl NewWindowRequestedHandlerVtbl;
ICoreWebView2AcceleratorKeyPressedEventHandlerVtbl AcceleratorKeyPressedHandlerVtbl;

// ICoreWebView2CreateCoreWebView2EnvironmentCompletedHandler
HRESULT STDMETHODCALLTYPE EnvironmentCompletedHandler_Invoke(ICoreWebView2CreateCoreWebView2EnvironmentCompletedHandler *This, HRESULT result,
                                                             ICoreWebView2Environment *environment) {
    if (FAILED(result)) {
        FatalError(L"Failed to create ICoreWebView2Environment");
    }

    ICoreWebView2Environment3 *environment3;
    ICoreWebView2_QueryInterface(environment, &IID_ICoreWebView2Environment3, (void **)&environment3);

    ICoreWebView2CreateCoreWebView2CompositionControllerCompletedHandler *controllerCompletedHandler =
        malloc(sizeof(ICoreWebView2CreateCoreWebView2CompositionControllerCompletedHandler));
    controllerCompletedHandler->lpVtbl = &ControllerCompletedHandlerVtbl;
    ICoreWebView2Environment3_CreateCoreWebView2CompositionController(environment3, window_hwnd, controllerCompletedHandler);

    ICoreWebView2Environment3_Release(environment3);
    return S_OK;
}

ICoreWebView2CreateCoreWebView2EnvironmentCompletedHandlerVtbl EnvironmentCompletedHandlerVtbl = {
    (HRESULT(STDMETHODCALLTYPE *)(ICoreWebView2CreateCoreWebView2EnvironmentCompletedHandler * This, REFIID riid, void **ppvObject)) Unknown_QueryInterface,
    (ULONG(STDMETHODCALLTYPE *)(ICoreWebView2CreateCoreWebView2EnvironmentCompletedHandler * This)) Unknown_AddRef,
    (ULONG(STDMETHODCALLTYPE *)(ICoreWebView2CreateCoreWebView2EnvironmentCompletedHandler * This)) Unknown_Release, EnvironmentCompletedHandler_Invoke};

// ICoreWebView2CreateCoreWebView2CompositionControllerCompletedHandler
HRESULT STDMETHODCALLTYPE ControllerCompletedHandler_Invoke(ICoreWebView2CreateCoreWebView2CompositionControllerCompletedHandler *This, HRESULT result,
                                                            ICoreWebView2CompositionController *new_composition_controller) {
    if (FAILED(result)) {
        FatalError(L"Failed to create ICoreWebView2CompositionController");
    }
    ICoreWebView2CompositionController_AddRef(new_composition_controller);
    composition_controller = new_composition_controller;

    ICoreWebView2CompositionController_put_RootVisualTarget(composition_controller, (void *)webview_visual);
    composition_device->lpVtbl->Commit(composition_device);

    ICoreWebView2CompositionController_QueryInterface(composition_controller, &IID_ICoreWebView2Controller, (void **)&controller);
    ICoreWebView2Controller_get_CoreWebView2(controller, &webview);

    ICoreWebView2Controller2 *controller2;
    ICoreWebView2Controller_QueryInterface(controller, &IID_ICoreWebView2Controller2, (void **)&controller2);
    ICoreWebView2Controller2_put_DefaultBackgroundColor(controller2, ((COREWEBVIEW2_COLOR){0xff, 0x0a, 0x0a, 0x0a}));

    ICoreWebView2_13 *webview13;
    ICoreWebView2_QueryInterface(webview, &IID_ICoreWebView2_13, (void **)&webview13);
    ICoreWebView2Profile *profile;
    ICoreWebView2_13_get_Profile(webview13, &profile);
    ICoreWebView2_13_Release(webview13);
    ICoreWebView2Profile_put_PreferredColorScheme(profile, COREWEBVIEW2_PREFERRED_COLOR_SCHEME_DARK);
    ICoreWebView2Profile_Release(profile);

    ICoreWebView2Settings *settings;
    ICoreWebView2_get_Settings(webview, &settings);
    ICoreWebView2Settings_put_AreDevToolsEnabled(settings, FALSE);
    ICoreWebView2Settings_put_AreDefaultContextMenusEnabled(settings, FALSE);
    ICoreWebView2Settings_put_IsStatusBarEnabled(settings, FALSE);
    ICoreWebView2Settings_Release(settings);

    ICoreWebView2Settings2 *settings2;
    ICoreWebView2_QueryInterface(settings, &IID_ICoreWebView2Settings2, (void **)&settings2);
    UINT app_version[4];
    GetAppVersion(app_version);
    wchar_t userAgent[255];
    wsprintf(userAgent, L"BassieMusic Windows App v%d.%d.%d.%d", app_version[0], app_version[1], app_version[2], app_version[3]);
    ICoreWebView2Settings2_put_UserAgent(settings2, userAgent);
    ICoreWebView2Settings2_Release(settings2);

    ICoreWebView2NewWindowRequestedEventHandler *newWindowRequestedHandler = malloc(sizeof(ICoreWebView2NewWindowRequestedEventHandler));
    newWindowRequestedHandler->lpVtbl = &NewWindowRequestedHandlerVtbl;
    ICoreWebView2_add_NewWindowRequested(webview, newWindowRequestedHandler, NULL);

    ICoreWebView2AcceleratorKeyPressedEventHandler *newAcceleratorKeyPressedHandler = malloc(sizeof(ICoreWebView2AcceleratorKeyPressedEventHandler));
    newAcceleratorKeyPressedHandler->lpVtbl = &AcceleratorKeyPressedHandlerVtbl;
    ICoreWebView2Controller_add_AcceleratorKeyPressed(controller, newAcceleratorKeyPressedHandler, NULL);

    ResizeBrowser(window_hwnd);
    ICoreWebView2Controller_put_IsVisible(controller, TRUE);
    ICoreWebView2_Navigate(webview, GetString(ID_STRING_WEBVIEW_URL));
    return S_OK;
}

ICoreWebView2CreateCoreWebView2CompositionControllerCompletedHandlerVtbl ControllerCompletedHandlerVtbl = {
    (HRESULT(STDMETHODCALLTYPE *)(ICoreWebView2CreateCoreWebView2CompositionControllerCompletedHandler * This, REFIID riid, void **ppvObject))
        Unknown_QueryInterface,
    (ULONG(STDMETHODCALLTYPE *)(ICoreWebView2CreateCoreWebView2CompositionControllerCompletedHandler * This)) Unknown_AddRef,
    (ULONG(STDMETHODCALLTYPE *)(ICoreWebView2CreateCoreWebView2CompositionControllerCompletedHandler * This)) Unknown_Release,
    ControllerCompletedHandler_Invoke,
};

// ICoreWebView2NewWindowRequestedEventHandler
HRESULT STDMETHODCALLTYPE NewWindowRequestedHandler_Invoke(ICoreWebView2NewWindowRequestedEventHandler *This, ICoreWebView2 *sender,
                                                           ICoreWebView2NewWindowRequestedEventArgs *args) {
    ICoreWebView2NewWindowRequestedEventArgs_put_Handled(args, TRUE);
    wchar_t *url;
    ICoreWebView2NewWindowRequestedEventArgs_get_Uri(args, &url);
    ShellExecute(window_hwnd, L"OPEN", url, NULL, NULL, SW_SHOWNORMAL);
    return S_OK;
}

ICoreWebView2NewWindowRequestedEventHandlerVtbl NewWindowRequestedHandlerVtbl = {
    (HRESULT(STDMETHODCALLTYPE *)(ICoreWebView2NewWindowRequestedEventHandler * This, REFIID riid, void **ppvObject)) Unknown_QueryInterface,
    (ULONG(STDMETHODCALLTYPE *)(ICoreWebView2NewWindowRequestedEventHandler * This)) Unknown_AddRef,
    (ULONG(STDMETHODCALLTYPE *)(ICoreWebView2NewWindowRequestedEventHandler * This)) Unknown_Release, NewWindowRequestedHandler_Invoke};

// ICoreWebView2AcceleratorKeyPressedEventHandler
HRESULT STDMETHODCALLTYPE AcceleratorKeyPressedHandler_Invoke(ICoreWebView2AcceleratorKeyPressedEventHandler *This, ICoreWebView2Controller *sender,
                                                              ICoreWebView2AcceleratorKeyPressedEventArgs *args) {
    COREWEBVIEW2_KEY_EVENT_KIND state;
    ICoreWebView2AcceleratorKeyPressedEventArgs_get_KeyEventKind(args, &state);
    UINT key;
    ICoreWebView2AcceleratorKeyPressedEventArgs_get_VirtualKey(args, &key);
    if (state == COREWEBVIEW2_KEY_EVENT_KIND_SYSTEM_KEY_DOWN) {
        if (HandleSystemKeyDown(key)) {
            ICoreWebView2AcceleratorKeyPressedEventArgs_put_Handled(args, TRUE);
        }
    }
    return S_OK;
}

ICoreWebView2AcceleratorKeyPressedEventHandlerVtbl AcceleratorKeyPressedHandlerVtbl = {
    (HRESULT(STDMETHODCALLTYPE *)(ICoreWebView2AcceleratorKeyPressedEventHandler * This, REFIID riid, void **ppvObject)) Unknown_QueryInterface,
    (ULONG(STDMETHODCALLTYPE *)(ICoreWebView2AcceleratorKeyPressedEventHandler * This)) Unknown_AddRef,
    (ULONG(STDMETHODCALLTYPE *)(ICoreWebView2AcceleratorKeyPressedEventHandler * This)) Unknown_Release, AcceleratorKeyPressedHandler_Invoke};

// Window code
LRESULT WINAPI WndProc(HWND hwnd, UINT msg, WPARAM wParam, LPARAM lParam) {
    // When window is created
    if (msg == WM_CREATE) {
        HMENU sysMenu = GetSystemMenu(hwnd, FALSE);
        InsertMenu(sysMenu, 5, MF_BYPOSITION | MF_SEPARATOR, 0, NULL);
        InsertMenu(sysMenu, 6, MF_BYPOSITION, ID_MENU_ABOUT, GetString(ID_STRING_ABOUT_MENU));

        // Create dxgi device
        ID3D11Device *d3d11Device;
        D3D11CreateDevice(NULL, D3D_DRIVER_TYPE_HARDWARE, NULL, D3D11_CREATE_DEVICE_BGRA_SUPPORT, NULL, 0, D3D11_SDK_VERSION, &d3d11Device, NULL, NULL);

        IDXGIDevice *dxgiDevice;
        IID IID_IDXGIDevice = {0x54ec77fa, 0x1377, 0x44e6, {0x8c, 0x32, 0x88, 0xfd, 0x5f, 0x44, 0xc8, 0x4c}};
        ID3D11Device_QueryInterface(d3d11Device, &IID_IDXGIDevice, (void **)&dxgiDevice);

        // Create Direct2D Factory
        IID IID_ID2D1Factory = {0xbb12d362, 0xdaee, 0x4b9a, {0xaa, 0x1d, 0x14, 0xba, 0x40, 0x1c, 0xfa, 0x1f}};
        D2D1CreateFactory(D2D1_FACTORY_TYPE_SINGLE_THREADED, &IID_ID2D1Factory, NULL, (void **)&d2d_factory);

        // Create composition device
        HMODULE hDcomp = LoadLibrary(L"dcomp.dll");
        _DCompositionCreateDevice DCompositionCreateDevice = (_DCompositionCreateDevice)GetProcAddress(hDcomp, "DCompositionCreateDevice");
        IID IID_IDCompositionDevice = {0xC37EA93A, 0xE7AA, 0x450D, {0xB1, 0x6F, 0x97, 0x46, 0xCB, 0x04, 0x07, 0xF3}};
        DCompositionCreateDevice(dxgiDevice, &IID_IDCompositionDevice, (void **)&composition_device);

        IDCompositionTarget *hwndRenderTarget;
        composition_device->lpVtbl->CreateTargetForHwnd(composition_device, hwnd, TRUE, &hwndRenderTarget);

        // Create composition visuals
        IDCompositionVisual *root_visual;
        composition_device->lpVtbl->CreateVisual(composition_device, &root_visual);
        hwndRenderTarget->lpVtbl->SetRoot(hwndRenderTarget, root_visual);

        composition_device->lpVtbl->CreateVisual(composition_device, &webview_visual);
        root_visual->lpVtbl->AddVisual(root_visual, webview_visual, FALSE, NULL);

        // Create titlebar visual
        IDCompositionVisual *titlebar_visual;
        composition_device->lpVtbl->CreateVisual(composition_device, &titlebar_visual);
        root_visual->lpVtbl->AddVisual(root_visual, titlebar_visual, FALSE, NULL);

        RECT client_rect;
        GetClientRect(hwnd, &client_rect);
        composition_device->lpVtbl->CreateVirtualSurface(composition_device, client_rect.right, client_rect.bottom, DXGI_FORMAT_B8G8R8A8_UNORM,
                                                         DXGI_ALPHA_MODE_PREMULTIPLIED, &titlebar_surface);
        titlebar_visual->lpVtbl->SetContent(titlebar_visual, (IUnknown *)titlebar_surface);
        composition_device->lpVtbl->Commit(composition_device);

        // Check for updates
        CreateThread(NULL, 0, (LPTHREAD_START_ROUTINE)CheckForUpdates, NULL, 0, NULL);
        return 0;
    }

    // Menu commands
    if (msg == WM_SYSCOMMAND || msg == WM_COMMAND) {
        UINT id = LOWORD(wParam);
        if (id == ID_MENU_ABOUT) {
            OpenAboutWindow();
            return 0;
        }
        if (msg == WM_COMMAND && id >= 0xF000) {
            return DefWindowProc(hwnd, WM_SYSCOMMAND, wParam, lParam);
        }
    }

    // Handle dpi changes
    if (msg == WM_DPICHANGED) {
        window_dpi = HIWORD(wParam);
        RECT *window_rect = (RECT *)lParam;
        SetWindowPos(hwnd, NULL, window_rect->left, window_rect->top, window_rect->right - window_rect->left, window_rect->bottom - window_rect->top,
                     SWP_NOZORDER | SWP_NOACTIVATE);
        return 0;
    }

    // Make window have custom window decoration
    if (msg == WM_ACTIVATE) {
        MARGINS borderless = {1, 1, 1, 1};
        DwmExtendFrameIntoClientArea(hwnd, &borderless);
        return 0;
    }
    if (msg == WM_NCCALCSIZE) {
        if (wParam) {
            WINDOWPLACEMENT placement;
            GetWindowPlacement(hwnd, &placement);
            if (placement.showCmd == SW_MAXIMIZE) {
                NCCALCSIZE_PARAMS *params = (NCCALCSIZE_PARAMS *)lParam;
                HMONITOR monitor = MonitorFromWindow(hwnd, MONITOR_DEFAULTTONULL);
                if (!monitor) return 0;
                MONITORINFO monitor_info;
                monitor_info.cbSize = sizeof(MONITORINFO);
                if (!GetMonitorInfo(monitor, &monitor_info)) {
                    return 0;
                }
                params->rgrc[0] = monitor_info.rcWork;
            }
        }
        return 0;
    }
    if (msg == WM_NCHITTEST) {
        int x = GET_X_LPARAM(lParam);
        int y = GET_Y_LPARAM(lParam);

        RECT window_rect;
        GetWindowRect(hwnd, &window_rect);

        int border_horizontal = GetSystemMetrics(SM_CXSIZEFRAME);
        int border_vertical = GetSystemMetrics(SM_CYSIZEFRAME);

        if (y >= window_rect.top && y < window_rect.bottom) {
            if (x >= window_rect.left && x < window_rect.left + border_horizontal) {
                if (y < window_rect.top + border_vertical) {
                    return HTTOPLEFT;
                }
                if (y > window_rect.bottom - border_vertical) {
                    return HTBOTTOMLEFT;
                }
                return HTLEFT;
            }
            if (x >= window_rect.right - border_horizontal && x < window_rect.right) {
                if (y < window_rect.top + border_vertical) {
                    return HTTOPRIGHT;
                }
                if (y > window_rect.bottom - border_vertical) {
                    return HTBOTTOMRIGHT;
                }
                return HTRIGHT;
            }
        }

        if (x >= window_rect.left && x < window_rect.right) {
            if (y >= window_rect.top && y < window_rect.top + border_vertical) {
                return HTTOP;
            }
            if (y >= window_rect.bottom - border_vertical && y < window_rect.bottom) {
                return HTBOTTOM;
            }
        }

        if (x >= window_rect.left && y >= window_rect.top && x < window_rect.right && y < window_rect.bottom) {
            return HTCLIENT;
        }

        return HTNOWHERE;
    }

    // Resize browser
    if (msg == WM_SIZE) {
        // Resize browser
        ResizeBrowser(hwnd);

        // Resize titlebar surface
        titlebar_surface->lpVtbl->Resize(titlebar_surface, LOWORD(lParam), HIWORD(lParam));
        return 0;
    }

    // Set window min size
    if (msg == WM_GETMINMAXINFO) {
        RECT window_rect = {0, 0, MulDiv(480, window_dpi, 96), MulDiv(480, window_dpi, 96)};
        AdjustWindowRectExForDpi(&window_rect, WINDOW_STYLE, FALSE, 0, window_dpi);
        MINMAXINFO *minMaxInfo = (MINMAXINFO *)lParam;
        minMaxInfo->ptMinTrackSize.x = window_rect.right - window_rect.left;
        minMaxInfo->ptMinTrackSize.y = window_rect.bottom - window_rect.top;
        return 0;
    }

    // Handle keydown messages
    if (msg == WM_SYSKEYDOWN) {
        HandleSystemKeyDown(wParam);
        return 0;
    }

    // Window titlebar events
    if (msg == WM_MOUSEMOVE || msg == WM_LBUTTONDBLCLK || msg == WM_RBUTTONUP) {
        POINT point;
        POINTSTOPOINT(point, lParam);
        RECT client_rect;
        GetClientRect(hwnd, &client_rect);

        // Mobile
        if (client_rect.right < MulDiv(1024, window_dpi, 96)) {
            if (point.y < MulDiv(TITLEBAR_BUTTON_HEIGHT_MOBILE, window_dpi, 96) && point.x > MulDiv(TITLEBAR_BUTTON_HEIGHT_MOBILE, window_dpi, 96) &&
                point.x < client_rect.right - MulDiv(3 * TITLEBAR_BUTTON_WIDTH, window_dpi, 96)) {
                if ((GetKeyState(VK_LBUTTON) & 0x100) != 0 && msg == WM_MOUSEMOVE) {
                    ReleaseCapture();
                    SendMessage(hwnd, WM_NCLBUTTONDOWN, HTCAPTION, 0);
                }
                if (msg == WM_LBUTTONDBLCLK) {
                    WindowToggleMaximize();
                }
                if (msg == WM_RBUTTONUP) {
                    ClientToScreen(hwnd, &point);
                    TrackPopupMenu(GetSystemMenu(hwnd, FALSE), TPM_TOPALIGN | TPM_LEFTALIGN, point.x, point.y, 0, hwnd, NULL);
                }
            }
        }
        // Desktop
        else {
            if (point.y < MulDiv(TITLEBAR_BUTTON_HEIGHT_DESKTOP, window_dpi, 96)) {
                if (point.x < client_rect.right - MulDiv(3 * TITLEBAR_BUTTON_WIDTH, window_dpi, 96)) {
                    if ((GetKeyState(VK_LBUTTON) & 0x100) != 0 && msg == WM_MOUSEMOVE) {
                        ReleaseCapture();
                        SendMessage(hwnd, WM_NCLBUTTONDOWN, HTCAPTION, 0);
                    }
                    if (msg == WM_LBUTTONDBLCLK) {
                        WindowToggleMaximize();
                    }
                    if (msg == WM_RBUTTONUP) {
                        ClientToScreen(hwnd, &point);
                        TrackPopupMenu(GetSystemMenu(hwnd, FALSE), TPM_TOPALIGN | TPM_LEFTALIGN, point.x, point.y, 0, hwnd, NULL);
                    }
                }
            }
        }

        // Titlebar buttons hover
        if (msg == WM_MOUSEMOVE) {
            int titlebar_height = client_rect.right < MulDiv(1024, window_dpi, 96) ? MulDiv(TITLEBAR_BUTTON_HEIGHT_MOBILE, window_dpi, 96)
                                                                                   : MulDiv(TITLEBAR_BUTTON_HEIGHT_DESKTOP, window_dpi, 96);
            BOOL new_minimize_hover = point.y < titlebar_height && point.x >= client_rect.right - MulDiv(3 * TITLEBAR_BUTTON_WIDTH, window_dpi, 96) &&
                                      point.x < client_rect.right - MulDiv(2 * TITLEBAR_BUTTON_WIDTH, window_dpi, 96);
            BOOL new_maximize_hover = point.y < titlebar_height && point.x >= client_rect.right - MulDiv(2 * TITLEBAR_BUTTON_WIDTH, window_dpi, 96) &&
                                      point.x < client_rect.right - MulDiv(TITLEBAR_BUTTON_WIDTH, window_dpi, 96);
            BOOL new_close_hover = point.y < titlebar_height && point.x >= client_rect.right - MulDiv(TITLEBAR_BUTTON_WIDTH, window_dpi, 96);
            if (new_minimize_hover != minimize_hover || new_maximize_hover != maximize_hover || new_close_hover != close_hover) {
                minimize_hover = new_minimize_hover;
                maximize_hover = new_maximize_hover;
                close_hover = new_close_hover;
                InvalidateRect(hwnd, NULL, TRUE);
            }
        }
    }

    // Window titlebar button actions
    if (msg == WM_LBUTTONUP) {
        int x = GET_X_LPARAM(lParam);
        int y = GET_Y_LPARAM(lParam);
        RECT client_rect;
        GetClientRect(hwnd, &client_rect);

        int titlebar_height = client_rect.right < MulDiv(1024, window_dpi, 96) ? MulDiv(TITLEBAR_BUTTON_HEIGHT_MOBILE, window_dpi, 96)
                                                                               : MulDiv(TITLEBAR_BUTTON_HEIGHT_DESKTOP, window_dpi, 96);
        if (y < titlebar_height) {
            if (x >= client_rect.right - MulDiv(3 * TITLEBAR_BUTTON_WIDTH, window_dpi, 96) &&
                x < client_rect.right - MulDiv(2 * TITLEBAR_BUTTON_WIDTH, window_dpi, 96)) {
                minimize_hover = FALSE;
                ShowWindow(hwnd, SW_MINIMIZE);
            }
            if (x >= client_rect.right - MulDiv(2 * TITLEBAR_BUTTON_WIDTH, window_dpi, 96) &&
                x < client_rect.right - MulDiv(TITLEBAR_BUTTON_WIDTH, window_dpi, 96)) {
                maximize_hover = FALSE;
                WindowToggleMaximize();
            }
            if (x >= client_rect.right - MulDiv(TITLEBAR_BUTTON_WIDTH, window_dpi, 96)) {
                WindowClose();
            }
        }
    }

    // Send mouse events to webview2
    if (
        // clang-format off
        msg == WM_LBUTTONDOWN || msg == WM_MBUTTONDOWN || msg == WM_RBUTTONDOWN || msg == WM_XBUTTONDOWN ||
        msg == WM_LBUTTONUP || msg == WM_MBUTTONUP || msg == WM_RBUTTONUP || msg == WM_XBUTTONUP ||
        msg == WM_LBUTTONDBLCLK || msg == WM_MBUTTONDBLCLK || msg == WM_RBUTTONDBLCLK || msg == WM_XBUTTONDBLCLK ||
        msg == WM_MOUSEWHEEL || msg == WM_MOUSEHWHEEL || msg == WM_MOUSEMOVE
        // clang-format on
    ) {
        if (composition_controller != NULL) {
            POINT point;
            POINTSTOPOINT(point, lParam);
            int mouseData = 0;
            if (msg == WM_MOUSEWHEEL || msg == WM_MOUSEHWHEEL) {
                mouseData = GET_WHEEL_DELTA_WPARAM(wParam);
                ScreenToClient(hwnd, &point);
            } else {
                mouseData = GET_XBUTTON_WPARAM(wParam);
            }

            ICoreWebView2CompositionController_SendMouseInput(composition_controller, msg, GET_KEYSTATE_WPARAM(wParam), mouseData, point);
        }
        return 0;
    }

    // Paint window titlebar
    if (msg == WM_PAINT) {
        PAINTSTRUCT ps;
        BeginPaint(hwnd, &ps);
        RECT client_rect;
        GetClientRect(hwnd, &client_rect);

        POINT offset;
        IDXGISurface *dxgiSurface;
        IID IID_IDXGISurface = {0xcafcb56c, 0x6ac3, 0x4889, {0xbf, 0x47, 0x9e, 0x23, 0xbb, 0xd2, 0x60, 0xec}};
        titlebar_surface->lpVtbl->BeginDraw(titlebar_surface, &client_rect, &IID_IDXGISurface, (void **)&dxgiSurface, &offset);

        D2D1_RENDER_TARGET_PROPERTIES props = {
            D2D1_RENDER_TARGET_TYPE_DEFAULT, {DXGI_FORMAT_B8G8R8A8_UNORM, DXGI_ALPHA_MODE_PREMULTIPLIED}, 96, 96, D2D1_RENDER_TARGET_USAGE_NONE,
            D2D1_FEATURE_LEVEL_DEFAULT};
        ID2D1RenderTarget *render_target;
        ID2D1Factory_CreateDxgiSurfaceRenderTarget(d2d_factory, dxgiSurface, &props, (ID2D1RenderTarget **)&render_target);

        ID2D1RenderTarget_BeginDraw(render_target);
        ID2D1RenderTarget_Clear(render_target, (&(D2D1_COLOR_F){0, 0, 0, 0}));

        // Titlebar
        int titlebar_height = client_rect.right < MulDiv(1024, window_dpi, 96) ? TITLEBAR_BUTTON_HEIGHT_MOBILE : TITLEBAR_BUTTON_HEIGHT_DESKTOP;
        Direct2d_FillRect(render_target,
                          &(CanvasRect){client_rect.right - MulDiv(3 * TITLEBAR_BUTTON_WIDTH, window_dpi, 96), 0,
                                        MulDiv(3 * TITLEBAR_BUTTON_WIDTH, window_dpi, 96), MulDiv(titlebar_height, window_dpi, 96)},
                          client_rect.right < MulDiv(1024, window_dpi, 96) ? 0xff121212 : 0xff0a0a0a);

        // Minimize button
        if (minimize_hover) {
            Direct2d_FillRect(render_target,
                              &(CanvasRect){client_rect.right - MulDiv(3 * TITLEBAR_BUTTON_WIDTH, window_dpi, 96), 0,
                                            MulDiv(TITLEBAR_BUTTON_WIDTH, window_dpi, 96), MulDiv(titlebar_height, window_dpi, 96)},
                              0x20ffffff);
        }
        Direct2d_FillPath(
            d2d_factory, render_target,
            &(CanvasRect){client_rect.right - MulDiv(3 * TITLEBAR_BUTTON_WIDTH, window_dpi, 96) + MulDiv((TITLEBAR_BUTTON_WIDTH - 10) / 2, window_dpi, 96),
                          MulDiv((titlebar_height - 10) / 2, window_dpi, 96), MulDiv(10, window_dpi, 96), MulDiv(10, window_dpi, 96)},
            2048, 2048, "M2048 1229v-205h-2048v205h2048z", 0xffffffff);

        // Maximize button
        if (maximize_hover) {
            Direct2d_FillRect(render_target,
                              &(CanvasRect){client_rect.right - MulDiv(2 * TITLEBAR_BUTTON_WIDTH, window_dpi, 96), 0,
                                            MulDiv(TITLEBAR_BUTTON_WIDTH, window_dpi, 96), MulDiv(titlebar_height, window_dpi, 96)},
                              0x20ffffff);
        }
        WINDOWPLACEMENT placement;
        GetWindowPlacement(hwnd, &placement);
        if (placement.showCmd == SW_MAXIMIZE) {
            Direct2d_FillPath(
                d2d_factory, render_target,
                &(CanvasRect){client_rect.right - MulDiv(2 * TITLEBAR_BUTTON_WIDTH, window_dpi, 96) + MulDiv((TITLEBAR_BUTTON_WIDTH - 10) / 2, window_dpi, 96),
                              MulDiv((titlebar_height - 10) / 2 + 10, window_dpi, 96), MulDiv(10, window_dpi, 96), MulDiv(10, window_dpi, 96)},
                2048, -2048, "M2048 410h-410v-410h-1638v1638h410v410h1638v-1638zM1434 1434h-1229v-1229h1229v1229zM1843 1843h-1229v-205h1024v-1024h205v1229z",
                0xffffffff);
        } else {
            Direct2d_FillPath(
                d2d_factory, render_target,
                &(CanvasRect){client_rect.right - MulDiv(2 * TITLEBAR_BUTTON_WIDTH, window_dpi, 96) + MulDiv((TITLEBAR_BUTTON_WIDTH - 10) / 2, window_dpi, 96),
                              MulDiv((titlebar_height - 10) / 2, window_dpi, 96), MulDiv(10, window_dpi, 96), MulDiv(10, window_dpi, 96)},
                2048, 2048, "M2048 2048v-2048h-2048v2048h2048zM1843 1843h-1638v-1638h1638v1638z", 0xffffffff);
        }

        // Close button
        if (close_hover) {
            Direct2d_FillRect(render_target,
                              &(CanvasRect){client_rect.right - MulDiv(TITLEBAR_BUTTON_WIDTH, window_dpi, 96), 0,
                                            MulDiv(TITLEBAR_BUTTON_WIDTH, window_dpi, 96), MulDiv(titlebar_height, window_dpi, 96)},
                              0xa00000ff);
        }
        Direct2d_FillPath(
            d2d_factory, render_target,
            &(CanvasRect){client_rect.right - MulDiv(TITLEBAR_BUTTON_WIDTH, window_dpi, 96) + MulDiv((TITLEBAR_BUTTON_WIDTH - 10) / 2, window_dpi, 96),
                          MulDiv((titlebar_height - 10) / 2, window_dpi, 96), MulDiv(10, window_dpi, 96), MulDiv(10, window_dpi, 96)},
            2048, 2048, "M1169 1024l879 -879l-145 -145l-879 879l-879 -879l-145 145l879 879l-879 879l145 145l879 -879l879 879l145 -145z", 0xffffffff);

        ID2D1RenderTarget_EndDraw(render_target, NULL, NULL);
        ID2D1RenderTarget_Release(render_target);

        titlebar_surface->lpVtbl->EndDraw(titlebar_surface);
        composition_device->lpVtbl->Commit(composition_device);
        EndPaint(hwnd, &ps);
        return 0;
    }

    // Save window state
    if (msg == WM_CLOSE) {
        WindowClose();
        return 0;
    }

    // Quit application
    if (msg == WM_DESTROY) {
        PostQuitMessage(0);
        return 0;
    }

    // Open update window
    if (msg == WM_SHOW_UPDATE_WINDOW) {
        OpenUpdateWindow();
        return 0;
    }

    return DefWindowProc(hwnd, msg, wParam, lParam);
}

void _start(void) {
    // Only allow own app instance
    wchar_t *window_class_name = L"bassiemusic";
    HANDLE mutex = CreateMutex(NULL, TRUE, window_class_name);
    if (mutex == NULL || GetLastError() == ERROR_ALREADY_EXISTS) {
        HWND window = FindWindow(window_class_name, NULL);
        if (window != NULL) {
            if (IsIconic(window)) {
                ShowWindow(window, SW_SHOW);
            }
            SetForegroundWindow(window);
        }
        ExitProcess(0);
    }

    // Register window class
    WNDCLASSEX wc = {0};
    wc.cbSize = sizeof(WNDCLASSEX);
    wc.style = CS_HREDRAW | CS_VREDRAW | CS_DBLCLKS;
    wc.lpfnWndProc = WndProc;
    wc.hInstance = GetModuleHandle(NULL);
    wc.hIcon = (HICON)LoadImage(wc.hInstance, MAKEINTRESOURCE(ID_ICON_APP), IMAGE_ICON, 0, 0, LR_DEFAULTSIZE | LR_DEFAULTCOLOR | LR_SHARED);
    wc.hCursor = LoadCursor(NULL, IDC_ARROW);
    wc.lpszClassName = window_class_name;
    wc.hIconSm = (HICON)LoadImage(wc.hInstance, MAKEINTRESOURCE(ID_ICON_APP), IMAGE_ICON, GetSystemMetrics(SM_CXSMICON), GetSystemMetrics(SM_CYSMICON),
                                  LR_DEFAULTCOLOR | LR_SHARED);
    RegisterClassEx(&wc);

    // Create centered window
    window_dpi = GetPrimaryDesktopDpi();
    UINT window_width = MulDiv(1280, window_dpi, 96);
    UINT window_height = MulDiv(720, window_dpi, 96);
    RECT window_rect;
    window_rect.left = (GetSystemMetrics(SM_CXSCREEN) - window_width) / 2;
    window_rect.top = (GetSystemMetrics(SM_CYSCREEN) - window_height) / 2;
    window_rect.right = window_rect.left + window_width;
    window_rect.bottom = window_rect.top + window_height;
    AdjustWindowRectExForDpi(&window_rect, WINDOW_STYLE, FALSE, 0, window_dpi);
    window_hwnd = CreateWindowEx(WS_EX_NOREDIRECTIONBITMAP, wc.lpszClassName, GetString(ID_STRING_APP_NAME), WINDOW_STYLE, window_rect.left, window_rect.top,
                                 window_rect.right - window_rect.left, window_rect.bottom - window_rect.top, HWND_DESKTOP, NULL, wc.hInstance, NULL);

    // Enable dark window decoration
    BOOL enabled = TRUE;
    if (FAILED(DwmSetWindowAttribute(window_hwnd, DWMWA_USE_IMMERSIVE_DARK_MODE, &enabled, sizeof(BOOL)))) {
        DwmSetWindowAttribute(window_hwnd, DWMWA_USE_IMMERSIVE_DARK_MODE_BEFORE_20H1, &enabled, sizeof(BOOL));
    }

    // Restore old window state
    wchar_t windowStatePath[MAX_PATH];
    SHGetFolderPath(NULL, CSIDL_LOCAL_APPDATA, NULL, 0, windowStatePath);
    wcscat(windowStatePath, L"\\BassieMusic\\window");
    HANDLE windowStateFile = CreateFile(windowStatePath, GENERIC_READ, 0, NULL, OPEN_EXISTING, FILE_ATTRIBUTE_NORMAL, NULL);
    if (windowStateFile != INVALID_HANDLE_VALUE) {
        WINDOWPLACEMENT windowState;
        ReadFile(windowStateFile, &windowState, sizeof(WINDOWPLACEMENT), NULL, NULL);
        SetWindowPlacement(window_hwnd, &windowState);
        CloseHandle(windowStateFile);
    } else {
        ShowWindow(window_hwnd, window_width >= GetSystemMetrics(SM_CXSCREEN) ? SW_SHOWMAXIMIZED : SW_SHOWDEFAULT);
    }
    UpdateWindow(window_hwnd);

    // Load webview2 laoder
    HMODULE hWebview2Loader = LoadLibrary(L"WebView2Loader.dll");
    _CreateCoreWebView2EnvironmentWithOptions CreateCoreWebView2EnvironmentWithOptions =
        (_CreateCoreWebView2EnvironmentWithOptions)GetProcAddress(hWebview2Loader, "CreateCoreWebView2EnvironmentWithOptions");
    if (CreateCoreWebView2EnvironmentWithOptions != NULL) {
        // Find app data path
        wchar_t appDataPath[MAX_PATH];
        SHGetFolderPath(NULL, CSIDL_LOCAL_APPDATA, NULL, 0, appDataPath);
        wcscat(appDataPath, L"\\BassieMusic");

        // Init webview2 stuff
        ICoreWebView2CreateCoreWebView2EnvironmentCompletedHandler *environmentCompletedHandler =
            malloc(sizeof(ICoreWebView2CreateCoreWebView2EnvironmentCompletedHandler));
        environmentCompletedHandler->lpVtbl = &EnvironmentCompletedHandlerVtbl;
        if (FAILED(CreateCoreWebView2EnvironmentWithOptions(NULL, appDataPath, NULL, environmentCompletedHandler))) {
            FatalError(L"Failed to call CreateCoreWebView2EnvironmentWithOptions");
        }
    } else {
        FatalError(L"Failed to load WebView2Loader.dll");
    }

    // Main window event loop
    MSG message;
    while (GetMessage(&message, NULL, 0, 0) > 0) {
        TranslateMessage(&message);
        DispatchMessage(&message);
    }
    ExitProcess(message.wParam);
}
