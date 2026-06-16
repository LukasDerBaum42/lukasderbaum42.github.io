package main
import (
	"os"
	"path/filepath"
	"io/fs"
	"os/exec"
	"strings"
	"fmt"
)

func devCreate(path string, dev bool) (*os.File, error) {
	if dev {
		return os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
	}
	return os.Create(path)
}

func devCopyFS(dst string, src fs.FS, dev bool) error {
	if dev {
		return fs.WalkDir(src, ".", func(path string, d fs.DirEntry, err error) error {
			if err != nil {
				return err
			}
			target := filepath.Join(dst, path)
			if d.IsDir() {
				return os.MkdirAll(target, 0755)
			}
			os.Remove(target)
			data, err := fs.ReadFile(src, path)
			if err != nil {
				return err
			}
			return os.WriteFile(target, data, 0644)
		})
	}
	return os.CopyFS(dst, src)
}

func compilTS(dev bool) error {
    os.MkdirAll("build/scripts", 0755)
    cmd := exec.Command("tsc", "--outDir", "build/scripts", "--target", "ES2020", "--module", "ESNext")

    // Glob all .ts files in scripts/
    entries, err := os.ReadDir("scripts")
    if err != nil {
        return err
    }
    for _, e := range entries {
        if !e.IsDir() && strings.HasSuffix(e.Name(), ".ts") {
            cmd.Args = append(cmd.Args, "scripts/"+e.Name())
        }
    }

    out, err := cmd.CombinedOutput()
    if err != nil {
        return fmt.Errorf("%s\n%s", err, out)
    }
    return nil
}