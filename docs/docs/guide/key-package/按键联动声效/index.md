# Key Sound Linkage

After the previous steps—whether directly using the original audio from **Load Audio Files (Step 1)**, the carefully processed clips from **Define Sounds by Trimming (Step 2)**, or the advanced sound combinations created through **Forge Premium Key Sounds (Step 3)**—now, it’s time to link these wonderful sounds to your actual key operations. This is the core task of **Key Sound Linkage (Step 4)**.

## What is the purpose of "Key Sound Linkage"?

The purpose of this step is straightforward: it allows you to bind your carefully prepared various sound resources (including the original **audio source files**, the trimmed **defined sounds**, and the dynamically varied **premium key sounds**) to a specific key or to all keys globally.

This way, when you press or release these keys, KeyTone will play the corresponding sound effects you have set.

## Flexible Binding Strategies

KeyTone offers multiple binding methods to meet the needs of different users and allows them to coexist harmoniously:

* **Global Binding:**
  * If you want to save trouble or pursue a unified key sound style, you can directly use global binding.
  * In global binding mode, you can select a sound for the press or release event of all keys globally. This sound can be an "audio source file" that meets the key sound duration standard, a "trimmed defined sound," or a more advanced "premium key sound."
  * **Benefits**: Extremely easy to operate; you don’t need to tediously set sounds for each key, allowing for quick global sound effect replacement.
  * **Common Usage**: A very practical and effective usage is to configure a "premium key sound" in "random mode" for **press/release sounds** and perform global key binding. This way, each of your keystrokes will have subtle, unpredictable sound variations, which are both unified and not monotonous, greatly enhancing the pleasure of typing.

* **Individual Key Configuration:**
  * For users who pursue ultimate personalization, KeyTone allows you to configure the press and release sounds for each individual key separately.
  * **Benefits**: Maximum personalization, covering every key. You can set one sound for the `A` key, another for the `Space` key, and even completely different sound combinations for the press and release of the `Shift` key(or any key), achieving the highest level of customization.
  * **Common Usage**: A very common usage is to clone the sounds of your actual keyboard to a key sound album. You can record the actual key press sounds of your keyboard (either directly as **audio source files (Step 1)** that meet the key sound duration standard or roughly record and further trim them into **sounds from (Step 2)**), and then bind them one-to-one to the corresponding keys to achieve the ultimate goal of cloning the actual keyboard sounds.

* **Priority Override Logic:**
  * You can **configure both global key bindings and individual key bindings simultaneously**.
  * Their working logic is based on **priority override**:
    * If you set a global key sound, by default, all keys will use this global sound effect.
    * However, if you also configure a sound for a specific key (e.g., the `Enter` key), when this `Enter` key is triggered, **it will prioritize playing the individually configured sound instead of the global sound effect**, meaning the global sound effect is overridden by the individually configured key sound for that specific key.
    * For other keys that do not have individual configurations, they will continue to use the initially bound global sound effect.
  * This design allows you to perform "special care" for a few keys you particularly care about on the basis of a unified global sound effect, achieving both efficient and personalized configuration.

## Conclusion

Through **Key Sound Linkage (Step 4)**, all your preparation work will be transformed into actual auditory feedback at this moment. Whether it’s simple and unified, extremely personalized, or a clever combination of both, KeyTone strives to provide you with the best configuration experience, making your keys "sound" alive from now on!