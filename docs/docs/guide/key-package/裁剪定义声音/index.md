# Define and Trim Sounds

After completing **Loading Audio Files (Step 1)**, you may have some raw audio materials. If these materials are longer recordings or you wish to extract multiple specific segments from a single audio file to use as independent sounds, **Define and Trim Sounds (Step 2)** will be your powerful tool.

## **What Is the Purpose of "Define and Trim Sounds"?**

The core purpose of this step is to take the audio source files uploaded in **Loading Audio Files (Step 1)** and, through precise trimming and parameter settings, define them into new, independent, and directly usable sound segments. These defined sounds will become key elements in building your personalized key sound album.

## **How to Trim and Define?**

* **Select Source File**: First, you need to choose an audio source file from the "Audio Source Files" list as the dependency for the defined sound. This source file is the audio file you uploaded in **Loading Audio Files (Step 1)**.
* **Set Time Segment**: Next, you need to specify the precise **start time** and **end time** for the sound segment you want to extract. KeyTone will extract the corresponding audio content from the source file based on the time range you set.
* **Adjust Initial Volume**: You can set an **initial volume parameter** for the newly defined sound. This volume parameter applies independently to the sound being defined and does not affect the volume of the original "audio source file" or other sounds defined from the same source file. In other words, each sound definition is independent and does not interfere with others. This provides great flexibility, allowing you to set different default volumes for different sound segments as needed.

Through **Define and Trim Sounds (Step 2)**, you can transform lengthy recordings or complex audio tracks into a series of concise, precise, and readily usable sound units, laying a solid foundation for subsequent creations.

## **Waveform-Based Visual Trimming** <Badge type="tip" text="Enhanced Feature" />

In addition to manually entering start/end timestamps, KeyTone now includes a built-in **waveform visualization component** in the Create/Edit Sound dialog, making the trimming process more intuitive and precise.

### Waveform & Selection

- The waveform of the selected audio source file is rendered directly inside the dialog.  
  Use the mouse wheel to zoom (with `Ctrl`) and scroll horizontally (with `Shift`) to navigate longer audio files.
- Right-click and drag on the waveform to create a **selection range** and visually define the trimming boundaries.  
  The selection remains **bi-directionally synchronized** with the numeric start/end time input fields.

### Frontend Preview Playback Bar

* A built-in **frontend preview playback bar** is placed above the waveform (fully running on the frontend layer, without relying on the SDK preview), supporting:
  * Play / Pause with a draggable playhead.
  * **Full playback** (entire source file) or **Selection playback** (only the defined trimmed range).
* When the dialog is closed, any ongoing preview playback will automatically stop, preventing audio from continuing after the window is dismissed.

### Volume Indicator Bar (dB)

* A **horizontal volume indicator bar** is integrated within the waveform. Drag it up or down to adjust the initial volume of the trimmed segment in real time:
  * Visible range: **±18 dB**, with scale markings at 18 / 12 / 6 / 0 / −6 / −12 / −18 displayed on the left.
  * The current dB value is shown on the right side of the waveform.
  * At 0 dB, the indicator rests on the waveform’s center baseline.
* This design allows you to fine-tune volume while previewing, eliminating the need for repetitive manual numeric input.

::: tip
The original numeric trimming workflow remains fully intact. The waveform component is an **enhancement**, not a replacement—you can freely choose whichever method best fits your preference.
:::

## **Next Steps**

<!-- **What Are the Uses of Defined Sounds?** -->

The sounds carefully trimmed and defined in this step have their own independent identities. They can be further used in:

* [**Craft Premium Key Sounds (Step 3)**](../../key-package/铸造至臻键音/index.md): As materials for creating advanced key sounds, whether for key press, release, or integration into different playback modes.
* [**Link Sounds to Key Actions (Step 4)**](../../key-package/按键联动声效/index.md): Directly bind these defined sounds to specific keys, reaching the final step of sound effect configuration.