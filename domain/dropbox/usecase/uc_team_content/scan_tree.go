package uc_team_content

import (
	"errors"
	"github.com/watermint/toolbox/domain/dropbox/api/dbx_context"
	"github.com/watermint/toolbox/domain/dropbox/model/mo_path"
	"github.com/watermint/toolbox/domain/dropbox/model/mo_sharedfolder"
	"github.com/watermint/toolbox/domain/dropbox/service/sv_file"
	"github.com/watermint/toolbox/essentials/kvs/kv_kvs"
	"github.com/watermint/toolbox/essentials/kvs/kv_storage"
	"github.com/watermint/toolbox/essentials/log/esl"
	"github.com/watermint/toolbox/infra/control/app_control"
)

type Tree struct {
	NamespaceId   string `json:"namespace_id"`
	NamespaceName string `json:"namespace_name"`
	RelativePath  string `json:"relative_path"`
}

type TeamFolderNestedFolderScanWorker struct {
	ctl          app_control.Control
	ctx          dbx_context.Context
	metadata     kv_storage.Storage
	tree         kv_storage.Storage
	teamFolderId string
	descendants  []string
}

func (z *TeamFolderNestedFolderScanWorker) addTree(t *Tree) error {
	return z.tree.Update(func(kvs kv_kvs.Kvs) error {
		return kvs.PutJsonModel(t.NamespaceId, t)
	})
}

func (z *TeamFolderNestedFolderScanWorker) Exec() error {
	// Breadth first search for nested folders
	l := z.ctl.Log().With(esl.String("teamFolderId", z.teamFolderId))
	teamFolderName := ""
	l.Debug("lookup team folder name")
	err := z.metadata.View(func(kvs kv_kvs.Kvs) error {
		tf := &mo_sharedfolder.SharedFolder{}
		if err := kvs.GetJsonModel(z.teamFolderId, tf); err != nil {
			return err
		}
		teamFolderName = tf.Name
		return nil
	})
	if err != nil {
		return err
	}
	l = l.With(esl.String("teamFolderName", teamFolderName))

	err = z.tree.Update(func(kvs kv_kvs.Kvs) error {
		return kvs.PutJsonModel(z.teamFolderId, &Tree{
			NamespaceId:   z.teamFolderId,
			NamespaceName: teamFolderName,
			RelativePath:  mo_path.NewDropboxPath(teamFolderName).Path(),
		})
	})
	if err != nil {
		return err
	}

	l.Debug("search nested folders", esl.Strings("descendants", z.descendants))
	traverse := make(map[string]bool)
	for _, d := range z.descendants {
		traverse[d] = false
	}
	completed := func() bool {
		for _, t := range traverse {
			if !t {
				return false
			}
		}
		return true
	}
	ErrorScanCompleted := errors.New("scan completed")

	ctx := z.ctx.WithPath(dbx_context.Namespace(z.teamFolderId))
	var scan func(path mo_path.DropboxPath) error
	scan = func(path mo_path.DropboxPath) error {
		entries, err := sv_file.NewFiles(ctx).List(path)
		if err != nil {
			l.Debug("Unable to retrieve entries", esl.Error(err), esl.String("path", path.Path()))
			return err
		}

		// Mark nested folders
		for _, entry := range entries {
			if f, ok := entry.Folder(); ok {
				if f.EntrySharedFolderId != "" {
					traverse[f.EntrySharedFolderId] = true
					rp := path.ChildPath(f.EntryName)
					err := z.addTree(&Tree{
						NamespaceId:   f.EntrySharedFolderId,
						NamespaceName: f.EntryName,
						RelativePath:  teamFolderName + rp.Path(),
					})
					if err != nil {
						return err
					}
				}
			}
		}

		// Return if the scan completed
		if completed() {
			return ErrorScanCompleted
		}

		// Dive into descendants
		for _, entry := range entries {
			if f, ok := entry.Folder(); ok {
				if err := scan(path.ChildPath(f.Name())); err != nil {
					return err
				}
			}
		}
		return nil
	}

	if err := scan(mo_path.NewDropboxPath("")); err != nil && err != ErrorScanCompleted {
		l.Debug("The error occurred on scanning team folder", esl.Error(err))
		return err
	}

	l.Debug("Scan completed")

	return nil
}

// Use breadth first search for file tree
type TeamFolderScanner struct {
	Ctl      app_control.Control
	Ctx      dbx_context.Context
	Metadata kv_storage.Storage
	Tree     kv_storage.Storage
}

func (z *TeamFolderScanner) parentChildRelationship() (relation map[string]string, err error) {
	l := z.Ctl.Log()
	l.Debug("Making mapping of parent-child relationship")

	// namespace_id -> parent namespace_id
	relation = make(map[string]string)

	err = z.Metadata.View(func(kvs kv_kvs.Kvs) error {
		return kvs.ForEachModel(&mo_sharedfolder.SharedFolder{}, func(key string, m interface{}) error {
			ns := m.(*mo_sharedfolder.SharedFolder)
			relation[ns.SharedFolderId] = ns.ParentSharedFolderId
			return nil
		})
	})
	if err != nil {
		return nil, err
	}
	l.Debug("Relation", esl.Any("relation", relation))

	return relation, nil
}

func (z *TeamFolderScanner) namespaceToTopNamespaceId() (top map[string]string, err error) {
	l := z.Ctl.Log()
	// namespace_id -> top level namespace_id
	top = make(map[string]string)

	relation, err := z.parentChildRelationship()
	if err != nil {
		return nil, err
	}

	l.Debug("Making child - top level folder namespace mapping")
	for nsid := range relation {
		ll := l.With(esl.String("namespace_id", nsid))
		top[nsid] = ""
		chain := make([]string, 0)
		parent := relation[nsid]
		current := parent
		for {
			chain = append(chain, parent)

			var ok bool
			current, ok = relation[parent]
			if current == "" || !ok {
				break
			}
			parent = current
		}
		top[nsid] = parent
		ll.Debug("Top folder", esl.String("top", parent))
	}

	return top, nil
}

func (z *TeamFolderScanner) nestedFolderNamespaceIds() (nested map[string][]string, others []string, err error) {
	l := z.Ctl.Log()

	// team_folder's namespace_id -> array of nested team folder namespace_id
	nested = make(map[string][]string)

	// other un-nested namespaces
	others = make([]string, 0)

	top, err := z.namespaceToTopNamespaceId()
	if err != nil {
		return nil, nil, err
	}

	l.Debug("Aggregate nested folders")
	for nsid, t := range top {
		if t == "" {
			others = append(others, nsid)
			continue
		}
		if _, ok := nested[t]; !ok {
			nested[t] = make([]string, 0)
		}
		nested[t] = append(nested[t], nsid)
	}

	l.Debug("Team folders and nested folders", esl.Any("nested", nested))
	l.Debug("Others", esl.Strings("others", others))

	return nested, others, nil
}

func (z *TeamFolderScanner) Scan() error {
	l := z.Ctl.Log()
	queue := z.Ctl.NewQueue()
	nested, others, err := z.nestedFolderNamespaceIds()
	if err != nil {
		return err
	}
	for nsid, descendants := range nested {
		queue.Enqueue(&TeamFolderNestedFolderScanWorker{
			ctl:          z.Ctl,
			ctx:          z.Ctx,
			metadata:     z.Metadata,
			tree:         z.Tree,
			teamFolderId: nsid,
			descendants:  descendants,
		})
	}

	var lastErr error

	for _, nsid := range others {
		err := z.Metadata.View(func(kvs kv_kvs.Kvs) error {
			meta := &mo_sharedfolder.SharedFolder{}
			if err := kvs.GetJsonModel(nsid, meta); err != nil {
				return err
			}
			return z.Tree.Update(func(kvs kv_kvs.Kvs) error {
				return kvs.PutJsonModel(nsid, &Tree{
					NamespaceId:   nsid,
					NamespaceName: meta.Name,
					RelativePath:  mo_path.NewDropboxPath(meta.Name).Path(),
				})
			})
		})
		if err != nil {
			l.Debug("Unable to convert namespace_id to tree", esl.String("nsid", nsid), esl.Error(err))
			lastErr = err
		}
	}

	queue.Wait()

	if lastErr != nil {
		return lastErr
	}
	return nil
}