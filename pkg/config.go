package projector

// Config contains all information required to generate project.
type Config struct {
	WorkingDirectory string
	ProjectAuthor    string
	ProjectName      string
	ProjectPackage   string
	Manifest         *TemplateManifest
}
