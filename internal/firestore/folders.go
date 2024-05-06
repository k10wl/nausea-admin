package firestore

import (
	"context"
	"errors"
	"slices"

	"nausea-admin/internal/models"

	"cloud.google.com/go/firestore"
)

func (f *Firestore) GetFolderByID(id string) (models.Folder, error) {
	var folder models.Folder
	doc, err := f.collectionFolders().Doc(id).Get(f.ctx)
	if err != nil {
		return folder, err
	}
	err = doc.DataTo(&folder)
	folder.ID.ID = id
	return folder, err
}

func (f *Firestore) CreateFolder(
	folder models.Folder,
) (models.Folder, models.FolderContent, error) {
	folderDoc := f.collectionFolders().Doc(folder.ID.ID)
	parentDoc := f.collectionFolders().Doc(folder.ParentID)
	asContent, err := folder.AsContent()
	if err != nil {
		return folder, asContent, err
	}
	err = f.client.RunTransaction(f.ctx, func(_ context.Context, tx *firestore.Transaction) error {
		parentDocumentSnapshot, err := tx.Get(parentDoc)
		if err != nil {
			return err
		}
		var parentFolder models.Folder
		err = parentDocumentSnapshot.DataTo(&parentFolder)
		if err != nil {
			return err
		}
		if parentFolder.ProhibitNested {
			return errors.New("cannot nest folders")
		}
		err = tx.Create(folderDoc, folder)
		if err != nil {
			return err
		}
		contents := append(parentFolder.FolderContents, asContent)
		err = tx.Update(
			parentDoc,
			[]firestore.Update{{Path: "folders", Value: contents}},
		)

		return err
	})
	return folder, asContent, err
}

func (f *Firestore) MarkFolderDeletedByID(id string) (models.Folder, error) {
	var folder models.Folder
	err := f.client.RunTransaction(f.ctx, func(_ context.Context, tx *firestore.Transaction) error {
		toDeleteDoc := f.collectionFolders().Doc(id)
		toDeleteSnapshot, err := tx.Get(toDeleteDoc)
		if err != nil {
			return err
		}
		var toDeleteFolder models.Folder
		err = toDeleteSnapshot.DataTo(&toDeleteFolder)
		if err != nil {
			return err
		}
		if toDeleteFolder.Protected {
			return errors.New("cannot delete protected folder")
		}
		parentDoc := f.collectionFolders().Doc(toDeleteFolder.ParentID)
		parentSnapshot, err := tx.Get(parentDoc)
		if err != nil {
			return err
		}
		var parentFolder models.Folder
		err = parentSnapshot.DataTo(&parentFolder)
		if err != nil {
			return err
		}
		parentFolder.MarkDeletedFolderContents(id)
		err = tx.Update(
			parentDoc,
			[]firestore.Update{{Path: "folders", Value: parentFolder.FolderContents}},
		)
		if err != nil {
			return err
		}
		toDeleteFolder.Delete()
		toDeleteFolder.UpdatedAt = *toDeleteFolder.DeletedAt
		err = tx.Update(
			toDeleteDoc,
			[]firestore.Update{
				{Path: "deletedAt", Value: toDeleteFolder.DeletedAt},
				{Path: "updatedAt", Value: toDeleteFolder.UpdatedAt},
			},
		)
		folder = toDeleteFolder
		folder.ID.ID = id
		return err
	})
	return folder, err
}

func (f *Firestore) MarkFolderRestoredByID(id string) (models.Folder, error) {
	var folder models.Folder
	err := f.client.RunTransaction(f.ctx, func(_ context.Context, tx *firestore.Transaction) error {
		toDeleteDoc := f.collectionFolders().Doc(id)
		toDeleteSnapshot, err := tx.Get(toDeleteDoc)
		if err != nil {
			return err
		}
		var toRestoreFolder models.Folder
		err = toDeleteSnapshot.DataTo(&toRestoreFolder)
		if err != nil {
			return err
		}
		parentDoc := f.collectionFolders().Doc(toRestoreFolder.ParentID)
		parentSnapshot, err := tx.Get(parentDoc)
		if err != nil {
			return err
		}
		var parentFolder models.Folder
		err = parentSnapshot.DataTo(&parentFolder)
		if err != nil {
			return err
		}
		parentFolder.MarkRestoredFolderContents(id)
		err = tx.Update(
			parentDoc,
			[]firestore.Update{
				{Path: "folders", Value: parentFolder.FolderContents},
			},
		)
		if err != nil {
			return err
		}
		toRestoreFolder.Restore()
		toRestoreFolder.Update()
		err = tx.Update(
			toDeleteDoc,
			[]firestore.Update{
				{Path: "deletedAt", Value: toRestoreFolder.DeletedAt},
				{Path: "updatedAt", Value: toRestoreFolder.UpdatedAt},
			},
		)
		folder = toRestoreFolder
		folder.ID.ID = id
		return err
	})
	return folder, err
}

func (f *Firestore) PatchFolder(folderID string, patch models.Folder) (models.Folder, error) {
	var folder models.Folder
	err := f.client.RunTransaction(f.ctx, func(_ context.Context, tx *firestore.Transaction) error {
		toPatchDoc := f.collectionFolders().Doc(folderID)
		toPatchSnapshot, err := tx.Get(toPatchDoc)
		if err != nil {
			return err
		}
		var toPatchFolder models.Folder
		err = toPatchSnapshot.DataTo(&toPatchFolder)
		if err != nil {
			return err
		}
		toPatchFolder.Override(patch)
		parentDoc := f.collectionFolders().Doc(toPatchFolder.ParentID)
		parentSnapshot, err := tx.Get(parentDoc)
		if err != nil {
			return err
		}
		var parentFolder models.Folder
		err = parentSnapshot.DataTo(&parentFolder)
		if err != nil {
			return err
		}
		for i, val := range parentFolder.FolderContents {
			if val.RefID == folderID {
				asContent, err := toPatchFolder.AsContent()
				if err != nil {
					return err
				}
				parentFolder.FolderContents[i] = asContent
				break
			}
		}
		err = tx.Update(
			parentDoc,
			[]firestore.Update{
				{Path: "folders", Value: parentFolder.FolderContents},
			},
		)
		if err != nil {
			return err
		}
		toPatchFolder.Update()
		tx.Set(toPatchDoc, toPatchFolder)
		folder = toPatchFolder
		folder.ID.ID = folderID
		return err
	})
	return folder, err
}

func (f *Firestore) UploadMediaToFolder(media []models.MediaContent, folderId string) error {
	var folder models.Folder
	folderRef := f.collectionFolders().Doc(folderId)
	doc, _ := folderRef.Get(context.TODO())
	doc.DataTo(&folder)
	if folder.ProhibitMedia {
		return errors.New("cannot upload media to this folder")
	}
	_, err := folderRef.Update(context.TODO(), []firestore.Update{
		{Path: "media", Value: append(folder.MediaContents, media...)},
	})
	return err
}

func (f *Firestore) MarkMediaAsDeletedInFolder(
	mediaID string,
	folderID string,
) (models.MediaContent, error) {
	var folder models.Folder
	folderRef := f.collectionFolders().Doc(folderID)
	doc, _ := folderRef.Get(context.TODO())
	doc.DataTo(&folder)
	i := slices.IndexFunc(folder.MediaContents, func(m models.MediaContent) bool {
		return m.ID.ID == mediaID
	})
	folder.MediaContents[i].Delete()
	_, err := folderRef.Update(context.TODO(), []firestore.Update{
		{Path: "media", Value: folder.MediaContents},
	})
	return folder.MediaContents[i], err
}

func (f *Firestore) MarkMediaAsRestoredInFolder(
	mediaID string,
	folderID string,
) (models.MediaContent, error) {
	var folder models.Folder
	folderRef := f.collectionFolders().Doc(folderID)
	doc, _ := folderRef.Get(context.TODO())
	doc.DataTo(&folder)
	i := slices.IndexFunc(folder.MediaContents, func(m models.MediaContent) bool {
		return m.ID.ID == mediaID
	})
	folder.MediaContents[i].Restore()
	_, err := folderRef.Update(context.TODO(), []firestore.Update{
		{Path: "media", Value: folder.MediaContents},
	})
	return folder.MediaContents[i], err
}

func (f *Firestore) UpdateMediaInFolder(patch models.MediaContent) (models.MediaContent, error) {
	var folder models.Folder
	folderDoc := f.collectionFolders().Doc(patch.ParentID)
	folderSnapshot, err := folderDoc.Get(context.TODO())
	if err != nil {
		return models.MediaContent{}, err
	}
	err = folderSnapshot.DataTo(&folder)
	if err != nil {
		return models.MediaContent{}, err
	}
	if folder.ProhibitMedia {
		return models.MediaContent{}, errors.New("cannot upload media to this folder")
	}
	var media models.MediaContent
	for i, v := range folder.MediaContents {
		if v.ID.ID == patch.ID.ID {
			v.Override(patch)
			folder.MediaContents[i] = v
			media = v
			break
		}
	}
	_, err = folderDoc.Update(context.TODO(), []firestore.Update{
		{Path: "media", Value: folder.MediaContents},
	})
	return media, err
}

func (f *Firestore) PermanentlyDeleteMedia(
	folderID string,
	mediaID string,
) (models.MediaContent, error) {
	var content models.MediaContent
	err := f.client.RunTransaction(
		f.ctx,
		func(ctx context.Context, tx *firestore.Transaction) error {
			var folder models.Folder
			folderDoc := f.collectionFolders().Doc(folderID)
			folderSnapshot, err := folderDoc.Get(ctx)
			if err != nil {
				return err
			}
			err = folderSnapshot.DataTo(&folder)
			i := slices.IndexFunc(folder.MediaContents, func(media models.MediaContent) bool {
				return media.ID.ID == mediaID
			})
			if i == -1 {
				return errors.New("not found")
			}
			content = folder.MediaContents[i]
			wihtoutMedia := slices.Delete(folder.MediaContents, i, i+1)
			err = tx.Update(folderDoc, []firestore.Update{{Path: "media", Value: wihtoutMedia}})
			if err != nil {
				return err
			}
			tx.Delete(f.collectionMedia().Doc(content.RefID))
			return err
		},
	)
	return content, err
}

func (f *Firestore) PermanentlyDeleteFolder(
	folderID string,
) ([]models.MediaContent, error) {
	media := []models.MediaContent{}
	err := f.client.RunTransaction(
		f.ctx,
		func(ctx context.Context, t *firestore.Transaction) error {
			m, err := f.permanentlyDeleteFoldedInTransaction(folderID, ctx, t)
			media = slices.Concat(media, m)
			return err
		},
	)
	return media, err
}

func (f *Firestore) permanentlyDeleteFoldedInTransaction(
	folderID string,
	ctx context.Context,
	t *firestore.Transaction,
) ([]models.MediaContent, error) {
	media := []models.MediaContent{}
	// get folder
	var folder models.Folder
	folderDoc := f.collectionFolders().Doc(folderID)
	folderSnapshot, err := folderDoc.Get(ctx)
	if err != nil {
		return media, err
	}
	err = folderSnapshot.DataTo(&folder)
	if err != nil {
		return media, err
	}
	media = folder.MediaContents
	// delete from parent
	var parent models.Folder
	parentDoc := f.collectionFolders().Doc(folder.ParentID)
	parentSnapshot, err := parentDoc.Get(ctx)
	if err != nil {
		return media, err
	}
	parentSnapshot.DataTo(&parent)
	i := slices.IndexFunc(parent.FolderContents, func(f models.FolderContent) bool {
		return f.RefID == folderID
	})
	if i == -1 {
		return media, err
	}
	folders := slices.Delete(parent.FolderContents, i, i+1)
	t.Update(parentDoc, []firestore.Update{{Path: "folders", Value: folders}})
	// delete media
	for _, m := range folder.MediaContents {
		err := t.Delete(f.collectionMedia().Doc(m.RefID))
		if err != nil {
			return media, err
		}
	}
	for _, fc := range folder.FolderContents {
		m, err := f.permanentlyDeleteFoldedInTransaction(fc.RefID, ctx, t)
		media = slices.Concat(media, m)
		if err != nil {
			return media, err
		}
	}
	// delete itself
	err = t.Delete(folderDoc)
	return media, err
}

func (f *Firestore) ReorderMedia(folderID string, from int, to int) error {
	var folder models.Folder
	doc := f.collectionFolders().Doc(folderID)
	snap, err := doc.Get(f.ctx)
	if err != nil {
		return err
	}
	snap.DataTo(&folder)
	el := folder.MediaContents[from]
	removed := slices.Delete(folder.MediaContents, from, from+1)
	swapped := slices.Insert(removed, to, el)
	_, err = doc.Update(f.ctx, []firestore.Update{{
		Path: "media", Value: swapped,
	}})
	return err
}
