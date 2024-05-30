package main

import (
	"archive/zip"
	"bytes"
	"flag"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"sync"

	"golang.org/x/net/html"
	"golang.org/x/net/html/atom"
)

func listFiles(dir string, recursive bool) []string {
	files := []string{}
	filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() && strings.HasSuffix(info.Name(), ".epub") {
			files = append(files, path)
		}
		if !recursive && info.IsDir() && path != dir {
			return filepath.SkipDir
		}
		return nil
	})
	return files
}

func parseNCX(content []byte) []string {
	htmlNode, err := html.Parse(strings.NewReader(string(content)))
	if err != nil {
		fmt.Println("Error parsing NCX file:", err)
		return nil
	}

	var htmlList []string
	var f func(*html.Node)
	f = func(n *html.Node) {
		if n.Type == html.ElementNode && n.DataAtom == atom.Content {
			for _, a := range n.Attr {
				if a.Key == "src" {
					htmlList = append(htmlList, a.Val[3:])
				}
			}
		}
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			f(c)
		}
	}
	f(htmlNode)
	return htmlList
}

func parseHTMLForImage(content []byte) string {
	htmlNode, err := html.Parse(strings.NewReader(string(content)))
	if err != nil {
		fmt.Println("Error parsing HTML file:", err)
		return ""
	}

	var imgURL string
	var f func(*html.Node)
	f = func(n *html.Node) {
		if n.Type == html.ElementNode && n.DataAtom == atom.Img {
			for _, a := range n.Attr {
				if a.Key == "src" {
					imgURL = a.Val[3:]
					break
				}
			}
		}
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			f(c)
		}
	}
	f(htmlNode)
	return imgURL
}

func processFile(file string) error {
	fmt.Println("Processing file:", file)
	r, err := zip.OpenReader(file)
	if err != nil {
		fmt.Println("Error opening zip file:", err)
		return err
	}
	defer r.Close()

	//var volOpf, volNcx string
	var volNcx string
	for _, f := range r.File {
		if strings.HasSuffix(f.Name, ".opf") {
			// volOpf = f.Name
		}
		if strings.HasSuffix(f.Name, ".ncx") {
			volNcx = f.Name
		}
	}

	var htmlList []string
	for _, f := range r.File {
		if f.Name == volNcx {
			rc, err := f.Open()
			if err != nil {
				fmt.Println("Error opening NCX file:", err)
				break
			}
			content, err := io.ReadAll(rc)
			if err != nil {
				fmt.Println("Error reading NCX file:", err)
				rc.Close()
				break
			}
			rc.Close()

			htmlList = parseNCX(content)
			break
		}
	}

	num := 1
	var buffer bytes.Buffer
	zipWriter := zip.NewWriter(&buffer)

	for _, f := range r.File {
		if f.Name == "image/cover.jpg" {
			rc, err := f.Open()
			if err != nil {
				fmt.Println("Error opening cover image:", err)
				break
			}
			w, err := zipWriter.Create("cover.jpg")
			if err != nil {
				fmt.Println("Error creating cover image in zip:", err)
				rc.Close()
				return err
			}
			_, err = io.Copy(w, rc)
			if err != nil {
				fmt.Println("Error copying cover image:", err)
				return err
			}
			rc.Close()
			break
		}
	}

	for _, htmlFile := range htmlList {
		for _, f := range r.File {
			if f.Name == htmlFile {
				rc, err := f.Open()
				if err != nil {
					fmt.Println("Error opening HTML file:", err)
					return err
				}
				defer rc.Close()
				content, err := io.ReadAll(rc)
				if err != nil {
					fmt.Println("Error reading HTML file:", err)
					return err
				}

				imgURL := parseHTMLForImage(content)
				format := filepath.Ext(imgURL)
				var name string

				switch {
				case strings.Contains(imgURL, "cover"):
					name = "cover"
				case strings.Contains(imgURL, "createby"):
					name = "createby"
				default:
					name = fmt.Sprintf("%04d", num)
					num++
				}

				for _, imgFile := range r.File {
					if imgFile.Name == imgURL {
						rc, err := imgFile.Open()
						if err != nil {
							fmt.Println("Error opening image file:", err)
							return err
						}
						w, err := zipWriter.Create(name + format)
						if err != nil {
							fmt.Println("Error creating image file in zip:", err)
							return err
						}
						_, err = io.Copy(w, rc)
						if err != nil {
							fmt.Println("Error copying image file:", err)
							return err
						}
						break
					}
				}
				break
			}
		}
	}

	err = zipWriter.Close()
	if err != nil {
		return err
	}

	outputFile, err := os.Create(strings.TrimSuffix(file, filepath.Ext(file)) + "_images.zip")
	if err != nil {
		return err
	}
	defer outputFile.Close()

	_, err = buffer.WriteTo(outputFile)
	if err != nil {
		return err
	}

	return nil
}

func main() {
	recursive := flag.Bool("r", false, "recursively search directories for epub files")
	deleteOriginal := flag.Bool("d", false, "delete original epub files after processing")
	flag.Parse()

	files := listFiles(".", *recursive)
	var wg sync.WaitGroup

	for _, file := range files {
		wg.Add(1)
		go func(file string) {
			defer wg.Done()
			err := processFile(file)
			if err != nil {
				fmt.Printf("Error processing file %s: %v\n", file, err)
			} else if *deleteOriginal {
				err = os.Remove(file)
				if err != nil {
					fmt.Printf("Error deleting original file %s: %v\n", file, err)
				}
			}
		}(file)
	}

	wg.Wait()
}
