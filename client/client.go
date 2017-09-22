package client

import (
	"encoding/json"
	"fmt"
	"github.com/sequoiia/arma3mission-dl/model"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
)

const DRIVEFILES_ENDPOINT string = "https://clients6.google.com/drive/v2beta/files"
const DRIVEDLFILE_ENDPOINT string = "https://drive.google.com/uc"
const QUERY string = "trashed%20%3D%20false%20and%20%270B8j-xMQtDZvwVjN6R25sWF94dG8%27%20in%20parents"
const KEY string = "AIzaSyC1qbk75NzWBvSaDh6KnsjjA9pIrP4lYIE"
const REFERER string = "https://drive.google.com/drive/folders/0B8j-xMQtDZvwVjN6R25sWF94dG8"

type DlClient struct {
	HttpClient *http.Client
}

func (d *DlClient) GetFiles() *model.GoogleDriveGetFilesResponse {
	u, err := url.Parse(DRIVEFILES_ENDPOINT)
	if err != nil {
		log.Fatal(err)
	}

	values := u.Query()
	values.Set("openDrive", "true")
	values.Set("reason", "102")
	values.Set("syncType", "0")
	values.Set("errorRecovery", "false")
	values.Set("fields", "kind,nextPageToken,items(kind,title,mimeType,createdDate,modifiedDate,modifiedByMeDate,lastViewedByMeDate,fileSize,owners(kind,permissionId,displayName,picture),lastModifyingUser(kind,permissionId,displayName,picture),hasThumbnail,thumbnailVersion,iconLink,id,shared,sharedWithMeDate,userPermission(role),explicitlyTrashed,quotaBytesUsed,shareable,copyable,fileExtension,sharingUser(kind,permissionId,displayName,picture),spaces,editable,version,teamDriveId,hasAugmentedPermissions,trashingUser(kind,permissionId,displayName,picture),trashedDate,parents(id),labels(starred,hidden,trashed,restricted,viewed),capabilities(canCopy,canDownload,canEdit,canAddChildren,canDelete,canRemoveChildren,canShare,canTrash,canRename,canReadTeamDrive,canMoveTeamDriveItem)),incompleteSearch")
	values.Set("appDataFilter", "NO_APP_DATA")
	values.Set("spaces", "drive")
	values.Set("maxResults", "50")
	values.Set("orderBy", "folder,title_natural asc")
	values.Set("key", KEY)
	u.RawQuery = values.Encode()

	u.RawQuery = u.RawQuery + "&q=" + QUERY

	req, err := http.NewRequest("GET", u.String(), nil)

	req.Header.Set("Referer", REFERER)

	resp, err := d.HttpClient.Do(req)

	response := &model.GoogleDriveGetFilesResponse{}

	jdecoder := json.NewDecoder(resp.Body)
	jdecoder.Decode(response)

	return response
}

func (d *DlClient) DownloadFiles(files *model.GoogleDriveGetFilesResponse, path string) {
	for i := 0; i < len(files.Items); i++ {
		file := files.Items[i]
		fmt.Printf("Downloading %s\n", file.Filename)
		err := d.DownloadFile(file.Id, file.Filename, path)
		if err != nil {
			log.Fatal(err)
		}
	}
}

func (d *DlClient) DownloadFile(id string, filename string, path string) error {
	u, err := url.Parse(DRIVEDLFILE_ENDPOINT)
	if err != nil {
		return err
	}

	values := u.Query()
	values.Set("id", id)
	values.Set("authuser", "0")
	values.Set("export", "download")
	u.RawQuery = values.Encode()

	req, err := http.NewRequest("GET", u.String(), nil)
	if err != nil {
		return err
	}

	out, err := os.Create(path + filename)
	if err != nil {
		return err
	}
	defer out.Close()

	resp, err := d.HttpClient.Do(req)
	defer resp.Body.Close()

	_, err = io.Copy(out, resp.Body)
	if err != nil {
		return err
	}

	return nil
}
