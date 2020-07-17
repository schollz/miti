# saps

*what is saps?*

- *saps* sees instruments.
- *saps* speaks MIDI.
- *saps* reads notes.
- *saps* feels time.
- *saps* returns.
- *saps* is a *simple as possible sequencer*.

## Install

First install `portmidi`, following [these directions](https://schollz.com/blog/portmidi/).

Next install [Go](https://golang.org/dl/) and then in a terminal:

	> go install github.com/schollz/saps

That's it! `saps` is now available from the command-line.

## Usage

### First steps

To get started, first plugin your instruments to your computer. If you run `saps` you will see which instruments are available to you.

```
> saps
Available instruments:

***********
* NTS-1   *
***********
```

You can use these instruments to build and chain patterns of notes.

### Musical notation

*saps* reads a `.saps` file, which is a high-level musical notation developed for *saps*. The musical notation is simple and powerful, allowing you to create patterns of notes that can be played on many instruments simultaneously.

#### Basic pattern

The basic unit is the *pattern*. A *pattern* contains a collection of *instruments*. Each *instrument* contains a collection of notes.
For example, here is a simple pattern that plays Cmaj followed by Fmaj and then repeats.

```bash
pattern a
instruments <instrument1>
CEG
FAC
```

The `pattern a` designates the pattern name, "a". This pattern has a single instrument, `<instrument1>` 
(normally you will fill in the name of an instrument that matches one given above). 

Each line under `instruments` designates a different measure. This is where you put notes. Notes without spaces are considered a chord and will be played simultaneously. The first measure plays C, E, and G (C major) and the second measure plays F, A and C (F major). This pattern will repeat indefinetly when played.

To add more patterns, simply add a line with `pattern X` and again add `instruments` and their notes. The patterns will be played in order.

#### Add instruments and subdivisions

You can easily add a second instrument to this section by adding another line with the instrument name:

```bash
pattern a 
instruments <instrument1>
CEG 
FAC
instruments <instrument2>
A F E C A F E C
```

This second instrument will play arpeggios. 
It consists of a single repeated measure which eight notes. 
Since each note is separated by a space, they are not played together as a chord (unlike in instrument1) and are automatically subidivided according to the number of notes in that measure. In this case they are subdivided into 1/8th notes since there are eight of them in that measure. Since there is only one measure for the `instrument2`, it will repeat over every chord of `instrument1`.

You can add rests in using `.` to create specific subdivisions. For example, we can change instrument1 to play chords on the off-beat of beat 1 and beat 2:

```
instruments <instrument1>
. CEG . . . CEG . . 
. FAC . . . FAC . . 
```

#### Specific notes

By default, the note played will be the note closest to the previous. If you want to specify the exact note you can add a suffix to include the octave. For example, instead of writing `CEG` you could instead write `C3E4G5` which will span the chord over three octaves.


## To Do

- [ ] Allow chaining patterns in different ways `chain: a a b b a a`
- [ ] Add legato control `legato: 90`
- [ ] Understand chords `Bmin Gmaj`

## License 

MIT
