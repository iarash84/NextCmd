param(
    [Parameter(Mandatory = $true)]
    [string]$Executable
)

$ErrorActionPreference = 'Stop'
$path = (Resolve-Path -LiteralPath $Executable).Path

# Notify Explorer about this file only. Do not delete global icon caches or
# restart Explorer: both would disrupt unrelated applications and windows.
Add-Type -TypeDefinition @'
using System;
using System.Runtime.InteropServices;

public static class NextCmdShellIcon
{
    [StructLayout(LayoutKind.Sequential, CharSet = CharSet.Unicode)]
    public struct SHFILEINFO
    {
        public IntPtr hIcon;
        public int iIcon;
        public uint attributes;
        [MarshalAs(UnmanagedType.ByValTStr, SizeConst = 260)]
        public string displayName;
        [MarshalAs(UnmanagedType.ByValTStr, SizeConst = 80)]
        public string typeName;
    }

    [DllImport("shell32.dll", CharSet = CharSet.Unicode)]
    public static extern IntPtr SHGetFileInfo(string path, uint attributes,
        out SHFILEINFO info, uint size, uint flags);

    [DllImport("shell32.dll", CharSet = CharSet.Unicode)]
    public static extern void SHUpdateImage(string hashItem, int iconIndex,
        uint flags, int imageIndex);

    [DllImport("shell32.dll", CharSet = CharSet.Unicode)]
    public static extern void SHChangeNotify(
        uint eventId, uint flags,
        [MarshalAs(UnmanagedType.LPWStr)] string item1,
        IntPtr item2);
}
'@

# UPDATEITEM alone can leave the per-path system image-list entry unchanged.
# Resolve and invalidate that entry before notifying open Explorer windows.
$info = New-Object NextCmdShellIcon+SHFILEINFO
$size = [Runtime.InteropServices.Marshal]::SizeOf($info)
$list = [NextCmdShellIcon]::SHGetFileInfo($path, 0, [ref]$info, $size, 0x00004000)
if ($list -ne [IntPtr]::Zero) {
    [NextCmdShellIcon]::SHUpdateImage($path, 0, 0, $info.iIcon)
}

# SHCNE_UPDATEITEM, SHCNF_PATHW | SHCNF_FLUSH.
[NextCmdShellIcon]::SHChangeNotify(0x00002000, 0x00001005, $path, [IntPtr]::Zero)
# Refresh the parent directory's cached listing as well (SHCNE_UPDATEDIR).
[NextCmdShellIcon]::SHChangeNotify(0x00001000, 0x00001005, (Split-Path -Parent $path), [IntPtr]::Zero)
# Windows 10/11 may retain additional Explorer icon entries beyond UPDATEITEM.
# Ask the built-in shell utility to refresh them without deleting cache files
# or terminating Explorer. Wait so the build does not finish before the request.
$refreshTool = Join-Path $env:SystemRoot 'System32/ie4uinit.exe'
if (Test-Path -LiteralPath $refreshTool) {
    $process = Start-Process -FilePath $refreshTool -ArgumentList '-show' -Wait -PassThru
    if ($process.ExitCode -ne 0) {
        Write-Warning "Windows icon refresh utility returned $($process.ExitCode)"
    }
}
Write-Output "Invalidated Windows image-list entry and requested icon refresh: $path"
