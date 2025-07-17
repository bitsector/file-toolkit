# Text Analysis API

A FastAPI application for text analysis with OCR capabilities using pytesseract.

## Features

- **Text Extraction**: Extract text from images using OCR
- **Confidence Scoring**: Get confidence levels for extracted text
- **Language Detection**: Detect the language of text in images
- **Pattern Matching**: Check if extracted text contains specific patterns

## Prerequisites

- Python 3.9+ (tested with Python 3.13.2)
- Poetry for dependency management
- Tesseract OCR engine (required by pytesseract)

## Installation

### 1. Install Poetry

If Poetry is not already installed:

```bash
curl -sSL https://install.python-poetry.org | python3 -
```

Add Poetry to your PATH:
```bash
export PATH="$HOME/.local/bin:$PATH"
```

Verify installation:
```bash
poetry --version
```

### 2. Install Dependencies

Navigate to the project directory:
```bash
cd /path/to/file-toolkit/cmd/text_analysis
```

Install all dependencies:
```bash
poetry install
```

This will create a virtual environment and install:
- FastAPI
- Uvicorn (ASGI server)
- pytesseract (OCR library)
- Pillow (image processing)
- python-multipart (file upload support)
- Development tools (pytest, black, flake8)

### 3. Install Tesseract OCR Engine

**Ubuntu/Debian:**
```bash
sudo apt-get update
sudo apt-get install tesseract-ocr
```

**macOS:**
```bash
brew install tesseract
```

**Windows:**
Download and install from: https://github.com/UB-Mannheim/tesseract/wiki

## Running the Application

### Method 1: Using Poetry (Recommended)

From the project directory:
```bash
poetry run python main.py
```

### Method 2: Direct Python Execution

Get the virtual environment Python path:
```bash
poetry env info --executable
```

Run with full path:
```bash
/home/ak/.cache/pypoetry/virtualenvs/text-analysis-api-tXzkVpFG-py3.13/bin/python main.py
```

### Method 3: Activate Virtual Environment

```bash
poetry shell
python main.py
```

## Setting up VS Code

### 1. Get the Python Interpreter Path

Run this command in your terminal:
```bash
poetry env info --executable
```

This will output something like:
```
/home/ak/.cache/pypoetry/virtualenvs/text-analysis-api-tXzkVpFG-py3.13/bin/python
```

### 2. Configure VS Code Python Interpreter

**Option A: Command Palette**
1. Open VS Code in the project directory
2. Press `Ctrl+Shift+P` (or `Cmd+Shift+P` on macOS)
3. Type "Python: Select Interpreter"
4. Click on "Python: Select Interpreter"
5. Click "Enter interpreter path..."
6. Paste the path from step 1
7. Press Enter

**Option B: Settings UI**
1. Open VS Code settings (`Ctrl+,` or `Cmd+,`)
2. Search for "python default interpreter"
3. Set "Python › Default Interpreter Path" to the path from step 1

**Option C: Workspace Settings**
1. Create `.vscode/settings.json` in your project root:
```json
{
    "python.defaultInterpreterPath": "/home/ak/.cache/pypoetry/virtualenvs/text-analysis-api-tXzkVpFG-py3.13/bin/python"
}
```

### 3. Verify Setup

1. Open a Python file in VS Code
2. Check the bottom-left status bar - it should show the Python version and virtual environment name
3. Open a terminal in VS Code (`Ctrl+``) - it should automatically activate the virtual environment

## API Endpoints

Once running, the API will be available at `http://localhost:8000`

### Available Endpoints:

- `GET /` - Root endpoint
- `GET /health` - Health check
- `POST /contains/` - Check if text contains specific patterns *(coming soon)*
- `POST /extract/` - Basic text extraction *(coming soon)*
- `POST /extract-with-confidence/` - OCR with confidence scores *(coming soon)*
- `POST /detect-language/` - Language detection *(coming soon)*

### API Documentation

- **Interactive docs**: http://localhost:8000/docs
- **ReDoc**: http://localhost:8000/redoc
- **OpenAPI JSON**: http://localhost:8000/openapi.json

## Development

### Code Formatting
```bash
poetry run black .
```

### Linting
```bash
poetry run flake8 .
```

### Running Tests
```bash
poetry run pytest
```

### Adding Dependencies
```bash
poetry add package_name
```

### Development Dependencies
```bash
poetry add --group dev package_name
```

## Virtual Environment Management

### Get environment info:
```bash
poetry env info
```

### Get environment path:
```bash
poetry env info --path
```

### Get Python executable path:
```bash
poetry env info --executable
```

### Remove environment:
```bash
poetry env remove python
```

### List environments:
```bash
poetry env list
```

## Troubleshooting

### Poetry not found
Make sure Poetry is in your PATH:
```bash
export PATH="$HOME/.local/bin:$PATH"
```

### Virtual environment issues
Remove and recreate the environment:
```bash
poetry env remove python
poetry install
```

### Tesseract not found
Install tesseract-ocr system package and ensure it's in your PATH.

### VS Code not using correct interpreter
1. Restart VS Code
2. Check that the interpreter path is correct
3. Try manually selecting the interpreter again

## Project Structure

```
cmd/text_analysis/
├── main.py              # FastAPI application
├── pyproject.toml       # Poetry configuration
├── poetry.lock          # Locked dependencies
└── README.md           # This file
```

## License

This project is part of the file-toolkit repository.
