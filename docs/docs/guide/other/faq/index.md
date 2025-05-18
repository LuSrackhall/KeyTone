# Frequently Asked Questions

::: info
For open-source projects, it's a good practice to search for solutions or submit new issues in [GitHub Issues](https://github.com/LuSrackhall/KeyTone/issues).
:::

Exiting via the system tray and reopening the application can resolve 99.99% of uncommon bugs.

**Common non-bug-related questions are listed below:**

::: tip
Click on a question to view the solution.
:::

> ::: details In Windows, the "Auto-Start" feature is enabled but does not work. (This issue does not occur in the Microsoft Store version)
> > **This may be due to Windows system settings managing auto-start applications. The detailed solution is as follows:**
> > * Open the system 'Settings' > 'Apps' > 'Startup' page.(Note: This refers to opening the operating system's 'Settings'.)
> > * Check if the application is allowed to auto-start.
> > * If disabled, set it to enabled.
> >
> > If this solution does not work for you, consider submitting an issue on [GitHub Issues](https://github.com/LuSrackhall/KeyTone/issues).
> :::
> ::: details In macOS, every time the software is opened, an "Accessibility Access" request pop-up appears, and the software functions do not work properly.
>
> > **This may be caused by an update to the Apple system, though the exact reason is unclear. The detailed solution is as follows:**
> >
> > Try granting "Accessibility" permissions to `/Applications/KeyTone.app/Contents/MacOS/KeyTone` (note: not `/Applications/KeyTone.app`).
> >
> > Detailed steps:
> > * Exit `KeyTone.app` via the system tray icon.
> > * Open “System Settings” and navigate to the “Accessibility” settings:
> >   * Click the Apple menu () in the top-left corner of the screen and select “System Settings” (or “System Preferences,” depending on the macOS version).
> >   * In “System Settings,” find and click “Privacy & Security” > “Accessibility.” Keep this window open, as you’ll need to drag an icon into the permissions settings box later.
> > * Use a shortcut to open the specified path in Finder:
> >   * Open “Finder.”
> >   * Press `Command + Shift + G` to open the “Go to Folder” window.
> >   * In the input box, copy and paste the following path: `/Applications/KeyTone.app/Contents/MacOS/KeyTone`, then press “Go” or the Enter key.
> > * Drag the KeyTone executable file to the “Accessibility” settings window:
> >   * In Finder, locate the `KeyTone` file under the specified path (this is an executable file, typically displayed with a terminal icon).
> >   * Drag the `KeyTone` file to the application list in the “Accessibility” window in “System Settings.”
> >   * If prompted for an administrator password, enter it to confirm.
> > * Ensure permissions are enabled:
> >   * In the “Accessibility” list, find `KeyTone` (the one you manually added, not `KeyTone.app`) and ensure the checkbox next to it is selected.
> >   * If unchecked, manually check it to enable the permission.
> > * Reopen `KeyTone.app`.
> >
> > If this solution does not work for you, consider submitting an issue on [GitHub Issues](https://github.com/LuSrackhall/KeyTone/issues).
> :::