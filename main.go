//line README_orig.md:66
package main

import (
//line README_orig.md:157
	"fmt"
	"os"
	"io"
//line README_orig.md:224
	"bufio"
//line README_orig.md:408
	"regexp"
//line README_orig.md:547
	"strings"
//line SubdirectoryFiles.md:37
	"path/filepath"
//line README_orig.md:70
)

//line LineNumbers.md:26
type File string
type CodeBlock []CodeLine
type BlockName string
type language string
//line LineNumbers.md:38
type CodeLine struct {
	text   string
	file   File
	lang   language
	number int
}
//line LineNumbers.md:31

var blocks map[BlockName]CodeBlock
var files map[File]CodeBlock
//line README_orig.md:427
var namedBlockRe *regexp.Regexp
//line README_orig.md:461
var fileBlockRe *regexp.Regexp
//line README_orig.md:554
var replaceRe *regexp.Regexp
//line README_orig.md:73

//line LineNumbers.md:125
// Updates the blocks and files map for the markdown read from r.
func ProcessFile(r io.Reader, inputfilename string) error {
//line LineNumbers.md:87
	scanner := bufio.NewReader(r)
	var err error
	
	var line CodeLine
	line.file = File(inputfilename)
	
	var inBlock, appending bool
	var bname BlockName
	var fname File
	var block CodeBlock
//line LineNumbers.md:105
	for {
		line.number++
		line.text, err = scanner.ReadString('\n')
		switch err {
		case io.EOF:
			return nil
		case nil:
			// Nothing special
		default:
			return err
		}
//line LineNumbers.md:154
		if inBlock {
			if strings.TrimSpace(line.text) == "```" {
//line LineNumbers.md:60
				inBlock = false
				// Update the files map if it's a file.
				if fname != "" {
					if appending {
						files[fname] = append(files[fname], block...)
					} else {
						files[fname] = block
					}
				}
				
				// Update the named block map if it's a named block.
				if bname != "" {
					if appending {
						blocks[bname] = append(blocks[bname], block...)
					} else {
						blocks[bname] = block
					}
				}
//line LineNumbers.md:157
				continue
			}
//line LineNumbers.md:51
			block = append(block, line)
//line LineNumbers.md:160
			continue
		}
//line h6_tags.md:3
		if line.text == "" {
			continue
		}
		
		switch line.text[0] {
		case '`':
//line h6_tags.md:19
			if len(line.text) >= 3 && (line.text[0:3] == "```") {
				inBlock = true
				// We were outside of a block, so just blindly reset it.
				block = make(CodeBlock, 0)
				codeStartRe := regexp.MustCompile("^`{3,}\\s?(\\w*)\\s*.*$")
				if matches := codeStartRe.FindStringSubmatch(strings.TrimSpace(line.text)); matches != nil {
					line.lang = language(matches[1])
				}
			}
//line h6_tags.md:10
		case '#':
			if len(line.text) >= 6 && (line.text[0:6] == "######") {
//line h6_tags.md:32
				fname, bname, appending = parseHeader(line.text)
//line h6_tags.md:13
			}
		}
//line LineNumbers.md:117
	}
//line LineNumbers.md:128
}
//line h6_tags.md:37
func parseHeader(line string) (File, BlockName, bool) {
	line = strings.TrimSpace(line)
//line h6_tags.md:55
	var matches []string
	if matches = fileBlockRe.FindStringSubmatch(line); matches != nil {
		return File(matches[1]), "", (strings.TrimSpace(matches[2]) == "+=")
	}
	if matches = namedBlockRe.FindStringSubmatch(line); matches != nil {
		if matches[2] == "+=" {println("Found +=", "'"+matches[1]+"'")}
		return "", BlockName(matches[1]), (strings.TrimSpace(matches[2]) == "+=")
	}
	return "", "", false
//line h6_tags.md:40
}
//line WhitespacePreservation.md:37
// Replace expands all macros in a CodeBlock and returns a CodeBlock with no
// references to macros.
func (c CodeBlock) Replace(prefix string) (ret CodeBlock) {
//line LineNumbers.md:270
	var line string
	for _, v := range c {
		line = v.text
//line LineNumbers.md:252
		matches := replaceRe.FindStringSubmatch(line)
		if matches == nil {
			if v.text != "\n" {
				v.text = prefix + v.text
			}
			ret = append(ret, v)
			continue
		}
//line LineNumbers.md:237
		bname := BlockName(matches[2])
		if val, ok := blocks[bname]; ok {
			ret = append(ret, val.Replace(prefix+matches[1])...)
		} else {
			fmt.Fprintf(os.Stderr, "Warning: Block named %s referenced but not defined.\n", bname)
			ret = append(ret, v)
		}
//line LineNumbers.md:274
	}
	return
//line WhitespacePreservation.md:41
}
//line LineNumbers.md:298

// Finalize extract the textual lines from CodeBlocks and (if needed) prepend a
// notice about "unexpected" filename or line changes, which is extracted from
// the contained CodeLines. The result is a string with newlines ready to be
// pasted into a file.
func (c CodeBlock) Finalize() (ret string) {
	var file File
	var formatstring string
	var linenumber int
	for _, l := range c {
		if linenumber+1 != l.number || file != l.file {
			switch l.lang {
			case "go", "golang":
				formatstring = "//line %[2]v:%[1]v\n"
			case "C", "c":
				formatstring = "#line %v \"%v\"\n"
			}
			ret += fmt.Sprintf(formatstring, l.number, l.file)
		}
		ret += l.text
		linenumber = l.number
		file = l.file
	}
	return
}
//line README_orig.md:75

func main() {
//line README_orig.md:166
	// Initialize the maps
	blocks = make(map[BlockName]CodeBlock)
	files = make(map[File]CodeBlock)
//line h6_tags.md:45
	namedBlockRe = regexp.MustCompile("^###### ([^+=]+)(\\s+[+][=])?$")
//line h6_tags.md:50
	fileBlockRe = regexp.MustCompile("^###### file:([\\w\\.\\-\\/]+)(\\s+[+][=])?$")
//line WhitespacePreservation.md:12
	replaceRe = regexp.MustCompile(`^([\s]*)<<<(.+)>>>[\s]*$`)
//line README_orig.md:143
	
	// os.Args[0] is the command name, "lmt". We don't want to process it.
	for _, file := range os.Args[1:] {
//line LineNumbers.md:135
		f, err := os.Open(file)
		if err != nil {
			fmt.Fprintln(os.Stderr, "error: ", err)
			continue
		}
		
		if err := ProcessFile(f, file); err != nil {
			fmt.Fprintln(os.Stderr, "error: ", err)
		}
		// Don't defer since we're in a loop, we don't want to wait until the function
		// exits.
		f.Close()
//line README_orig.md:147
	
	}
//line LineNumbers.md:330
	for filename, codeblock := range files {
		if dir := filepath.Dir(string(filename)); dir != "." {
			if err := os.MkdirAll(dir, 0775); err != nil {
				fmt.Fprintf(os.Stderr, "%v\n", err)
			}
		}
	
		f, err := os.Create(string(filename))
		if err != nil {
			fmt.Fprintf(os.Stderr, "%v\n", err)
			continue
		}
		fmt.Fprintf(f, "%s", codeblock.Replace("").Finalize())
		// We don't defer this so that it'll get closed before the loop finishes.
		f.Close()
	}
//line README_orig.md:78
}
