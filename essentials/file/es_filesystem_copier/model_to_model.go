package es_filesystem_copier

import (
	"errors"
	"github.com/watermint/toolbox/essentials/file/es_filesystem"
	"github.com/watermint/toolbox/essentials/file/es_filesystem_model"
	"github.com/watermint/toolbox/essentials/log/esl"
	"github.com/watermint/toolbox/essentials/model/em_file"
)

func NewModelToModel(l esl.Logger, sourceRoot, targetRoot em_file.Folder) es_filesystem.Connector {
	return &modelToModelCopier{
		l:          l,
		sourceRoot: sourceRoot,
		targetRoot: targetRoot,
	}
}

type modelToModelCopier struct {
	l          esl.Logger
	sourceRoot em_file.Folder
	targetRoot em_file.Folder
}

func (z modelToModelCopier) Copy(source es_filesystem.Entry, target es_filesystem.Path) (copied es_filesystem.Entry, err es_filesystem.FileSystemError) {
	l := z.l.With(esl.Any("source", source.AsData()), esl.String("target", target.Path()))
	l.Debug("Copy")

	sourceNode := em_file.ResolvePath(z.sourceRoot, source.Path().Path())
	if sourceNode == nil {
		l.Debug("Unable to find the source node")
		return nil, es_filesystem_model.NewError(errors.New("source node not found"), es_filesystem_model.ErrorTypePathNotFound)
	}

	sourceFile, ok := sourceNode.(em_file.File)
	if !ok || sourceNode.Type() != em_file.FileNode {
		l.Debug("Node is not a file")
		return nil, es_filesystem_model.NewError(errors.New("source node is not a file"), es_filesystem_model.ErrorTypeOther)
	}

	targetFolderPath := target.Ancestor()
	targetFolder := em_file.ResolvePath(z.targetRoot, targetFolderPath.Path())
	if targetFolder == nil {
		if !em_file.CreateFolder(z.targetRoot, targetFolderPath.Path()) {
			l.Debug("Unable to create folder", esl.Error(err))
			return nil, es_filesystem_model.NewError(err, es_filesystem_model.ErrorTypeOther)
		}
		targetFolder = em_file.ResolvePath(z.targetRoot, targetFolderPath.Path())
		if targetFolder == nil {
			l.Debug("Unable to resolve target path")
			return nil, es_filesystem_model.NewError(errors.New("unable to resolve target path"), es_filesystem_model.ErrorTypePathNotFound)
		}
	}

	if targetFolder.Type() != em_file.FolderNode {
		l.Debug("target nod is not a folder")
		return nil, es_filesystem_model.NewError(errors.New("target node is not a folder"), es_filesystem_model.ErrorTypeConflict)
	}

	targetFile := sourceFile.Clone()
	targetFolder.(em_file.Folder).Add(targetFile)

	return es_filesystem_model.NewEntry(target.Path(), targetFile), nil
}
