
# Oronoxyl

> Easily create XML sitemaps for your website.

Oronoxyl is a simple [sitemap](https://www.sitemaps.org) generator as a command line tool. It is capable of crawling web pages in a configurable number of parallel workers and it doesn't use any external dependencies.
## Table of contents
- [Install](#install)
- [Ussage](#ussage)
- [Options](#options)
- [License](#license)


## Install

Clone the code.

```bash
git clone https://github.com/Mihai22125/oronoxyl
```


Get into the source directory.

```bash
cd oronoxyl/cmd
```

Build the code from the root directory

```bash
go build -o oronoxyl
```

> This sends the output of `go build` to a file called `oronoxyl` in the same directory.

## Ussage

The crawler will fetch all links found in the <a> elements. The crawler is able to apply the `base` value to found links.

```BASH
oronoxyl [options]
```

When the crawler finished the XML Sitemap will be built and saved to your specified path.

Example:

```BASH
oronoxyl -url=http://example.com
```
## Options

```BASH
oronoxyl -help

  Usage: oronoxyl [options]

  Options:
    -max-depth   (int)                    max depth of url navigation recursion (default 3)
    -output-file (string)                 output file path (default "./temp.xml")
    -parallel    (int)                    number of parallel workers to navigate through site (Default 3)
    -url         (string)                 site url for sitemap generation
    -verbose     (bool)                   display detailed processing information (default true)
    -help        (bool)                   output usage information
```

### parallel

Sets the maximum number of requests the crawler will run simultaneously (default: 3).

### output-file

Path to file to write including the filename itself. Path can be absolute or relative. Default is `temp.xml`.

Examples:

- `sitemap.xml`
- `mymap.xml`
- `/var/www/sitemap.xml`

### maxDepth

Set a maximum distance from the original request to crawl URLs, useful for generating smaller `sitemap.xml` files. Defaults to 3.

### url

Specify the url for which the sitemap will be generated.

### verbose

Print debug messages during crawling process. Also prints out a summery when finished.

### help

Output usage information.
## License

[MIT](https://github.com/Mihai22125/oronoxyl/blob/main/LICENSE)
