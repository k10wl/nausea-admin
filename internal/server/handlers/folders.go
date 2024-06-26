package handlers

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"sync"

	"nausea-admin/internal/db"
	"nausea-admin/internal/models"
	"nausea-admin/internal/server/template"
	"nausea-admin/internal/server/utils"
	"nausea-admin/internal/storage"
)

type FoldersHandler struct {
	DB       db.DB
	Template template.Template
	Storage  storage.Storage
}

func NewFoldersHandler(
	db db.DB,
	template template.Template,
	storage storage.Storage,
) FoldersHandler {
	return FoldersHandler{
		DB:       db,
		Template: template,
		Storage:  storage,
	}
}

func (fh FoldersHandler) GetHomePage(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		_, _, data := utils.WithPageData(w, r, map[string]interface{}{})
		w.WriteHeader(http.StatusNotFound)
		fh.Template.ExecuteTemplate(w, "/404", data)
		return
	}
	fh.Template.ExecuteTemplate(utils.WithPageData(w, r, map[string]interface{}{}))
}

func (fh FoldersHandler) GetFoldersPage(w http.ResponseWriter, r *http.Request) {
	folderID := getFolderID(r)
	folder, err := fh.DB.GetFolderByID(folderID)
	// TODO interface as struct
	_, _, data := utils.WithPageData(
		w,
		r,
		map[string]interface{}{
			"Folder": folder,
			"Error":  err,
			"Title":  folder.Name,
		},
	)
	fh.Template.ExecuteTemplate(w, "/folders", data)
}

func (fh FoldersHandler) CreateFolder(w http.ResponseWriter, r *http.Request) {
	parentID := getFolderID(r)
	name := r.FormValue("name")
	if name == "" {
		w.Header().Set("HX-Reswap", "innerHTML")
		utils.ErrorResponse(
			w,
			r,
			http.StatusBadRequest,
			errors.New("insufficient information, missing name"),
		)
		return
	}
	if parentID == "" {
		w.Header().Set("HX-Reswap", "innerHTML")
		utils.ErrorResponse(
			w,
			r,
			http.StatusBadRequest,
			errors.New("insufficient information, missing parent folder id"),
		)
		return
	}
	folder, err := models.NewFolder(parentID, name)
	if err != nil {
		utils.ErrorResponse(
			w,
			r,
			http.StatusInternalServerError,
			errors.New("failed to create folder"),
		)
		return
	}
	_, asContent, err := fh.DB.CreateFolder(*folder)
	if err != nil {
		utils.ErrorResponse(
			w,
			r,
			http.StatusInternalServerError,
			errors.New("failed to create folder"),
		)
		return
	}
	fh.Template.ExecuteTemplate(w, "folder-list", asContent)
}

func (fh FoldersHandler) MarkFolderAsDeleted(w http.ResponseWriter, r *http.Request) {
	folderID := getFolderID(r)
	if folderID == "" {
		w.Header().Set("HX-Reswap", "innerHTML")
		utils.ErrorResponse(w, r, http.StatusBadRequest, errors.New("no folderID"))
		return
	}
	if folderID == models.RootFolderID {
		w.Header().Set("HX-Reswap", "innerHTML")
		utils.ErrorResponse(w, r, http.StatusBadRequest, errors.New("can't delete root folder"))
		return
	}
	folder, err := fh.DB.MarkFolderDeletedByID(folderID)
	if err != nil {
		utils.ErrorResponse(w, r, http.StatusInternalServerError, err)
		return
	}
	asContent, err := folder.AsContent()
	if err != nil {
		utils.ErrorResponse(
			w,
			r,
			http.StatusInternalServerError,
			errors.New("folder is deleted, reload to update UI"),
		)
		return
	}
	w.WriteHeader(http.StatusOK)
	fh.Template.ExecuteTemplate(w, "folder-list", asContent)
}

func (fh FoldersHandler) MarkMediaAsDeletedInFolder(w http.ResponseWriter, r *http.Request) {
	folderID := getFolderID(r)
	mediaID := r.PathValue("media_id")
	if folderID == "" {
		w.Header().Set("HX-Reswap", "innerHTML")
		utils.ErrorResponse(w, r, http.StatusBadRequest, errors.New("no folder id"))
		return
	}
	if mediaID == "" {
		w.Header().Set("HX-Reswap", "innerHTML")
		utils.ErrorResponse(w, r, http.StatusBadRequest, errors.New("no media id"))
		return
	}
	media, err := fh.DB.MarkMediaAsDeletedInFolder(mediaID, folderID)
	if err != nil {
		w.Header().Set("HX-Reswap", "innerHTML")
		utils.ErrorResponse(w, r, http.StatusBadRequest, errors.New("failed to restore"))
		return
	}
	fh.Template.ExecuteTemplate(w, "media-list", media)
}

func (fh FoldersHandler) RestoreMediaInFolder(w http.ResponseWriter, r *http.Request) {
	folderID := getFolderID(r)
	mediaID := r.PathValue("media_id")
	if folderID == "" {
		w.Header().Set("HX-Reswap", "innerHTML")
		utils.ErrorResponse(w, r, http.StatusBadRequest, errors.New("no folder id"))
		return
	}
	if mediaID == "" {
		w.Header().Set("HX-Reswap", "innerHTML")
		utils.ErrorResponse(w, r, http.StatusBadRequest, errors.New("no media id"))
		return
	}
	media, err := fh.DB.MarkMediaAsRestoredInFolder(mediaID, folderID)
	if err != nil {
		w.Header().Set("HX-Reswap", "innerHTML")
		utils.ErrorResponse(w, r, http.StatusBadRequest, errors.New("failed to restore"))
		return
	}
	fh.Template.ExecuteTemplate(w, "media-list", media)
}

func (fh FoldersHandler) RestoreFolder(w http.ResponseWriter, r *http.Request) {
	folderID := getFolderID(r)
	if folderID == "" {
		w.Header().Set("HX-Reswap", "innerHTML")
		utils.ErrorResponse(w, r, http.StatusBadRequest, errors.New("no folderID"))
		return
	}
	folder, err := fh.DB.MarkFolderRestoredByID(folderID)
	if err != nil {
		utils.ErrorResponse(w, r, http.StatusInternalServerError, err)
		return
	}
	asContent, err := folder.AsContent()
	if err != nil {
		utils.ErrorResponse(
			w,
			r,
			http.StatusInternalServerError,
			errors.New("folder is restored, reload to update UI"),
		)
		return
	}
	w.WriteHeader(http.StatusOK)
	fh.Template.ExecuteTemplate(w, "folder-list", asContent)
}

func (fh FoldersHandler) PatchFolder(w http.ResponseWriter, r *http.Request) {
	folderID := getFolderID(r)
	if folderID == "" {
		w.Header().Set("HX-Reswap", "innerHTML")
		utils.ErrorResponse(w, r, http.StatusBadRequest, errors.New("no folderID"))
		return
	}
	var patch models.Folder
	patch.Name = r.FormValue("name")
	folder, _ := fh.DB.PatchFolder(folderID, patch)
	asContent, err := folder.AsContent()
	if err != nil {
		w.Header().Set("HX-Reswap", "innerHTML")
		utils.ErrorResponse(
			w,
			r,
			http.StatusBadRequest,
			errors.New("cannot show update, please refresh"),
		)
		return
	}
	if r.URL.Query().Has("from-inside") {
		fmt.Fprintf(w, fmt.Sprintf("%[1]s<script>document.title=%[1]q</script>", asContent.Name))
		return
	}
	fh.Template.ExecuteTemplate(w, "folder-list", asContent)
}

func (fh FoldersHandler) EditFolderMedia(w http.ResponseWriter, r *http.Request) {
	folderID := getFolderID(r)
	mediaID := r.PathValue("media_id")
	if folderID == "" {
		w.Header().Set("HX-Reswap", "innerHTML")
		utils.ErrorResponse(w, r, http.StatusBadRequest, errors.New("no folderID"))
		return
	}
	name := r.FormValue("name")
	description := r.FormValue("description")
	update := models.MediaContent{
		ContentBase: models.ContentBase{ID: models.ID{ID: mediaID}},
		Name:        name,
		Description: description,
		ParentID:    folderID,
	}
	mediaContent, err := fh.DB.UpdateMediaInFolder(update)
	if err != nil {
		w.Header().Set("HX-Reswap", "innerHTML")
		utils.ErrorResponse(w, r, http.StatusBadRequest, errors.New("failed to update media"))
		return
	}
	fh.Template.ExecuteTemplate(w, "media-list", mediaContent)
}

func (fh FoldersHandler) DeleteFolderMedia(w http.ResponseWriter, r *http.Request) {
	folderID := getFolderID(r)
	mediaID := r.PathValue("media_id")
	if folderID == "" {
		w.Header().Set("HX-Reswap", "innerHTML")
		utils.ErrorResponse(w, r, http.StatusBadRequest, errors.New("no folderID"))
		return
	}
	content, err := fh.DB.PermanentlyDeleteMedia(folderID, mediaID)
	if err != nil {
		w.Header().Set("HX-Reswap", "innerHTML")
		utils.ErrorResponse(w, r, http.StatusBadRequest, errors.New("failed to update media"))
		return
	}
	err = fh.Storage.RemoveObject(fh.Storage.ParseURLKey(content.ThumbnailURL))
	if err != nil {
		fmt.Println(">", err)
		w.Header().Set("HX-Reswap", "innerHTML")
		utils.ErrorResponse(
			w,
			r,
			http.StatusBadRequest,
			errors.New("failed to remove image thumbnail from cloud"),
		)
		return
	}
	err = fh.Storage.RemoveObject(fh.Storage.ParseURLKey(content.URL))
	if err != nil {
		w.Header().Set("HX-Reswap", "innerHTML")
		utils.ErrorResponse(
			w,
			r,
			http.StatusBadRequest,
			errors.New("failed to remove image from cloud"),
		)
		return
	}
	w.WriteHeader(http.StatusOK)
}

func (fh FoldersHandler) DeleteFolder(w http.ResponseWriter, r *http.Request) {
	folderID := getFolderID(r)
	contents, err := fh.DB.PermanentlyDeleteFolder(folderID)
	if err != nil {
		w.Header().Set("HX-Reswap", "innerHTML")
		utils.ErrorResponse(w, r, http.StatusBadRequest, errors.New("failed to delete folder"))
		return
	}
	echan := make(chan error)
	wg := sync.WaitGroup{}
	wg.Add(len(contents))
	for _, m := range contents {
		go func() {
			err := fh.Storage.RemoveObject(fh.Storage.ParseURLKey(m.URL))
			if err != nil {
				echan <- err
			}
			err = fh.Storage.RemoveObject(fh.Storage.ParseURLKey(m.ThumbnailURL))
			if err != nil {
				echan <- err
			}
			wg.Done()
		}()
	}
	wg.Wait()
	close(echan)
	errs := []error{}
	for err := range echan {
		errs = append(errs, err)
	}
	if len(errs) > 0 {
		w.Header().Set("HX-Reswap", "innerHTML")
		utils.ErrorResponse(
			w,
			r,
			http.StatusBadRequest,
			errors.New(fmt.Sprintln("failed to remove some objects from storage:", errs)),
		)
		return
	}
	w.WriteHeader(http.StatusOK)
}

func (fh FoldersHandler) ReorderFolders(w http.ResponseWriter, r *http.Request) {
	folderID := getFolderID(r)
	from, err := strconv.Atoi(r.FormValue("from"))
	if err != nil {
		return
	}
	to, err := strconv.Atoi(r.FormValue("to"))
	if err != nil {
		return
	}
	fh.DB.ReorderFolders(folderID, from, to)
}

func (fh FoldersHandler) ReorderMedia(w http.ResponseWriter, r *http.Request) {
	folderID := getFolderID(r)
	from, err := strconv.Atoi(r.FormValue("from"))
	if err != nil {
		return
	}
	to, err := strconv.Atoi(r.FormValue("to"))
	if err != nil {
		return
	}
	fh.DB.ReorderMedia(folderID, from, to)
}
