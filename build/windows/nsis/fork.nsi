Unicode True
!include "MUI2.nsh"
!include "LogicLib.nsh"
!include "x64.nsh"
!include "WinVer.nsh"

!ifndef APP_VERSION
  !error "APP_VERSION is required"
!endif
!ifndef FILE_VERSION
  !error "FILE_VERSION is required"
!endif

!define APP_EXE "TcNo-Acc-Switcher.exe"
!define UNINSTALL_KEY "Software\Microsoft\Windows\CurrentVersion\Uninstall\BearSwitcher"
!define WEBVIEW_KEY "Software\Microsoft\EdgeUpdate\Clients\{F3017226-FE2A-4295-8BDF-00C3A9A7E4C5}"

Name "BearSwitcher"
OutFile "..\..\..\bin\BearSwitcher-Setup.exe"
InstallDir "$LOCALAPPDATA\Programs\BearSwitcher"
InstallDirRegKey HKCU "${UNINSTALL_KEY}" "InstallLocation"
RequestExecutionLevel user
ManifestDPIAware true
SetCompressor /SOLID lzma
ShowInstDetails show
ShowUninstDetails show
VIProductVersion "${FILE_VERSION}"
VIAddVersionKey "ProductName" "BearSwitcher"
VIAddVersionKey "FileDescription" "BearSwitcher Setup"
VIAddVersionKey "FileVersion" "${APP_VERSION}"
VIAddVersionKey "LegalCopyright" "TroubleChute (Wesley Pyburn) and contributors"

!define MUI_ICON "..\icon.ico"
!define MUI_UNICON "..\icon.ico"
!define MUI_ABORTWARNING
!define MUI_LANGDLL_REGISTRY_ROOT HKCU
!define MUI_LANGDLL_REGISTRY_KEY "${UNINSTALL_KEY}"
!define MUI_LANGDLL_REGISTRY_VALUENAME "InstallerLanguage"
!define MUI_LANGDLL_ALLLANGUAGES
!insertmacro MUI_PAGE_WELCOME
!insertmacro MUI_PAGE_LICENSE "..\..\..\LICENSE"
!insertmacro MUI_PAGE_DIRECTORY
!insertmacro MUI_PAGE_COMPONENTS
!insertmacro MUI_PAGE_INSTFILES
!define MUI_FINISHPAGE_RUN "$INSTDIR\${APP_EXE}"
!define MUI_FINISHPAGE_RUN_TEXT "$(LaunchApp)"
!insertmacro MUI_PAGE_FINISH
!insertmacro MUI_UNPAGE_CONFIRM
!insertmacro MUI_UNPAGE_INSTFILES
!insertmacro MUI_LANGUAGE "Arabic"
!insertmacro MUI_LANGUAGE "English"

LangString MainFiles ${LANG_ARABIC} "ملفات البرنامج واختصار قائمة ابدأ"
LangString MainFiles ${LANG_ENGLISH} "Application and Start menu shortcut"
LangString DesktopShortcut ${LANG_ARABIC} "اختصار سطح المكتب"
LangString DesktopShortcut ${LANG_ENGLISH} "Desktop shortcut"
LangString LaunchApp ${LANG_ARABIC} "تشغيل BearSwitcher"
LangString LaunchApp ${LANG_ENGLISH} "Launch BearSwitcher"
LangString CloseApp ${LANG_ARABIC} "أغلق BearSwitcher من قائمة أيقونته بجوار الساعة، ثم أعد المحاولة."
LangString CloseApp ${LANG_ENGLISH} "Exit BearSwitcher from its tray menu, then try again."
LangString UnsupportedOS ${LANG_ARABIC} "يتطلب البرنامج ويندوز 10 أو أحدث بمعمارية x64."
LangString UnsupportedOS ${LANG_ENGLISH} "Windows 10 or later on x64 is required."
LangString RuntimeFailed ${LANG_ARABIC} "تعذر تثبيت WebView2. تحقق من اتصال الإنترنت ثم أعد تشغيل المثبّت."
LangString RuntimeFailed ${LANG_ENGLISH} "WebView2 installation failed. Check your internet connection and run setup again."

; Do not kill the running app or replace an executable it still has open.
!macro RequireClosedApp
  System::Call 'kernel32::OpenMutexW(i 0x100000, i 0, w "TcNo-Acc-Switcher-Singleton") p.r0'
  ${If} $0 != 0
    System::Call 'kernel32::CloseHandle(p r0)'
    MessageBox MB_OK|MB_ICONEXCLAMATION "$(CloseApp)" /SD IDOK
    SetErrorLevel 2
    Quit
  ${EndIf}
!macroend

Function .onInit
  SetShellVarContext current
  SetRegView 64
  StrCpy $LANGUAGE ${LANG_ARABIC}
  !insertmacro MUI_LANGDLL_DISPLAY
  ${IfNot} ${AtLeastWin10}
    MessageBox MB_OK|MB_ICONSTOP "$(UnsupportedOS)" /SD IDOK
    SetErrorLevel 3
    Quit
  ${EndIf}
  ${IfNot} ${IsNativeAMD64}
    MessageBox MB_OK|MB_ICONSTOP "$(UnsupportedOS)" /SD IDOK
    SetErrorLevel 3
    Quit
  ${EndIf}
  !insertmacro RequireClosedApp
FunctionEnd

Function DetectWebView
  ReadRegStr $0 HKLM "Software\WOW6432Node\Microsoft\EdgeUpdate\Clients\{F3017226-FE2A-4295-8BDF-00C3A9A7E4C5}" "pv"
  ${If} $0 == ""
  ${OrIf} $0 == "0.0.0.0"
    ReadRegStr $0 HKCU "${WEBVIEW_KEY}" "pv"
  ${EndIf}
FunctionEnd

Section "$(MainFiles)" MainSection
  SectionIn RO
  Call DetectWebView
  ${If} $0 == ""
  ${OrIf} $0 == "0.0.0.0"
    InitPluginsDir
    SetOutPath "$PLUGINSDIR"
    File "MicrosoftEdgeWebview2Setup.exe"
    ExecWait '"$PLUGINSDIR\MicrosoftEdgeWebview2Setup.exe" /silent /install' $1
    Call DetectWebView
    ${If} $0 == ""
    ${OrIf} $0 == "0.0.0.0"
      MessageBox MB_OK|MB_ICONSTOP "$(RuntimeFailed)" /SD IDOK
      SetErrorLevel 4
      Abort
    ${EndIf}
  ${EndIf}

  SetOutPath "$INSTDIR"
  File "..\..\..\bin\${APP_EXE}"
  File "..\..\..\LICENSE"
  File "..\..\..\OPEN_SOURCE_LICENSES.txt"
  WriteUninstaller "$INSTDIR\Uninstall.exe"
  CreateDirectory "$SMPROGRAMS\BearSwitcher"
  CreateShortcut "$SMPROGRAMS\BearSwitcher\BearSwitcher.lnk" "$INSTDIR\${APP_EXE}"
  CreateShortcut "$SMPROGRAMS\BearSwitcher\Uninstall.lnk" "$INSTDIR\Uninstall.exe"
  WriteRegStr HKCU "${UNINSTALL_KEY}" "DisplayName" "BearSwitcher"
  WriteRegStr HKCU "${UNINSTALL_KEY}" "DisplayVersion" "${APP_VERSION}"
  WriteRegStr HKCU "${UNINSTALL_KEY}" "InstallLocation" "$INSTDIR"
  WriteRegStr HKCU "${UNINSTALL_KEY}" "DisplayIcon" "$INSTDIR\${APP_EXE}"
  WriteRegStr HKCU "${UNINSTALL_KEY}" "Publisher" "2-sa / Switcher contributors"
  WriteRegStr HKCU "${UNINSTALL_KEY}" "URLInfoAbout" "https://github.com/2-sa/Switcher"
  WriteRegStr HKCU "${UNINSTALL_KEY}" "UninstallString" '$\"$INSTDIR\Uninstall.exe$\"'
  WriteRegStr HKCU "${UNINSTALL_KEY}" "QuietUninstallString" '$\"$INSTDIR\Uninstall.exe$\" /S'
  WriteRegDWORD HKCU "${UNINSTALL_KEY}" "NoModify" 1
  WriteRegDWORD HKCU "${UNINSTALL_KEY}" "NoRepair" 1
SectionEnd

Section /o "$(DesktopShortcut)" DesktopSection
  CreateShortcut "$DESKTOP\BearSwitcher.lnk" "$INSTDIR\${APP_EXE}"
SectionEnd

Function un.onInit
  SetShellVarContext current
  SetRegView 64
  !insertmacro MUI_UNGETLANGUAGE
  !insertmacro RequireClosedApp
FunctionEnd

Section "Uninstall"
  ; Never remove user data, custom files, or recursively delete an install dir.
  ReadRegStr $0 HKCU "${UNINSTALL_KEY}" "InstallLocation"
  ${If} $0 == $INSTDIR
    Delete "$SMPROGRAMS\BearSwitcher\BearSwitcher.lnk"
    Delete "$SMPROGRAMS\BearSwitcher\Uninstall.lnk"
    RMDir "$SMPROGRAMS\BearSwitcher"
    Delete "$DESKTOP\BearSwitcher.lnk"
    DeleteRegKey HKCU "${UNINSTALL_KEY}"
  ${EndIf}
  Delete "$INSTDIR\${APP_EXE}"
  Delete "$INSTDIR\LICENSE"
  Delete "$INSTDIR\OPEN_SOURCE_LICENSES.txt"
  Delete "$INSTDIR\Uninstall.exe"
  RMDir "$INSTDIR"
SectionEnd
