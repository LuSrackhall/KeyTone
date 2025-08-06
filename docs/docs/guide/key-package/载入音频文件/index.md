# Load Audio Files

**Loading Audio Files (Step 1)** is the most critical first step in creating a key sound album in KeyTone. Without this step, none of the subsequent exciting possibilities would be possible.

## **Why Is "Loading Audio Files" Necessary?**

Audio source files are the direct source of all sounds in a key sound album and are an indispensable component. The audio files you load will serve as raw materials for processing and use in subsequent steps.

* **Direct Use**: You can load pre-processed audio files that are of appropriate length and ready to be used directly for a single key press or release. Such files can skip the **Define and Trim Sounds (Step 2)** process and be used directly in **Craft Premium Key Sounds (Step 3)** for advanced key sound creation or in **Link Sounds to Key Actions (Step 4)** for binding to keys.
* **Further Processing**: Alternatively, you can load a longer recording, such as an audio file containing various key press sounds. These files can be precisely trimmed and defined in the **Define and Trim Sounds (Step 2)** step to create multiple independent sound segments for use in **Craft Premium Key Sounds (Step 3)** or **Link Sounds to Key Actions (Step 4)**.

## **How Does KeyTone Handle Your Loaded Files?**

In KeyTone, each time you load a source file, a new, independent key sound source file entry is created in the list. This means that even if you load the same file multiple times, they will appear as separate, independent entries in the source file list. These entries are also completely independent in their subsequent use and management, allowing you to assign different names or descriptions to distinguish and use them.

You don’t need to worry about loading the same audio source file multiple times causing the key sound album’s size to increase. KeyTone was designed with this in mind. To optimize storage efficiency, **KeyTone employs a smart storage mechanism: identical audio files are stored only once at the underlying level**.

Thanks to this clever storage design, even though they rely on the same underlying file, identical audio files added through multiple load operations will still be displayed as separate, independent entries in the management list. This independence extends beyond display; each entry is treated as a fully isolated file entry in KeyTone’s usage. This design maximizes flexibility and clarity during configuration without sacrificing storage efficiency.

## **How to Obtain "Audio Source Files" Resources?**

Audio source files, which are the prerequisite for everything, are not provided in this project and will not be provided in the future.

   <blockquote>
   <details>

   <summary>You can obtain audio source files through self-recording, searching open-source communities, free audio resource sharing websites, AI generation, etc.</summary>

   >
   > `Generally, you are free to use these audio resources locally; however, if you need to share them further, please be sure to review their specific license agreements.`
   >
   > * [Nigh/OpenKeySound](https://github.com/Nigh/OpenKeySound) — This repository, provided by [Nigh](https://github.com/Nigh), contains **self-recorded** and edited mechanical keyboard switch sounds, along with related usage instructions.
   > * [Pixabay](https://pixabay.com/sound-effects/search/keyboard/), [Freesound](https://freesound.org/search/?q=keyboard), and other websites claiming to share free audio resources. (*Note: Third-party website resources have not been verified for content, please validate them yourself.*)
   > * **With the rapid development and continuous advancement of artificial intelligence, it may be possible in the future to use AI audio generation technology to create customized keyboard sound effects through prompts; or even inform it of the KeyTone sound album format to generate sound albums directly importable for use.**
   </details>
   </blockquote>

## **Next Steps**

After loading audio files:
* You can decide whether to proceed to [**Define and Trim Sounds (Step 2)**](../../key-package/裁剪定义声音/index.md) for further processing.
* If the loaded audio files meet the key sound duration standards, you can directly use them in [**Craft Premium Key Sounds (Step 3)**](../../key-package/铸造至臻键音/index.md) or [**Link Sounds to Key Actions (Step 4)**](../../key-package/按键联动声效/index.md).