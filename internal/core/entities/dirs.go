package entities

import (
	"github.com/excoriate/stiletto/internal/utils"
	"os"
)

func GetDirCfg() *DirCfg {
	currentDir, _ := os.Getwd()
	currentDirAbs, _ := utils.PathToAbsolute(currentDir)
	homeDir, _ := os.UserHomeDir()
	homeDirAbs, _ := utils.PathToAbsolute(homeDir)

	if err := utils.IsGitRepository(currentDirAbs); err != nil {
		return &DirCfg{
			BaseDir:    currentDir,
			BaseDirAbs: currentDirAbs,
			HomeDir:    homeDir,
			HomeDirAbs: homeDirAbs,
			IsGitRepo:  false,
			GitDirAbs:  "",
		}
	}

	return &DirCfg{
		BaseDir:    currentDir,
		BaseDirAbs: currentDirAbs,
		IsGitRepo:  true,
		HomeDir:    homeDir,
		HomeDirAbs: homeDirAbs,
		GitDirAbs:  "",
	}
}
