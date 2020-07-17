# saps

*what is saps?*

*saps* is a *simple as possible sequencer*.

- *saps* sees instruments.
- *saps* reads notes.
- *saps* feels time.
- *saps* returns.

## Install

First install `portmidi`, following [these directions](https://schollz.com/blog/portmidi/).

Next install [Go](https://golang.org/dl/) and then in a terminal:

	> go install github.com/schollz/saps

That's it! `saps` is now available from the command-line.

## 

- number of measures is biggest number of measures in each

```
section a

instruments NTS-1 digital kit 1 SOUND
C5 A E4 C5 A E4 C5 A E4 C5 A E4 C5 A E4 C5 A E4 C5 A E4 C5 A E4 C5 A E4 C5 A E4 C5 A E4 C5 A E4
C5 G E4 C5 G E4 C5 G E4 C5 G E4 C5 G E4 C5 G E4 C5 G E4 C5 G E4 C5 G E4 C5 G E4 C5 G E4 C5 G E4

instruments Boutique SH-01A
A3CE  
C4EG 
A3CE  
C4EG 

section b
	

instruments Boutique SH-01A
DF#A
DF#A
DF#A
DF#A

```
