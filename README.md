# webp-to-jpg-converter
Golang service to convert .webp files to .jpg format

### Usage example

```bash
curl -X POST -F "file=@samples/meme.webp" http://localhost:3000/convert -o output.jpg
```