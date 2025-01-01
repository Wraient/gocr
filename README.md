# GoOCR

GoOCR is a simple yet effective Optical Character Recognition (OCR) application built in Go. With GoOCR, you can select an image file, extract text from it, and interact with the extracted text. The app lets you click on text boxes to copy specific pieces or view all the extracted text below the image. While not 100% accurate, it provides a user-friendly OCR experience.

## Demo

![image](https://github.com/user-attachments/assets/98c6aac4-da2f-4a53-8ca3-db2afeb6554c)

## Features
- Open an image file for OCR scanning.
- Extract text using Tesseract OCR.
- Click on specific text boxes to copy text to the clipboard.
- View all extracted text in a scrollable section below the image.
- Requires Tesseract OCR engine for operation.

## Requirements
1. [Tesseract OCR](https://github.com/tesseract-ocr/tesseract)
2. GoOCR binary: Download the **latest release** from [GoOCR Releases](https://github.com/Wraient/gocr/releases/latest).

## Installing and Setup

> **Note**: GoOCR requires Tesseract to work. Installation instructions for Tesseract on various Linux distributions are included below.

### Linux
<details>
<summary>Arch Linux / Manjaro (AUR-based systems)</summary>

```bash
sudo pacman -S tesseract tesseract-data-eng
curl -Lo gocr https://github.com/Wraient/gocr/releases/latest/download/gocr
chmod +x gocr
sudo mv gocr /usr/bin/
```
</details>

<details>
<summary>Debian / Ubuntu (and derivatives)</summary>

```bash
sudo apt update
sudo apt install tesseract-ocr
curl -Lo gocr https://github.com/Wraient/gocr/releases/latest/download/gocr
chmod +x gocr
sudo mv gocr /usr/bin/
```
</details>

<details>
<summary>Fedora Installation</summary>

```bash
sudo dnf update
sudo dnf install tesseract
curl -Lo gocr https://github.com/Wraient/gocr/releases/latest/download/gocr
chmod +x gocr
sudo mv gocr /usr/bin/
```
</details>

<details>
<summary>openSUSE Installation</summary>

```bash
sudo zypper refresh
sudo zypper install tesseract
curl -Lo gocr https://github.com/Wraient/gocr/releases/latest/download/gocr
chmod +x gocr
sudo mv gocr /usr/bin/
```
</details>

<details>
<summary>Generic Linux Installation</summary>

```bash
sudo <package-manager> install tesseract
curl -Lo gocr https://github.com/Wraient/gocr/releases/latest/download/gocr
chmod +x gocr
sudo mv gocr /usr/bin/
```
</details>

## Usage

### Options

| Flag                      | Description                                                             |
|---------------------------|-------------------------------------------------------------------------|
| `-i [Image path]`         | Returns extracted text from image in terminal                           |
| `-g [Image path]`         | Opens gui with specified image                                          |

## Uninstallation

To remove GoOCR:

```bash
sudo rm /usr/bin/gocr
```
