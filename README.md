# file-toolbox
Golang service for file processing operations including WebP to JPG conversion

### Usage example

```bash
curl -X POST -F "file=@samples/meme.webp" http://localhost:3000/convert -o output.jpg
```