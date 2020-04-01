package melrose

const (
	Sharp = iota
	Flat
	Eight
	Quarter
	Half
	Whole
	Dot

	Left
	Right

	PrintAsSpecified

	Major
	Minor // Natural
	HarmonicMinor
	MelodicMinor

	// https://en.wikipedia.org/wiki/Chord_(music)#Common_types_of_chords
	MajorTriad
	MajorSeventh
	Triad
	Seventh

	Ground
	Inversion1
	Inversion2
	Inversion3

	// https://nl.wikipedia.org/wiki/Dynamiek_(muziek)
	Pianissimo // ---
	Piano      // --
	MezzoPiano // -
	MezzoForte // +
	Forte      // ++
	Fortissimo // +++
)

const (
	F_Pianissimo = 0.2 // ---
	F_Piano      = 0.4 // --
	F_MezzoPiano = 0.8 // -
	F_MezzoForte = 1.2 // +
	F_Forte      = 1.4 // ++
	F_Fortissimo = 1.8 // +++
)

const (
	// NonRestNoteNames maps a tone to an index (C=0)
	NonRestNoteNames = "CDEFGAB"
)
