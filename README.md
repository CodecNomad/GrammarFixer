# GrammarFixer

GrammarFixer is a lightweight, command-line tool that corrects the grammar of text in your clipboard using the Gemini API. It's designed to be seamlessly integrated into your workflow by binding it to a hotkey. Simply copy a piece of text, press your designated key combination, and the grammatically improved version will be ready to paste.

## Features

*   **Clipboard-based Grammar Correction**: Instantly fixes the grammar of any text you have copied.
*   **Hotkey Integration**: Intended for use with a global hotkey for quick and easy access.
*   **Customizable AI Control**: Dynamically alter the AI's behavior with in-text tags.

## Prerequisites

Before you begin, ensure you have the following installed on your system:

*   `Go`
*   `wl-clipboard`
*   `notify-send`

## Installation

1.  **Obtain a Gemini API Key**:
    *   Visit [Google AI Studio](https://makersuite.google.com/app/apikey) to generate your free API key.

2.  **Configure Environment Variable**:
    *   Set the `GEMINI_APIKEY` environment variable to your newly generated key. You can add this to your shell's configuration file (e.g., `.bashrc`, `.zshrc`) for persistence.
        ```bash
        export GEMINI_APIKEY="YOUR_API_KEY_HERE"
        ```

3.  **Build and Install**:
    *   Clone the repository and build the application using the following commands:
        ```bash
        go build
        sudo mv GrammarFixer /usr/local/bin/
        ```

## Usage

To use GrammarFixer, bind the `GrammarFixer` command to a hotkey in your desktop environment or window manager's settings.

Once configured, you can:

1.  Copy any text to your clipboard.
2.  Press your designated hotkey.
3.  A notification will appear to confirm that the text has been processed.
4.  Paste the corrected text.

## On-the-Fly AI Customization

You can control the AI model's behavior by embedding special tags directly within your copied text. These tags are automatically removed from the final output.

*   `[QUALITY]`: Employs a more advanced (and slower) model for higher-quality results.
*   `[FAST]`: Utilizes a faster (though less sophisticated) model for quicker corrections.
*   `[INST=your instruction]`: Provide a specific instruction to the AI. For example, `[INST=be more formal]`.
*   `[STYLE=desired style]`: Dictate a particular writing style. For instance, `[STYLE=like a 19th-century novelist]`.

**Example:**

```
[QUALITY] [INST=Translate to German] This is a test.
```

## Future Development

*   **Unit Tests**: Implement a comprehensive suite of unit tests to ensure code quality and reliability.
*   **Enhanced Configuration**: Introduce a configuration file for more advanced settings and customization.
