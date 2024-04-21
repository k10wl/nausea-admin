package firestore

import (
	"context"
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
			if val.ContentBase.RefID == folderID {
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
	_, err := folderRef.Update(context.TODO(), []firestore.Update{
		{Path: "media", Value: append(folder.MediaContents, media...)},
	})
	return err
}

func (f *Firestore) MarkMediaAsDeletedInFolder(mediaID string, folderID string) (models.MediaContent, error) {
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

func (f *Firestore) MarkMediaAsRestoredInFolder(mediaID string, folderID string) (models.MediaContent, error) {
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
