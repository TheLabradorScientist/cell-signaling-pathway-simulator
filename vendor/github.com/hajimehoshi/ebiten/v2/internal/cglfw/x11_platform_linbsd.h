// SPDX-License-Identifier: Apache-2.0
// SPDX-FileCopyrightText: 2002-2006 Marcus Geelnard
// SPDX-FileCopyrightText: 2006-2019 Camilla Löwy <elmindreda@glfw.org>
// SPDX-FileCopyrightText: 2023 The Ebitengine Authors

//go:build freebsd || linux || netbsd || openbsd

#include <unistd.h>
#include <signal.h>
#include <stdint.h>
#include <dlfcn.h>
#define XRRAllocGamma _glfw.x11.randr.AllocGamma
#define XRRFreeCrtcInfo _glfw.x11.randr.FreeCrtcInfo
#define XRRFreeGamma _glfw.x11.randr.FreeGamma
#define XRRFreeOutputInfo _glfw.x11.randr.FreeOutputInfo
#define XRRFreeScreenResources _glfw.x11.randr.FreeScreenResources
#define XRRGetCrtcGamma _glfw.x11.randr.GetCrtcGamma
#define XRRGetCrtcGammaSize _glfw.x11.randr.GetCrtcGammaSize
#define XRRGetCrtcInfo _glfw.x11.randr.GetCrtcInfo
#define XRRGetOutputInfo _glfw.x11.randr.GetOutputInfo
#define XRRGetOutputPrimary _glfw.x11.randr.GetOutputPrimary
#define XRRGetScreenResourcesCurrent _glfw.x11.randr.GetScreenResourcesCurrent
#define XRRQueryExtension _glfw.x11.randr.QueryExtension
#define XRRQueryVersion _glfw.x11.randr.QueryVersion
#define XRRSelectInput _glfw.x11.randr.SelectInput
#define XRRSetCrtcConfig _glfw.x11.randr.SetCrtcConfig
#define XRRSetCrtcGamma _glfw.x11.randr.SetCrtcGamma
#define XRRUpdateConfiguration _glfw.x11.randr.UpdateConfiguration
#define XShapeQueryExtension _glfw.x11.xshape.QueryExtension
#define XShapeQueryVersion _glfw.x11.xshape.QueryVersion
#define XShapeCombineRegion _glfw.x11.xshape.ShapeCombineRegion
#define XShapeCombineMask _glfw.x11.xshape.ShapeCombineMask

#include "posix_thread.h"
#include "posix_time_linbsd.h"
#include "xkb_unicode_linbsd.h"
#include "glx_context_linbsd.h"
#include "egl_context.h"
#include "osmesa_context.h"

#define _glfw_dlopen(name) dlopen(name, RTLD_LAZY | RTLD_LOCAL)
#define _glfw_dlclose(handle) dlclose(handle)
#define _glfw_dlsym(handle, name) dlsym(handle, name)

#define _GLFW_EGL_NATIVE_WINDOW  ((EGLNativeWindowType) window->x11.handle)
#define _GLFW_EGL_NATIVE_DISPLAY ((EGLNativeDisplayType) _glfw.x11.display)

#define _GLFW_PLATFORM_WINDOW_STATE         _GLFWwindowX11  x11
#define _GLFW_PLATFORM_LIBRARY_WINDOW_STATE _GLFWlibraryX11 x11
#define _GLFW_PLATFORM_MONITOR_STATE        _GLFWmonitorX11 x11
#define _GLFW_PLATFORM_CURSOR_STATE         _GLFWcursorX11  x11

void _glfwGrabErrorHandlerX11(void);
void _glfwReleaseErrorHandlerX11(void);
void _glfwInputErrorX11(int error, const char* message);

void _glfwPushSelectionToManagerX11(void);

