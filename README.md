# idim

*idim* is for *interfacing different instruments' midi*. It provides surprisingly simple sequencing for synthesizers or other instruments.

## Install

First install `portmidi`, following [these directions](https://schollz.com/blog/portmidi/).

Next install [Go](https://golang.org/dl/) and then in a terminal:

	> go install github.com/schollz/idim

That's it! `idim` is now available from the command-line.

## Usage

### First steps

To get started, first plugin your instruments to your computer. Open a command prompt and type `idim` to see which instruments are available to you.

```
> idim
+---------------------------+
|        INSTRUMENTS        |
+---------------------------+
| NTS-1 digital kit 1 SOUND |
+---------------------------+
```

You can use these instruments to build and chain patterns of notes.

Modify an example in the `examples` to make sure its set to the instrument that you have. Then to run, you can just do

```
> idim --file examples/song1.idim
[info]  2020/07/17 08:18:12 playing
```

And you'll hear some music!

## idim musical notation

*idim* reads a `.idim` file, which is a high-level musical notation developed for *idim*. The musical notation is simple and powerful, allowing you to create patterns of notes that can be played on many instruments simultaneously.

### Basic pattern

The basic unit is the *pattern*. A *pattern* contains a collection of *instruments*. Each *instrument* contains a collection of notes.
For example, here is a simple pattern that plays Cmaj followed by Fmaj and then repeats.

```bash
pattern a
instruments <instrument1>
CEG
FAC
```

The `pattern a` designates the pattern name, "a". This pattern has a single instrument, `<instrument1>` 
(normally you will fill in the name of an instrument, like `NTS-1 digital kit 1 SOUND` in the example above). 

Each line under `instruments` designates a different measure. This is where you put notes. Notes without spaces are considered a chord and will be played simultaneously. The first measure plays C, E, and G (C major) and the second measure plays F, A and C (F major). This pattern will repeat indefinitely when played.

To add more patterns, simply add a line with `pattern X` and again add `instruments` and their notes. The patterns will be played in order.

### Add instruments and subdivisions

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
Since each note is separated by a space, they are not played together as a chord (unlike in instrument1) and are automatically subdivided according to the number of notes in that measure. In this case they are subdivided into 1/8th notes since there are eight of them in that measure. Since there is only one measure for the `instrument2`, it will repeat over every chord of `instrument1`.

You can add rests in using `.` to create specific subdivisions. For example, we can change instrument1 to play chords on the off-beat of beat 1 and beat 2:

```
instruments <instrument1>
. CEG . . . CEG . . 
. FAC . . . FAC . . 
```


### Other specifications

#### Specific notes

By default, the note played will be the note closest to the previous. If you want to specify the exact note you can add a suffix to include the octave. For example, instead of writing `CEG` you could instead write `C3E4G5` which will span the chord over three octaves.


Here are other keywords you can use to modulate the song in the `.idim` file:

#### Setting the tempo

You can add a line to change the tempo, anywhere in the file.

```
tempo <10-300>
````

#### Changing the legato

The legato specifies how much to hold each note until releasing it. Full legato (100) holds to the very end, while the shortest legato (1) will release immediately after playing.

```
legato <1-100>
```

#### Sustain 

For a pedal note (sustain) add a `*` to the end of the note. For example, the following will sustain a C major chord for two measures:

```
CEG* 
CEG
```

#### Multiple instruments

You can assign multiple instruments to a pattern by separating each instrument by a comma. For example:

```
instruments <instrument1>, <instrumnet2>
C E G
```

will play the C, E, G arpeggio on both instruments 1 and 2.



## To Do

- [x] Add legato control `legato: 90`
- [x] Hot-reload file
- [x] in midi, create a channel for each instrument
- [x] in midi, each instrument keeps track of which notes are off
- [x] in midi, accept -1 to turn off all notes 
- [x] in midi, accept -2 to turn off all notes and shut down
- [x] Add `*` suffix for adding sustain
- [ ] Allow chaining patterns in different ways `chain: a a b b a a`
- [ ] Understand chords `Bmin Gmaj`



## License 

MIT
