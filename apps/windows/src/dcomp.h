#pragma once

// Custom dcomp headers because they only have C++ headers and no IDl file is published :(

typedef struct IDCompositionVisual IDCompositionVisual;
typedef struct IDCompositionVisualVtbl {
    IUnknownVtbl base;
    void *padding[13];
    void(__stdcall *AddVisual)(IDCompositionVisual *this, IDCompositionVisual* visual, BOOL insertAbove, IDCompositionVisual *referenceVisual);
} IDCompositionVisualVtbl;
struct IDCompositionVisual {
    IDCompositionVisualVtbl *lpVtbl;
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
    void *padding[2];
    void(__stdcall *CreateTargetForHwnd)(IDCompositionDevice *this, HWND hwnd, BOOL topmost, IDCompositionTarget **target);
    void(__stdcall *CreateVisual)(IDCompositionDevice *this, IDCompositionVisual **visual);
} IDCompositionDeviceVtbl;
struct IDCompositionDevice {
    IDCompositionDeviceVtbl *lpVtbl;
};

typedef HRESULT (*_DCompositionCreateDevice)(IDXGIDevice *dxgiDevice, REFIID iid, void **dcompositionDevice);
