package remote

import (
	"encoding/json"
	"errors"
	"fmt"
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

var pluginDir = filepath.Join(HomeDir(), ".mynvim")

type TagsResp []github.RepositoryTag

// todo: 处理分页
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
	tags, err := repoTags(fullName)
	if err != nil {
		neovim.Echomsg(err)
		return
	}
	//git clone --depth 1 --branch <tag_name> <repo_url>
	if len(tags) > 0 {
		_, err = git.PlainClone(filepath.Join(pluginDir, "github.com", fullName), false, &git.CloneOptions{
			URL:           fmt.Sprintf("https://github.com/%s", fullName),
			ReferenceName: plumbing.ReferenceName("refs/tags/" + tags[0]),
			Depth:         1,
		})
	} else {
		_, err = git.PlainClone(filepath.Join(pluginDir, "github.com", fullName), false, &git.CloneOptions{
			URL: fmt.Sprintf("https://github.com/%s", fullName),
		})
	}
	if err != nil {
		return
	}
}

func Push(nvm *nvim.Nvim, args []string) {
	repo := args[0]
	neovim.Echomsg(repo)
}

func HomeDir() string {
	return os.Getenv("HOME")
}
