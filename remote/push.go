package remote

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/beardnick/mynvim/check"
	"github.com/beardnick/mynvim/global"
	"github.com/beardnick/mynvim/neovim"
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/google/go-github/v33/github"
	"github.com/neovim/go-client/nvim"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

func PullDir(nvm *nvim.Nvim, args []string) {
	dir := args[0]
	err := nvm.SetVar("mynvim_plugin_dir", dir)
	if err != nil {
		neovim.EchoErrStack(err)
	}
}

var repos []string

type TagsResp []github.RepositoryTag

// todo: 处理tag分页
func repoTags(fullName string) (tags []string, err error) {
	tagsUrl := fmt.Sprintf("https://api.github.com/repos/%s/tags", fullName)
	rep, err := http.DefaultClient.Get(tagsUrl)
	if err != nil {
		return
	}
	if rep.StatusCode != 200 {
		err = errors.New("get tags failed")
		return
	}
	defer rep.Body.Close()
	buf, err := ioutil.ReadAll(rep.Body)
	if err != nil {
		return
	}
	t := TagsResp{}
	err = json.Unmarshal(buf, &t)
	if err != nil {
		return
	}
	for _, i := range t {
		tags = append(tags, *i.Name)
	}
	return
}

func Pull(nvm *nvim.Nvim, args []string) {
	fullName := args[0]
	fullName = strings.Trim(fullName, "'")
	fullName = strings.Trim(fullName, `"`)
	if !strings.Contains(fullName, "/") {
		neovim.Echomsg("invalid plug , / is needed")
		return
	}
	repos = append(repos, fullName)
	if !check.RepoExists(getStoreDir(fullName)) {
		pullLatest(fullName)
	}
	neovim.Echomsg("begin source")
	err := loadPlugin(nvm, fullName)
	if err != nil {
		neovim.Echomsg(err)
	}
	return
}

func getStoreDir(fullName string) string {
	return filepath.Join(GetPluginDir(), "github.com", fullName)
}

func GetPluginDir() (plugin string) {
	err := global.Nvm.Var("mynvim_plugin_dir", &plugin)
	if err != nil {
		plugin = filepath.Join(HomeDir(),".cache/mynvim/plugin")
		_ = global.Nvm.SetVar("mynvim_plugin_dir", plugin)
	}
	return
}

func Push(nvm *nvim.Nvim, args []string) {
	repo := args[0]
	neovim.Echomsg(repo)
}

func HomeDir() string {
	return os.Getenv("HOME")
}

func pullLatest(fullName string) (err error) {
	tags, err := repoTags(fullName)
	if err != nil {
		neovim.Echomsg(err)
		return
	}
	storeDir := getStoreDir(fullName)
	//git clone --depth 1 --branch <tag_name> <repo_url>
	if len(tags) > 0 {
		_, err = git.PlainClone(storeDir, false, &git.CloneOptions{
			URL:           fmt.Sprintf("https://github.com/%s", fullName),
			ReferenceName: plumbing.ReferenceName("refs/tags/" + tags[0]),
			Depth:         1,
		})
	} else {
		_, err = git.PlainClone(storeDir, false, &git.CloneOptions{
			URL: fmt.Sprintf("https://github.com/%s", fullName),
		})
	}
	if err != nil {
		return
	}
	return
}

func PullAll(nvm *nvim.Nvim) {
	err := EnsurePath(GetPluginDir())
	if err != nil {
		neovim.Echomsg(err)
		return
	}
	for _, repo := range repos {
		err := pullLatest(repo)
		if err != nil {
			neovim.Echomsg(err)
		}
	}
}

func loadPlugin(nvm *nvim.Nvim, fullName string) (err error) {
	dir := getStoreDir(fullName)
	plugins, err := filepath.Glob(filepath.Join(dir, "plugin/*.vim"))
	if err != nil {
		return
	}
	afters, err := filepath.Glob(filepath.Join(dir, "after/*.vim"))
	if err != nil {
		return
	}
	all := append(plugins, afters...)
	err = nvm.Command(fmt.Sprintf("source %s", strings.Join(all, " ")))
	return
}

func EnsurePath(path string) (err error) {
	_, err = os.Stat(path)
	if errors.Is(err,os.ErrNotExist) {
		err = os.MkdirAll(path,os.ModePerm)
		return
	}
	return
}