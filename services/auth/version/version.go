package version

var (
	Tag       = ""
	GitCommit = ""
)

type Version struct {
	Tag       string
	GitCommit string
}

func GetVersion() Version {
	return Version{
		Tag:       Tag,
		GitCommit: GitCommit,
	}
}
