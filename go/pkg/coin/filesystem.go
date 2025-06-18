package coin

import (
	"io"
	"os"
	"path/filepath"
	"runtime"
)

// ErrorAlreadyExists mirrors the C errno returned when a directory exists.
var ErrorAlreadyExists = os.ErrExist

// CreatePath creates all missing directories in the given path.
func CreatePath(path string) error {
	return os.MkdirAll(path, 0o755)
}

// CopyFile copies src to dst, truncating dst if it exists.
func CopyFile(src, dst string) error {
	in, err := os.Open(src)
	if err != nil {
		return err
	}
	defer in.Close()
	out, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer out.Close()
	if _, err := io.Copy(out, in); err != nil {
		return err
	}
	return out.Sync()
}

// PathContents returns the file and directory names within path.
func PathContents(path string) ([]string, error) {
	entries, err := os.ReadDir(path)
	if err != nil {
		return nil, err
	}
	names := make([]string, 0, len(entries))
	for _, e := range entries {
		names = append(names, e.Name())
	}
	return names, nil
}

// DataPath returns an OS-appropriate directory for application data.
func DataPath() string {
	name := ClientName
	if TestNet {
		name += "TestNet"
	}
	home := homePath()
	switch runtime.GOOS {
	case "windows":
		if app := os.Getenv("APPDATA"); app != "" {
			return filepath.Join(app, name)
		}
		return filepath.Join(home, name)
	case "darwin":
		return filepath.Join(home, "Library", "Application Support", name)
	case "android":
		return filepath.Join(home, "data")
	default:
		return filepath.Join(home, "."+name, "data")
	}
}

func homePath() string {
	if h, err := os.UserHomeDir(); err == nil {
		if !filepath.IsAbs(h) {
			if cwd, err := os.Getwd(); err == nil {
				return filepath.Join(cwd, h)
			}
		}
		return h
	}
	if dir := os.Getenv("HOME"); dir != "" {
		return dir
	}
	if dir := os.Getenv("USERPROFILE"); dir != "" {
		return dir
	}
	if drive, home := os.Getenv("HOMEDRIVE"), os.Getenv("HOMEPATH"); drive != "" && home != "" {
		return drive + home
	}
	return "."
}
