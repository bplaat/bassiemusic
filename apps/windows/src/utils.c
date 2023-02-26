#include "utils.h"

#define COBJMACROS
#include <wincodec.h>

// We don't like with Standard C Library so we have to write our own win32 implementations
void *malloc(size_t size) { return HeapAlloc(GetProcessHeap(), 0, size); }

void free(void *ptr) { HeapFree(GetProcessHeap(), 0, ptr); }

size_t wcslen(const wchar_t *string) {
    wchar_t *c = (wchar_t *)string;
    while (*c != '\0') c++;
    return c - string;
}

wchar_t *wcscpy(wchar_t *dest, const wchar_t *src) {
    wchar_t *start = dest;
    while ((*dest++ = *src++) != '\0')
        ;
    return start;
}

wchar_t *wcscat(wchar_t *dest, const wchar_t *src) {
    wchar_t *start = dest;
    while (*dest != '\0') dest++;
    wcscpy(dest, src);
    return start;
}

int wcscmp(const wchar_t *s1, const wchar_t *s2) {
    while (*s1 && (*s1 == *s2)) {
        s1++;
        s2++;
    }
    return *(const wchar_t *)s1 - *(const wchar_t *)s2;
}

// Win32 Helper functions
UINT GetPrimaryDesktopDpi(void) {
    HDC hdc = GetDC(HWND_DESKTOP);
    UINT dpi = GetDeviceCaps(hdc, LOGPIXELSY);
    ReleaseDC(HWND_DESKTOP, hdc);
    return dpi;
}

typedef BOOL(STDMETHODCALLTYPE *_AdjustWindowRectExForDpi)(RECT *lpRect, DWORD dwStyle, BOOL bMenu, DWORD dwExStyle, UINT dpi);

BOOL AdjustWindowRectExForDpi(RECT *lpRect, DWORD dwStyle, BOOL bMenu, DWORD dwExStyle, UINT dpi) {
    HMODULE huser32 = LoadLibrary(L"user32.dll");
    _AdjustWindowRectExForDpi AdjustWindowRectExForDpi = (_AdjustWindowRectExForDpi)GetProcAddress(huser32, "AdjustWindowRectExForDpi");
    if (AdjustWindowRectExForDpi) {
        return AdjustWindowRectExForDpi(lpRect, dwStyle, bMenu, dwExStyle, dpi);
    }
    return AdjustWindowRectEx(lpRect, dwStyle, bMenu, dwExStyle);
}

wchar_t *GetString(UINT id) {
    wchar_t *string;
    LoadString(GetModuleHandle(NULL), id, (wchar_t *)&string, 0);
    return string;
}

HBITMAP LoadPNGFromResource(wchar_t *type, wchar_t *name) {
    HRSRC hsrc = FindResource(NULL, name, type);

    CLSID CLSID_WICImagingFactory = {0xcacaf262, 0x9370, 0x4615, {0xa1, 0x3b, 0x9f, 0x55, 0x39, 0xda, 0x4c, 0x0a}};
    IID IID_IWICImagingFactory = {0xec5ec8a9, 0xc395, 0x4314, {0x9c, 0x77, 0x54, 0xd7, 0xa9, 0x35, 0xff, 0x70}};
    IWICImagingFactory *wicFactory;
    CoCreateInstance(&CLSID_WICImagingFactory, NULL, CLSCTX_INPROC_SERVER, &IID_IWICImagingFactory, (void **)&wicFactory);

    IWICStream *wicStream;
    IWICImagingFactory_CreateStream(wicFactory, &wicStream);
    IWICStream_InitializeFromMemory(wicStream, LockResource(LoadResource(NULL, hsrc)), SizeofResource(NULL, hsrc));

    IWICBitmapDecoder *wicDecoder;
    IWICImagingFactory_CreateDecoderFromStream(wicFactory, (IStream *)wicStream, NULL, WICDecodeMetadataCacheOnDemand, &wicDecoder);

    IWICBitmapFrameDecode *wicFrame;
    IWICBitmapDecoder_GetFrame(wicDecoder, 0, &wicFrame);
    UINT width, height;
    IWICBitmapSource_GetSize(wicFrame, &width, &height);

    IWICFormatConverter *wicConverter;
    IWICImagingFactory_CreateFormatConverter(wicFactory, &wicConverter);
    GUID GUID_WICPixelFormat24bppBGR = {0x6fddc324, 0x4e03, 0x4bfe, {0xb1, 0x85, 0x3d, 0x77, 0x76, 0x8d, 0xc9, 0x0c}};
    IWICFormatConverter_Initialize(wicConverter, (IWICBitmapSource *)wicFrame, &GUID_WICPixelFormat24bppBGR, WICBitmapDitherTypeNone, NULL, 0,
                                   WICBitmapPaletteTypeCustom);

    IID IID_IWICBitmapSource = {0x00000120, 0xa8f2, 0x4877, {0xba, 0x0a, 0xfd, 0x2b, 0x66, 0x45, 0xfb, 0x94}};
    IWICBitmapSource *wicConvertedSource;
    IWICFormatConverter_QueryInterface(wicConverter, &IID_IWICBitmapSource, (void **)&wicConvertedSource);

    HDC hdc = GetDC(HWND_DESKTOP);
    BITMAPINFO bitmapInfo = {0};
    bitmapInfo.bmiHeader.biSize = sizeof(BITMAPINFOHEADER);
    bitmapInfo.bmiHeader.biWidth = width;
    bitmapInfo.bmiHeader.biHeight = -height;
    bitmapInfo.bmiHeader.biPlanes = 1;
    bitmapInfo.bmiHeader.biBitCount = 24;
    bitmapInfo.bmiHeader.biCompression = BI_RGB;
    BYTE *bitmapBuffer = NULL;
    HBITMAP bitmap = CreateDIBSection(hdc, &bitmapInfo, DIB_RGB_COLORS, (void **)&bitmapBuffer, NULL, 0);
    IWICBitmapSource_CopyPixels(wicConvertedSource, NULL, width * 3, width * 3 * height, bitmapBuffer);
    ReleaseDC(HWND_DESKTOP, hdc);

    IWICBitmapSource_Release(wicConvertedSource);
    IWICFormatConverter_Release(wicConverter);
    IWICBitmapFrameDecode_Release(wicFrame);
    IWICBitmapFrameDecode_Release(wicDecoder);
    IWICStream_Release(wicStream);
    IWICImagingFactory_Release(wicFactory);
    return bitmap;
}

void GetAppVersion(UINT *version) {
    wchar_t file_name[MAX_PATH];
    GetModuleFileName(NULL, file_name, sizeof(file_name) / sizeof(wchar_t));
    DWORD version_info_size = GetFileVersionInfoSize(file_name, NULL);
    BYTE *version_info = malloc(version_info_size);
    GetFileVersionInfo(file_name, 0, version_info_size, version_info);

    VS_FIXEDFILEINFO *file_info;
    UINT file_info_size;
    VerQueryValue(version_info, L"\\", (LPVOID *)&file_info, &file_info_size);

    version[0] = HIWORD(file_info->dwProductVersionMS);
    version[1] = LOWORD(file_info->dwProductVersionMS);
    version[2] = HIWORD(file_info->dwProductVersionLS);
    version[3] = LOWORD(file_info->dwProductVersionLS);

    free(version_info);
}
