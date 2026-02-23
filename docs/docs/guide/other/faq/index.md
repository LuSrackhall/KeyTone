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
> ::: details Why are the configurations not synced between the Microsoft Store version and the Standard version?
> > **This behavior is caused by the "File System Virtualization (MSIX Container)" mechanism used by Windows for Store apps. Detailed explanation and solutions are as follows:**
> >
> > #### 1. Storage Path Differences
> > * **Standard Version (Win32):** Reads/Writes directly to the physical path:  
> >   `C:\Users\<Username>\AppData\Roaming\<Your_Config_Folder>`
> > * **Microsoft Store Version:** Windows automatically redirects operations to a private container path:  
> >   `C:\Users\<Username>\AppData\Local\Packages\<Package_Family_Name>\LocalCache\Roaming\<Your_Config_Folder>`
> >   *(Note: `<Package_Family_Name>` usually consists of the app name and a random string, e.g., `KeyTone_8wekyb3d8bbwe`)*
> >
> > #### 2. Why do configs sometimes sync, but other times stay independent?
> > This depends on your installation sequence (due to the Windows **"Merged View"** mechanism):
> > * **Standard First:** If you installed the Standard version first, the Store version will "see through" the container and read the existing physical configuration if its own container is empty.
> > * **Store First:** The Store version creates a new config folder inside its private container. Since the Standard version lacks container access permissions, it cannot "see" these files, resulting in independent configurations.
> >
> > #### 3. Uninstallation & Data Retention (Important)
> > * **Standard Version:** Uninstalling will **NOT** delete the configuration folder in `AppData\Roaming`. Your data remains safe.
> > * **Microsoft Store Version:** Follows the "Clean Uninstall" principle. Uninstalling the Store app will **permanently destroy** all private data within the `Packages` folder. **Please back up manually before uninstalling, or your data will be lost.**
> >   *(Note: If the Store version is "seeing through" and using the Standard version's config, uninstalling the Store version will not touch the external physical data.)*
> >
> > #### Solution: How to manually migrate configurations?
> > 1. Ensure the application is completely closed.
> > 2. Locate the source folder and copy its contents to the target folder:
> >    * **Standard ➔ Store:** Copy contents from `AppData\Roaming\<Your_Config_Folder>` to the Store version's `LocalCache\Roaming\<Your_Config_Folder>`.
> >    * **Store ➔ Standard:** Perform the reverse.
> >
> > If this solution does not work for you, feel free to submit an issue on [GitHub Issues](https://github.com/LuSrackhall/KeyTone/issues).
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
> ::: details Unable to open normally in macOS, with the prompt: `“KeyTone.app” is damaged and can’t be opened. You should move it to the Trash.`
> 
> #### Explanation of the Situation
> Thank you for using KeyTone.app! You may have encountered the error message from macOS stating `“KeyTone.app” is damaged and can’t be opened. You should move it to the Trash.` Below are the reasons and background explanation for this issue:
> 
> 1. **Incomplete Apple Developer Signing**:
>    - The development team of KeyTone.app currently does not hold a valid Apple Developer account, therefore, it is unable to perform Apple’s Notarization and Code Signing operations on the application.
>    - The Gatekeeper security mechanism in macOS by default prevents the running of unsigned or unnotarized applications to ensure user safety. This leads to the aforementioned error message.
> 
> 2. **macOS Security Mechanism**:
>    - Starting from macOS 10.15 (Catalina) and later versions, Apple requires that all applications downloaded from the internet must be notarized; otherwise, they will be marked by the system as “damaged” or “untrusted.”
> 
> 3. **Impact**:
>    - You may not be able to open KeyTone.app directly by double-clicking.
>    - This issue does not indicate a problem with the application itself, but rather that macOS’s security policy restricts the running of unsigned applications.
> 
> #### Solution
> To resolve this issue and run KeyTone.app normally, you can follow the steps below to manually remove the application’s quarantine attribute and allow it to run.
> 
> ##### Using the `xattr` Command to Remove Quarantine Attribute (Recommended)
> 1. **Open Terminal**:
>    - Press `Command + Space` to open Spotlight search.
>    - Type `Terminal` and open the Terminal application.
> 
> 2. **Execute Command**:
>    - Enter the following command and press Enter:
>      ```
>      xattr -cr /Applications/KeyTone.app
>      ```
>    - Explanation:
>      - `xattr -cr` will remove the extended attributes of KeyTone.app (including the quarantine flag), allowing macOS to run it.
>      - Ensure that the path `/Applications/KeyTone.app` is correct. If you have installed the application in another location, replace it with the actual path (e.g., `/Users/YourUsername/Downloads/KeyTone.app`).
> 
> 3. **Verification**:
>    - After executing the command, try double-clicking to open KeyTone.app; the application should run normally.
> 
> 4. **Notes**:
>    - Executing the command may require administrator privileges, and the system will prompt for a password.
>    - If the command is ineffective, ensure that the entered path is correct, or check if the application has been moved or renamed.
> 
> #### Additional Recommendations
> - **Check Application Source**: Ensure that KeyTone.app is downloaded from official or trusted channels to avoid potential security risks.
> - **Update the Application**: The development team may plan to apply for an Apple Developer account in the future. Therefore, future versions may resolve this issue. Please follow official channels (such as the official website or social media) for updates.
> - **Contact Support**: If the above methods do not resolve the issue, or if you have other questions, please visit [GitHub Issues](https://github.com/LuSrackhall/KeyTone/issues) to check for solutions or raise a new issue.
> 
> #### Disclaimer
> - Ensure that you download KeyTone.app from official or trusted sources.
> - Performing the above operations is at the user’s own risk; it is recommended to back up important data before proceeding.
> - If you are unsure about the operation steps, please consult a professional or contact the development team for support.
> 
> Thank you for your understanding and support! If you have further questions, please feel free to contact us.
> 
> ---
> 
> **KeyTone.app Development Team**  
> Date: May 20, 2025
> :::