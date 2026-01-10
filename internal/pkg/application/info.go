package application

var (
	Version string
	Commit  string
	BuildAt string
)

type Info struct {
	Version string
	Commit  string
	BuildAt string
}

func GetInfo() Info {
	return Info{
		Version: Version,
		Commit:  Commit,
		BuildAt: BuildAt,
	}
}
