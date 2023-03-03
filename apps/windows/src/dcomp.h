#pragma once

// Custom DComposition headers because they only have C++ headers and no IDL file is published :(
#include <dxgi1_3.h>

typedef struct IDCompositionVisual IDCompositionVisual;
typedef struct IDCompositionVisualVtbl {
    IUnknownVtbl base;
    void *padding[12];
    void(__stdcall *SetContent)(IDCompositionVisual *this, IUnknown *content);
    void(__stdcall *AddVisual)(IDCompositionVisual *this, IDCompositionVisual *visual, BOOL insertAbove, IDCompositionVisual *referenceVisual);
} IDCompositionVisualVtbl;
struct IDCompositionVisual {
    IDCompositionVisualVtbl *lpVtbl;
};

typedef struct IDCompositionVirtualSurface IDCompositionVirtualSurface;
typedef struct IDCompositionVirtualSurfaceVtbl {
    IUnknownVtbl base;
    void(__stdcall *BeginDraw)(IDCompositionVirtualSurface *this, const RECT *updateRect, REFIID iid, void **updateObject, POINT *updateOffset);
    void(__stdcall *EndDraw)(IDCompositionVirtualSurface *this);
    void *padding[3];
    void(__stdcall *Resize)(IDCompositionVirtualSurface *this, UINT width, UINT height);
} IDCompositionVirtualSurfaceVtbl;
struct IDCompositionVirtualSurface {
    IDCompositionVirtualSurfaceVtbl *lpVtbl;
};

typedef struct IDCompositionTarget IDCompositionTarget;
typedef struct IDCompositionTargetVtbl {
    IUnknownVtbl base;
    void(__stdcall *SetRoot)(IDCompositionTarget *this, IDCompositionVisual *visual);
} IDCompositionTargetVtbl;
struct IDCompositionTarget {
    IDCompositionTargetVtbl *lpVtbl;
};

typedef struct IDCompositionDevice IDCompositionDevice;
typedef struct IDCompositionDeviceVtbl {
    IUnknownVtbl base;
    void(__stdcall *Commit)(IDCompositionDevice *this);
    void *padding1[2];
    void(__stdcall *CreateTargetForHwnd)(IDCompositionDevice *this, HWND hwnd, BOOL topmost, IDCompositionTarget **target);
    void(__stdcall *CreateVisual)(IDCompositionDevice *this, IDCompositionVisual **visual);
    void *padding2[1];
    void(__stdcall *CreateVirtualSurface)(IDCompositionDevice *this, UINT width, UINT height, DXGI_FORMAT pixelFormat, DXGI_ALPHA_MODE alphaMode,
                                          IDCompositionVirtualSurface **surface);
} IDCompositionDeviceVtbl;
struct IDCompositionDevice {
    IDCompositionDeviceVtbl *lpVtbl;
};

typedef HRESULT (*_DCompositionCreateDevice)(IDXGIDevice *dxgiDevice, REFIID iid, void **dcompositionDevice);
