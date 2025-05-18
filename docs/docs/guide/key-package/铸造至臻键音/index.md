# Craft Premium Key Sounds

While directly obtaining sounds through **Loading Audio Files (Step 1)**, refining them via **Define and Trim Sounds (Step 2)**, and binding them in **Link Sounds to Key Actions (Step 4)** is sufficient to create a highly refined key sound album, KeyTone’s design goes beyond this. 

We aim to provide a more advanced, dynamic, and enriched key usage experience. Thus, we introduce **Craft Premium Key Sounds (Step 3)**, often referred to as “defining advanced key sounds.”

## **Core Concept of "Premium Key Sounds"**

The design philosophy of “Premium Key Sounds” (or “Advanced Key Sounds”) is based on defining a **complete key activity cycle** as the fundamental unit. This cycle includes not only the sound triggered when a key is **pressed** but also the sound triggered when a key is **released**. However, you are not required to configure both; you can **define only one part or none at all**. KeyTone allows you to create a blank “Premium Key Sound” configuration. For example, you can bind a blank “Premium Key Sound” to specific keys to mute them.

Of course, the core design of “Premium Key Sounds” also includes **high flexibility** and **diverse playback modes**, which will be detailed below.

## **High Flexibility**

When linked to keys in **Link Sounds to Key Actions (Step 4)**, “Premium Key Sounds” (or “Advanced Key Sounds”) offer great flexibility in how they are **used**.

For instance, you can choose to use only its “press sound,” only its “release sound,” or, by default, both the “press and release sounds” together. In short, the design of “Premium Key Sounds” fully considers the creative freedom needed, providing you with endless possibilities.

## **Diverse Playback Modes**

The “advanced” nature of “Premium Key Sounds” lies in the carefully designed variety of playback modes. These modes inject infinite possibilities into your key sound experience, including:

* **`Single` Mode:**
  * **Dependency**: In this mode, the press or release part of a “Premium Key Sound” can only depend on a single specific sound. This sound can be an **audio source file (from Step 1)** that meets key sound duration standards, a **defined sound (from Step 2)**, or even **another “Premium Key Sound”**!
  * **Playback Behavior**: Each time it is triggered (press or release), the specified sound is played consistently.
  * **Inheritance Behavior**: When the specified sound is another “Premium Key Sound,” it fully inherits all sound effects and playback modes of that “Premium Key Sound” and plays accordingly.
  * **Nesting Performance**: When the specified sound is another “Premium Key Sound,” one or more layers of nesting are inevitable. However, since KeyTone’s backend is developed in Go, you can confidently use multi-layer or even mutual nesting. No matter how complex your nesting combinations are, you won’t experience any perceptible lag, as Go’s high-performance fully support your “tinkering” spirit.

* **`Random` Mode:**
  * **Dependency**: In this mode, the press or release part of a “Premium Key Sound” can depend on **multiple** **audio source files** or **defined sounds** that meet key sound duration standards, or even **other “Premium Key Sounds”**.
  * **Playback Behavior**: Each time it is triggered, a sound is randomly selected from the dependent sound list for playback, adding surprise and variety to your typing experience.
  * **Inheritance Behavior**: Since the dependent sounds can include **other “Premium Key Sounds”**, when a “Premium Key Sound” is randomly selected from the list, it fully inherits all sound effects and playback modes of that “Premium Key Sound” and plays accordingly.
  * **Nesting Performance**: When the specified sounds include other “Premium Key Sounds,” one or more layers of nesting are inevitable. However, since KeyTone’s backend is developed in Go, you can confidently use multi-layer or even mutual nesting. No matter how complex your nesting combinations are, you won’t experience any perceptible lag, as Go’s high-performance fully support your “tinkering” spirit.

* **`Sequential` Mode:**
  * **Dependency**: Similar to the Random mode, the press or release part of a “Premium Key Sound” in this mode can depend on **multiple** **audio source files** or **defined sounds** that meet key sound duration standards, or even **other “Premium Key Sounds”**.
  * **Playback Behavior**: Each time it is triggered, the next sound in the dependent sound list is played in sequence. After reaching the last item in the list, the next trigger loops back to the first item, continuing the cycle.
  * **Inheritance Behavior**: Since the dependent sounds can include **other “Premium Key Sounds”**, when the sequence reaches a “Premium Key Sound” in the list, it fully inherits all sound effects and playback modes of that “Premium Key Sound” and plays accordingly.
  * **Nesting Performance**: When the specified sounds include other “Premium Key Sounds,” one or more layers of nesting are inevitable. However, since KeyTone’s backend is developed in Go, you can confidently use multi-layer or even mutual nesting. No matter how complex your nesting combinations are, you won’t experience any perceptible lag, as Go’s high-performance fully support your “tinkering” spirit.

## **Emphasis on Inheritance**

To reiterate the inheritance characteristics of “Premium Key Sounds”:
* The "press" inherits the "press". When the "press" part of one "Premium Key Sound" inherits from another "Premium Key Sound", it only inherits the "press" part of that other "Premium Key Sound", and does not inherit its "release" part, nor does it affect the "release" part of the current "Premium Key Sound", since the "press" part of the current "Premium Key Sound" is being configured.

* Similarly, “Release” inherits “Release.”

## **Next Steps**

Through **Craft Premium Key Sounds (Step 3)**, you can create truly unique and dynamically varied key sound schemes, making every keystroke in [**Link Sounds to Key Actions (Step 4)**](../../key-package/按键联动声效/index.md) full of fun.