###### Handle nonblock line
```go
if line.text == "" {
	continue
}

switch line.text[0] {
case '`':
	<<<Check block start>>>
case '#':
	if len(line.text) >= 6 && (line.text[0:6] == "######") {
		<<<Check block header>>>
	}
}
```

###### Check block start
```go
if len(line.text) >= 3 && (line.text[0:3] == "```") {
	inBlock = true
	// We were outside of a block, so just blindly reset it.
	block = make(CodeBlock, 0)
	codeStartRe := regexp.MustCompile("^`{3,}\\s?(\\w*)\\s*.*$")
	if matches := codeStartRe.FindStringSubmatch(strings.TrimSpace(line.text)); matches != nil {
		line.lang = language(matches[1])
	}
}
```

###### Check block header
```go
fname, bname, appending = parseHeader(line.text)
```

###### ParseHeader Declaration
```go
func parseHeader(line string) (File, BlockName, bool) {
	line = strings.TrimSpace(line)
	<<<parseHeader implementation>>>
}
```

###### Namedblock Regex
```go
namedBlockRe = regexp.MustCompile("^###### ([^+=]+)(\\s+[+][=])?$")
```

###### Fileblock Regex
```go
fileBlockRe = regexp.MustCompile("^###### file:([\\w\\.\\-\\/]+)(\\s+[+][=])?$")
```

###### parseHeader implementation
```go
var matches []string
if matches = fileBlockRe.FindStringSubmatch(line); matches != nil {
	return File(matches[1]), "", (strings.TrimSpace(matches[2]) == "+=")
}
if matches = namedBlockRe.FindStringSubmatch(line); matches != nil {
	if matches[2] == "+=" {println("Found +=", "'"+matches[1]+"'")}
	return "", BlockName(matches[1]), (strings.TrimSpace(matches[2]) == "+=")
}
return "", "", false
```
