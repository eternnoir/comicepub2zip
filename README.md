# comicepub2zip

This project provides a utility to convert EPUB files containing comics into ZIP files, making them easier to read with mainstream comic reader software. It supports the VOX/MOE comic websites and reorganizes and restores shuffled filenames.

## Features

- Convert EPUB files to ZIP archives.
- Recursively search directories for EPUB files.
- Extract and rename images from EPUB files based on VOX/MOE comic websites.
- Optionally delete the original EPUB files after processing.
- **Parallel processing for fast conversion of multiple files.**

## Requirements

- Go 1.16 or later

## Installation

1. Clone the repository:
    ```sh
    git clone https://github.com/eternnoir/comicepub2zip.git
    cd comicepub2zip
    ```

2. Build the project:
    ```sh
    go build -o comicepub2zip main.go
    ```

## Usage

```sh
./comicepub2zip [options]
```

### Options

- `-r` : Recursively search directories for EPUB files.
- `-d` : Delete original EPUB files after processing.

### Examples

1. Convert all EPUB files in the current directory:
    ```sh
    ./comicepub2zip
    ```

2. Recursively search directories for EPUB files and convert them:
    ```sh
    ./comicepub2zip -r
    ```

3. Convert all EPUB files and delete the original files after processing:
    ```sh
    ./comicepub2zip -d
    ```

4. Recursively search directories for EPUB files, convert them, and delete the original files:
    ```sh
    ./comicepub2zip -r -d
    ```

## Contributing

Contributions are welcome! Please open an issue or submit a pull request.

## License

This project is licensed under the MIT License. See the [LICENSE](LICENSE) file for details.

## Acknowledgements

This project uses the following Go packages:

- [archive/zip](https://golang.org/pkg/archive/zip/)
- [golang.org/x/net/html](https://pkg.go.dev/golang.org/x/net/html)

This project was inspired by [Kox-Moe-Epub-To-Zip](https://github.com/Dean-Zheng/Kox-Moe-Epub-To-Zip). Special thanks to Dean-Zheng for the initial idea and implementation.

## Contact

If you have any questions or suggestions, feel free to open an issue or contact the maintainer.


