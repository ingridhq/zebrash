package encoder

type Options struct {
	Mode             *Mode
	CharacterSetName string
	AppendGS1        bool
	VersionNumber    int
	MaskPattern      *int
	QuietZone        int
}
