# saps

*what is saps?*

- *saps* sees instruments.
- *saps* reads notes.
- *saps* feels time.
- *saps* returns.

*saps* is a *simple as possible sequencer*.

## Install

First install `portmidi`, following [these directions](https://schollz.com/blog/portmidi/).

Next install [Go](https://golang.org/dl/) and then in a terminal:

	> go install github.com/schollz/saps

That's it! `saps` is now available from the command-line.

## Usage

*saps* reads a `.saps` file, which is a collection of sections, instruments, and notes.

If you run `saps` you will see which instruments are available to you.

```
> saps
Available instruments:

***********
* NTS-1   *
***********
```

You can use these instruments to build sequences.

Sequences are made of chained patterns. You can define a pattern, the instruments and the notes. For example, here is a simple pattern that plays Cmaj followed by Fmaj and then repeats.

```bash
pattern a
instruments <instrument1>
CEG
FAC
```

The `pattern a` designates the pattern "a". This pattern has a single instrument, `<instrument1>` (normally you will fill in the name of an instrument that matches one given. Each line under `instruments` designates a different measure. The first measure plays Cmaj (C, E, and G) and the second measure plays Fmaj (F, A, and C).

You can easily add a second instrument by just including another section.

```bash
pattern a 
instruments <instrument1>
CEG 
FAC
instruments <instrument2>
A F E C A F E C
```

This second instrument will play arpeggios. It consists of a single repeated measure which eight notes. Since each note is separated by a space, they are not played together (unlike in instrument1) and are automatically subidivided according to how many their are. In this case they are subdivided into 1/8th notes since there are eight of them.

You can add rests in using `.` to keep a subdivision. For example, we can change instrument1 to play chords on the off-beat:

```
instruments <instrument1>
. CEG . . . CEG . . .
. FAC . . . FAC . . .
```

By default, the chords will be composed of the nearest cluster of notes. If you want to specify the exact note you can add a suffix to include the octave. For example, `CEG` could instead be `C3E4G5` which will span the chord over three octaves.

To add more patterns, simply add a line with `pattern X` and again add `instruments` and their notes. The patterns will be played in order.

## To Do

- [ ] Allow chaining patterns in different ways `chain: a a b b a a`
- [ ] Add legato control `legato: 90`
- [ ] Understand chords `Bmin Gmaj`

## License 

MIT
