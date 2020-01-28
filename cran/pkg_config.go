package cran


type PkgConfig interface {
	GetOrigin() RepoURL // Returns repo ordinarily, but can return other things
	SetOrigin(string,  string)
	GetSourceType() string
	SetSourceType(string)
	GetSourceType2() SourceType
	IsDownloadable() bool
}

//PkgConfigImpl stores configuration information about a given package
type PkgConfigImpl struct {
	Repo RepoURL
	Type SourceType
}

func (pc *PkgConfigImpl) GetOrigin() RepoURL {
	return pc.Repo
}

func (pc *PkgConfigImpl) SetOrigin(name, url string) {
	pc.Repo = RepoURL{URL: url, Name: name}
}

func (pc *PkgConfigImpl) GetSourceType() string {
	return pc.Type.String()
}

func (pc *PkgConfigImpl) GetSourceType2() SourceType {
	return pc.Type
}

func (pc *PkgConfigImpl) SetSourceType(st string) {
	var t SourceType
	switch st {
	case "source":
		t = Source
		break
	case "binary":
		t = Binary
		break
	default:
		t = Default
		break
	}

	pc.Type = t
}

func (pc *PkgConfigImpl) IsDownloadable() bool {
	return true
}



type PkgConfigTarball struct {
	Path RepoURL
	Type SourceType
}

func (pc *PkgConfigTarball) GetOrigin() RepoURL {
	return pc.Path
}

func (pc *PkgConfigTarball) SetOrigin(name, url string) {
	pc.Path = RepoURL{URL: url, Name: name}
}

func (pc *PkgConfigTarball) GetSourceType() string {
	return pc.Type.String()
}

func (pc *PkgConfigTarball) GetSourceType2() SourceType {
	return pc.Type
}

func (pc *PkgConfigTarball) SetSourceType(st string) {
	var t SourceType
	switch st {
	case "source":
		t = Source
		break
	case "binary":
		t = Binary
		break
	default:
		t = Default
		break
	}

	pc.Type = t
}

func (pc *PkgConfigTarball) IsDownloadable() bool {
	return false
}