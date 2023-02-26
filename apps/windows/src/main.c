// window decoration styling
// remove webview2 messaging

#define UNICODE
#include <dwmapi.h>
#include <objbase.h>
#include <shlobj.h>
#include <windows.h>
#define COBJMACROS
#include <d2d1.h>
#include <d3d11.h>
#include <d3d11_2.h>
#include <dxgi1_3.h>

#include "../res/resource.h"
#include "WebView2.h"
#include "about.h"
#include "dcomp.h"
#include "utils.h"

#define ID_MENU_ABOUT 2
#define WINDOW_STYLE (WS_POPUP | WS_THICKFRAME | WS_CAPTION | WS_SYSMENU | WS_MAXIMIZEBOX | WS_MINIMIZEBOX)

HWND window_hwnd;
UINT window_dpi;
BOOL window_titlebar_drag = FALSE;
ICoreWebView2CompositionController *composition_controller = NULL;
ICoreWebView2Controller *controller = NULL;
ICoreWebView2 *webview = NULL;
IDCompositionDevice *composition_device = NULL;
IDCompositionVisual *webview_visual = NULL;

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

// Default IUnknown method wrappers
HRESULT STDMETHODCALLTYPE Unknown_QueryInterface(IUnknown *This, REFIID riid, void **ppvObject) { return E_NOINTERFACE; }

ULONG STDMETHODCALLTYPE Unknown_AddRef(IUnknown *This) { return E_NOTIMPL; }

ULONG STDMETHODCALLTYPE Unknown_Release(IUnknown *This) { return E_NOTIMPL; }

// Forward interface reference
ICoreWebView2CreateCoreWebView2EnvironmentCompletedHandlerVtbl EnvironmentCompletedHandlerVtbl;
ICoreWebView2CreateCoreWebView2CompositionControllerCompletedHandlerVtbl ControllerCompletedHandlerVtbl;
ICoreWebView2WebMessageReceivedEventHandlerVtbl WebMessageReceivedEventHandlerVtbl;
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

    ICoreWebView2WebMessageReceivedEventHandler *webMessageReceivedEventHandler = malloc(sizeof(ICoreWebView2WebMessageReceivedEventHandler));
    webMessageReceivedEventHandler->lpVtbl = &WebMessageReceivedEventHandlerVtbl;
    ICoreWebView2_add_WebMessageReceived(webview, webMessageReceivedEventHandler, NULL);

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

// ICoreWebView2WebMessageReceivedEventHandler
HRESULT STDMETHODCALLTYPE WebMessageReceivedEventHandler_Invoke(ICoreWebView2WebMessageReceivedEventHandler *This, ICoreWebView2 *sender,
                                                                ICoreWebView2WebMessageReceivedEventArgs *args) {
    wchar_t *message;
    ICoreWebView2WebMessageReceivedEventArgs_get_WebMessageAsJson(args, &message);
    if (!wcscmp(message, L"\"minimize\"")) {
        ShowWindow(window_hwnd, SW_MINIMIZE);
    }
    if (!wcscmp(message, L"\"maximize\"")) {
        WINDOWPLACEMENT placement;
        GetWindowPlacement(window_hwnd, &placement);
        if (placement.showCmd == SW_MAXIMIZE) {
            ShowWindow(window_hwnd, SW_RESTORE);
        } else {
            ShowWindow(window_hwnd, SW_MAXIMIZE);
        }
    }
    if (!wcscmp(message, L"\"close\"")) {
        WindowClose();
    }
    return S_OK;
}

ICoreWebView2WebMessageReceivedEventHandlerVtbl WebMessageReceivedEventHandlerVtbl = {
    (HRESULT(STDMETHODCALLTYPE *)(ICoreWebView2WebMessageReceivedEventHandler * This, REFIID riid, void **ppvObject)) Unknown_QueryInterface,
    (ULONG(STDMETHODCALLTYPE *)(ICoreWebView2WebMessageReceivedEventHandler * This)) Unknown_AddRef,
    (ULONG(STDMETHODCALLTYPE *)(ICoreWebView2WebMessageReceivedEventHandler * This)) Unknown_Release, WebMessageReceivedEventHandler_Invoke};

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

        // Create D3D11Devic
        ID3D11Device *d3d11Device;
        D3D11CreateDevice(NULL, D3D_DRIVER_TYPE_HARDWARE, NULL, D3D11_CREATE_DEVICE_BGRA_SUPPORT, NULL, 0, D3D11_SDK_VERSION, &d3d11Device, NULL, NULL);

        // https://learn.microsoft.com/en-us/archive/msdn-magazine/2014/june/windows-with-c-high-performance-window-layering-using-the-windows-composition-engine
        IDXGIDevice *dxgiDevice;
        IID IID_IDXGIDevice = {0x54ec77fa, 0x1377, 0x44e6, {0x8c, 0x32, 0x88, 0xfd, 0x5f, 0x44, 0xc8, 0x4c}};
        ID3D11Device_QueryInterface(d3d11Device, &IID_IDXGIDevice, (void **)&dxgiDevice);

        // Create composition device
        HMODULE hDcomp = LoadLibrary(L"dcomp.dll");
        _DCompositionCreateDevice DCompositionCreateDevice = (_DCompositionCreateDevice)GetProcAddress(hDcomp, "DCompositionCreateDevice");
        IID IID_IDCompositionDevice = {0xC37EA93A, 0xE7AA, 0x450D, {0xB1, 0x6F, 0x97, 0x46, 0xCB, 0x04, 0x07, 0xF3}};
        DCompositionCreateDevice(dxgiDevice, &IID_IDCompositionDevice, (void **)&composition_device);

        IDCompositionTarget *hwndRenderTarget;
        composition_device->lpVtbl->CreateTargetForHwnd(composition_device, hwnd, TRUE, &hwndRenderTarget);

        // Create composition visuals
        IDCompositionVisual *rootVisual;
        composition_device->lpVtbl->CreateVisual(composition_device, &rootVisual);
        hwndRenderTarget->lpVtbl->SetRoot(hwndRenderTarget, rootVisual);

        composition_device->lpVtbl->CreateVisual(composition_device, &webview_visual);
        rootVisual->lpVtbl->AddVisual(rootVisual, webview_visual, TRUE, NULL);

        composition_device->lpVtbl->Commit(composition_device);
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
        ResizeBrowser(hwnd);
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

    // Window titlebar drag
    if (msg == WM_LBUTTONDOWN) {
        int x = GET_X_LPARAM(lParam);
        int y = GET_Y_LPARAM(lParam);
        RECT client_rect;
        GetClientRect(hwnd, &client_rect);

        // Mobile
        if (client_rect.right < MulDiv(1024, window_dpi, 96)) {
            if (y < MulDiv(52, window_dpi, 96) && x > MulDiv(52, window_dpi, 96) && x < client_rect.right - MulDiv(3 * 48, window_dpi, 96)) {
                window_titlebar_drag = TRUE;
            }
        }
        // Desktop
        else {
            if (y < MulDiv(30, window_dpi, 96) && x < client_rect.right - MulDiv(3 * 48, window_dpi, 96)) {
                window_titlebar_drag = TRUE;
            }
        }
    }
    if (msg == WM_MOUSEMOVE) {
        if (window_titlebar_drag) {
            ReleaseCapture();
            SendMessage(hwnd, WM_NCLBUTTONDOWN, HTCAPTION, 0);
        }
    }
    if (msg == WM_LBUTTONUP) {
        window_titlebar_drag = FALSE;
    }

    // Window titlebar popup system menu
    if (msg == WM_RBUTTONUP) {
        POINT point;
        POINTSTOPOINT(point, lParam);
        RECT client_rect;
        GetClientRect(hwnd, &client_rect);

        // Mobile
        if (client_rect.right < MulDiv(1024, window_dpi, 96)) {
            if (point.y < MulDiv(52, window_dpi, 96) && point.x > MulDiv(52, window_dpi, 96)) {
                ClientToScreen(hwnd, &point);
                TrackPopupMenu(GetSystemMenu(hwnd, FALSE), TPM_TOPALIGN | TPM_LEFTALIGN, point.x, point.y, 0, hwnd, NULL);
            }
        }
        // Desktop
        else {
            if (point.y < MulDiv(30, window_dpi, 96) && point.x < client_rect.right - MulDiv(3 * 48, window_dpi, 96)) {
                ClientToScreen(hwnd, &point);
                TrackPopupMenu(GetSystemMenu(hwnd, FALSE), TPM_TOPALIGN | TPM_LEFTALIGN, point.x, point.y, 0, hwnd, NULL);
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

    return DefWindowProc(hwnd, msg, wParam, lParam);
}

void _start(void) {
    // Only allow own app instance
    wchar_t *window_class_name = L"bassiemusic";
    HANDLE mutex = CreateMutex(NULL, TRUE, window_class_name);
    if (mutex == NULL) {
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
    wc.style = CS_HREDRAW | CS_VREDRAW;
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
