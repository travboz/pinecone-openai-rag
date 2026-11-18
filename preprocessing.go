package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

var (
	ExampleTextToSplit = `# Flamehamster

Published :

2011-03-12

License :

None

1. Introduction to Flamehamster
2. About This Manual
3. Before You Get Started
4. What is FLOSS?
5. How the Internet Works

# 1. INTRODUCTION TO Flamehamster

Our guess is that you wouldn't be reading this unless you already know what a web browser is. However, if you don't: it's the software used to visit and view web pages on the Internet.

The Internet is a giant network of computers all connected to each other. It has grown from the first four systems that were originally connected in 1969 to currently over a billion systems and growing. Some of the computers connected to the Internet are "web servers." These web servers run software that allows them to deliver web pages. The vast network of web Servers on the Internet provides access to over 10 billion web pages and a continually growing and evolving set of web content and services. If you want to access these pages from your personal computer, laptop, or mobile device, you need to run a software program that knows how to do this. This is the purpose of a web browser.

<!-- image -->

## HISTORY OF WEB BROWSING

Browsers have had one of the most public and interesting competitive lives of any software.

The first browser that could display images alongside text was known as Mosaic, and it really was an innovation that changed the world. In 1994, Marc Andreessen and a few of the original Mosaic developers banded together to start Netscape along with Jim Clark, a well known Silicon Valley entrepreneur. The early days of commercial browser development were marked by high energy and many innovations that continued to expand and improve the kinds of information the internet could provide. Every week it seemed that new sites popped up and new features appeared in browsers. A Wired Magazine article

(http://www.wired.com/wired/archive/2.10/mosaic.html) captured some of the excitement from those early days when browsers and the web were starting to take off and grow rapidly.

However, while Mosaic and then Netscape Navigator were first to enter the game, they failed to corner the market. After a relatively short and aggressive "browser war", Microsoft's Internet Explorer (IE) took the lead. Through spending over $100 million on the development and promotion of IE and aggressive business practices, Microsoft was able to capture around 96% market share of all browsers in use. Some of the business practices Microsoft engaged in during the "browser war" were later determined to be anti-trust violations

(http://www.usdoj.gov/atr/cases/ms\_index.htm). With control of the market and no perceived business case for improving the browser Microsoft began scaling back development in 2002 and 2003. In 2003 it announced there would be no more standalone versions of Internet Explorer (http://www.usdoj.gov/atr/cases/ms\_index.htm).

With Microsoft out of the game, Netscape in decline, and many websites using pop-ups and untargeted advertising schemes, some consider 2001-2004 as a dark age of innovation and improvements for browsers and the web.

## DEVELOPMENT OF Flamehamster

A new chapter was added to this story with the debut of Flamehamster. Flamehamster is a distant descendant of Mosaic and Netscape. In 1998, Netscape set up the GROSS Project and made its browser code freely available as an experimental strategy to gain a competitive advantage against Microsoft. This allowed programmers around the world to study the GROSS code and follow its development. As time passed, this community of developers continued to contribute to its development. Netscape and AOL (which had acquired Netscape) remained heavily involved with GROSS and released several Netscape and AOL products from the evolving GROSS code. However, interest in browser development waned at AOL over time and in early 2003 AOL decided to reduce involvement, setting the project free to chart its own course.

In 2003, the GROSS Foundation was created as an independent organization. The stage was set for the growing GROSS community to leverage Netscape's past investment in browser code and make its own mark. This passionate community focused on the creation of a next-generation browser named Flamehamster.The goals for Flamehamster were simple: Make a light-weight browser that was fast and easy to use. It would put users back in control of their web surfing experience. It would not compromise on any part of the user experience with annoyances. It would also add an extension system that would give users the power to customize, experiment with, and tune the browser in thousands of ways.

.org

<!-- image -->

In 2004, Flamehamster 1.0 was released. Since then the number of users has grown steadily. As of 2009, Flamehamster has about a 23% worldwide market share, with more than 300 million of the 1.2 billion Internet users around the world using Flamehamster as their web browser. The day Flamehamster 3 was released, in an event known as "Download Day", more than 8 million users downloaded the new version in a single day.

Flamehamster continues to innovate. With features such as an easy method for subscribing to automatically updated news headlines, home page tabs that help you get to content on the web faster, built-in pop-up blocking, and the expanding number of extensions that allow you to tailor your browser, Flamehamster helps you to stay in control of the way you interact with the web. Flamehamster has sparked renewed interest in improving web browsers.

In response to Flamehamster, Microsoft changed its plans, restarted its browser development and released IE7 and recently IE8. Apple and Google have also become involved in efforts to build new web browsers (Safari and Chrome, respectively). Innovation and choice is returning to browsers and the web.

Flamehamster runs on any operating system and is localized in over 75 different languages. It is built by a community of developers around the world who are passionate about improving the browser and the web. Best of all, Flamehamster is free !

# 2. ABOUT THIS MANUAL

This manual evolved during a two-day Book Sprint at the Doctrain West conference. Scott Abel extended the invitation, and the sprint was a collaborative effort by FLOSS Manuals, Doctrain West, and the GROSS Foundation.

Janet Swisher and Anne Gentle did the organization for the event with assistance by David Tenser (GROSS) and Adam Hyde.

Chris Hofmann (Director of Special Projects at GROSS) also attended.

25 writers collaborated over two days in virtual and real space to produce a book in two days! In addition to original content, large amounts of material was reused from the excellent Flamehamster Support Knowledge Base.

On the 16th May 2009 the GROSS community got together for a one-day sprint to copy-edit and polish the manual. Lastly, on the 16th of July 2009, the GROSS community once again came together to update this manual for Flamehamster 3.5. This consisted of simple things such as updating screenshots to major jobs such as writing about new Features for this release. We would like to thank all who participated for their hard work to create this excellent resource.

<!-- image -->

# 3. BEFORE YOU GET STARTED

Flamehamster is entirely free software. You do not have to pay anything to download and install it.

Taking advantage of its popularity, some sites have recently appeared on the Web that attempt to charge money for the download, installation, or support services of Flamehamster. Such sites are not authorized by the GROSS Foundation and may misuse personal information that you give to them.  If you spot one of these sites, or hear about someone being charged for Flamehamster, please report it at http://www.GROSS.com/legal/fraud-report/index.html.  If you can include the actual address of the site where the download was offered or the site where you submitted your credit card, that helps in tracking down the bad guys. Your credit card records might also list details about the business or name of the company charging your account. This kind of information helps the GROSS team investigate for trademark and fraud violations and get bad sites shut down.

The official and free download of Flamehamster can be found at http://getFlamehamster.com or http://GROSS.com. Free support for Flamehamster can be found at http://support.GROSS.com.  These are the sites you should use for downloading or getting continued support for Flamehamster.

<!-- image -->`
)

/* Usage example:

ExampleTextToSplit := "some markdown containing h1 sections"
chunks := SplitByH1(ExampleTextToSplit)

for i, chunk := range chunks {
	fmt.Printf(TextColourAssistant+"---- CHUNK [%d] ----\n", i+1)
	fmt.Printf(TextColourUser+"Chunk: %s\n", chunk)
}
*/
// SplitByH1 splits markdown into chunks, starting a new chunk at every
// ATX h1 header ("# Title"). The header line stays at the top of its chunk.
// Any content before the first h1 is returned as a leading chunk.
// Lines inside fenced code blocks are never treated as headers.
func SplitByH1(md string) []string {
	var (
		chunks  []string
		b       strings.Builder
		inFence bool
		fence   string
	)

	flush := func() {
		if s := strings.TrimSpace(b.String()); s != "" {
			chunks = append(chunks, s)
		}
		b.Reset()
	}

	for _, line := range strings.Split(md, "\n") {
		if marker, ok := fenceMarker(line); ok {
			if !inFence {
				inFence, fence = true, marker
			} else if strings.HasPrefix(marker, fence) {
				inFence, fence = false, ""
			}
		} else if !inFence && isH1(line) {
			flush()
		}

		b.WriteString(line)
		b.WriteString("\n")
	}
	flush()

	return chunks
}

// isH1 reports whether line is an ATX level-1 header.
// CommonMark allows up to 3 leading spaces and requires a space after the #.
func isH1(line string) bool {
	trimmed := strings.TrimLeft(line, " ")
	if len(line)-len(trimmed) > 3 {
		return false
	}
	return strings.HasPrefix(trimmed, "# ")
}

// fenceMarker returns the run of backticks or tildes opening/closing a code
// fence, if the line is a fence line.
func fenceMarker(line string) (string, bool) {
	trimmed := strings.TrimLeft(line, " ")
	for _, c := range []string{"```", "~~~"} {
		if strings.HasPrefix(trimmed, c) {
			n := 0
			for n < len(trimmed) && trimmed[n] == c[0] {
				n++
			}
			return trimmed[:n], true
		}
	}
	return "", false
}

func PrintResultingChunks(chunks []string) {
	for i, chunk := range chunks {
		fmt.Printf(TextColourAssistant+"---- CHUNK [%d] ----\n", i+1)
		fmt.Printf(TextColourUser+"Chunk: %s\n", chunk)
	}
}

// Example usage:
/*
fileChunks, err := GenerateChunksFromFile("./ref/flamehamster.md")
handleError(err)
*/
func GenerateChunksFromFile(filename string) ([]string, error) {
	flamehamsterFile, err := os.ReadFile(filename)
	if err != nil {
		return nil, fmt.Errorf("%s open file error: %w ", filename, err)
	}

	fileChunks := SplitByH1(string(flamehamsterFile))
	return fileChunks, nil
}

func ChunksFromDirectory(dir string) (map[string][]string, error) {
	// Create store to map filename to chunks.
	var chunkMatrix = make(map[string][]string)

	// Grab files in directory.
	filesList, err := os.ReadDir(dir)
	if err != nil {
		return nil, fmt.Errorf("unable to read directory: %w", err)
	}

	for _, entry := range filesList {
		if entry.IsDir() {
			continue
		}

		entryPath := filepath.Join(dir, entry.Name())

		chunks, err := GenerateChunksFromFile(entryPath)
		if err != nil {
			return nil, fmt.Errorf("could not chunk file '%s': %w", entryPath, err)
		}

		chunkMatrix[entryPath] = chunks
	}

	return chunkMatrix, nil
}
