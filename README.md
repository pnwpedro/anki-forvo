# anki-forvo

VSCode launch config looks like:
    "configurations": [
        {
            "name": "Launch Package",
            "type": "go",
            "request": "launch",
            "mode": "auto",
            "program": "${fileDirname}",
            "env": { "FORVO_API_KEY": "YOUR_API_KEY_HERE" }
        }