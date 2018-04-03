package stress

type PushConfig struct {
	Method    string
	Namespace string
	RepoName  string

	TagName               string
	TagPattern            string
	PushesPerBatch        int
	Batches               int
	LayerMinBytes         int
	LayerMaxBytes         int
	LayerSizeScalePattern string
	ImageMinLayers        int
	ImageMaxLayers        int
	LayersScalePattern    string
	Retries               int
}

type PullConfig struct {
	Method     string
	Namespace  string
	RepoName   string
	TagName    string
	TagPattern string
	Duration   string
}

// The mass-population config is a superset of the push config:
// It has a user/repo range that it will randomly populate within, instead of using a single namepsace/repository
// Mass-population will only run if UsersToPopulate > 0 and ReposToPopulate > 0,
// (note that for both of them being 1, you can just use a regular push and specify which namespace/repo you want to populate)
type MassPopulateConfig struct {
	Push *PushConfig

	// Requirement is that these users and repos exist and they need to follow a similar pattern to
	// the the creation of such in the setup phase (although they can have been generated differently,
	// for example the user pattern is identical to our big LDAP test server)
	// Mass-populate will randomly choose UsersToPopulate users and ReposToPopulate per user repos to populate with the
	// pushconfig settings (replacing the namespace/repoName in it)
	TotalNumUsers        int
	TotalNumReposPerUser int
	UsersToPopulate      int
	ReposToPopulate      int
}

type Config struct {
	DTRURL         string
	DTRCA          string
	DTRInsecureTLS bool
	Username       string
	Password       string
	RefreshToken   string
	Seed           int
	Push           *PushConfig
	Pull           *PullConfig
	MassPopulate   *MassPopulateConfig
}
