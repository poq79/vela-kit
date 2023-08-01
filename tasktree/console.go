package tasktree

type Console struct {
	Path   string `json:"path"`
	Script string `json:"script"`
}

/*

path: "/usr/local/ssoc",
Script: file.cat("daemon.log")

path: "/usr/local/ssoc",
Script: file.dir("."),

path: "/usr/local/ssoc"
Script: vela.ps()




*/
