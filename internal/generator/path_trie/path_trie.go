package path_trie

type Path []string

type PathDir struct {
	children      map[string]PathDir
	pathStartMeta interface{}
}

func New() PathDir {
	return PathDir{children: map[string]PathDir{}}
}

var pathStartFlagMeta = struct{}{} // for the case meta is not provided to LoadPath

func (dir *PathDir) LoadPath(path Path, meta interface{}) {
	if len(path) == 0 {
		if meta == nil {
			meta = pathStartFlagMeta
		}
		dir.pathStartMeta = meta
		return
	}
	lastIndex := len(path) - 1
	child, ok := dir.children[path[lastIndex]]
	if !ok {
		child = New()
	}
	child.LoadPath(path[:lastIndex], meta)
	dir.children[path[lastIndex]] = child
}

type ReducedPath struct {
	Rest Path
	Meta interface{}
}

func (dir PathDir) ReducedPaths() []ReducedPath {
	paths := []ReducedPath(nil)
	if dir.pathStartMeta != nil {
		meta := dir.pathStartMeta
		if meta == pathStartFlagMeta {
			meta = nil
		}
		paths = append(paths, ReducedPath{
			Rest: nil,
			Meta: meta,
		})
	}
	for childName, childDir := range dir.children {
		childPaths := childDir.ReducedPaths()
		// if len(childPaths) == 0 {
		// 	panic("impossible")
		// }
		if len(childPaths) == 1 {
			paths = append(paths, ReducedPath{
				Rest: Path{childName},
				Meta: childPaths[0].Meta,
			})
			continue
		}
		for _, childPath := range childPaths {
			paths = append(paths, ReducedPath{
				Rest: append(childPath.Rest, childName),
				Meta: childPath.Meta,
			})
		}
	}
	return paths
}
