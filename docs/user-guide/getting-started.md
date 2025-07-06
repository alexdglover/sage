# Getting Started

Sage is a local-only application. No account or login is required. All your data stays on your device.

## Installation

### Installation walkthrough video

[![Installation Walkthrough Video](https://img.youtube.com/vi/lfyVE-FLz_g/0.jpg)](https://www.youtube.com/watch?v=lfyVE-FLz_g)

Installation is easy, except Apple made it a bit harder because of how [application signing](https://support.apple.com/guide/security/app-code-signing-process-sec3ad8e6e53/web)
is enforced on newer versions of MacOS.

### For Apple/Mac users:

1. Create a new folder where you want to run Sage - `mkdir ~/sage && cd ~/sage` for example
1. Browse to `https://github.com/alexdglover/sage/releases/latest`
1. Select the right download for your computer, based on what operating system you use and whether
your computer uses an Intel/x86 CPU or an ARM CPU.
1. Right click on the link and copy the URL
1. Download the file - `curl -L https://github.com/alexdglover/sage/releases/download/v1.1.2/sage_Darwin_arm64.tar.gz -o sage.tar.gz`
    * Downloading the file this way bypasses MacOS' application signing requirement. You can't
    just click the link in your browser
1. Unpack the archive - `tar xvzf sage.tar.gz`

### For Windows and Linux users

1. Create a new folder where you want to run Sage
1. Browse to `https://github.com/alexdglover/sage/releases/latest`
1. Select the right download for your computer, based on what operating system you use and whether
your computer uses an Intel/x86 CPU or an ARM CPU.
1. Download the file (either via your browser or CLI) and unpack the archive

### Launching the app

You can either launch the app from a CLI (`~/sage/sage`) or double click the icon. Launching the
application will create a `sage.db` file in the same directory. All of your financial data is
stored locally in the `sage.db` file. Remember to back up this file periodically to avoid data loss.

## Main Workflows
- **Add new accounts**: Set up your financial accounts in Sage.
- **Import statements**: Bring in your transaction history using CSV files.
- **Update balances**: Keep your account balances current.
- **Categorize transactions**: Organize your spending for better insights.
- **View reports**: Analyze your finances with built-in reports.
